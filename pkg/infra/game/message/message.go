package message

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/google/uuid"
)

type Payload interface {
	isPayload()
}

type Type int64

const (
	Close Type = iota
	Request
	Inform
)

type Message struct {
	mType   Type
	payload Payload
}

func NewMessage(mType Type, payload Payload) *Message {
	return &Message{mType: mType, payload: payload}
}

func (m Message) MType() Type {
	return m.mType
}

func (m Message) Payload() Payload {
	return m.payload
}

type TaggedMessage struct {
	sender  commons.ID
	message Message
	mId     uuid.UUID
}

func NewTaggedMessage(sender commons.ID, message Message, mId uuid.UUID) *TaggedMessage {
	return &TaggedMessage{sender: sender, message: message, mId: mId}
}

func (t TaggedMessage) Sender() commons.ID {
	return t.sender
}

func (t TaggedMessage) Message() Message {
	return t.message
}

type ActionMessage struct {
	Action decision.FightAction
	Sender commons.ID
}
