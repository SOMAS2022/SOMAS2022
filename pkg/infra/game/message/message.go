package message

import (
	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
)

/*
Each message should either be an info message ort a request messages,
request messages require the receiver to give a response in a certain amount of time

Each message should implement either a ProcessRequestMessage or ProcessInfoMessage method
These methods should call a method in the strategy interface on the strategy passed to them
This means that message types do not have to be switched by the agent implementations
*/

// TODO need to add a mechanism for agents to send messages

type Strategy interface {
	ProcessStartOfRound(view *state.View, log *immutable.Map[commons.ID, decision.FightAction])
	GenerateActionDecision() decision.FightAction
	ProcessFightDecisionRequestMessage(FightDecisionRequestMessage Message) FightDecisionMessage
	ProcessFightDecisionMessage(FightDecisionMessage)
}

type Message interface {
	GenerateUUID()
	GetUUID() uuid.UUID
	SetUUID(UUID uuid.UUID)
}

type RequestMessageInterface interface {
	GenerateUUID()
	SetUUID(UUID uuid.UUID)
	GetUUID() uuid.UUID
	ProcessRequestMessage(Strategy Strategy,
		view *state.View,
		log *immutable.Map[commons.ID, decision.FightAction]) InfoMessageInterface
}

type InfoMessageInterface interface {
	GenerateUUID()
	SetUUID(UUID uuid.UUID)
	GetUUID() uuid.UUID
	ProcessInfoMessage(Strategy Strategy,
		view *state.View,
		log *immutable.Map[commons.ID, decision.FightAction])
}

// Base Message

type BaseMessage struct {
	UUID uuid.UUID
}

func (b *BaseMessage) GenerateUUID() {
	b.UUID = uuid.New()
}

func (b *BaseMessage) GetUUID() uuid.UUID {
	return b.UUID
}

func (b *BaseMessage) SetUUID(UUID uuid.UUID) {
	b.UUID = UUID
}

// Request Message force the receiver to respond in a given amount of time

type RequestMessage struct {
	*BaseMessage
}

func (RequestMessage) ProcessRequestMessage(Strategy Strategy,
	view *state.View,
	log *immutable.Map[commons.ID, decision.FightAction]) InfoMessageInterface {
	return InfoMessage{}
}

// Info messages only send information from one agent to another
type InfoMessage struct {
	*BaseMessage
}

func (InfoMessage) ProcessInfoMessage(Strategy Strategy,
	view *state.View,
	log *immutable.Map[commons.ID, decision.FightAction]) {
}

/*
Message Structs
*/

// Message sent to agents to signify the start of a round
type FightRoundStartMessage struct {
	InfoMessage
}

// The start round message passes the current state to the agent
func (m FightRoundStartMessage) ProcessInfoMessage(Strategy Strategy,
	view *state.View,
	log *immutable.Map[commons.ID, decision.FightAction]) {
	Strategy.ProcessStartOfRound(view, log)
}

type FightDecisionMessage struct {
	// Info message from one agent to another indicating what its current fight decision is
	InfoMessage
	FightDecision decision.FightAction
}

func (m FightDecisionMessage) ProcessInfoMessage(Strategy Strategy,
	view *state.View,
	log *immutable.Map[commons.ID, decision.FightAction]) {
	Strategy.ProcessFightDecisionMessage(m)
}

type FightDecisionRequestMessage struct {
	// Request from one agent to another asking what its current fight decision is
	RequestMessage
	FightDecision decision.FightAction
}

func (m FightDecisionRequestMessage) ProcessRequestMessage(Strategy Strategy,
	view *state.View,
	log *immutable.Map[commons.ID, decision.FightAction]) InfoMessageInterface {
	return Strategy.ProcessFightDecisionRequestMessage(m)
}

/*
Associated messaging structs
*/

type TaggedMessage struct {
	Sender  commons.ID
	Message Message
}

type ActionDecision struct {
	Action decision.FightAction
	Sender commons.ID
}
