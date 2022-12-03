package message

import (
	"infra/game/commons"

	"github.com/google/uuid"
)

type Message interface {
	sealedMessage()
}

type Inform interface {
	Message
	sealedInform()
}

type Request interface {
	Message
	sealedRequest()
}

type Proposal interface {
	Message
	sealedProposal()
}

type FightRequest interface {
	Request
	sealedFightRequest()
}

type LootRequest interface {
	Request
	sealedLootRequest()
}

type FightInform interface {
	Inform
	sealedFightInform()
}

type TaggedMessage struct {
	sender  commons.ID
	message Message
	mID     uuid.UUID
}

func NewTaggedMessage(sender commons.ID, message Message, mID uuid.UUID) *TaggedMessage {
	return &TaggedMessage{sender: sender, message: message, mID: mID}
}

func (t TaggedMessage) Sender() commons.ID {
	return t.sender
}

func (t TaggedMessage) Message() Message {
	return t.message
}

func (t TaggedMessage) MID() uuid.UUID {
	return t.mID
}

type TaggedRequestMessage[R Request] struct {
	sender  commons.ID
	message R
	mID     uuid.UUID
}

func NewTaggedRequestMessage[R Request](sender commons.ID, message R, mID uuid.UUID) *TaggedRequestMessage[R] {
	return &TaggedRequestMessage[R]{sender: sender, message: message, mID: mID}
}

type TaggedInformMessage[I Inform] struct {
	sender  commons.ID
	message I
	mID     uuid.UUID
}

func NewTaggedInformMessage[I Inform](sender commons.ID, message I, mID uuid.UUID) *TaggedInformMessage[I] {
	return &TaggedInformMessage[I]{sender: sender, message: message, mID: mID}
}

type StartFight struct{}

func (s StartFight) sealedMessage() {
	// TODO implement me
	panic("implement me")
}

func (s StartFight) sealedInform() {
	panic("implement me")
}

func (s StartFight) sealedFightInform() {
	panic("implement me")
}
