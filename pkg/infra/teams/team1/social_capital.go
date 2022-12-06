package team1

import (
	"infra/game/agent"
	"infra/game/decision"
	"infra/game/message"
	"math"
	"sort"
)

// Function which defines how an agent perceives an action
func actionSentiment(action decision.FightAction) [4]float64 {
	switch action {
	case decision.Cower:
		return [4]float64{0.0, 0.0, -0.1, -0.1}
	case decision.Attack:
		return [4]float64{0.0, 0.0, 0.1, 0.1}
	case decision.Defend:
		return [4]float64{0.0, 0.0, 0.1, 0.1}
	default:
		return [4]float64{0.0, 0.0, 0.0, 0.0}
	}
}

// Ensures a float is between -1 and 1
func boundFloat(inputNumber float64) float64 {
	if inputNumber > 1.0 {
		return 1.0
	} else if inputNumber < -1.0 {
		return -1.0
	} else {
		return inputNumber
	}
}

// Ensures array values are between -1 and 1
func boundArray(inputArray [4]float64) [4]float64 {
	return [4]float64{
		boundFloat(inputArray[0]),
		boundFloat(inputArray[1]),
		boundFloat(inputArray[2]),
		boundFloat(inputArray[3]),
	}
}

// Add two arrays
func addArrays(A [4]float64, B [4]float64) [4]float64 {
	return [4]float64{
		A[0] + B[0],
		A[1] + B[1],
		A[2] + B[2],
		A[3] + B[3],
	}
}

func decayNumber(inputNumber float64) float64 {
	if inputNumber < 0 {
		return 0.70 * inputNumber
	} else {
		return 0.90 * inputNumber
	}
}

func decayArray(inputArray [4]float64) [4]float64 {
	return [4]float64{
		decayNumber(inputArray[0]),
		decayNumber(inputArray[1]),
		decayNumber(inputArray[2]),
		decayNumber(inputArray[3]),
	}
}

func softmax(inputArray [3]float64) [3]float64 {
	expValues := [3]float64{
		math.Exp(inputArray[0]),
		math.Exp(inputArray[1]),
		math.Exp(inputArray[2]),
	}

	// Sum exponential array
	sum := 0.0
	for i := 0; i < 3; i++ {
		sum += expValues[i]
	}

	// Divide each element in input array by sum
	for i := 0; i < 3; i++ {
		expValues[i] /= sum
	}

	return expValues
}

func makeIncremental(inputArray [3]float64) [3]float64 {
	var outputArray [3]float64

	outputArray[0] = inputArray[0]

	for i := 1; i < 3; i++ {
		outputArray[i] = outputArray[i-1] + inputArray[i]
	}

	return outputArray
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

	// Update socialCapital values
	for agentID := range s.socialCapital {
		// Decay existing socialCapital values
		s.socialCapital[agentID] = decayArray(s.socialCapital[agentID])

		// Update socialCapital based on agent action
		// TODO: Update of socialCaptial should be dependent on the agents own action (especially for honour)
		action, exists := choices.Get(agentID)
		if exists {
			s.socialCapital[agentID] = addArrays(s.socialCapital[agentID], boundArray(actionSentiment(action)))
		}
	}

	// Ensure all socialCapital values are between -1 and 1
	for key := range s.socialCapital {
		s.socialCapital[key] = boundArray(s.socialCapital[key])
	}
}

/**
 * Tell other trusted agents above a threshold T about the A agents they admire the most,
 * and the H agents they hate the most.
 * Currently messages may be sent to recipients praising or denouncing the recipient
 */
func (r *SocialAgent) sendGossip(agent agent.BaseAgent) {
	selfID := agent.Name()

	sortedSCTrustworthiness := make([]SocialCapInfo, 0, len(r.socialCapital))
	for k, sc := range r.socialCapital {
		if k == selfID { // Exclude self
			continue
		}
		sci := SocialCapInfo{ID: k, arr: sc}
		sortedSCTrustworthiness = append(sortedSCTrustworthiness, sci)
	}

	sort.Slice(sortedSCTrustworthiness, func(i int, j int) bool {
		return sortedSCTrustworthiness[i].arr[2] > sortedSCTrustworthiness[j].arr[2]
	})

	numAdmire := int(r.propAdmire * float64(len(sortedSCTrustworthiness)))
	numHate := int(r.propHate * float64(len(sortedSCTrustworthiness)))

	admiredAgents := make([]string, 0, numAdmire)
	hatedAgents := make([]string, 0, numHate)
	for i, scit := range sortedSCTrustworthiness {
		if i >= numAdmire {
			break
		}

		admiredAgents = append(admiredAgents, scit.ID)
	}

	for i := len(sortedSCTrustworthiness) - 1; i >= 0; i-- {
		if len(sortedSCTrustworthiness)-1-i >= numHate {
			break
		}

		hatedAgents = append(hatedAgents, sortedSCTrustworthiness[i].ID)
	}

	for _, sci := range sortedSCTrustworthiness {
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
