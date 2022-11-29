package main

import (
	"flag"
	"fmt"
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/example"
	gamemath "infra/game/math"
	"infra/game/message"
	"infra/game/stage/election"
	"infra/game/stage/fight"
	"infra/game/stages"
	"infra/game/state"
	"infra/logging"
	"math"
	"math/rand"

	"golang.org/x/exp/constraints"

	"github.com/benbjohnson/immutable"
	"github.com/joho/godotenv"
)

var InitAgentMap = map[commons.ID]func() agent.Strategy{
	"RANDOM": example.NewRandomAgent,
}

/*
Each agent is initialised with these modes of communication
Whereby the peer channels enable bidirectional p2p communication
Enables peers to send and receive messages with broadcasting possible via non-blocking sends
*/

func main() {
	// define flags
	useJSONFormatter := flag.Bool("j", false, "Whether to output logs in JSON")
	debug := flag.Bool("d", false, "Whether to run in debug mode. If false, only logs with level info or above will be shown")
	flag.Parse()

	logging.InitLogger(*useJSONFormatter, *debug)
	viewPtr := initViews()
	agentMap, globalState, gameConfig := initGame(viewPtr)
	gameLoop(globalState, agentMap, gameConfig, viewPtr)
}

func initViews() *state.View {
	view := &state.View{}
	return view
}

func gameLoop(globalState state.State, agentMap map[commons.ID]agent.Agent, gameConfig config.GameConfig, ptr *state.View) {
	var decisionMap map[commons.ID]decision.FightAction
	var channelsMap map[commons.ID]chan message.TaggedMessage
	var termLeft uint
	channelsMap = addCommsChannels(agentMap, ptr)
	*ptr = globalState.ToView()

	updateView := func() {
		*ptr = globalState.ToView()
	}

	for globalState.CurrentLevel = 1; globalState.CurrentLevel < (gameConfig.NumLevels + 1); globalState.CurrentLevel++ {
		for globalState.CurrentLevel = 0; globalState.CurrentLevel < gameConfig.NumLevels; globalState.CurrentLevel++ {
			updateView()

			if termLeft == 0 {
				termLeft = runElection(&globalState, agentMap, gameConfig)
			} else {
				termLeft = runConfidenceVote(&globalState, agentMap, gameConfig, termLeft)
			}

			for globalState.MonsterHealth != 0 {
				decisionMapView := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
				for u, action := range decisionMap {
					decisionMapView.Set(u, action)
				}
				fightRoundResult := decision.FightResult{Choices: stages.AgentFightDecisions(&globalState, agentMap, *decisionMapView.Map(), channelsMap)}
				updateView()
				fight.HandleFightRound(&globalState, gameConfig.StartingHealthPoints, &fightRoundResult)
				updateView()

				logging.Log(logging.Info, logging.LogField{
					"currLevel":     globalState.CurrentLevel,
					"monsterHealth": globalState.MonsterHealth,
					"monsterDamage": globalState.MonsterAttack,
					"numCoward":     len(fightRoundResult.CoweringAgents),
					"attackSum":     fightRoundResult.AttackSum,
					"shieldSum":     fightRoundResult.ShieldSum,
					"numAgents":     len(agentMap),
				}, "Battle Summary")

				damageCalculation(&globalState, agentMap, fightRoundResult)
				updateView()

				channelsMap = addCommsChannels(agentMap, ptr)

				if float64(len(agentMap)) < math.Ceil(float64(gameConfig.ThresholdPercentage)*float64(gameConfig.InitialNumAgents)) {
					logging.Log(logging.Info, nil, fmt.Sprintf("Lost on level %d  with %d remaining", globalState.CurrentLevel, len(agentMap)))
					return
				}
			}

			//todo: Results in infinite game run-through
			globalState.MonsterHealth = gamemath.CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.Stamina, gameConfig.NumLevels, globalState.CurrentLevel+1)
			globalState.MonsterAttack = gamemath.CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.Stamina, gameConfig.ThresholdPercentage, gameConfig.NumLevels, globalState.CurrentLevel+1)

			// TODO: End of Level looting and trading
			// FIXME: This loot allocation should not stay for long!
			weaponLoot := make([]uint, len(agentMap))
			shieldLoot := make([]uint, len(agentMap))

			for i := range weaponLoot {
				weaponLoot[i] = globalState.CurrentLevel * uint(rand.Intn(3))
				shieldLoot[i] = globalState.CurrentLevel * uint(rand.Intn(3))
			}

			logging.Log(logging.Info, nil, fmt.Sprintf("------------------------------ Level %d Ended ----------------------------", globalState.CurrentLevel))
			newGlobalState := stages.AgentLootDecisions(globalState, agentMap, weaponLoot, shieldLoot)
			updateView()

			// TODO: Add verification if needed
			globalState = newGlobalState
			termLeft--
		}
		logging.Log(logging.Info, nil, fmt.Sprintf("Congratulations, The Peasants have escaped the pit with %d remaining.", len(agentMap)))
	}
}

