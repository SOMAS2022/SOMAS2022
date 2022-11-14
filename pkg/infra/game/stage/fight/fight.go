package fight

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"math"
	"sync"

	"github.com/benbjohnson/immutable"
)

func DealDamage(attack uint, agentMap map[commons.ID]agent.Agent, globalState *state.State) {
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

func AgentFightDecisions(state *state.View, agents map[commons.ID]agent.Agent, previousDecisions *immutable.Map[commons.ID, decision.FightAction], channelsMap map[commons.ID]chan message.TaggedMessage) map[string]decision.FightAction {
	decisionMap := make(map[commons.ID]decision.FightAction)
	channel := make(chan message.ActionMessage, 100)

	var wg sync.WaitGroup

	for _, a := range agents {
		a := a
		wg.Add(1)
		startAgentFightHandlers(*state, &a, previousDecisions, channel, &wg)
	}

	for _, messages := range channelsMap {
		messages <- message.TaggedMessage{
			Sender:  "server",
			Message: *message.NewMessage(message.Something, nil),
		}
	}

	go func(group *sync.WaitGroup) {
		group.Wait()
		for _, messages := range channelsMap {
			close(messages)
		}
		close(channel)
	}(&wg)

	for actionMessage := range channel {
		decisionMap[actionMessage.Sender] = actionMessage.Action
	}

	wg.Wait()

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

func startAgentFightHandlers(view state.View, a *agent.Agent, decisionLog *immutable.Map[commons.ID, decision.FightAction], channel chan message.ActionMessage, wg *sync.WaitGroup) {
	go a.HandleFight(view, *decisionLog, channel, wg)
}
