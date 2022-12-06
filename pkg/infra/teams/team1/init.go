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
	"infra/game/commons"
)

var InitAgentMap = map[commons.ID]func() agent.Strategy{
	"CollaborativeAgent": CreateCollaborativeAgent,
	"SelfishAgent":       CreateSelfishAgent,
}
