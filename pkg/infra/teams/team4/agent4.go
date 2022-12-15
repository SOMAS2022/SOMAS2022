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

// Struct for AgentFour
type AgentFour struct {
	HP           int
	ST           int
	AT           int
	SH           int
	C            int
	bravery      int
	uR           map[commons.ID]int
	uP           map[commons.ID]int
	uC           map[commons.ID]int
	utilityScore map[commons.ID]int
	TSN          []commons.ID
}

// bravery and utility score defined
func NewAgentFour() agent.Strategy {
	return &AgentFour{
		bravery:      rand.Intn(10),
		utilityScore: make(map[string]int),
	}
}

// *********************************** STRATEGY INTERFACE FUNCTIONS ***********************************

// we always pick our best shield
func (a *AgentFour) HandleUpdateShield(baseAgent agent.BaseAgent) decision.ItemIdx {
	view := baseAgent.View()
	agentState := view.AgentState()

	shields := agentState.Shields

	if shields.Len() > 0 { // if we have shields

		largestItemValue := 0
		itemIndex := 0

		for i := 0; i < shields.Len(); i++ {
			if shields.Get(i).Value() > largestItemValue {
				largestItemValue = shields.Get(i).Value()
				itemIndex = i
			}
		}

		return decision.ItemIdx(itemIndex)

	}

	// what to do if we have no shields?
	return decision.ItemIdx(0)
}

// we always pick the best weapon
func (a *AgentFour) HandleUpdateWeapon(baseAgent agent.BaseAgent) decision.ItemIdx {
	view := baseAgent.View()
	agentState := view.AgentState()

	weapons := agentState.Weapons

	if weapons.Len() > 0 { // if we have weapons

		largestItemValue := 0
		itemIndex := 0

		for i := 0; i < weapons.Len(); i++ {
			if weapons.Get(i).Value() > largestItemValue {
				largestItemValue = weapons.Get(i).Value()
				itemIndex = i
			}
		}

		return decision.ItemIdx(itemIndex)

	}

	// what to do if we have no weapons?
	return decision.ItemIdx(0)
}

// Define and update the attributes for agent four
func (a *AgentFour) UpdateInternalState(baseAgent agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	a.UpdateUtility(baseAgent)
	a.HP = int(baseAgent.AgentState().Hp)
	a.ST = int(baseAgent.AgentState().Stamina)
	a.AT = int(baseAgent.AgentState().Attack)
	//a.SH = int(baseAgent.AgentState().Shields)
	a.C = 0
}

// *********************************** ELECTION INTERFACE FUNCTIONS ***********************************

func (a *AgentFour) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	//
}

func (a *AgentFour) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	//
}

func (a *AgentFour) HandleElectionBallot(baseAgent agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	//
}

// *********************************** FIGHT INTERFACE FUNCTIONS ***********************************

func (a *AgentFour) HandleFightInformation(_ message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	//
}

func (a *AgentFour) HandleFightRequest(_ message.TaggedRequestMessage[message.FightRequest], _ *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	//
}

func (a *AgentFour) FightResolution(baseAgent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]]) immutable.Map[commons.ID, decision.FightAction] {
	//
}

func (a *AgentFour) HandleFightProposal(m message.Proposal[decision.FightAction], baseAgent agent.BaseAgent) decision.Intent {
	//
}

func (a *AgentFour) HandleFightProposalRequest(_ message.Proposal[decision.FightAction], _ agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) bool {
	//
}

func (a *AgentFour) FightActionNoProposal(baseAgent agent.BaseAgent) decision.FightAction {
	//
}

func (a *AgentFour) FightAction(baseAgent agent.BaseAgent, proposedAction decision.FightAction, acceptedProposal message.Proposal[decision.FightAction]) decision.FightAction {
	//
}

// *********************************** LOOT INTERFACE FUNCTIONS ***********************************

func (a *AgentFour) HandleLootInformation(m message.TaggedInformMessage[message.LootInform], agent agent.BaseAgent) {
	//
}

func (a *AgentFour) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	//
}

func (a *AgentFour) HandleLootProposal(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) decision.Intent {
	//
}

func (a *AgentFour) HandleLootProposalRequest(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	//
}

