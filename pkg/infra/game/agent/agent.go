package agent

import (
	"fmt"
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"
	"sync"
)

type Strategy interface {
	HandleFightMessage(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) decision.FightAction
	Default() decision.FightAction
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
}

func (a *Agent) HandleFight(view state.View, log immutable.Map[commons.ID, decision.FightAction], decisionChan chan message.ActionMessage, wg *sync.WaitGroup) {
	for m := range a.BaseAgent.communication.receipt {
		action := a.handleMessage(&view, &log, m)
		if action != decision.Undecided {
			go func() {
				<-a.BaseAgent.communication.receipt
			}()
			decisionChan <- message.ActionMessage{Action: action, Sender: a.BaseAgent.Id}
			wg.Done()
			return
		}
	}
	decisionChan <- message.ActionMessage{Action: a.Strategy.Default(), Sender: a.BaseAgent.Id}
}

func (a *Agent) handleMessage(view *state.View, log *immutable.Map[commons.ID, decision.FightAction], m message.TaggedMessage) decision.FightAction {
	switch m.Message.MType() {
	case message.Close:
	default:
		fightMessage := a.Strategy.HandleFightMessage(m, view, a.BaseAgent, log)
		return fightMessage
	}
	return decision.Undecided
}

func (ba *BaseAgent) log(lvl logging.Level, fields logging.LogField, msg string) {
	agentFields := logging.LogField{
		"agentName": ba.AgentName,
		"agentID":   ba.Id,
	}

	logging.Log(lvl, logging.CombineFields(agentFields, fields), msg)
}

type BaseAgent struct {
	communication Communication
	Id            commons.ID
	AgentName     string
}

func NewBaseAgent(communication Communication, id commons.ID, agentName string) BaseAgent {
	return BaseAgent{communication: communication, Id: id, AgentName: agentName}
}

type Communication struct {
	receipt <-chan message.TaggedMessage
	peer    immutable.Map[commons.ID, chan<- message.TaggedMessage]
}

func NewCommunication(receipt <-chan message.TaggedMessage, peer immutable.Map[commons.ID, chan<- message.TaggedMessage]) Communication {
	return Communication{receipt: receipt, peer: peer}
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
