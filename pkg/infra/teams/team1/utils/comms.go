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
	"infra/game/agent"
	"infra/game/message"
)

const (
	Praise int = iota
	Denounce
)

func Gossip(BA agent.BaseAgent, recipients string, mtype int, about []string) {
	msg := message.ArrayInfo{Num: mtype, StringArr: about}
	BA.SendBlockingMessage(recipients, msg)
}
