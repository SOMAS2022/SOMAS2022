package message

import "infra/game/commons"

type Payload interface {
	isPayload()
}

type Type int64

const (
	Close Type = iota
	Something
	SomethingElse
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
	Sender  commons.ID
	Message Message
}
