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
	receipt <-chan message.Message
	peer    *immutable.Map[commons.ID, chan<- message.Message]
}

func NewCommunication(receipt <-chan message.Message, peer *immutable.Map[commons.ID, chan<- message.Message]) *Communication {
	return &Communication{receipt: receipt, peer: peer}
}

func (c Communication) sendMessage(id commons.ID, message message.Message) error {
	m := message
	value, ok := c.peer.Get(id)
	if ok {
		value <- m
		return nil
	} else {
		return fmt.Errorf("agent %s not available for messaging, either dead or submitted", id)
	}
}
