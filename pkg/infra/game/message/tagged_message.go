package message

import (
	"infra/game/commons"

	"github.com/google/uuid"
)

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

func (t TaggedRequestMessage[R]) Sender() commons.ID {
	return t.sender
}

func (t TaggedRequestMessage[R]) Message() R {
	return t.message
}

func (t TaggedRequestMessage[R]) MID() uuid.UUID {
	return t.mID
}

type TaggedInformMessage[I Inform] struct {
	sender  commons.ID
	message I
	mID     uuid.UUID
}

func NewTaggedInformMessage[I Inform](sender commons.ID, message I, mID uuid.UUID) *TaggedInformMessage[I] {
	return &TaggedInformMessage[I]{sender: sender, message: message, mID: mID}
}

func (t TaggedInformMessage[I]) Sender() commons.ID {
	return t.sender
}

func (t TaggedInformMessage[I]) Message() I {
	return t.message
}

func (t TaggedInformMessage[I]) MID() uuid.UUID {
	return t.mID
}
