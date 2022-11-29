package main

import (
	"fmt"
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

var viewPtr *state.View = &state.View{}
var globalState *state.State
var agentMap *map[commons.ID]agent.Agent
var gameConfig *config.GameConfig

/*
	Init Helpers
*/

func updateView(ptr *state.View, globalState *state.State) {
	*ptr = globalState.ToView()
}

func initGame() {
	err := godotenv.Load()
	if err != nil {
		logging.Log(logging.Error, nil, "No .env file located, using defaults")
	}

	stages.Mode = config.EnvToString("MODE", "default")

	config := stages.InitGameConfig()
	gameConfig = &config
	defStrategyMap := stages.ChooseDefaultStrategyMap(InitAgentMap)
	numAgents, agents, agentStateMap := stages.InitAgents(defStrategyMap, config, viewPtr)
	gameConfig.InitialNumAgents = numAgents

	globalState = &state.State{
		MonsterHealth: gamemath.CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.Stamina, gameConfig.NumLevels, 1),
		MonsterAttack: gamemath.CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.Stamina, gameConfig.ThresholdPercentage, gameConfig.NumLevels, 1),
		AgentState:    agentStateMap,
	}
	agentMap = &agents
}

/*
	Communication Helpers
*/

func addCommsChannels() (res map[commons.ID]chan message.TaggedMessage) {
	keys := make([]commons.ID, len(*agentMap))
	res = make(map[commons.ID]chan message.TaggedMessage)
	i := 0
	for k := range *agentMap {
		keys[i] = k
		i++
	}

	for _, key := range keys {
		res[key] = make(chan message.TaggedMessage, 100)
	}
	immutableMap := createImmutableMapForChannels(res)
	for id, a := range *agentMap {
		a.BaseAgent = agent.NewBaseAgent(agent.NewCommunication(res[id], *immutableMap.Delete(id)), id, a.BaseAgent.Name(), viewPtr)
		(*agentMap)[id] = a
	}
	return
}

func createImmutableMapForChannels[K constraints.Ordered, V any](peerChannels map[K]chan V) immutable.Map[K, chan<- V] {
	builder := immutable.NewMapBuilder[K, chan<- V](nil)
	for pId, channel := range peerChannels {
		builder.Set(pId, channel)
	}
	return *builder.Map()
}

/*
	Election Helpers
*/

func runElection() uint {
	electedAgent, manifesto, percentage := election.HandleElection(globalState, *agentMap, decision.VotingStrategy(gameConfig.VotingStrategy), gameConfig.VotingPreferences)
	termLeft := manifesto.TermLength()
	globalState.LeaderManifesto = manifesto
	globalState.CurrentLeader = electedAgent
	logging.Log(logging.Info, nil, fmt.Sprintf("[%d] New leader has been elected %s with %d%% of the vote", globalState.CurrentLevel, electedAgent, percentage))
	updateView(viewPtr, globalState)
	return termLeft
}

func runConfidenceVote(termLeft uint) uint {
	votes := make(map[decision.Intent]uint)
	for _, a := range *agentMap {
		votes[a.Strategy.HandleConfidencePoll(a.BaseAgent)]++
	}
	leader := (*agentMap)[globalState.CurrentLeader]
	leaderName := leader.BaseAgent.Name()
	logging.Log(logging.Info, logging.LogField{
		"positive":  votes[decision.Positive],
		"negative":  votes[decision.Negative],
		"abstain":   votes[decision.Abstain],
		"threshold": globalState.LeaderManifesto.OverthrowThreshold(),
		"leader":    globalState.CurrentLeader,
		"team":      leaderName,
	}, "Confidence Vote")
	if 100*votes[decision.Negative]/(votes[decision.Negative]+votes[decision.Positive]) > globalState.LeaderManifesto.OverthrowThreshold() {
		logging.Log(logging.Info, nil, fmt.Sprintf("%s got ousted", globalState.CurrentLeader))
		termLeft = runElection()
	}
	return termLeft
}

/*
	Fight Helpers
*/

func damageCalculation(fightRoundResult decision.FightResult) {
	if len(fightRoundResult.CoweringAgents) != len(*agentMap) {
		globalState.MonsterHealth = commons.SaturatingSub(globalState.MonsterHealth, fightRoundResult.AttackSum)
		if globalState.MonsterHealth > 0 && fightRoundResult.ShieldSum < globalState.MonsterAttack {
			agentsFighting := append(fightRoundResult.AttackingAgents, fightRoundResult.ShieldingAgents...)
			damageTaken := globalState.MonsterAttack - fightRoundResult.ShieldSum
			fight.DealDamage(damageTaken, agentsFighting, *agentMap, globalState)
			// TODO: Monster disruptive ability
		}
	} else {
		damageTaken := globalState.MonsterAttack
		fight.DealDamage(damageTaken, fightRoundResult.CoweringAgents, *agentMap, globalState)
	}
	*viewPtr = globalState.ToView()
}