func damageCalculation(globalState *state.State, agentMap map[commons.ID]agent.Agent, fightRoundResult decision.FightResult) {
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
}

func runConfidenceVote(globalState *state.State, agentMap map[commons.ID]agent.Agent, gameConfig config.GameConfig, termLeft uint) uint {
	votes := make(map[decision.Intent]uint)
	for _, a := range agentMap {
		votes[a.Strategy.HandleConfidencePoll(a.BaseAgent)]++
	}
	logging.Log(logging.Info, logging.LogField{
		"positive":  votes[decision.Positive],
		"negative":  votes[decision.Negative],
		"abstain":   votes[decision.Abstain],
		"threshold": globalState.LeaderManifesto.OverthrowThreshold(),
		"agent":     globalState.CurrentLeader,
	}, "Confidence Vote")
	if 100*votes[decision.Negative]/(votes[decision.Negative]+votes[decision.Positive]) > globalState.LeaderManifesto.OverthrowThreshold() {
		logging.Log(logging.Info, nil, fmt.Sprintf("%s got ousted", globalState.CurrentLeader))
		termLeft = runElection(globalState, agentMap, gameConfig)
	}
	return termLeft
}

func runElection(globalState *state.State, agentMap map[commons.ID]agent.Agent, gameConfig config.GameConfig) uint {
	electedAgent, manifesto, percentage := election.HandleElection(globalState, agentMap, decision.VotingStrategy(gameConfig.VotingStrategy), gameConfig.VotingPreferences)
	termLeft := manifesto.TermLength()
	globalState.LeaderManifesto = manifesto
	globalState.CurrentLeader = electedAgent
	logging.Log(logging.Info, nil, fmt.Sprintf("[%d] New leader has been elected %s with %d%% of the vote", globalState.CurrentLevel, electedAgent, percentage))
	return termLeft
}

func initGame(vPtr *state.View) (map[commons.ID]agent.Agent, state.State, config.GameConfig) {
	err := godotenv.Load()
	if err != nil {
		logging.Log(logging.Error, nil, "No .env file located, using defaults")
	}

	stages.Mode = config.EnvToString("MODE", "default")

	gameConfig := stages.InitGameConfig()
	defStrategyMap := stages.ChooseDefaultStrategyMap(InitAgentMap)
	numAgents, agentMap, agentStateMap := stages.InitAgents(defStrategyMap, gameConfig, vPtr)
	gameConfig.InitialNumAgents = numAgents

	globalState := state.State{
		MonsterHealth: gamemath.CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.Stamina, gameConfig.NumLevels, 1),
		MonsterAttack: gamemath.CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.Stamina, gameConfig.ThresholdPercentage, gameConfig.NumLevels, 1),
		AgentState:    agentStateMap,
	}

	return agentMap, globalState, gameConfig
}

func addCommsChannels(agentMap map[commons.ID]agent.Agent, viewPtr *state.View) (res map[commons.ID]chan message.TaggedMessage) {
	keys := make([]commons.ID, len(agentMap))
	res = make(map[commons.ID]chan message.TaggedMessage)
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
		agentMap[id] = a
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
