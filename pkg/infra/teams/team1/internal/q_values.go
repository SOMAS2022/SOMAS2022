/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/
package internal

// Estimate reward for agent actions
func CooperationQ(state QState) [3]float64 {
	reward := QFunction(state, true)
	// fmt.Println(reward)
	return reward
}

func SelfishQ(state QState) [3]float64 {
	return QFunction(state, false)
}
