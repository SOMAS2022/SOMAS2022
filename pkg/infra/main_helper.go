package main

import (
	"fmt"
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/stage/election"
	"infra/game/stage/fight"
	"infra/game/stages"
	"infra/game/state"
	"infra/logging"
	"sync"

	gamemath "infra/game/math"

	"github.com/benbjohnson/immutable"
	"github.com/joho/godotenv"
	"golang.org/x/exp/constraints"
)

/*
	Package Variables
*/

var (
	viewPtr     = &state.View{}
	globalState *state.State
	agentMap    map[commons.ID]agent.Agent
	gameConfig  *config.GameConfig
)

/*
	Init Helpers
*/

func updateView(ptr *state.View, globalState *state.State) {
	*ptr = globalState.ToView()
}

func initGame() {
	if godotenv.Load() != nil {
		logging.Log(logging.Error, nil, "No .env file located, using defaults")
	}

	stages.Mode = config.EnvToString("MODE", "default")

	initGameConfig := stages.InitGameConfig()
	gameConfig = &initGameConfig
	defStrategyMap := stages.ChooseDefaultStrategyMap(InitAgentMap)
	numAgents, agents, agentStateMap, inventoryMap := stages.InitAgents(defStrategyMap, initGameConfig, viewPtr)
	gameConfig.InitialNumAgents = numAgents

	globalState = &state.State{
		MonsterHealth: gamemath.CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.Stamina, gameConfig.NumLevels, 1),
		MonsterAttack: gamemath.CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.Stamina, gameConfig.ThresholdPercentage, gameConfig.NumLevels, 1),
		AgentState:    agentStateMap,
		InventoryMap:  inventoryMap,
	}
	agentMap = agents
}

/*
	Communication Helpers
*/

func addCommsChannels() map[commons.ID]chan message.TaggedMessage {
	keys := make([]commons.ID, len(agentMap))
	res := make(map[commons.ID]chan message.TaggedMessage)
	i := 0
	for k := range agentMap {
		keys[i] = k
		i++
	}

	for _, key := range keys {
		res[key] = make(chan message.TaggedMessage, 100)
	}
	immutableMap := createImmutableMapForChannels(res)
	for id, a := range agentMap {
		a.BaseAgent = agent.NewBaseAgent(agent.NewCommunication(res[id], *immutableMap.Delete(id)), id, a.BaseAgent.Name(), viewPtr)
		(agentMap)[id] = a
	}
	return res
}

func createImmutableMapForChannels[K constraints.Ordered, V any](peerChannels map[K]chan V) immutable.Map[K, chan<- V] {
	builder := immutable.NewMapBuilder[K, chan<- V](nil)
	for pID, channel := range peerChannels {
		builder.Set(pID, channel)
	}
	return *builder.Map()
}

/*
	Election Helpers
*/

func runElection() uint {
	electedAgent, manifesto := election.HandleElection(globalState, agentMap, decision.VotingStrategy(gameConfig.VotingStrategy), gameConfig.VotingPreferences)
	termLeft := manifesto.TermLength()
	globalState.LeaderManifesto = manifesto
	globalState.CurrentLeader = electedAgent
	updateView(viewPtr, globalState)
	return termLeft
}

func runConfidenceVote(termLeft uint) (uint, map[decision.Intent]uint) {
	votes := make(map[decision.Intent]uint)
	for _, a := range agentMap {
		votes[a.Strategy.HandleConfidencePoll(a.BaseAgent)]++
	}
	leader := agentMap[globalState.CurrentLeader]
	leaderName := leader.BaseAgent.Name()

	logging.Log(logging.Info, logging.LogField{
		"positive":  votes[decision.Positive],
		"negative":  votes[decision.Negative],
		"abstain":   votes[decision.Abstain],
		"threshold": globalState.LeaderManifesto.OverthrowThreshold(),
		"leader":    globalState.CurrentLeader,
		"team":      leaderName,
	}, "Confidence Vote")

	if votes[decision.Negative]+votes[decision.Positive] == 0 {
		return termLeft, votes
	} else if 100*votes[decision.Negative]/(votes[decision.Negative]+votes[decision.Positive]) > globalState.LeaderManifesto.OverthrowThreshold() {
		logging.Log(logging.Info, nil, fmt.Sprintf("%s got ousted", globalState.CurrentLeader))
		termLeft = runElection()
	}
	return termLeft, votes
}

