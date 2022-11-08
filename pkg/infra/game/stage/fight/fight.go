package fight

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

func DealDamage(attack uint, agentMap map[uint]agent.Agent, globalState state.State, stateChannels map[uint]chan<- state.State, decisionChannels map[uint]<-chan decision.Decision) {
	splitDamage := attack / uint(len(agentMap))
	for id, agentState := range globalState.AgentState {
		newHp := commons.SaturatingSub(agentState.Hp, splitDamage)
		if newHp == 0 {
			// kill agent
			// todo: prune peer channels somehow...
			delete(globalState.AgentState, id)
			delete(agentMap, id)
			delete(stateChannels, id)
			delete(decisionChannels, id)
		} else {
			globalState.AgentState[id] = state.AgentState{
				Hp:            newHp,
				Attack:        agentState.Attack,
				Defense:       agentState.Defense,
				AbilityPoints: agentState.AbilityPoints,
				BonusAttack:   agentState.BonusAttack,
				BonusDefense:  agentState.BonusDefense,
			}
		}
	}
}

func HandleFightRound(state *state.State, agents map[uint]agent.Agent, decisionChannels map[uint]<-chan decision.Decision) (uint, uint, uint) {
	for _, a := range agents {
		go a.Strategy.HandleFight(*state, a.BaseAgent)
	}
	decisions := make(map[uint]decision.FightAction)

	for agentID, decisionC := range decisionChannels {
		handleFightDecision(decisionC, decisions, agentID, state.AgentState[agentID])
	}

	var coweringAgents uint
	var attackSum uint
	var shieldSum uint

	for agentID, d := range decisions {
		switch v := d.(type) {
		case decision.Fight:
			attackSum += v.Attack
			shieldSum += v.Defend
		case decision.Cower:
			coweringAgents++
			if entry, ok := state.AgentState[agentID]; ok {
				entry.Hp += 1
				state.AgentState[agentID] = entry
			}
		}
	}

	return coweringAgents, attackSum, shieldSum
}

func handleFightDecision(decisionC <-chan decision.Decision, decisions map[uint]decision.FightAction, agentID uint, s state.AgentState) {
	for {
		received := <-decisionC
		switch d := received.(type) {
		case decision.FightDecision:
			switch f := d.Action.(type) {
			case decision.Fight:
				if f.Attack <= s.TotalAttack() && f.Defend <= s.TotalDefense() && f.Attack+f.Defend <= s.AbilityPoints {
					decisions[agentID] = f
				} else {
					decisions[agentID] = decision.Cower{}
				}
			case decision.Cower:
				decisions[agentID] = f
			}
			return
		default:
			continue
		}
	}
}
