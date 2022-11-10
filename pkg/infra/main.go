package main

import (
	"flag"
	"fmt"
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	game_math "infra/game/math"
	"infra/game/stage/fight"
	"infra/game/state"
	"infra/logging"
	"math"
	"math/rand"

	"github.com/joho/godotenv"
)

/*
Each agent is initialised with these modes of communication
Whereby the peer channels enable bidirectional p2p communication
Enables peers to send and receive messages with broadcasting possible via non-blocking sends
*/

func main() {
	// define flags
	useJSONFormatter := flag.Bool("j", false, "whether to use JSONFormatter for logging")
	flag.Parse()

	logging.InitLogger(*useJSONFormatter)

	agentMap, decisionChannels, globalState, gameConfig := initialise()
	gameLoop(globalState, agentMap, decisionChannels, gameConfig)
}

func gameLoop(globalState state.State, agentMap map[uint]agent.Agent, decisionChannels map[uint]<-chan decision.Decision, gameConfig config.GameConfig) {
	for globalState.CurrentLevel = 0; globalState.CurrentLevel < gameConfig.NumLevels; globalState.CurrentLevel++ {
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		for globalState.MonsterHealth != 0 {
			coweringAgents, attackSum, shieldSum := fight.HandleFightRound(&globalState, agentMap, decisionChannels, gameConfig.StartingHealthPoints)
			logging.Log.WithFields(logging.LogField{
				"currLevel":     globalState.CurrentLevel,
				"monsterHealth": globalState.MonsterHealth,
				"monsterDamage": globalState.MonsterAttack,
				"numCoward":     coweringAgents,
				"attackSum":     attackSum,
				"shieldSum":     shieldSum,
				"numAgents":     len(agentMap),
			}).Info(fmt.Sprintf("Battle summary"))
			if coweringAgents == uint(len(agentMap)) {
				attack := globalState.MonsterAttack
				fight.DealDamage(attack, agentMap, globalState, decisionChannels)
			} else {
				globalState.MonsterHealth = commons.SaturatingSub(globalState.MonsterHealth, attackSum)
				if globalState.MonsterHealth > 0 {
					damageTaken := globalState.MonsterAttack - shieldSum
					fight.DealDamage(damageTaken, agentMap, globalState, decisionChannels)
					// TODO: Monster disruptive ability
				}
			}
			if float64(len(agentMap)) < math.Ceil(float64(gameConfig.ThresholdPercentage)*float64(gameConfig.InitialNumAgents)) {
				logging.Log.Infof("Lost on level %d  with %d remaining", globalState.CurrentLevel, len(agentMap))
				return
			}
		}
		//todo: fix this
		globalState.MonsterHealth = game_math.CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.StartingAttackStrength, 0.8, gameConfig.NumLevels, globalState.CurrentLevel)
		globalState.MonsterAttack = game_math.CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.StartingShieldStrength, 0.8, gameConfig.ThresholdPercentage, gameConfig.NumLevels, globalState.CurrentLevel)

		// TODO: End of Level looting and trading
		// FIXME: This loot allocation should not stay for long!
		weaponLoot := make([]uint, len(agentMap))
		shieldLoot := make([]uint, len(agentMap))

		for i := range weaponLoot {
			weaponLoot[i] = globalState.CurrentLevel * uint(rand.Intn(3))
			shieldLoot[i] = globalState.CurrentLevel * uint(rand.Intn(3))
		}

		for _, agentState := range globalState.AgentState {
			allocatedWeapon := rand.Intn(len(weaponLoot))
			allocatedShield := rand.Intn(len(shieldLoot))

			agentState.BonusAttack = weaponLoot[allocatedWeapon]
			agentState.BonusDefense = shieldLoot[allocatedShield]
			weaponLoot, _ = commons.DeleteElFromSlice(weaponLoot, allocatedWeapon)
			shieldLoot, _ = commons.DeleteElFromSlice(shieldLoot, allocatedShield)
		}
	}
}

func initialise() (map[uint]agent.Agent, map[uint]<-chan decision.Decision, state.State, config.GameConfig) {
	agentMap := make(map[uint]agent.Agent)

	agentStateMap := make(map[uint]state.AgentState)
	decisionChannels := make(map[uint]<-chan decision.Decision)

	err := godotenv.Load()
	if err != nil {
		logging.Log.Fatalln("No .env file located, using defaults")
	}

	gameConfig := config.GameConfig{
		NumLevels:              config.EnvToUint("LEVELS", 60),
		StartingHealthPoints:   config.EnvToUint("STARTING_HP", 1000),
		StartingAttackStrength: config.EnvToUint("STARTING_ATTACK", 1000),
		StartingShieldStrength: config.EnvToUint("STARTING_SHIELD", 1000),
		ThresholdPercentage:    config.EnvToFloat("THRESHOLD_PCT", 0.6),
		AgentRandomQty:         config.EnvToUint("AGENT_RANDOM_QUANTITY", 100),
		InitialNumAgents:       uint(0),
		Stamina:                config.EnvToUint("BASE_STAMINA", 2000),
	}

	instantiateAgent(gameConfig, decisionChannels, agentMap, agentStateMap, agent.RandomAgent{})

	gameConfig.InitialNumAgents = gameConfig.AgentRandomQty

	globalState := state.State{
		MonsterHealth: game_math.CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.StartingAttackStrength, 0.8, gameConfig.NumLevels, 0),
		MonsterAttack: game_math.CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.StartingShieldStrength, 0.8, gameConfig.ThresholdPercentage, gameConfig.NumLevels, 0),
		AgentState:    agentStateMap,
	}
	return agentMap, decisionChannels, globalState, gameConfig
}

func instantiateAgent[S agent.Strategy](gameConfig config.GameConfig,
	decisionChannels map[uint]<-chan decision.Decision,
	agentMap map[uint]agent.Agent,
	agentStateMap map[uint]state.AgentState,
	strategy S) {
	for i := uint(0); i < gameConfig.AgentRandomQty; i++ {
		// TODO: add peer channels
		decisionChan := make(chan decision.Decision)

		decisionChannels[i] = decisionChan

		agentMap[i] = agent.Agent{
			BaseAgent: agent.BaseAgent{
				Communication: commons.Communication{
					Peer:   nil,
					Sender: decisionChan,
				},
				Id: i,
			},
			Strategy: strategy,
		}
		agentStateMap[i] = state.AgentState{
			Hp:           gameConfig.StartingHealthPoints,
			Stamina:      gameConfig.Stamina,
			Attack:       gameConfig.StartingAttackStrength,
			Defense:      gameConfig.StartingShieldStrength,
			BonusAttack:  0,
			BonusDefense: 0,
		}
	}
}
