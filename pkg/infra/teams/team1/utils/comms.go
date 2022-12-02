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
)

const (
	Praise int = 2048
	Denounce
)

type MessageContent struct {
	mtype  int
	agents []string
}

func (mc *MessageContent) sealedMessage() {
}

func Gossip(BA agent.BaseAgent, recepient string, mtype int, about []string) {

	msg := MessageContent{mtype: mtype, agents: about}

	BA.SendBlockingMessage(recepient, msg)

}
