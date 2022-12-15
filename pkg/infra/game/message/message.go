package message

import "infra/game/state"

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

type LootInform interface {
	Inform
	sealedLootInform()
}

type StartLoot struct {
	state.LootPool
}

func NewStartLoot(lootPool state.LootPool) *StartLoot {
	return &StartLoot{LootPool: lootPool}
}

func (s StartLoot) sealedMessage() {
	//TODO implement me
	panic("implement me")
}

func (s StartLoot) sealedInform() {
	//TODO implement me
	panic("implement me")
}

func (s StartLoot) sealedLootInform() {
	//TODO implement me
	panic("implement me")
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
