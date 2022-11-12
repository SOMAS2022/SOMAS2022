package fight

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"math"

	"github.com/benbjohnson/immutable"
)

func DealDamage(attack uint, agentMap map[string]agent.Agent, globalState *state.State) {
	splitDamage := attack / uint(len(agentMap))
	for id, agentState := range globalState.AgentState {
		newHp := commons.SaturatingSub(agentState.Hp, splitDamage)
		if newHp == 0 {
			// kill agent
			// todo: prune peer channels somehow...
			delete(globalState.AgentState, id)
			delete(agentMap, id)
		} else {
			globalState.AgentState[id] = state.AgentState{
				Hp:           newHp,
				Attack:       agentState.Attack,
				Defense:      agentState.Defense,
				BonusAttack:  agentState.BonusAttack,
				BonusDefense: agentState.BonusDefense,
			}
		}
	}
}

func AgentFightDecisions(state *state.View, agents map[string]agent.Agent, previousDecisions *immutable.Map[string, decision.FightAction]) map[string]decision.FightAction {
	decisionMap := make(map[string]decision.FightAction)
	channels := make(map[string]chan decision.FightAction)

	for i, a := range agents {
		channels[i] = startAgentFightHandlers(state, a, previousDecisions)
	}
	for i, dChan := range channels {
		decisionMap[i] = <-dChan
		close(dChan)
	}

	return decisionMap
}

func HandleFightRound(state *state.State, baseHealth uint, decisionMap map[string]decision.FightAction) (uint, uint, uint) {

	var coweringAgents uint
	var attackSum uint
	var shieldSum uint

	for agentID, d := range decisionMap {
		agentState := state.AgentState[agentID]

		switch d {
		case decision.Attack:
			if agentState.Stamina > agentState.BonusAttack {
				attackSum += agentState.TotalAttack()
				agentState.Stamina = commons.SaturatingSub(agentState.Stamina, agentState.BonusAttack)
			} else {
				coweringAgents++
				decisionMap[agentID] = decision.Cower
				agentState.Hp += uint(math.Ceil(0.05 * float64(baseHealth)))
				agentState.Stamina += 1
			}
		case decision.Defend:
			if agentState.Stamina > agentState.BonusDefense {
				shieldSum += agentState.TotalDefense()
				agentState.Stamina = commons.SaturatingSub(agentState.Stamina, agentState.BonusDefense)
			} else {
				coweringAgents++
				decisionMap[agentID] = decision.Cower
				agentState.Hp += uint(math.Ceil(0.05 * float64(baseHealth)))
				agentState.Stamina += 1
			}
		case decision.Cower:
			coweringAgents++
			agentState.Hp += uint(math.Ceil(0.05 * float64(baseHealth)))
			agentState.Stamina += 1
		}
		state.AgentState[agentID] = agentState
	}

	return coweringAgents, attackSum, shieldSum
}

func startAgentFightHandlers(view *state.View, a agent.Agent, decisionLog *immutable.Map[string, decision.FightAction]) chan decision.FightAction {
	decisionChan := make(chan decision.FightAction)
	go a.Strategy.HandleFight(view, a.BaseAgent, decisionChan, decisionLog)
	return decisionChan
}
