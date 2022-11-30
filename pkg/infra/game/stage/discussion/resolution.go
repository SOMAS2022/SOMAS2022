package discussion

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/tally"
)

func ResolveFightDiscussion(agentMap *map[commons.ID]agent.Agent, currentLeader agent.Agent, manifesto decision.Manifesto, tally *tally.Tally[decision.FightAction]) decision.FightResult {
	var fightActions map[commons.ID]decision.FightAction

	if manifesto.FightImposition() {
		currentLeader.Strategy.FightResolution(currentLeader.BaseAgent)
	} else {
		// get proposal with most votes
		winningProp := tally.GetMax().Proposal()
		for id, a := range *agentMap {
			if val, ok := winningProp.Get(id); ok {
				fightActions[id] = val
			} else {
				fightActions[id] = a.Strategy.CurrentAction()
			}
		}
	}

	return decision.FightResult{
		Choices:         fightActions,
		AttackingAgents: nil,
		ShieldingAgents: nil,
		CoweringAgents:  nil,
		AttackSum:       0,
		ShieldSum:       0,
	}
}
