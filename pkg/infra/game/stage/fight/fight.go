package fight

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

func DealDamage(attack uint, agentMap map[uint]agent.Agent, globalState state.State) {
	splitDamage := attack / uint(len(agentMap))
	for id, agent := range agentMap {
		newHp := commons.SaturatingSub(agent.State.Hp, splitDamage)
		if newHp == 0 {
			// kill agent
			// todo: prune peer channels somehow...
			delete(agentMap, id)
		} else {
			agent.State.Hp = newHp
			agentMap[id] = agent
		}
	}
}

func HandleFightRound(state *state.State, agents map[uint]agent.Agent, prevDecisions map[uint]decision.FightDecision) (uint, uint, uint, map[uint]decision.FightDecision) {

	for _, a := range agents {
		go a.Strategy.HandleFight(*state, a, prevDecisions)
	}
	decisions := make(map[uint]decision.FightDecision)

	for agentID, agent := range agents {
		handleFightDecision(agent.FightDecisionChannel, decisions, agentID, agent.State)
	}
	var coweringAgents uint
	var attackSum uint
	var shieldSum uint

	for agentID, d := range decisions {
		if d.Cower {
			coweringAgents++
			if entry, ok := agents[agentID]; ok {
				entry.State.Hp += 1
				agents[agentID] = entry
			}
		} else {
			attackSum += d.Attack
			shieldSum += d.Defend
		}
	}

	return coweringAgents, attackSum, shieldSum, decisions
}

func handleFightDecision(decisionC chan decision.FightDecision, decisions map[uint]decision.FightDecision, agentID uint, s state.AgentState) {
	receivedDecision := <-decisionC
	receivedDecision.ValidateDecision(s)
	decisions[agentID] = receivedDecision
}
