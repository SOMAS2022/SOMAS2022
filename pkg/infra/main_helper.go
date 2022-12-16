package main

import (
	"fmt"
	"sort"

	"github.com/google/uuid"

	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	gamemath "infra/game/math"
	"infra/game/message"
	"infra/game/stage/election"
	"infra/game/stage/fight"
	"infra/game/stages"
	"infra/game/state"
	"infra/logging"

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
		Defection:     gameConfig.Defection,
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
		a.SetCommunication(agent.NewCommunication(res[id], *immutableMap.Delete(id)))
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
		votes[a.Strategy.HandleConfidencePoll(*a.BaseAgent)]++
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

/*
	Hp Pool Helpers
*/

func checkHpPool() bool {
	if globalState.HpPool >= globalState.MonsterHealth {
		logging.Log(logging.Info, logging.LogField{
			"Original HP Pool":  globalState.HpPool,
			"Monster Health":    globalState.MonsterHealth,
			"HP Pool Remaining": globalState.HpPool - globalState.MonsterHealth,
		}, fmt.Sprintf("Skipping level %d through HP Pool", globalState.CurrentLevel))

		globalState.HpPool -= globalState.MonsterHealth
		globalState.MonsterHealth = 0
		return true
	}
	return false
}

func generateLootPool(numAgents uint) *state.LootPool {
	nWeapons, nShields := gamemath.GetEquipmentDistribution(numAgents)
	nHealthPotions, nStaminaPotions := gamemath.GetPotionDistribution(numAgents)

	makeItems := func(nItems uint, stats uint) *commons.ImmutableList[state.Item] {
		items := make([]state.Item, nItems)
		for i := uint(0); i < nItems; i++ {
			items[i] = *state.NewItem(uuid.NewString(), stats)
		}
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Value() > items[j].Value()
		})
		return commons.NewImmutableList(items)
	}

	recalculatedMonsterHealth := gamemath.CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.Stamina, gameConfig.NumLevels, globalState.CurrentLevel)

	return state.NewLootPool(
		// Weapons
		makeItems(nWeapons, gamemath.GetWeaponDamage(recalculatedMonsterHealth, numAgents)),
		// Shields
		makeItems(nShields, gamemath.GetShieldProtection(globalState.MonsterAttack, numAgents)),
		// Health Potions
		makeItems(nHealthPotions, gamemath.GetHealthPotionValue(globalState.MonsterAttack, numAgents)),
		// Stamina Potions
		makeItems(nStaminaPotions, gamemath.GetStaminaPotionValue(recalculatedMonsterHealth, numAgents)),
	)
}
