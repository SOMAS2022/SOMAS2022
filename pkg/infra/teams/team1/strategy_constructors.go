/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/

package team1

import (
	"infra/game/agent"
)

// Demonstrate creating a strategy with input parameters
func CreateAggressiveAgent() agent.Strategy {
	return NewProbabilisticAgent(0.1, 0.8, 0.1)
}

func CreateDefensiveAgent() agent.Strategy {
	return NewProbabilisticAgent(0.1, 0.1, 0.8)
}
