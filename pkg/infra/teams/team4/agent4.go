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

// HP pool donation
func (a *AgentFour) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	C := 0
	C_thresh_HP := 1
	// If our health is > 50% and we feel generous then donate some (max 20%) HP
	if a.HP > 0.8 && C < C_thresh_HP {
		return uint(rand.Intn((a.HP * 20) / 100))
		C +=1
	}
	return 0
}

// Replenish Health
func (a *AgentFour) RepenlishHealth(baseAgent agent.BaseAgent) uint {
	aux_var := (Y/(0.5*N_surv) - a.SH)
	for (a.HP < aux_var && we have a HP potion) {
		//#use health potion -> health_potion = health_potion - 1
	}
}

// Replenish Stamina
func (a *AgentFour) RepenlishStamina(baseAgent agent.BaseAgent) uint {
	aux_var := (Y/(0.5*N_surv) - a.SH)
	for ((a.ST < TotalAttack || a.ST < TotalDefense) && (we have a ST potion)) {
		#use stamina potion
	}
}

// FUNCTIONS COPIED //

func (a *AgentFour) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	//
	threshold_fight_HP = ratio_agents_HPLow*(250) + ratio_agents_HPNormal*(500) + ratio_agents_HPHigh*(750)
	threshold_fight_ST = ratio_agents_STLow*(500) + ratio_agents_STNormal*(1000) + ratio_agents_STHigh*(1500)

	if (HP > threshold_fight_HP AND ST > threshold_fight_ST){
		fight
	}
	else{
		cowers
	}

	if our decision == cower and manifesto decision for us == fight{
	if the random_prob < Thresh_fight{            # bad boy
			threshold_fight_HP = HP + 10
			U -= 1
	}
	else{                                                         # accept decision from manifesto
			U += 1
	}

	}
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
	a.UpdateUtility(baseAgent)
	a.HP = int(baseAgent.AgentState().Hp)
	a.ST = int(baseAgent.AgentState().Stamina)
	a.AT = int(baseAgent.AgentState().Attack)
	a.SH = int(baseAgent.AgentState().Shields)
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
