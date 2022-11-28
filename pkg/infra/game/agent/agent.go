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
	HandleFightInformation(m message.TaggedMessage, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction])
	HandleFightRequest(m message.TaggedMessage, log *immutable.Map[commons.ID, decision.FightAction]) message.Payload
	CurrentAction() decision.FightAction
	CreateManifesto(baseAgent BaseAgent) *decision.Manifesto
	HandleConfidencePoll(baseAgent BaseAgent) decision.Intent
	HandleElectionBallot(baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot
	HandleFightProposal(proposal *message.FightProposalMessage, baseAgent BaseAgent) decision.Intent
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
}

func (a *Agent) SubmitManifesto(agentState state.AgentState) *decision.Manifesto {
	a.BaseAgent.latestState = agentState
	return a.Strategy.CreateManifesto(a.BaseAgent)
}

// HandleNoConfidenceVote todo: do we need to send the baseAgent here? I.e. is communication necessary at this point?
func (a *Agent) HandleNoConfidenceVote(agentState state.AgentState) decision.Intent {
	a.BaseAgent.latestState = agentState
	return a.Strategy.HandleConfidencePoll(a.BaseAgent)
}

func (a *Agent) HandleElection(agentState state.AgentState, params *decision.ElectionParams) decision.Ballot {
	a.BaseAgent.latestState = agentState
	return a.Strategy.HandleElectionBallot(a.BaseAgent, params)
}

func (a *Agent) HandleFight(agentState state.AgentState, log immutable.Map[commons.ID, decision.FightAction], decisionChan chan<- message.ActionMessage, wg *sync.WaitGroup) {
	a.BaseAgent.latestState = agentState
	for m := range a.BaseAgent.communication.receipt {
		a.handleMessage(&log, m)
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

func (a *Agent) handleMessage(log *immutable.Map[commons.ID, decision.FightAction], m message.TaggedMessage) decision.FightAction {
	switch m.Message().MType() {
	case message.Close:
	case message.Request:
		payload := a.Strategy.HandleFightRequest(m, log)
		err := a.BaseAgent.SendBlockingMessage(m.Sender(), *message.NewMessage(message.Inform, payload))
		logging.Log(logging.Error, nil, err.Error())
	case message.Inform:
		a.Strategy.HandleFightInformation(m, a.BaseAgent, log)
	case message.Proposal:
		proposalMessage := message.NewFightProposalMessage(m)
		a.Strategy.HandleFightProposal(proposalMessage, a.BaseAgent)
	default:
		a.Strategy.HandleFightInformation(m, a.BaseAgent, log)
	}
	return a.Strategy.CurrentAction()
}
