package team1

import (
	"infra/game/agent"
	"infra/game/decision"
	"infra/game/message"
	"infra/teams/team1/internal"
	"sort"
)

// Called by own InitAgent function when running team experiment, before game starts
func (s *SocialAgent) initSocialCapital(allAgents []string) {
	// Create empty map
	s.socialCapital = map[string][4]float64{}
	for _, id := range allAgents {
		s.socialCapital[id] = [4]float64{0.0, 0.0, 0.0, 0.0}
	}
}

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
	}

	// Extract agentState from base agent
	view := self.View()
	agentState := view.AgentState()

	// Calculate how cooperative agents own action was
	cooperativeQ := internal.CooperationQ(internal.BaseAgentToQState(self))
	cooperationScale := internal.Normalise(cooperativeQ)
	selfAction, _ := choices.Get(self.ID())
	selfCooperation := cooperationScale[int(selfAction)]

	// Update socialCapital values
	for agentID := range s.socialCapital {
		// Decay existing socialCapital values
		s.socialCapital[agentID] = internal.DecayArray(s.socialCapital[agentID])

		// If agent did an action, update socialCapital based on action
		action, exists := choices.Get(agentID)
		if exists {
			// Get hidden state of agent
			otherAgentState, _ := agentState.Get(agentID)

			// Calculate how cooperative each action is in other agents current state
			cooperativeQ := internal.CooperationQ(internal.HiddenAgentToQState(otherAgentState, view))

			// Put actions on linear scale from -1 (least cooperative) to 1 (most cooperative)
			cooperationScale := internal.Normalise(cooperativeQ)

			// Calculate update of trustworthiness based on how cooperative action was
			deltaTrust := 0.1 * cooperationScale[int(action)]

			// Calculate update of based on how cooperative action was compared to the agents own action
			deltaHonour := 0.1 * (cooperationScale[int(action)] - selfCooperation)

			// Update the socialCapital array based on calculated delta for trustworthiness and honour
			s.socialCapital[agentID] = internal.BoundArray(internal.AddArrays(s.socialCapital[agentID], [4]float64{0.0, 0.0, deltaTrust, deltaHonour}))
		}
	}
}

/**
 * Tell other trusted agents above a threshold T about the A agents they admire the most,
 * and the H agents they hate the most.
 * Currently messages may be sent to recipients praising or denouncing the recipient
 */
func (s *SocialAgent) sendGossip(agent agent.BaseAgent) {
	// This type will make it easier to extract from map, sort, and retrieve agent ID
	type SocialCapInfo struct {
		ID  string
		arr [4]float64
	}
	selfID := agent.ID()

	sortedSCTrustHonor := make([]SocialCapInfo, 0, len(s.socialCapital))
	for k, sc := range s.socialCapital {
		if k == selfID { // Exclude self
			continue
		}
		sci := SocialCapInfo{ID: k, arr: sc}
		sortedSCTrustHonor = append(sortedSCTrustHonor, sci)
	}

	sort.Slice(sortedSCTrustHonor, func(i int, j int) bool {
		return (sortedSCTrustHonor[i].arr[2] + sortedSCTrustHonor[i].arr[3]) > (sortedSCTrustHonor[j].arr[2] + sortedSCTrustHonor[j].arr[3])
	})

	numAdmire := int(s.propAdmire * float64(len(sortedSCTrustHonor)))
	numHate := int(s.propHate * float64(len(sortedSCTrustHonor)))

	admiredAgents := make([]string, 0, numAdmire)
	hatedAgents := make([]string, 0, numHate)
	for i, scit := range sortedSCTrustHonor {
		if i >= numAdmire {
			break
		}

		admiredAgents = append(admiredAgents, scit.ID)
	}

	for i := len(sortedSCTrustHonor) - 1; i >= 0; i-- {
		if len(sortedSCTrustHonor)-1-i >= numHate {
			break
		}

		hatedAgents = append(hatedAgents, sortedSCTrustHonor[i].ID)
	}

	for _, sci := range sortedSCTrustHonor {
		if sci.arr[1] < s.gossipThreshold {
			continue
		}
		internal.Gossip(agent, sci.ID, internal.MessagePraise, admiredAgents)
		internal.Gossip(agent, sci.ID, internal.MessageDenounce, hatedAgents)
	}
}

/**
 * On receiving gossip, scale the network value up/down by a constant and
 * the senders overall perception to this agent. This means that someone
 * with 0.2 network would become 0.22, and someone with 0.7 network would
 * become 0.77, with a 10% increase
 *
 */
func (s *SocialAgent) receiveGossip(m message.ArrayInfo, sender string) {
	// Will reverse if sender's perception is negative
	senderPerception := OverallPerception(s.socialCapital[sender])
	mtype := m.GetNum()
	var sign float64
	switch mtype {
	case internal.MessagePraise:
		sign = 1.0
	case internal.MessageDenounce:
		sign = -1.0
	}

	for _, about := range m.GetStringArr() {
		sc := s.socialCapital[about]
		sc[1] += sign * senderPerception * 0.1 * sc[1]
		sc = internal.BoundArray(sc)
		s.socialCapital[about] = sc
	}
}

func OverallPerception(inputArray [4]float64) float64 {
	return (inputArray[0] + inputArray[1] + inputArray[2] + inputArray[3]) * 0.25
}
