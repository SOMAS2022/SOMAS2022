package agent

import (
	"fmt"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

type Agent struct {
	*BaseAgent
	Strategy
}

func (a *Agent) HandleDonateToHpPool(agentState state.AgentState) uint {
	a.BaseAgent.latestState = agentState

	return a.Strategy.DonateToHpPool(*a.BaseAgent)
}

func (a *Agent) HandleUpdateInternalState(agentState state.AgentState, fightResults *commons.ImmutableList[decision.ImmutableFightResult], voteResults *immutable.Map[decision.Intent, uint], logChan chan<- logging.AgentLog) {
	a.BaseAgent.latestState = agentState

	a.Strategy.UpdateInternalState(*a.BaseAgent, fightResults, voteResults, logChan)
}

func (a *Agent) HandleUpdateWeapon(agentState state.AgentState) decision.ItemIdx {
	a.BaseAgent.latestState = agentState

	return a.Strategy.HandleUpdateWeapon(*a.BaseAgent)
}

func (a *Agent) HandleUpdateShield(agentState state.AgentState) decision.ItemIdx {
	a.BaseAgent.latestState = agentState
	return a.Strategy.HandleUpdateShield(*a.BaseAgent)
}

func (a *Agent) SubmitManifesto(agentState state.AgentState) *decision.Manifesto {
	a.BaseAgent.latestState = agentState

	return a.Strategy.CreateManifesto(*a.BaseAgent)
}

// HandleNoConfidenceVote todo: do we need to send the baseAgent here? I.e. is communication necessary at this point?
func (a *Agent) HandleNoConfidenceVote(agentState state.AgentState) decision.Intent {
	a.BaseAgent.latestState = agentState

	return a.Strategy.HandleConfidencePoll(*a.BaseAgent)
}

func (a *Agent) HandleElection(agentState state.AgentState, params *decision.ElectionParams) decision.Ballot {
	a.BaseAgent.latestState = agentState

	return a.Strategy.HandleElectionBallot(*a.BaseAgent, params)
}

func (a *Agent) HandleFight(agentState state.AgentState,
	log immutable.Map[commons.ID, decision.FightAction],
	votes chan commons.ProposalID,
	submission chan message.Proposal[decision.FightAction],
	closure <-chan struct{},
) {
	a.BaseAgent.latestState = agentState
	for {
		select {
		case taggedMessage := <-a.BaseAgent.communication.receipt:
			a.handleFightRoundMessage(&log, taggedMessage, votes, submission)
		case <-closure:
			return
		}
	}
}

func (a *Agent) isLeader() bool {
	return a.BaseAgent.ID() == a.BaseAgent.view.CurrentLeader()
}

func (a *Agent) SetCommunication(communication *Communication) {
	a.BaseAgent.setCommunication(communication)
}

func (a *Agent) handleFightRoundMessage(log *immutable.Map[commons.ID, decision.FightAction],
	m message.TaggedMessage,
	votes chan commons.ProposalID,
	submission chan message.Proposal[decision.FightAction],
) {
	switch r := m.Message().(type) {
	case message.FightRequest:
		req := *message.NewTaggedRequestMessage[message.FightRequest](m.Sender(), r, m.MID())
		resp := a.Strategy.HandleFightRequest(req, log)
		err := a.BaseAgent.SendBlockingMessage(m.Sender(), resp)
		logging.Log(logging.Error, nil, err.Error())
	case message.FightInform:
		inf := *message.NewTaggedInformMessage[message.FightInform](m.Sender(), r, m.MID())
		a.Strategy.HandleFightInformation(inf, *a.BaseAgent, log)

	case message.Proposal[decision.FightAction]:
		if a.isLeader() {
			if a.Strategy.HandleFightProposalRequest(r, *a.BaseAgent, log) {
				submission <- r
				iterator := a.BaseAgent.communication.peer.Iterator()
				for !iterator.Done() {
					_, value, _ := iterator.Next()
					value <- m
				}
			}
		}
		switch a.Strategy.HandleFightProposal(r, *a.BaseAgent) {
		case decision.Positive:
			votes <- r.ProposalID()
		default:
		}
	default:
		logging.Log(logging.Warn, nil, fmt.Sprintf("Unknown type, %T", r))
	}
}

func (a *Agent) HandleLoot(agentState state.AgentState, votes chan commons.ProposalID, submission chan message.Proposal[decision.LootAction], closure chan struct{}, start <-chan message.StartLoot) {
	a.BaseAgent.latestState = agentState
	for {
		select {
		case loot := <-start:
			a.addLoot(loot.LootPool)
		case taggedMessage := <-a.BaseAgent.communication.receipt:
			a.handleLootRoundMessage(taggedMessage, votes, submission)
		case <-closure:
			return
		}
	}
}

func (a *Agent) handleLootRoundMessage(
	m message.TaggedMessage,
	votes chan commons.ProposalID,
	submission chan message.Proposal[decision.LootAction],
) {
	switch r := m.Message().(type) {
	case message.LootRequest:
		req := *message.NewTaggedRequestMessage[message.LootRequest](m.Sender(), r, m.MID())
		resp := a.Strategy.HandleLootRequest(req)
		err := a.BaseAgent.SendBlockingMessage(m.Sender(), resp)
		logging.Log(logging.Error, nil, err.Error())
	case message.LootInform:
		inf := *message.NewTaggedInformMessage[message.LootInform](m.Sender(), r, m.MID())
		a.Strategy.HandleLootInformation(inf, *a.BaseAgent)
	case message.Proposal[decision.LootAction]:
		if a.isLeader() {
			if a.Strategy.HandleLootProposalRequest(r, *a.BaseAgent) {
				submission <- r
				iterator := a.BaseAgent.communication.peer.Iterator()
				for !iterator.Done() {
					_, value, _ := iterator.Next()
					value <- m
				}
			}
		}
		switch a.Strategy.HandleLootProposal(r, *a.BaseAgent) {
		case decision.Positive:
			votes <- r.ProposalID()
		default:
		}
	default:
		logging.Log(logging.Warn, nil, fmt.Sprintf("Unknown type, %T", r))
	}
}

func (a *Agent) addLoot(pool state.LootPool) {
	a.BaseAgent.loot = pool
}

func (a *Agent) HandleTrade(
	agentState state.AgentState,
	info message.TradeInfo,
	next <-chan interface{},
	closure <-chan interface{},
	responseChannel chan<- message.TradeMessage,
) {
	a.BaseAgent.latestState = agentState
	for {
		select {
		case <-closure:
			return
		case <-next:
			tradeMessage := a.Strategy.HandleTradeNegotiation(*a.BaseAgent, info)
			responseChannel <- tradeMessage
		}
	}
}
