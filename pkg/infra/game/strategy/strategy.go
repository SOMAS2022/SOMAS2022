package strategy

import (
	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"

	"infra/game/state"
)

type Strategy interface {
	ProcessStartOfRound(ba *message.BaseAgent, view *state.View, log *immutable.Map[commons.ID, decision.FightAction])
	GenerateActionDecision() decision.FightAction
	ProcessFightDecisionRequestMessage(ba *message.BaseAgent, FightDecisionRequestMessage message.Message, view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) FightDecisionMessage
	ProcessFightDecisionMessage(ba *message.BaseAgent, m FightDecisionMessage, view *state.View, log *immutable.Map[commons.ID, decision.FightAction])
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

type RequestMessage struct {
	*message.BaseMessage
}

func (RequestMessage) ProcessRequestMessage(Strategy Strategy, ba *message.BaseAgent, view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) InfoMessage {
	return InfoMessage{}
}

// InfoMessage Info messages only send information from one agent to another
type InfoMessage struct {
	*message.BaseMessage
}

func (InfoMessage) ProcessInfoMessage(Strategy Strategy,
	ba *message.BaseAgent,
	view *state.View,
	log *immutable.Map[commons.ID, decision.FightAction]) {
}

// FightRoundStartMessage Message sent to agents to signify the start of a round
type FightRoundStartMessage struct {
	InfoMessage
}

// ProcessInfoMessage The start round message passes the current state to the agent
func (m FightRoundStartMessage) ProcessInfoMessage(Strategy Strategy,
	ba *message.BaseAgent,
	view *state.View,
	log *immutable.Map[commons.ID, decision.FightAction]) {
	Strategy.ProcessStartOfRound(ba, view, log)
}

type FightDecisionMessage struct {
	// Info message from one agent to another indicating what its current fight decision is
	InfoMessage
	FightDecision decision.FightAction
}

func (m FightDecisionMessage) ProcessInfoMessage(Strategy Strategy,
	ba *message.BaseAgent,
	view *state.View,
	log *immutable.Map[commons.ID, decision.FightAction]) {
	Strategy.ProcessFightDecisionMessage(ba, m, view, log)
}

type FightDecisionRequestMessage struct {
	// Request from one agent to another asking what its current fight decision is
	RequestMessage
	FightDecision decision.FightAction
}

func (m FightDecisionRequestMessage) ProcessRequestMessage(Strategy Strategy,
	ba *message.BaseAgent,
	view *state.View,
	log *immutable.Map[commons.ID, decision.FightAction]) FightDecisionMessage {
	return Strategy.ProcessFightDecisionRequestMessage(ba, m, view, log)
}

type ActionDecision struct {
	Action decision.FightAction
	Sender commons.ID
}
