package team5

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/teams/team5/commons5"

	"github.com/benbjohnson/immutable"
)

type agent5 struct {
	//FightInformation
	// InternalState   internalState
	LootInformation commons5.Loot
	SocialNetwork   SocialNetwork
	LeaderCombo     uint
}

func (t *agent5) CreateManifesto(baseAgent agent.BaseAgent) *decision.Manifesto {
	returnType := decision.NewManifesto(false, true, t.LeaderCombo+1, 50)
	return returnType
}
func (t *agent5) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	var returnType decision.Intent
	return returnType
}
func (t *agent5) HandleElectionBallot(baseAgent agent.BaseAgent, params *decision.ElectionParams) decision.Ballot {
	var returnType decision.Ballot
	return returnType
}

func (t *agent5) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], ba agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
}
func (t *agent5) HandleFightRequest(m message.TaggedRequestMessage[message.FightRequest], log *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	var returnType message.FightInform
	return returnType
}
func (t *agent5) FightResolution(agent agent.BaseAgent) commons.ImmutableList[proposal.Rule[decision.FightAction]] {
	var returnType commons.ImmutableList[proposal.Rule[decision.FightAction]]
	return returnType
}
func (t *agent5) HandleFightProposal(proposal message.Proposal[decision.FightAction], ba agent.BaseAgent) decision.Intent {
	var returnType decision.Intent
	return returnType
}

// HandleFightProposalRequest only called as leader
func (t *agent5) HandleFightProposalRequest(proposal message.Proposal[decision.FightAction], ba agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool {
	var returnType bool
	return returnType
}

//This signature function is gone in the new infra update

// func (t *agent5) FightAction() decision.FightAction {
// 	var returnType decision.FightAction
// 	return returnType
// }

func (t *agent5) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent agent.BaseAgent) {
}
func (t *agent5) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	var returnType message.LootInform
	return returnType
}
func (t *agent5) HandleLootProposal(r message.Proposal[decision.LootAction], agent agent.BaseAgent) decision.Intent {
	var rr decision.Intent
	return rr
}
func (t *agent5) HandleLootProposalRequest(proposal message.Proposal[decision.LootAction], agent agent.BaseAgent) bool {
	var returnType bool
	return returnType
}
func (t *agent5) LootAllocation(agent agent.BaseAgent) immutable.Map[commons.ID, immutable.List[commons.ItemID]] {
	var returnType immutable.Map[commons.ID, immutable.List[commons.ItemID]]
	return returnType
}

func (t *agent5) LootAction() immutable.List[commons.ItemID] {
	var returnType immutable.List[commons.ItemID]
	return returnType
}

// HandleUpdateWeapon return the index of the weapon you want to use in AgentState.weapons
func (t *agent5) HandleUpdateWeapon(baseAgent agent.BaseAgent) decision.ItemIdx {
	var returnType decision.ItemIdx
	return returnType
}

// HandleUpdateShield return the index of the shield you want to use in AgentState.Shields
func (t *agent5) HandleUpdateShield(baseAgent agent.BaseAgent) decision.ItemIdx {
	var returnType decision.ItemIdx
	return returnType
}

func (t *agent5) UpdateInternalState(baseAgent agent.BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], voteResult *immutable.Map[decision.Intent, uint]) {
	//update leader combo
}

func (t *agent5) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	var returnType uint
	return returnType
}

// type internalState struct {
// 	AllAgents commons5.Agents
// 	//......
// }

func NewAgent5() agent.Strategy {
	var agent5 agent.Strategy

	return agent5
}
