package agent

import (
	"fmt"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/stages"
	"infra/game/state"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

type Agent struct {
	BaseAgent
	Strategy
}

func (a *Agent) HandleDonateToHpPool(agentState state.AgentState) uint {
	a.BaseAgent.latestState = agentState

	return a.Strategy.DonateToHpPool(a.BaseAgent)
}

func (a *Agent) HandleUpdateInternalState(agentState state.AgentState, fightResults *commons.ImmutableList[decision.ImmutableFightResult], voteResults *immutable.Map[decision.Intent, uint]) {
	a.BaseAgent.latestState = agentState

	a.Strategy.UpdateInternalState(a.BaseAgent, fightResults, voteResults)
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
		req := *message.NewTaggedRequestMessage[message.FightRequest](m.Sender(), r, m.MID())
		resp := a.Strategy.HandleFightRequest(req, log)
		err := a.BaseAgent.SendBlockingMessage(m.Sender(), resp)
		logging.Log(logging.Error, nil, err.Error())
	case message.FightInform:
		inf := *message.NewTaggedInformMessage[message.FightInform](m.Sender(), r, m.MID())
		a.Strategy.HandleFightInformation(inf, a.BaseAgent, log)
	case message.MapProposal[decision.FightAction]:
		//todo: Refactor this type to be similar to the types above
		v := *message.NewFightProposalMessage(m.Sender(), r.Proposal(), r.ProposalID())
		if a.isLeader() {
			if a.Strategy.HandleFightProposalRequest(v, a.BaseAgent, log) {
				submission <- *message.NewProposal(v.ProposalID(), v.Proposal())
				iterator := a.BaseAgent.communication.peer.Iterator()
				for !iterator.Done() {
					_, value, _ := iterator.Next()
					value <- m
				}
			}
		}
		switch a.Strategy.HandleFightProposal(v, a.BaseAgent) {
		case decision.Positive:
			votes <- v.ProposalID()
		default:
		}
	case message.CustomInform:
		if stages.Mode != "default" {
			inf := *message.NewTaggedInformMessage[message.CustomInform](m.Sender(), r, m.MID())
			a.Strategy.HandleCustomInformation(inf, a.BaseAgent, log)
		}
	default:
		logging.Log(logging.Warn, nil, fmt.Sprintf("Unknown type, %T", r))
	}
}
