package team1

import (
	"infra/game/agent"
	"infra/game/decision"
	"math"
)

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

// Make lowest value -1, highest 1 and everything else interpolation between
func normalise(array [3]float64) [3]float64 {
	max := max(array[:])
	min := min(array[:])

	var normArray [3]float64

	for index, value := range array {
		if value == max {
			normArray[index] = 1.0
		} else if value == min {
			normArray[index] = -1.0
		} else {
			// Interpolate between -1 and 1
			normArray[index] = (value - (max+min)/2) / (max - min)
		}
	}

	return normArray
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
