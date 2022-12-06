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
	"infra/game/message"
)

const (
	MessagePraise int = iota
	MessageDenounce
)

func Gossip(BA agent.BaseAgent, recipients string, mtype int, about []string) {
	message := message.NewArrayInfo(mtype, about)
	BA.SendBlockingMessage(recipients, message)
}
