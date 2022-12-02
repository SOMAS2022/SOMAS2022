package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

type Agent struct {
	BaseAgent
	Strategy
}

func (a *Agent) HandleUpdateWeapon(agentState state.AgentState, view state.View) decision.ItemIdx {
	a.BaseAgent.latestState = agentState

	return a.Strategy.HandleUpdateWeapon(&view, a.BaseAgent)
}

func (a *Agent) HandleUpdateShield(agentState state.AgentState, view state.View) decision.ItemIdx {
	a.BaseAgent.latestState = agentState

	return a.Strategy.HandleUpdateShield(&view, a.BaseAgent)
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
	submission chan message.MapProposal[decision.FightAction],
	closure <-chan struct{},
) {
	a.BaseAgent.latestState = agentState
	for {
		select {
		case taggedMessage := <-a.BaseAgent.communication.receipt:
			a.handleMessage(&log, taggedMessage, votes, submission)
		case <-closure:
			return
		}
	}
}

func (a *Agent) isLeader() bool {
	return a.BaseAgent.ID() == a.BaseAgent.view.CurrentLeader()
}

func (a *Agent) handleMessage(log *immutable.Map[commons.ID, decision.FightAction],
	m message.TaggedMessage,
	votes chan commons.ProposalID,
	submission chan message.MapProposal[decision.FightAction],
) {
	switch r := m.Message().(type) {
	case message.FightRequest:
		resp := a.Strategy.HandleFightRequest(m, log)
		err := a.BaseAgent.SendBlockingMessage(m.Sender(), resp)
		logging.Log(logging.Error, nil, err.Error())
	case message.FightInform:
		a.Strategy.HandleFightInformation(m, a.BaseAgent, log)
	case message.FightProposalMessage:
		if a.isLeader() {
			if a.Strategy.HandleFightProposalRequest(r, a.BaseAgent, log) {
				submission <- *message.NewProposal(r.ProposalID(), r.Proposal())
				iterator := a.BaseAgent.communication.peer.Iterator()
				for !iterator.Done() {
					_, value, _ := iterator.Next()
					value <- m
				}
			}
		}
		switch a.Strategy.HandleFightProposal(r, a.BaseAgent) {
		case decision.Positive:
			votes <- r.ProposalID()
		default:
		}

	default:
		a.Strategy.HandleFightInformation(m, a.BaseAgent, log)
	}
}
