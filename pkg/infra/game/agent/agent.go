package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"
	"sync"

	"github.com/benbjohnson/immutable"
)

type Strategy interface {
	HandleFightInformation(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction])
	HandleFightRequest(m message.TaggedMessage, view *state.View, log *immutable.Map[commons.ID, decision.FightAction]) message.Payload
	CurrentAction() decision.FightAction
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
}

func (a *Agent) HandleFight(view state.View, log immutable.Map[commons.ID, decision.FightAction], decisionChan chan message.ActionMessage, wg *sync.WaitGroup) {
	for m := range a.BaseAgent.communication.receipt {
		a.handleMessage(&view, &log, m)
		action := a.Strategy.CurrentAction()
		if action != decision.Undecided {
			go func() {
				<-a.BaseAgent.communication.receipt
			}()
			decisionChan <- message.ActionMessage{Action: action, Sender: a.BaseAgent.id}
			wg.Done()
			return
		}
	}
	decisionChan <- message.ActionMessage{Action: a.Strategy.CurrentAction(), Sender: a.BaseAgent.id}
}

func (a *Agent) handleMessage(view *state.View, log *immutable.Map[commons.ID, decision.FightAction], m message.TaggedMessage) decision.FightAction {
	switch m.Message().MType() {
	case message.Close:
	case message.Request:
		payload := a.Strategy.HandleFightRequest(m, view, log)
		err := a.BaseAgent.SendBlockingMessage(m.Sender(), *message.NewMessage(message.Inform, payload))
		logging.Log(logging.Error, nil, err.Error())
	case message.Inform:
		a.Strategy.HandleFightInformation(m, view, a.BaseAgent, log)
	default:
		a.Strategy.HandleFightInformation(m, view, a.BaseAgent, log)
	}
	return a.Strategy.CurrentAction()
}
