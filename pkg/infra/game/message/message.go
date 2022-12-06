package message

import "infra/game/state"

type Message interface {
	SealedMessage()
}

type Inform interface {
	Message
	SealedInform()
}

type Request interface {
	Message
	SealedRequest()
}

type FightRequest interface {
	Request
	SealedFightRequest()
}

type LootRequest interface {
	Request
	SealedLootRequest()
}

type FightInform interface {
	Inform
	SealedFightInform()
}

type LootInform interface {
	Inform
	SealedLootInform()
}

type StartLoot struct {
	state.LootPool
}

func NewStartLoot(lootPool state.LootPool) *StartLoot {
	return &StartLoot{LootPool: lootPool}
}

func (s StartLoot) SealedMessage() {
	//TODO implement me
	panic("implement me")
}

func (s StartLoot) SealedInform() {
	//TODO implement me
	panic("implement me")
}

func (s StartLoot) SealedLootInform() {
	//TODO implement me
	panic("implement me")
}

type StartFight struct{}

func (s StartFight) SealedMessage() {
	// TODO implement me
	panic("implement me")
}

func (s StartFight) SealedInform() {
	panic("implement me")
}

func (s StartFight) SealedFightInform() {
	panic("implement me")
}
