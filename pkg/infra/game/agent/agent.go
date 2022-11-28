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

	HandleLootInformation(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.LootDecision])
	HandleLootRequest(m message.TaggedMessage, view *state.View, log *immutable.Map[commons.ID, decision.LootDecision]) message.Payload

	HandleHpInformation(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.HpPoolDecision])
	HandleHpRequest(m message.TaggedMessage, view *state.View, log *immutable.Map[commons.ID, decision.HpPoolDecision]) message.Payload

	CurrentAction() decision.FightAction
	CurrentLootDescion() decision.LootDecision
	CurrentHpDecision() decision.HpPoolDecision

	CreateManifesto(view *state.View, baseAgent BaseAgent) *decision.Manifesto
	HandleConfidencePoll(view *state.View, baseAgent BaseAgent) decision.Intent
	HandleElectionBallot(view *state.View, baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
}

func (a *Agent) SubmitManifesto(agentState state.AgentState, view *state.View, baseAgent BaseAgent) *decision.Manifesto {
	a.BaseAgent.latestState = agentState
	return a.Strategy.CreateManifesto(view, baseAgent)
}

// HandleNoConfidenceVote todo: do we need to send the baseAgent here? I.e. is communication necessary at this point?
func (a *Agent) HandleNoConfidenceVote(agentState state.AgentState, view *state.View, baseAgent BaseAgent) decision.Intent {
	a.BaseAgent.latestState = agentState
	return a.Strategy.HandleConfidencePoll(view, baseAgent)
}

func (a *Agent) HandleElection(agentState state.AgentState, view *state.View, baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot {
	a.BaseAgent.latestState = agentState
	return a.Strategy.HandleElectionBallot(view, baseAgent, params)
}

func (a *Agent) HandleFight(agentState state.AgentState, view state.View, log immutable.Map[commons.ID, decision.FightAction], decisionChan chan<- message.ActionMessage, wg *sync.WaitGroup) {
	a.BaseAgent.latestState = agentState
	for m := range a.BaseAgent.communication.receipt {
		a.handleFightActionMessage(&view, &log, m)
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

func (a *Agent) handleFightActionMessage(view *state.View, log *immutable.Map[commons.ID, decision.FightAction], m message.TaggedMessage) decision.FightAction {
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

func (a *Agent) handleLootDecisionMessage(view *state.View, log *immutable.Map[commons.ID, decision.LootDecision], m message.TaggedMessage) decision.LootDecision {
	switch m.Message().MType() {
	case message.Close:
	case message.Request:
		payload := a.Strategy.HandleLootRequest(m, view, log)
		err := a.BaseAgent.SendBlockingMessage(m.Sender(), *message.NewMessage(message.Inform, payload))
		logging.Log(logging.Error, nil, err.Error())
	case message.Inform:
		a.Strategy.HandleLootInformation(m, view, a.BaseAgent, log)
	default:
		a.Strategy.HandleLootInformation(m, view, a.BaseAgent, log)
	}
	return a.Strategy.CurrentLootDescion()
}

func (a *Agent) handleHpDecisionMessage(view *state.View, log *immutable.Map[commons.ID, decision.HpPoolDecision], m message.TaggedMessage) decision.HpPoolDecision {
	switch m.Message().MType() {
	case message.Close:
	case message.Request:
		payload := a.Strategy.HandleHpRequest(m, view, log)
		err := a.BaseAgent.SendBlockingMessage(m.Sender(), *message.NewMessage(message.Inform, payload))
		logging.Log(logging.Error, nil, err.Error())
	case message.Inform:
		a.Strategy.HandleHpInformation(m, view, a.BaseAgent, log)
	default:
		a.Strategy.HandleHpInformation(m, view, a.BaseAgent, log)
	}
	return a.Strategy.CurrentHpDecision()
}
