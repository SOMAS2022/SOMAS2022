package team4

import (
	"math/rand"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

type AgentFour struct {
	HP           int
	ST           int
	AT           int
	bravery      int
	uR           map[commons.ID]int
	uP           map[commons.ID]int
	uC           map[commons.ID]int
	utilityScore map[commons.ID]int
	TSN          []commons.ID
}

func (a *AgentFour) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	//
}

// FUNCTIONS COPIED //

func (a *AgentFour) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	//
}

func (a *AgentFour) FightAction(baseAgent agent.BaseAgent, proposedAction decision.FightAction, acceptedProposal message.Proposal[decision.FightAction]) decision.FightAction {
	//
}

func (a *AgentFour) FightActionNoProposal(baseAgent agent.BaseAgent) decision.FightAction {
	//
}

func (a *AgentFour) FightResolution(baseAgent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]]) immutable.Map[commons.ID, decision.FightAction] {
	//
}

func (a *AgentFour) CurrentAction() decision.FightAction {
	//
}

func (a *AgentFour) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	//
}

func (a *AgentFour) HandleElectionBallot(baseAgent agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	//
}

func (a *AgentFour) HandleFightInformation(_ message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	//
}

func (a *AgentFour) HandleFightProposal(m message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	//
}

func (a *AgentFour) HandleFightProposalRequest(_ message.Proposal[decision.FightAction], _ agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) bool {
	//
}

func (a *AgentFour) HandleFightRequest(_ message.TaggedRequestMessage[message.FightRequest], _ *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	//
}

func (a *AgentFour) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent agent.BaseAgent) {
	//
}

func (a *AgentFour) HandleLootProposal(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	//
}

func (a *AgentFour) HandleLootProposalRequest(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	//
}

func (a *AgentFour) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	//
}

func (a *AgentFour) HandleTradeNegotiation(_ agent.BaseAgent, _ message.TradeInfo) message.TradeMessage {
	//
}

func (a *AgentFour) HandleUpdateShield(baseAgent agent.BaseAgent) decision.ItemIdx {
	//
}

func (a *AgentFour) HandleUpdateWeapon(baseAgent agent.BaseAgent) decision.ItemIdx {
	//
}

func (a *AgentFour) LootAction(baseAgent agent.BaseAgent, proposedLoot immutable.SortedMap[commons.ItemID, struct{}], acceptedProposal message.Proposal[decision.LootAction]) immutable.SortedMap[commons.ItemID, struct{}] {
	//
}

func (a *AgentFour) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
	//
}

func (a *AgentFour) LootAllocation(baseAgent agent.BaseAgent, proposal message.Proposal[decision.LootAction]) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	//
}

func allocateRandomly(iterator commons.Iterator[state.Item], ids []commons.ID, lootAllocation map[commons.ID][]commons.ItemID) {
	//
}

func (a *AgentFour) UpdateInternalState(baseAgent agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	//
}

func (a *AgentFour) UpdateUtility(baseAgent agent.BaseAgent) {
	//
}

func NewAgentFour() agent.Strategy {
	return &AgentFour{
		bravery:      rand.Intn(10),
		utilityScore: make(map[string]int),
	}
}
