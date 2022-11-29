package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/game/tally"
	"infra/logging"
)

type Strategy interface {
	HandleFightInformation(m message.TaggedMessage, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction])
	HandleFightRequest(m message.TaggedMessage, log *immutable.Map[commons.ID, decision.FightAction]) message.Payload
	CurrentAction() decision.FightAction
	CreateManifesto(baseAgent BaseAgent) *decision.Manifesto
	HandleConfidencePoll(baseAgent BaseAgent) decision.Intent
	HandleElectionBallot(baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot
	HandleFightProposal(proposal *message.FightProposalMessage, baseAgent BaseAgent) decision.Intent
	FightResolution(agent BaseAgent) tally.Proposal[decision.FightAction]
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

func (a *Agent) HandleFight(agentState state.AgentState,
	log immutable.Map[commons.ID, decision.FightAction],
	votes chan commons.ProposalID,
	submission chan tally.Proposal[decision.FightAction]) {
	a.BaseAgent.latestState = agentState
	for m := range a.BaseAgent.communication.receipt {
		if m.Message().MType() == message.Close {
			break
		}
		a.handleMessage(&log, m, votes, submission)
	}
}

func (a *Agent) handleMessage(log *immutable.Map[commons.ID, decision.FightAction],
	m message.TaggedMessage,
	votes chan commons.ProposalID,
	submission chan tally.Proposal[decision.FightAction]) {
	switch m.Message().MType() {
	case message.Close:
	case message.Request:
		payload := a.Strategy.HandleFightRequest(m, log)
		err := a.BaseAgent.SendBlockingMessage(m.Sender(), *message.NewMessage(message.Inform, payload))
		logging.Log(logging.Error, nil, err.Error())
	case message.Inform:
		a.Strategy.HandleFightInformation(m, a.BaseAgent, log)
	case message.Proposal:
		//todo: if I am the leader then decide whether to broadcast
		// todo: if broadcast then send and send to tally
		proposalMessage := message.NewFightProposalMessage(m)
		switch a.Strategy.HandleFightProposal(proposalMessage, a.BaseAgent) {
		case decision.Positive:
			votes <- proposalMessage.ProposalId()
		default:
		}
	default:
		a.Strategy.HandleFightInformation(m, a.BaseAgent, log)
	}
}