/*
	Fight Helpers
*/

func damageCalculation(fightRoundResult decision.FightResult) {
	if len(fightRoundResult.CoweringAgents) != len(agentMap) {
		globalState.MonsterHealth = commons.SaturatingSub(globalState.MonsterHealth, fightRoundResult.AttackSum)
		if globalState.MonsterHealth > 0 && fightRoundResult.ShieldSum < globalState.MonsterAttack {
			agentsFighting := append(fightRoundResult.AttackingAgents, fightRoundResult.ShieldingAgents...)
			damageTaken := globalState.MonsterAttack - fightRoundResult.ShieldSum
			fight.DealDamage(damageTaken, agentsFighting, agentMap, globalState)
			// TODO: Monster disruptive ability
		}
	} else {
		damageTaken := globalState.MonsterAttack
		fight.DealDamage(damageTaken, fightRoundResult.CoweringAgents, agentMap, globalState)
	}
	*viewPtr = globalState.ToView()
}

func updateInternalStates(immutableFightRounds *commons.ImmutableList[decision.ImmutableFightResult], votesResult *immutable.Map[decision.Intent, uint]) {
	var wg sync.WaitGroup
	for id, a := range agentMap {
		id := id
		a := a
		wg.Add(1)
		go func(wait *sync.WaitGroup) {
			a.HandleUpdateInternalState(globalState.AgentState[id], immutableFightRounds, votesResult)
			wait.Done()
		}(&wg)
	}
	wg.Wait()
}

/*
	Hp Pool Helpers
*/

type HpPoolDonation struct {
	agentId  commons.ID
	donation uint
}

func checkHpPool() {
	if globalState.HpPool >= globalState.MonsterAttack {
		logging.Log(logging.Info, logging.LogField{
			"Original HP Pool":  globalState.HpPool,
			"Monster Damage":    globalState.MonsterAttack,
			"HP Pool Remaining": globalState.HpPool - globalState.MonsterAttack,
		}, fmt.Sprintf("Skipping level %d through HP Pool", globalState.CurrentLevel))

		globalState.HpPool -= globalState.MonsterAttack
		globalState.MonsterHealth = 0
	}
}

func updateHpPool() {
	var wg sync.WaitGroup
	donationChan := make(chan HpPoolDonation, len(agentMap))
	for id, a := range agentMap {
		// fmt.Print(a.BaseAgent.AgentState().Hp)
		id := id
		a := a
		wg.Add(1)
		go func(wait *sync.WaitGroup, donationChan chan HpPoolDonation) {
			donation := a.HandleDonateToHpPool(globalState.AgentState[id])
			donationChan <- HpPoolDonation{id, donation}
			wait.Done()
		}(&wg, donationChan)
	}
	wg.Wait()

	sum := uint(0)
	for i := 0; i < len(agentMap); i++ {
		agentDonation := <-donationChan

		agentHp := globalState.AgentState[agentDonation.agentId].Hp
		if agentDonation.donation > agentHp {
			agentDonation.donation = agentHp
			delete(globalState.AgentState, agentDonation.agentId)
			delete(agentMap, agentDonation.agentId)
		}

		logging.Log(logging.Trace, logging.LogField{
			"Agent Donation": agentDonation,
			"Old Sum":        sum,
			"New Sum":        sum + agentDonation.donation,
		}, "HP Pool Donation")

		sum += agentDonation.donation
		if a, ok := globalState.AgentState[agentDonation.agentId]; ok {
			a.Hp = agentHp - agentDonation.donation
			globalState.AgentState[agentDonation.agentId] = a
		}
	}

	logging.Log(logging.Info, logging.LogField{
		"Old HP Pool":           globalState.HpPool,
		"HP Donated This Round": sum,
		"New Hp Pool":           globalState.HpPool + sum,
	}, "HP Pool Donation")

	globalState.HpPool += sum
}
