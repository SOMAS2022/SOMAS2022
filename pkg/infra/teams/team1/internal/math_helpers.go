/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/
package internal

import (
	"infra/game/agent"
	"infra/game/state"
	"math"
)

func Argmax(array []float64) int {
	maxIndex := 0

	for index, element := range array {
		if element > array[maxIndex] {
			maxIndex = index
		}
	}

	return maxIndex
}

func Argmin(array []float64) int {
	minIndex := 0

	for index, element := range array {
		if element < array[minIndex] {
			minIndex = index
		}
	}

	return minIndex
}

func Max(array []float64) float64 {
	if len(array) == 0 {
		return 0
	}

	max := array[0]

	for _, element := range array {
		if element > max {
			max = element
		}
	}

	return max
}

func Min(array []float64) float64 {
	if len(array) == 0 {
		return 0
	}

	min := array[0]

	for _, element := range array {
		if element < min {
			min = element
		}
	}

	return min
}

// Ensures a float is between -1 and 1
func BoundFloat(inputNumber float64) float64 {
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
		BoundFloat(inputArray[0]),
		BoundFloat(inputArray[1]),
		BoundFloat(inputArray[2]),
		BoundFloat(inputArray[3]),
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

func DecayNumber(inputNumber float64) float64 {
	if inputNumber < 0 {
		return 0.70 * inputNumber
	} else {
		return 0.90 * inputNumber
	}
}

func DecayArray(inputArray [4]float64) [4]float64 {
	return [4]float64{
		DecayNumber(inputArray[0]),
		DecayNumber(inputArray[1]),
		DecayNumber(inputArray[2]),
		DecayNumber(inputArray[3]),
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

// Make lowest value -1, highest 1 and everything else interpolation between
func Normalise(array [3]float64) [3]float64 {
	max := Max(array[:])
	min := Min(array[:])

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

func QStateToArray(state QState) [8]float64 {
	return [8]float64{
		1,
		state.Hp,
		state.Stamina,
		state.TotalAttack,
		state.TotalDefense,
		state.LevelsToWin,
		state.MonsterHealth,
		state.MonsterAttack,
	}
}

func BaseAgentToQState(agent agent.BaseAgent) QState {
	// Get agentState from baseAgent
	agentState := agent.AgentState()

	// Get view from baseAgent
	view := agent.View()

	return QState{
		float64(agentState.Hp),
		float64(agentState.Stamina),
		float64(agentState.Attack),
		float64(agentState.Defense),
		float64(view.CurrentLevel()),
		float64(view.MonsterHealth()),
		float64(view.MonsterAttack()),
	}
}

func HiddenAgentToQState(agent state.HiddenAgentState, view state.View) QState {
	return QState{
		float64(agent.Hp),
		float64(agent.Stamina),
		float64(agent.Attack),
		float64(agent.Defense),
		float64(view.CurrentLevel()),
		float64(view.MonsterHealth()),
		float64(view.MonsterAttack()),
	}
}

func ConstMulSlice(constant float64, input []float64) []float64 {
	multipliedSlice := make([]float64, len(input))

	for idx, element := range input {
		multipliedSlice[idx] = constant * element
	}

	return multipliedSlice
}

func AddSlices(a []float64, b []float64) []float64 {
	if len(a) != len(b) {
		panic("Slices are not of same size")
	}

	for idx, element := range b {
		a[idx] += element
	}

	return a
}
