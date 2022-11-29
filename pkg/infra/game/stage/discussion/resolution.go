package discussion

import (
	"github.com/benbjohnson/immutable"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
)

func ResolveFightDiscussion(agentMap *map[commons.ID]agent.Agent, currentLeader agent.Agent, manifesto decision.Manifesto, proposals immutable.Map[commons.ProposalID, uint]) decision.FightResult {
	var fightActions map[commons.ID]decision.FightAction

	if manifesto.FightImposition() {
		currentLeader.Strategy.FightResolution(currentLeader.BaseAgent)
	} else {
		// get proposal with most votes
		var winningProposal immutable.Map[commons.ID, decision.FightAction]
		for id, a := range *agentMap {
			if val, ok := winningProposal.Get(id); ok {
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
