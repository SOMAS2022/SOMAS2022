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

type GossipInform struct {
	mtype  int
	agents []string
}

func (mc *GossipInform) sealedMessage() {
}

func (mc *GossipInform) sealedInform() {
}

// func newGossipInform(mtype int, agents []string) message.Message {
// 	return GossipInform{mtype: mtype, agents: agents}
//}

func Gossip(BA agent.BaseAgent, recepient string, mtype int, about []string) {

	// msg := MessageContent{mtype: mtype, agents: about}
	// message.NewFightProposalMessage()

	// BA.SendBlockingMessage(recepient, newGossipInform(mtype, about))

}
