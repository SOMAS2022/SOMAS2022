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
	"infra/game/decision"
)

func CollaborativeFightDecision() decision.FightAction {
	// TODO: Select from collaborative Q-Table
	return decision.Attack
}

func SelfishFightDecision() decision.FightAction {
	// TODO: Select from selfish Q-Table
	return decision.Cower
}
