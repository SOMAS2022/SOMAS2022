package team1

import (
	"infra/game/agent"
	"infra/game/decision"
	"math"
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
