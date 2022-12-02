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
	"infra/game/message"

	"github.com/google/uuid"
)

// This type will make it easier to extract from map, sort, and retrieve agent ID

const (
	Praise message.Type = 2048
	Denounce
)

type MessageContent struct {
	agents []string
}

func (mc MessageContent) isPayload() {
}

func Gossip(from string, recepients []string, mtype message.Type, about []string) {

	msg := message.NewMessage(mtype, MessageContent{agents: about})

	for _, to := range recepients {
		mId, _ := uuid.NewUUID()
		message.NewTaggedMessage(from, msg, mId)
	}

}
