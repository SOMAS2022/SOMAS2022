package team1

import (
	"infra/game/agent"
	"infra/game/decision"
	"infra/game/message"
	"sort"
)

// Called by own InitAgent function when running team experiment, before game starts
func (r *SocialAgent) initSocialCapital(selfID string, allAgents []string) {
	// Create empty map
	r.socialCapital = map[string][4]float64{}
	for _, s := range allAgents {
		r.socialCapital[s] = [4]float64{0.0, 0.0, 0.0, 0.0}
	}

	// Delete the agents own id from the socialCapital array
	delete(r.socialCapital, selfID)
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

/**
 * Tell other trusted agents above a threshold T about the A agents they admire the most,
 * and the H agents they hate the most.
 * Currently messages may be sent to recipients praising or denouncing the recipient
 */
func (r *SocialAgent) sendGossip(agent agent.BaseAgent) {
	// This type will make it easier to extract from map, sort, and retrieve agent ID
	type SocialCapInfo struct {
		ID  string
		arr [4]float64
	}
	selfID := agent.ID()

	sortedSCTrustHonor := make([]SocialCapInfo, 0, len(r.socialCapital))
	for k, sc := range r.socialCapital {
		if k == selfID { // Exclude self
			continue
		}
		sci := SocialCapInfo{ID: k, arr: sc}
		sortedSCTrustHonor = append(sortedSCTrustHonor, sci)
	}

	sort.Slice(sortedSCTrustHonor, func(i int, j int) bool {
		return (sortedSCTrustHonor[i].arr[2] + sortedSCTrustHonor[i].arr[3]) > (sortedSCTrustHonor[j].arr[2] + sortedSCTrustHonor[j].arr[3])
	})

	numAdmire := int(r.propAdmire * float64(len(sortedSCTrustHonor)))
	numHate := int(r.propHate * float64(len(sortedSCTrustHonor)))

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
		if sci.arr[1] < r.gossipThreshold {
			continue
		}
		Gossip(agent, sci.ID, MessagePraise, admiredAgents)
		Gossip(agent, sci.ID, MessageDenounce, hatedAgents)
	}
}

/**
 * On receiving gossip, scale the network value up/down by a constant and
 * the senders overall perception to this agent. This means that someone
 * with 0.2 network would become 0.22, and someone with 0.7 network would
 * become 0.77, with a 10% increase
 *
 */
func (r *SocialAgent) receiveGossip(m message.ArrayInfo, sender string) {
	// Will reverse if sender's perception is negative
	senderPerception := OverallPerception(r.socialCapital[sender])
	mtype := m.GetNum()
	var sign float64
	switch mtype {
	case MessagePraise:
		sign = 1.0
	case MessageDenounce:
		sign = -1.0
	}

	for _, about := range m.GetStringArr() {
		sc := r.socialCapital[about]
		sc[1] += sign * senderPerception * 0.1 * sc[1]
		sc = boundArray(sc)
		r.socialCapital[about] = sc
	}
}

func OverallPerception(inputArray [4]float64) float64 {
	return (inputArray[0] + inputArray[1] + inputArray[2] + inputArray[3]) * 0.25
}
