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

func HiddenCooperationQ(state state.HiddenAgentState) [3]float64 {
	// TODO: Implement with q-function
	return [3]float64{0.0, 0.0, 0.0}
}

func HiddenSelfishQ(state state.HiddenAgentState) [3]float64 {
	// TODO: Implement with q-function
	return [3]float64{0.0, 0.0, 0.0}
}

func CooperationQ(state state.AgentState) [3]float64 {
	// TODO: Implement with q-function
	return [3]float64{0.0, 0.0, 0.0}
}

func SelfishQ(state state.AgentState) [3]float64 {
	// TODO: Implement with q-function
	return [3]float64{0.0, 0.0, 0.0}
}