func (a *AgentFour) LootAllocation(baseAgent agent.BaseAgent, proposal message.Proposal[decision.LootAction]) immutable.Map[commons.ID, immutable.SortedMap[commons.ItemID, struct{}]] {
	//
}

func (a *AgentFour) LootActionNoProposal(baseAgent agent.BaseAgent) immutable.SortedMap[commons.ItemID, struct{}] {
	//
}

func (a *AgentFour) LootAction(baseAgent agent.BaseAgent, proposedLoot immutable.SortedMap[commons.ItemID, struct{}], acceptedProposal message.Proposal[decision.LootAction]) immutable.SortedMap[commons.ItemID, struct{}] {
	//
}

// *********************************** HPPOOL INTERFACE FUNCTIONS ***********************************

// HP pool donation
func (a *AgentFour) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	C_thresh_HP := 1
	donation := 0
	// If our health is > 50% and we feel generous then donate some (max 20%) HP
	if a.HP > 0.8 && a.C < C_thresh_HP {
		donation = (a.HP * 20) / 100
		a.C += 1
	} else {
		donation = 0
	}
	return uint(donation)
}

// *********************************** TRADE INTERFACE FUNCTIONS ***********************************

func (a *AgentFour) HandleTradeNegotiation(theAgent agent.BaseAgent, m message.TradeInfo) message.TradeMessage {
	return message.TradeRequest{}
	// respond to requests

	// make requests

}

// *********************************** OTHER FUNCTIONS ***********************************

// Random value generator
func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func (a *AgentFour) HPLevels(agent agent.BaseAgent) {
	view := agent.View()
	agentState := view.AgentState()
	agentStateIterator := agentState.Iterator()
	sliceOfAgentsWithLowHealth := make([]commons.ID, 0)
	sliceOfAgentsWithMidHealth := make([]commons.ID, 0)
	sliceOfAgentsWithHighHealth := make([]commons.ID, 0)
	for !agentStateIterator.Done() {
		key, value, _ := agentStateIterator.Next()
		switch {
		case value.Hp == state.HealthRange(state.LowHealth):
			sliceOfAgentsWithLowHealth = append(sliceOfAgentsWithLowHealth, key)
		case value.Hp == state.HealthRange(state.MidHealth):
			sliceOfAgentsWithMidHealth = append(sliceOfAgentsWithMidHealth, key)
		case value.Hp == state.HealthRange(state.HighHealth):
			sliceOfAgentsWithHighHealth = append(sliceOfAgentsWithHighHealth, key)
		}
	}
}

func (a *AgentFour) STLevels(agent agent.BaseAgent) {
	view := agent.View()
	agentState := view.AgentState()
	agentStateIterator := agentState.Iterator()
	sliceOfAgentsWithLowST := make([]commons.ID, 0)
	sliceOfAgentsWithMidST := make([]commons.ID, 0)
	sliceOfAgentsWithHighST := make([]commons.ID, 0)
	for !agentStateIterator.Done() {
		key, value, _ := agentStateIterator.Next()
		switch {
		case value.Hp == state.HealthRange(state.LowStamina):
			sliceOfAgentsWithLowST = append(sliceOfAgentsWithLowST, key)
		case value.Hp == state.HealthRange(state.MidStamina):
			sliceOfAgentsWithMidST = append(sliceOfAgentsWithMidST, key)
		case value.Hp == state.HealthRange(state.HighStamina):
			sliceOfAgentsWithHighST = append(sliceOfAgentsWithHighST, key)
		}
	}
}

// Attack-Defend-Cower Strat
func (a *AgentFour) AttackDefendCower(baseAgent agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	var action decision.FightAction
	damage := logging.LevelStats.MonsterAttack / (len(decision.NewImmutableFightResult().AttackingAgents()) + len(decision.NewImmutableFightResult().ShieldingAgents())) //might be incorrect formula, change!
	if a.HP > (damage + 1) {
		if state.TotalAttack >= state.TotalDefense()*0.8 {
			action = decision.Attack
		} else {
			action = decision.Defend
		}
	} else {
		action = decision.Cower
	}
}

