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
	"infra/game/state"
)

// Hidden: Estimate reward of other player for action they took
func HiddenCooperationQ(state state.HiddenAgentState) [3]float64 {
	return QFunction(getQStateOther(state), true)
}

func HiddenSelfishQ(state state.HiddenAgentState) [3]float64 {
	return QFunction(getQStateOther(state), false)
}

// Estimate reward for agent actions
func CooperationQ(state state.AgentState) [3]float64 {
	reward := QFunction(getQState(state), true)
	// fmt.Println(reward)
	return reward
}

func SelfishQ(state state.AgentState) [3]float64 {
	return QFunction(getQState(state), false)
}
