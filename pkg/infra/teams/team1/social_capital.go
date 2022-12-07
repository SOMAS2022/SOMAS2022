package team1

import (
	"infra/game/agent"
	"infra/game/decision"
)

// Called any time a message is received, initialises or updates the socialCapital map
func (s *SocialAgent) updateSocialCapital(self agent.BaseAgent, fightDecisions decision.ImmutableFightResult) {
	// For some reason had to split .Choices() and .Len() for Golang not to complain
	choices := fightDecisions.Choices()

	// Initialize socialCapital map if it hasn't already
	if len(s.socialCapital) == 0 && choices.Len() > 1 {
		// Create empty map
		s.socialCapital = map[string][4]float64{}

		// Populate map with every currently living agent
		itr := choices.Iterator()
		for !itr.Done() {
			agentID, _, _ := itr.Next()

			s.socialCapital[agentID] = [4]float64{0.0, 0.0, 0.0, 0.0}
		}

		// Delete the agents own id from the socialCapital array
		delete(s.socialCapital, self.ID())
	}

	// Extract agentState from base agent
	view := self.View()
	agentState := view.AgentState()

	// Calculate how cooperative agents own action was
	cooperativeQ := cooperationQ(self.AgentState())
	cooperationScale := normalise(cooperativeQ)
	selfAction, _ := choices.Get(self.ID())
	selfCooperation := cooperationScale[int(selfAction)]

	// Update socialCapital values
	for agentID := range s.socialCapital {
		// Decay existing socialCapital values
		s.socialCapital[agentID] = decayArray(s.socialCapital[agentID])

		// If agent did an action, update socialCapital based on action
		action, exists := choices.Get(agentID)
		if exists {
			// Get hidden state of agent
			otherAgentState, _ := agentState.Get(agentID)

			// Calculate how cooperative each action is in other agents current state
			cooperativeQ := hiddenCooperationQ(otherAgentState)

			// Put actions on linear scale from -1 (least cooperative) to 1 (most cooperative)
			cooperationScale := normalise(cooperativeQ)

			// Calculate update of trustworthiness based on how cooperative action was
			deltaTrust := 0.1 * cooperationScale[int(action)]

			// Calculate update of based on how cooperative action was compared to the agents own action
			deltaHonour := 0.1 * (cooperationScale[int(action)] - selfCooperation)

			// Update the socialCapital array based on calculated delta for trustworthiness and honour
			s.socialCapital[agentID] = boundArray(addArrays(s.socialCapital[agentID], [4]float64{0.0, 0.0, deltaTrust, deltaHonour}))
		}
	}
}
