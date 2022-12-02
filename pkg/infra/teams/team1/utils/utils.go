/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/

package utils

import (
	"infra/config"
	"infra/game/decision"
	"infra/game/state"
	"math"
)

var Config config.GameConfig

/**
 * This function converts agents stats into a single 'usefullness-to-battle' value
 *
 * Calculated based on a(curr)
 */
func AgentBattleUtility(agentState state.AgentState) float64 {
	attackPortion := float64(agentState.TotalAttack() / Config.StartingAttackStrength)
	defensePortion := float64(agentState.TotalDefense() / Config.StartingShieldStrength)
	healthPortion := float64(agentState.Hp / Config.StartingHealthPoints)
	staminaPortion := float64(agentState.Stamina / Config.Stamina)

	return 0.25*attackPortion + 0.25*defensePortion + 0.25*healthPortion + 0.25*staminaPortion
}

// Function which defines how an agent perceives an action
func ActionSentiment(action decision.FightAction) [4]float64 {
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
func BoundArray(inputArray [4]float64) [4]float64 {
	return [4]float64{
		boundFloat(inputArray[0]),
		boundFloat(inputArray[1]),
		boundFloat(inputArray[2]),
		boundFloat(inputArray[3]),
	}
}

// Add two arrays
func AddArrays(A [4]float64, B [4]float64) [4]float64 {
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

func DecayArray(inputArray [4]float64) [4]float64 {
	return [4]float64{
		decayNumber(inputArray[0]),
		decayNumber(inputArray[1]),
		decayNumber(inputArray[2]),
		decayNumber(inputArray[3]),
	}
}

func Softmax(inputArray [3]float64) [3]float64 {
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

func MakeIncremental(inputArray [3]float64) [3]float64 {
	var outputArray [3]float64

	outputArray[0] = inputArray[0]

	for i := 1; i < 3; i++ {
		outputArray[i] = outputArray[i-1] + inputArray[i]
	}

	return outputArray
}
