package agent

import (
	"fmt"
	"infra/game/commons"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"

	"github.com/google/uuid"
)

type BaseAgent struct {
	communication *Communication
	id            commons.ID
	name          string
	latestState   state.AgentState
}

func (ba *BaseAgent) Id() commons.ID {
	return ba.id
}

func (ba *BaseAgent) Name() string {
	return ba.name
}

func NewBaseAgent(communication *Communication, id commons.ID, agentName string) BaseAgent {
	return BaseAgent{communication: communication, id: id, name: agentName}
}

func (ba *BaseAgent) BroadcastBlockingMessage(m message.Message) {
	iterator := ba.communication.peer.Iterator()
	mId, _ := uuid.NewUUID()
	tm := message.NewTaggedMessage(ba.id, m, mId)
	//	sender:  ba.id,
	//	message: m,
	//}
	for !iterator.Done() {
		_, c, ok := iterator.Next()
		if ok {
			c <- *tm
		}
	}
}

func (ba *BaseAgent) SendBlockingMessage(id commons.ID, m message.Message) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("agent %s not available for messaging, submitted", id)
		}
	}()

	channel, ok := ba.communication.peer.Get(id)
	if ok {
		mId, _ := uuid.NewUUID()
		channel <- *message.NewTaggedMessage(ba.id, m, mId)
	} else {
		e = fmt.Errorf("agent %s not available for messaging, dead", id)
	}
	return
}

func (ba *BaseAgent) Log(lvl logging.Level, fields logging.LogField, msg string) {
	agentFields := logging.LogField{
		"agentName": ba.name,
		"agentID":   ba.id,
	}

	logging.Log(lvl, logging.CombineFields(agentFields, fields), msg)
}

func (ba *BaseAgent) ViewState() state.AgentState {
	return ba.latestState
}
