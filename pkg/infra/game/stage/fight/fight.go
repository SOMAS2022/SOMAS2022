package fight

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"math"
)

func DealDamage(attack uint, agentMap map[uint]agent.Agent, globalState state.State) {
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

func HandleFightRound(state *state.State, agents map[uint]agent.Agent, baseHealth uint) (uint, uint, uint) {
	decisionMap := make(map[uint]decision.FightAction)
	for i, a := range agents {
		processAgentDecision(*state, a, decisionMap, i)
	}
	var coweringAgents uint
	var attackSum uint
	var shieldSum uint

	for agentID, d := range decisionMap {
		switch v := d.(type) {
		case decision.Fight:
			attackSum += v.Attack
			shieldSum += v.Defend
		case decision.Cower:
			coweringAgents++
			if entry, ok := state.AgentState[agentID]; ok {
				entry.Hp += uint(math.Ceil(0.05 * float64(baseHealth)))
				state.AgentState[agentID] = entry
			}
		}
	}

	return coweringAgents, attackSum, shieldSum
}

func processAgentDecision(state state.State, a agent.Agent, decisionMap map[uint]decision.FightAction, i uint) {
	decisionChan := make(chan decision.FightAction)
	defer close(decisionChan)
	go a.Strategy.HandleFight(state, a.BaseAgent, decisionChan)
	decisionMap[i] = (<-decisionChan).HandleAction(state.AgentState[i])
}
