package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/game/tally"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

type Strategy interface {
	HandleFightInformation(m message.TaggedMessage, baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction])
	HandleFightRequest(m message.TaggedMessage, log *immutable.Map[commons.ID, decision.FightAction]) message.Payload
	CurrentAction() decision.FightAction
	CreateManifesto(baseAgent BaseAgent) *decision.Manifesto
	HandleConfidencePoll(baseAgent BaseAgent) decision.Intent
	HandleElectionBallot(baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot
	HandleFightProposal(proposal *message.FightProposalMessage, baseAgent BaseAgent) decision.Intent
	FightResolution(agent BaseAgent) tally.Proposal[decision.FightAction]
	// HandleFightProposalRequest only called as leader
	HandleFightProposalRequest(proposal *message.FightProposalMessage, baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool
	// return the index of the weapon you want to use in AgentState.Weapons
	HandleUpdateWeapon(view *state.View, baseAgent BaseAgent) decision.ItemIdx
	// return the index of the shield you want to use in AgentState.Shields
	HandleUpdateShield(view *state.View, baseAgent BaseAgent) decision.ItemIdx
}

type Agent struct {
	BaseAgent BaseAgent
	Strategy  Strategy
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
	submission chan tally.Proposal[decision.FightAction],
	closure <-chan struct{},
) {
	a.BaseAgent.latestState = agentState
	for {
		select {
		case taggedMessage := <-a.BaseAgent.communication.receipt:
			a.handleMessage(&log, taggedMessage, votes, submission)
		case _ = <-closure:
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
	submission chan tally.Proposal[decision.FightAction],
) {
	switch m.Message().MType() {
	case message.Request:
		payload := a.Strategy.HandleFightRequest(m, log)
		err := a.BaseAgent.SendBlockingMessage(m.Sender(), *message.NewMessage(message.Inform, payload))
		logging.Log(logging.Error, nil, err.Error())
	case message.Inform:
		a.Strategy.HandleFightInformation(m, a.BaseAgent, log)
	case message.Proposal:
		// todo: make this generic for all future proposals
		proposalMessage := message.NewFightProposalMessage(m)
		if a.isLeader() {
			if a.Strategy.HandleFightProposalRequest(proposalMessage, a.BaseAgent, log) {
				submission <- *tally.NewProposal[decision.FightAction](proposalMessage.ProposalID(), proposalMessage.Proposal())
				iterator := a.BaseAgent.communication.peer.Iterator()
				for !iterator.Done() {
					_, value, _ := iterator.Next()
					value <- m
				}
			}
		}
		switch a.Strategy.HandleFightProposal(proposalMessage, a.BaseAgent) {
		case decision.Positive:
			votes <- proposalMessage.ProposalID()
		default:
		}
	default:
		a.Strategy.HandleFightInformation(m, a.BaseAgent, log)
	}
}

func (a *Agent) broadcastProposal() {

}
