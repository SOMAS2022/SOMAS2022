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
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
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

func (b BaseAgent) sendBlockingMessage(id commons.ID, m message.Message) error {
	value, ok := b.communication.peer.Get(id)
	if ok {
		value <- message.TaggedMessage{
			Sender:  b.Id,
			Message: m,
		}
		return nil
	} else {
		return fmt.Errorf("agent %s not available for messaging, either dead or submitted", id)
	}
}

func (b BaseAgent) receiveAllMessages() (res []message.TaggedMessage) {
	res = make([]message.TaggedMessage, 0)
	for {
		receive, b := nonBlockingReceive(b.communication.receipt)
		if b {
			res = append(res, receive)
		} else {
			return
		}
	}
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
