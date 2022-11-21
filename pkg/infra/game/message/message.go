package message

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

type Strategy interface {
	ProcessStartOfRound(view state.View, log immutable.Map[commons.ID, decision.FightAction])
	GenerateActionDecision() decision.FightAction
	ProcessFightDecisionRequestMessage(FightDecisionRequestMessage Message) FightDecisionMessage
	ProcessFightDecisionMessage(FightDecisionMessage)
}

type Message interface {
	RequiresResponse() bool
}

type RequestMessageInterface interface {
	ProcessRequestMessage(Strategy Strategy) InfoMessageInterface
	RequiresResponse() bool
}

type InfoMessageInterface interface {
	ProcessInfoMessage(Strategy Strategy)
	RequiresResponse() bool
}

// Request Message force the receiver to respond in a given amount of time

type RequestMessage struct{}

func (RequestMessage) RequiresResponse() bool {
	return true
}

func (RequestMessage) ProcessRequestMessage(Strategy Strategy) InfoMessageInterface {
	return InfoMessage{}
}

// Info messages only send information from one agent to another
type InfoMessage struct{}

func (InfoMessage) RequiresResponse() bool {
	return false
}

func (InfoMessage) ProcessInfoMessage(Strategy Strategy) {
}

/*
Examples of possible message structs
*/

type FightDecisionMessage struct {
	// Info message from one agent to another indicating what its current fight decision is
	InfoMessage
	FightDecision decision.FightAction
}

func (m FightDecisionMessage) ProcessInfoMessage(Strategy Strategy) {
	Strategy.ProcessFightDecisionMessage(m)
}

type FightDecisionRequestMessage struct {
	// Request from one agent to another asking what its current fight decision is
	RequestMessage
	FightDecision decision.FightAction
}

func (m FightDecisionRequestMessage) ProcessRequestMessage(Strategy Strategy) InfoMessageInterface {
	return Strategy.ProcessFightDecisionRequestMessage(m)
}

/*
Associated message structs
*/

type TaggedMessage struct {
	Sender  commons.ID
	Message Message
}

type ActionDecision struct {
	Action decision.FightAction
	Sender commons.ID
}

// type Payload interface {
// 	isPayload()
// }

// type InfoMessageType int64

// const (
// 	Close InfoMessageType = iota
// 	Something
// 	SomethingElse
// )

// type InfoMessage struct {
// 	mType   InfoMessageType
// 	payload Payload
// }

// type RequestMessageType int64

// const (
// 	RequestTrade RequestMessageType = iota
// )

// type RequestMessage struct {
// 	mType   RequestMessageType
// 	payload Payload
// }

// func NewMessage(mType Type, payload Payload) *Message {
// 	return &Message{mType: mType, payload: payload}
// }

// func (m Message) MType() Type {
// 	return m.mType
// }

// func (m Message) Payload() Payload {
// 	return m.payload
// }

// type TaggedMessage struct {
// 	Sender          commons.ID
// 	Message         Message
// 	ResponseMessage bool
// }
