package discussion

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message/proposal"
	"infra/game/state"
	"infra/game/tally"
)

func ResolveFightDiscussion(gs state.State, agentMap map[commons.ID]agent.Agent, currentLeader agent.Agent, manifesto decision.Manifesto, tally *tally.Tally[decision.FightAction]) decision.FightResult {
	fightActions := make(map[commons.ID]decision.FightAction)
	// todo: cleanup the nil check that acts to check if the leader died in combat
	if manifesto.FightImposition() && currentLeader.Strategy != nil {
		resolution := currentLeader.Strategy.FightResolution(*currentLeader.BaseAgent).Rules()
		predicate := proposal.ToPredicate(resolution)
		for id, a := range agentMap {
			fightActions[id] = predicate(gs, a.AgentState())
		}
	} else {
		// get proposal with most votes
		winningProp := tally.GetMax().Rules()
		predicate := proposal.ToPredicate(winningProp)
		for id, a := range agentMap {
			fightActions[id] = predicate(gs, a.AgentState())
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