// FightManifesto
func (a *AgentFour) FightManifesto(agent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]]) immutable.Map[commons.ID, decision.FightAction] {
	view := agent.View()
	builder := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
	var manifesto_decision decision.FightAction
	rand_prob := randFloats(0, 1, 1)
	get_HP_levels := a.HPLevels()
	get_ST_levels := a.STLevels()
	thresh_fight := 0.3

	ratio_agents_HPLow := get_HP_levels.sliceOfAgentsWithLowHealth / logging.LevelStats.NumberOfAgents
	ratio_agents_HPNormal := get_HP_levels.sliceOfAgentsWithMidHealth / logging.LevelStats.NumberOfAgents
	ratio_agents_HPHigh := get_HP_levels.sliceOfAgentsWithHighHealth / logging.LevelStats.NumberOfAgents

	ratio_agents_STLow := get_ST_levels.sliceOfAgentsWithLowST / logging.LevelStats.NumberOfAgents
	ratio_agents_STNormal := get_ST_levels.sliceOfAgentsWithMidST / logging.LevelStats.NumberOfAgents
	ratio_agents_STHigh := get_ST_levels.sliceOfAgentsWithHighST / logging.LevelStats.NumberOfAgents

	thresh_attack := decision.NewImmutableFightResult().AttackSum() / logging.LevelStats.NumberOfAgents
	thresh_defend := decision.NewImmutableFightResult().ShieldSum() / logging.LevelStats.NumberOfAgents

	threshold_fight_HP := ratio_agents_HPLow*(250) + ratio_agents_HPNormal*(500) + ratio_agents_HPHigh*(750)
	threshold_fight_ST := ratio_agents_STLow*(500) + ratio_agents_STNormal*(1000) + ratio_agents_STHigh*(1500)

	FightMethod := a.AttackDefendCower()

	for _, id := range commons.ImmutableMapKeys(view.AgentState()) {
		if a.HP > threshold_fight_HP && a.ST > threshold_fight_ST {
			if state.TotalAttack() > thresh_attack && state.TotalDefense() > thresh_defend {
				switch {
				case rand_prob >= 0.4:
					manifesto_decision = decision.Defend
				case rand_prob <= 0.6:
					manifesto_decision = decision.Attack
				}
			}
			if state.TotalAttack() > thresh_attack {
				manifesto_decision = decision.Attack
			}
			if state.TotalDefense() > thresh_defend {
				manifesto_decision = decision.Defend
			}
		} else {
			manifesto_decision = decision.Cower
		}

		if FightMethod.action == decision.Cower && manifesto_decision == decision.Attack {
			if rand_prob < thresh_fight {
				threshold_fight_HP = a.HP + 10
				a.C -= 1
			} else {
				a.C += 1
			}
			return *builder.Map()
		}
	}
}

//Alternative FightManifesto method
// func FightManifesto(agents map[commons.ID]agent.Agent) map[commons.ID]decision.FightAction {
// 	decisionMap := make(map[commons.ID]decision.FightAction)

// 	for i := range agents {
// 		decisionMap[i] = decision.Defend
// 	}

// 	return decisionMap
// }

// Equip weapon and shield
func EquipWeaponShield(iterator commons.Iterator[state.Item], ids []commons.ID, lootAllocation map[commons.ID][]commons.ItemID) {
	//
}

// Replenish Health
// func (a *AgentFour) RepenlishHealth(baseAgent agent.BaseAgent) uint {
// 	aux_var := (Y/(0.5*N_surv) - a.SH)
// 	for (a.HP < aux_var && we have a HP potion) {
// 		//#use health potion -> health_potion = health_potion - 1
// 	}
// }

// Replenish Stamina
// func (a *AgentFour) RepenlishStamina(baseAgent agent.BaseAgent) uint {
// 	aux_var := (Y/(0.5*N_surv) - a.SH)
// 	for ((a.ST < TotalAttack || a.ST < TotalDefense) && (we have a ST potion)) {
// 		#use stamina potion
// 	}
// }

// // FUNCTIONS COPIED //

// func (a *AgentFour) CurrentAction() decision.FightAction {
// 	//
// }

// func allocateRandomly(iterator commons.Iterator[state.Item], ids []commons.ID, lootAllocation map[commons.ID][]commons.ItemID) {
// 	//
// }

func (a *AgentFour) UpdateUtility(baseAgent agent.BaseAgent) {
	//???????
}
