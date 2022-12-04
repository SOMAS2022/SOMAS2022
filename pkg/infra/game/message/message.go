package message

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
