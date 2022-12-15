/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/
package internal

// Deprecated: To remove

// Estimate reward for agent actions
func CooperationQ(state QState, strat [3]ActionStrategy) [3]float64 {
	reward := QFunction(state, strat)
	// fmt.Println(reward)
	return reward
}

func SelfishQ(state QState, strat [3]ActionStrategy) [3]float64 {
	return QFunction(state, strat)
}
