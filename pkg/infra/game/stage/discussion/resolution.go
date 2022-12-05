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
	var prop commons.ImmutableList[proposal.Rule[decision.FightAction]]
	if manifesto.FightImposition() && currentLeader.Strategy != nil {
		prop = currentLeader.Strategy.FightResolution(*currentLeader.BaseAgent)
	} else {
		// get proposal with most votes
		prop = tally.GetMax().Rules()
	}

	predicate := proposal.ToPredicate(prop)
	if predicate == nil {
		for id, a := range agentMap {
			fightActions[id] = a.CurrentAction()
		}
	} else {
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
