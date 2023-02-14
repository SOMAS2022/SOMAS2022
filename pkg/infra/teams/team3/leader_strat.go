package team3

import (
	"math/rand"

	cmdline "infra/cmdLine"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"
	"infra/game/state"

	"github.com/benbjohnson/immutable"
)

type SanctionActivity struct {
	sanctionActive bool
	duration       int
}

func (s *SanctionActivity) makeSanction(length int) {
	s.sanctionActive = true
	s.duration = length
}

func (s *SanctionActivity) initialiseSanction() {
	s.sanctionActive = false
	s.duration = 0
}

func (s *SanctionActivity) agentIsSanctioned() bool {
	return s.sanctionActive
}

func (s *SanctionActivity) updateSanction() {
	if !s.sanctionActive {
		return
	}

	if s.duration > 0 {
		s.duration--
	} else {
		s.sanctionActive = false
	}

}

// Manifesto
func (a *AgentThree) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	// Submit Manifesto?
	submitManifesto := rand.Intn(100)
	if submitManifesto < a.personality {
		// Enter manifesto logic here for creating a manifesto
		manifesto := decision.NewManifesto(false, false, 10, 50)
		return manifesto
	} else {
		// submit a manifesto with a term length of 0
		manifesto := decision.NewManifesto(false, false, 0, 0)
		return manifesto
	}
}

// Leader function to grant the floor?
func (a *AgentThree) HandleFightProposalRequest(_ message.Proposal[decision.FightAction], _ agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

// Fight imposition
func (a *AgentThree) FightResolution(baseAgent agent.BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]], proposedActions immutable.Map[string, decision.FightAction]) immutable.Map[string, decision.FightAction] {
	view := baseAgent.View()
	// AS := baseAgent.AgentState()
	builder := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)

	for _, id := range commons.ImmutableMapKeys(view.AgentState()) {
		var fightAction decision.FightAction

		// Check for our agent and assign what we want to do
		if id == baseAgent.ID() {
			action := a.FightActionNoProposal(baseAgent)
			fightAction = action
		} else {
			switch rand.Intn(3) {
			case 0:
				fightAction = decision.Attack
			case 1:
				fightAction = decision.Defend
			default:
				fightAction = decision.Cower
			}
		}
		builder.Set(id, fightAction)
	}
	return *builder.Map()
}

func (a *AgentThree) HandleFightRequest(_ message.TaggedRequestMessage[message.FightRequest], _ *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (a *AgentThree) HandleLootRequest(m message.TaggedRequestMessage[message.LootRequest]) message.LootInform {
	//TODO implement me
	panic("implement me")
}

func (a *AgentThree) HandleLootProposalRequest(_ message.Proposal[decision.LootAction], _ agent.BaseAgent) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

func (a *AgentThree) LootAllocation(baseAgent agent.BaseAgent, proposal message.Proposal[decision.LootAction], proposedAllocations map[commons.ID]map[commons.ItemID]struct{}) immutable.Map[string, immutable.SortedMap[string, struct{}]] {
	lootAllocation := make(map[commons.ID][]commons.ItemID)
	view := baseAgent.View()
	ids := commons.ImmutableMapKeys(view.AgentState())
	iterator := baseAgent.Loot().Weapons().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = baseAgent.Loot().Shields().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = baseAgent.Loot().HpPotions().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	iterator = baseAgent.Loot().StaminaPotions().Iterator()
	allocateRandomly(iterator, ids, lootAllocation)
	mMapped := make(map[commons.ID]immutable.SortedMap[commons.ItemID, struct{}])
	for id, itemIDS := range lootAllocation {
		mMapped[id] = commons.ListToImmutableSortedSet(itemIDS)
	}
	return commons.MapToImmutable(mMapped)
}

// func allocateRandomly(iterator commons.Iterator[state.Item], ids []commons.ID, lootAllocation map[commons.ID][]commons.ItemID) {
// 	for !iterator.Done() {
// 		next, _ := iterator.Next()
// 		toBeAllocated := ids[rand.Intn(len(ids))]
// 		if l, ok := lootAllocation[toBeAllocated]; ok {
// 			l = append(l, next.Id())
// 			lootAllocation[toBeAllocated] = l
// 		} else {
// 			l := make([]commons.ItemID, 0)
// 			l = append(l, next.Id())
// 			lootAllocation[toBeAllocated] = l
// 		}
// 	}
// }

func allocateRandomly(iterator commons.Iterator[state.Item], ids []commons.ID, lootAllocation map[commons.ID][]commons.ItemID) {
	for !iterator.Done() {
		next, _ := iterator.Next()
		toBeAllocated := ids[rand.Intn(len(ids))]
		l, ok := lootAllocation[toBeAllocated]
		if !ok {
			l = make([]commons.ItemID, 0)
		}
		l = append(l, next.Id())
		lootAllocation[toBeAllocated] = l
	}
}

func (a *AgentThree) willSanctionConstant(agent agent.Agent) int {
	A := 0.8
	B := 0.2

	AS := agent.AgentState()
	id := agent.BaseAgent.ID()
	// were they a defector?
	D := float64(BoolToInt(AS.Defector.IsDefector()))

	// shift and scale the agent reputation
	S := (a.reputationMap[id] - 50) / 3
	// S := float64(rand.Intn(30) - 15)

	sanction := int(D * (A*float64(a.personality) + B*S))
	// fmt.Println(sanction)

	return sanction
}

func (a *AgentThree) sanctioningGraduated(agent agent.Agent) int {
	agentId := agent.ID()
	prevSanctions := a.sanctionHistory[agentId]
	mostRecentSanction := prevSanctions[len(prevSanctions)-1]
	return mostRecentSanction + 1
}

func (a *AgentThree) sanctioningDynamic(agent agent.Agent) int {
	return rand.Intn(5) + 1
}

func (a *AgentThree) updateSanctionHistory(agent agent.Agent, sanctionDuration int) {
	agentId := agent.ID()
	prevSanctions := a.sanctionHistory[agentId]
	updatedSanctions := append(prevSanctions, sanctionDuration)
	a.sanctionHistory[agentId] = updatedSanctions
}

func (a *AgentThree) createSanction(agent agent.Agent, length int) {
	agentId := agent.ID()
	sanction := SanctionActivity{}
	sanction.makeSanction(length)
	a.activeSanctionMap[agentId] = sanction
}

func (a *AgentThree) PruneAgentList(agentMap map[commons.ID]agent.Agent) map[commons.ID]agent.Agent {

	cmdParams := cmdline.CmdLineInits

	pruned := make(map[commons.ID]agent.Agent)
	for id, agent := range agentMap {

		currentSanction := a.activeSanctionMap[id]
		if currentSanction.agentIsSanctioned() {
			currentSanction.updateSanction()
			a.activeSanctionMap[id] = currentSanction
			continue
		}

		// Compare to 50 in order to sanction
		toSanctionOrNot := rand.Intn(100)
		if toSanctionOrNot > a.willSanctionConstant(agent) {
			pruned[id] = agent
		} else {
			// agent has been pruned. Choose sanction duration
			var sanctionDuration int
			if cmdParams.DynamicSanctions {
				sanctionDuration = a.sanctioningDynamic(agent)
			} else if cmdParams.GraduatedSanctions {
				sanctionDuration = a.sanctioningGraduated(agent)
			} else {
				sanctionDuration = cmdParams.FixedSanctionDuration
			}
			// register sanction to local map
			a.createSanction(agent, sanctionDuration)
			// update agent's sanction history
			a.updateSanctionHistory(agent, sanctionDuration)

		}
	}
	return pruned
}
