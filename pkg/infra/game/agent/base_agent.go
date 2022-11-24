package agent

import (
	"fmt"
	"infra/game/commons"
	"infra/game/message"
	"infra/logging"
)

type BaseAgent struct {
	communication Communication
	Id            commons.ID
	AgentName     string
}

func NewBaseAgent(communication Communication, id commons.ID, agentName string) BaseAgent {
	return BaseAgent{communication: communication, Id: id, AgentName: agentName}
}

func (ba *BaseAgent) BroadcastBlockingMessage(m message.Message) {
	iterator := ba.communication.peer.Iterator()
	tm := message.TaggedMessage{
		Sender:  ba.Id,
		Message: m,
	}
	for !iterator.Done() {
		_, c, ok := iterator.Next()
		if ok {
			c <- tm
		}
	}
}

func (ba *BaseAgent) SendBlockingMessage(id commons.ID, m message.Message) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("agent %s not available for messaging, submitted", id)
		}
	}()

	value, ok := ba.communication.peer.Get(id)
	if ok {
		value <- message.TaggedMessage{
			Sender:  ba.Id,
			Message: m,
		}
	} else {
		e = fmt.Errorf("agent %s not available for messaging, dead", id)
	}
	return
}

func (ba *BaseAgent) Log(lvl logging.Level, fields logging.LogField, msg string) {
	agentFields := logging.LogField{
		"agentName": ba.AgentName,
		"agentID":   ba.Id,
	}

	logging.Log(lvl, logging.CombineFields(agentFields, fields), msg)
}
