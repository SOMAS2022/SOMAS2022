package main

import (
	"flag"
	"fmt"
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	gamemath "infra/game/math"
	"infra/game/message"
	"infra/game/stage/fight"
	"infra/game/stages"
	"infra/game/state"
	"infra/logging"
	"math"
	"math/rand"

	"github.com/benbjohnson/immutable"
	"github.com/joho/godotenv"
)

var InitAgentMap = map[commons.ID]agent.Strategy{
	"RANDOM": agent.NewRandomAgent(),
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

	agentMap, globalState, gameConfig := init_game()
	gameLoop(globalState, agentMap, gameConfig)
}

func gameLoop(globalState state.State, agentMap map[commons.ID]agent.Agent, gameConfig config.GameConfig) {
	var decisionMap map[commons.ID]decision.FightAction
	var channelsMap map[commons.ID]chan message.TaggedMessage
	channelsMap = addCommsChannels(agentMap)
	for globalState.CurrentLevel = 1; globalState.CurrentLevel < (gameConfig.NumLevels + 1); globalState.CurrentLevel++ {
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		for globalState.MonsterHealth != 0 {
			decisionMapView := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
			for u, action := range decisionMap {
				decisionMapView.Set(u, action)
			}
			fightRoundResult := decision.FightResult{Choices: stages.AgentFightDecisions(globalState.ToView(), agentMap, *decisionMapView.Map(), channelsMap)}
			fight.HandleFightRound(&globalState, gameConfig.StartingHealthPoints, &fightRoundResult)
			// decisionMap = dMap

			logging.Log(logging.Info, logging.LogField{
				"currLevel":     globalState.CurrentLevel,
				"monsterHealth": globalState.MonsterHealth,
				"monsterDamage": globalState.MonsterAttack,
				"numCoward":     len(fightRoundResult.CoweringAgents),
				"attackSum":     fightRoundResult.AttackSum,
				"shieldSum":     fightRoundResult.ShieldSum,
				"numAgents":     len(agentMap),
			}, "Battle Summary")

			if len(fightRoundResult.CoweringAgents) != len(agentMap) {
				globalState.MonsterHealth = commons.SaturatingSub(globalState.MonsterHealth, fightRoundResult.AttackSum)
				if globalState.MonsterHealth > 0 && fightRoundResult.ShieldSum < globalState.MonsterAttack {
					agentsFighting := append(fightRoundResult.AttackingAgents, fightRoundResult.ShieldingAgents...)
					damageTaken := globalState.MonsterAttack - fightRoundResult.ShieldSum
					fight.DealDamage(damageTaken, agentsFighting, agentMap, &globalState)
					// TODO: Monster disruptive ability
				}
			} else {
				damageTaken := globalState.MonsterAttack
				fight.DealDamage(damageTaken, fightRoundResult.CoweringAgents, agentMap, &globalState)
			}

			channelsMap = addCommsChannels(agentMap)

			if float64(len(agentMap)) < math.Ceil(float64(gameConfig.ThresholdPercentage)*float64(gameConfig.InitialNumAgents)) {
				logging.Log(logging.Info, nil, fmt.Sprintf("Lost on level %d  with %d remaining", globalState.CurrentLevel, len(agentMap)))
				return
			}
		}
		logging.Log(logging.Info, nil, fmt.Sprintf("------------------------------ Level %d Ended ----------------------------", globalState.CurrentLevel))
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

		newGlobalState := stages.AgentLootDecisions(globalState, agentMap, weaponLoot, shieldLoot)
		// TODO: Add verification if needed
		globalState = newGlobalState
	}
	logging.Log(logging.Info, nil, fmt.Sprintf("Congratulations, The Peasants have escaped the pit with %d remaining.", len(agentMap)))
}

func init_game() (map[commons.ID]agent.Agent, state.State, config.GameConfig) {
	err := godotenv.Load()
	if err != nil {
		logging.Log(logging.Error, nil, "No .env file located, using defaults")
	}

	stages.Mode = config.EnvToString("MODE", "default")

	gameConfig := stages.InitGameConfig()
	defStratergyMap := stages.ChooseDefaultStratergyMap(InitAgentMap)
	numAgents, agentMap, agentStateMap := stages.InitAgents(defStratergyMap, gameConfig)
	gameConfig.InitialNumAgents = numAgents

	globalState := state.State{
		MonsterHealth: gamemath.CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.Stamina, gameConfig.NumLevels, 1),
		MonsterAttack: gamemath.CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.Stamina, gameConfig.ThresholdPercentage, gameConfig.NumLevels, 1),
		AgentState:    agentStateMap,
	}

	return agentMap, globalState, gameConfig
}

func addCommsChannels(agentMap map[commons.ID]agent.Agent) (res map[commons.ID]chan message.TaggedMessage) {
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
	immutableMap := createImmutableMap(res)
	for id, a := range agentMap {
		a.BaseAgent = agent.NewBaseAgent(agent.NewCommunication(res[id], *immutableMap.Delete(id)), id, a.BaseAgent.AgentName)
		agentMap[id] = a
	}
	return
}

func createImmutableMap(peerChannels map[commons.ID]chan message.TaggedMessage) immutable.Map[commons.ID, chan<- message.TaggedMessage] {
	builder := immutable.NewMapBuilder[commons.ID, chan<- message.TaggedMessage](nil)
	for pId, channel := range peerChannels {
		builder.Set(pId, channel)
	}
	return *builder.Map()
}
