package message

import "infra/game/commons"

type Payload interface {
	isPayload()
}

type Type int64

const (
	Proposal Type = iota
	Something
	SomethingElse
)

type Message struct {
	mType   Type
	payload Payload
}

type TaggedMessage struct {
	Sender  commons.ID
	Message Message
}
