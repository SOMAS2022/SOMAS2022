package agent

import (
	"fmt"
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
)

type Strategy interface {
	HandleFight(view *state.View, baseAgent BaseAgent, decisionC chan<- decision.FightAction, log *immutable.Map[commons.ID, decision.FightAction])
	HandleFightMessage(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) *decision.FightAction
	Default() decision.FightAction
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
}

func (a *Agent) HandleFight(view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) decision.FightAction {
	for m := range a.BaseAgent.communication.receipt {
		action := a.handleMessage(view, log, m)
		if action != nil {
			return *action
		}
	}
	return a.Strategy.Default()
}

func (a *Agent) handleMessage(view *state.View, log *immutable.Map[commons.ID, decision.FightAction], m message.TaggedMessage) *decision.FightAction {
	switch m.Message.MType() {
	case message.Close:
		a.BaseAgent.communication.peer = a.BaseAgent.communication.peer.Delete(m.Sender)
	default:
		fightMessage := a.Strategy.HandleFightMessage(m, view, a.BaseAgent, log)
		if fightMessage != nil {
			a.handleDecisionMaking()
		}
		return fightMessage
	}
	return nil
}

func (a *Agent) handleDecisionMaking() {
	iterator := a.BaseAgent.communication.peer.Iterator()
	go func() {
		<-a.BaseAgent.communication.receipt
	}()
	for !iterator.Done() {
		_, c, _ := iterator.Next()
		c <- message.TaggedMessage{
			Sender:  a.BaseAgent.Id,
			Message: *message.NewMessage(message.Close, nil),
		}
	}
}

type BaseAgent struct {
	communication *Communication
	Id            commons.ID
}

func NewBaseAgent(communication *Communication, id commons.ID) BaseAgent {
	return BaseAgent{communication: communication, Id: id}
}

type Communication struct {
	receipt <-chan message.TaggedMessage
	peer    *immutable.Map[commons.ID, chan<- message.TaggedMessage]
}

func NewCommunication(receipt <-chan message.TaggedMessage, peer *immutable.Map[commons.ID, chan<- message.TaggedMessage]) *Communication {
	return &Communication{receipt: receipt, peer: peer}
}

func (b BaseAgent) broadcastBlockingMessage(m message.Message) {
	iterator := b.communication.peer.Iterator()
	tm := message.TaggedMessage{
		Sender:  b.Id,
		Message: m,
	}
	for !iterator.Done() {
		_, c, ok := iterator.Next()
		if ok {
			c <- tm
		}
	}
}

func (b BaseAgent) sendBlockingMessage(id commons.ID, m message.Message) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("agent %s not available for messaging, submitted", id)
		}
	}()

	value, ok := b.communication.peer.Get(id)
	if ok {
		value <- message.TaggedMessage{
			Sender:  b.Id,
			Message: m,
		}
	} else {
		e = fmt.Errorf("agent %s not available for messaging, dead", id)
	}
	return
}

func nonBlockingReceive[T any](c <-chan T) (T, bool) {
	select {
	case msg := <-c:
		return msg, true
	default:
		var result T
		return result, false
	}
}
