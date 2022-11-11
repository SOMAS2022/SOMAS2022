package main

import (
	"flag"
	"github.com/benbjohnson/immutable"
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	gamemath "infra/game/math"
	"infra/game/stage/fight"
	"infra/game/state"
	"infra/logging"
	"math"
	"math/rand"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var InitAgentMap = map[string]agent.Strategy{
	"RANDOM": agent.RandomAgent{},
}

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

	agentMap, globalState, gameConfig := initialise()
	gameLoop(globalState, agentMap, gameConfig)
}

func gameLoop(globalState state.State, agentMap map[uint]agent.Agent, gameConfig config.GameConfig) {
	var decisionMap map[uint]decision.FightAction
	for globalState.CurrentLevel = 0; globalState.CurrentLevel < gameConfig.NumLevels; globalState.CurrentLevel++ {
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		for globalState.MonsterHealth != 0 {
			decisionMapView := immutable.NewMapBuilder[uint, decision.FightAction](nil)
			for u, action := range decisionMap {
				decisionMapView.Set(u, action)
			}
			coweringAgents, attackSum, shieldSum, dMap := fight.HandleFightRound(&globalState, agentMap, gameConfig.StartingHealthPoints, decisionMapView.Map())
			decisionMap = dMap

			logging.Log.WithFields(logging.LogField{
				"currLevel":     globalState.CurrentLevel,
				"monsterHealth": globalState.MonsterHealth,
				"monsterDamage": globalState.MonsterAttack,
				"numCoward":     coweringAgents,
				"attackSum":     attackSum,
				"shieldSum":     shieldSum,
				"numAgents":     len(agentMap),
			}).Info("Battle summary")
			if coweringAgents == uint(len(agentMap)) {
				attack := globalState.MonsterAttack
				fight.DealDamage(attack, agentMap, &globalState)
			} else {
				globalState.MonsterHealth = commons.SaturatingSub(globalState.MonsterHealth, attackSum)
				if globalState.MonsterHealth > 0 {
					damageTaken := globalState.MonsterAttack - shieldSum
					fight.DealDamage(damageTaken, agentMap, &globalState)
					// TODO: Monster disruptive ability
				}
			}
			if float64(len(agentMap)) < math.Ceil(float64(gameConfig.ThresholdPercentage)*float64(gameConfig.InitialNumAgents)) {
				logging.Log.Infof("Lost on level %d  with %d remaining", globalState.CurrentLevel, len(agentMap))
				return
			}
		}
		//todo: fix this
		//todo: There is a weird bug due to the mathematics that agents gain more health from cowering than monster attack when all cower
		//todo: Results in infinite game run-through
		globalState.MonsterHealth = gamemath.CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.StartingAttackStrength, 0.8, gameConfig.NumLevels, globalState.CurrentLevel)
		globalState.MonsterAttack = gamemath.CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.StartingShieldStrength, 0.8, gameConfig.ThresholdPercentage, gameConfig.NumLevels, globalState.CurrentLevel)

		// TODO: End of Level looting and trading
		// FIXME: This loot allocation should not stay for long!
		weaponLoot := make([]uint, len(agentMap))
		shieldLoot := make([]uint, len(agentMap))

		for i := range weaponLoot {
			weaponLoot[i] = globalState.CurrentLevel * uint(rand.Intn(3))
			shieldLoot[i] = globalState.CurrentLevel * uint(rand.Intn(3))
		}

		for i, agentState := range globalState.AgentState {
			allocatedWeapon := rand.Intn(len(weaponLoot))
			allocatedShield := rand.Intn(len(shieldLoot))

			agentState.BonusAttack = weaponLoot[allocatedWeapon]
			agentState.BonusDefense = shieldLoot[allocatedShield]
			weaponLoot, _ = commons.DeleteElFromSlice(weaponLoot, allocatedWeapon)
			shieldLoot, _ = commons.DeleteElFromSlice(shieldLoot, allocatedShield)

			globalState.AgentState[i] = agentState
		}
	}
}

func initialise() (map[uint]agent.Agent, state.State, config.GameConfig) {
	agentMap := make(map[uint]agent.Agent)

	agentStateMap := make(map[uint]state.AgentState)

	err := godotenv.Load()
	if err != nil {
		logging.Log.Warnln("No .env file located, using defaults")
	}

	gameConfig := config.GameConfig{
		NumLevels:              config.EnvToUint("LEVELS", 60),
		StartingHealthPoints:   config.EnvToUint("STARTING_HP", 1000),
		StartingAttackStrength: config.EnvToUint("STARTING_ATTACK", 1000),
		StartingShieldStrength: config.EnvToUint("STARTING_SHIELD", 1000),
		ThresholdPercentage:    config.EnvToFloat("THRESHOLD_PCT", 0.6),
		InitialNumAgents:       uint(0),
		Stamina:                config.EnvToUint("BASE_STAMINA", 2000),
	}

	for agentName, strategy := range InitAgentMap {
		expected_env_name := "AGENT_" + agentName + "_QUANTITY"
		quantity := config.EnvToUint(expected_env_name, 0)

		gameConfig.InitialNumAgents += quantity
		instantiateAgent(gameConfig, agentMap, agentStateMap, quantity, strategy)
	}

	globalState := state.State{
		MonsterHealth: gamemath.CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.StartingAttackStrength, 0.8, gameConfig.NumLevels, 0),
		MonsterAttack: gamemath.CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.StartingShieldStrength, 0.8, gameConfig.ThresholdPercentage, gameConfig.NumLevels, 0),
		AgentState:    agentStateMap,
	}
	return agentMap, globalState, gameConfig
}

func instantiateAgent[S agent.Strategy](gameConfig config.GameConfig,
	agentMap map[uint]agent.Agent,
	agentStateMap map[uint]state.AgentState,
	quantity uint,
	strategy S) {
	for i := uint(0); i < quantity; i++ {
		// TODO: add peer channels
		agentMap[i] = agent.Agent{
			BaseAgent: agent.BaseAgent{
				Communication: commons.Communication{
					Peer: nil,
				},
				Id: uuid.New().String(),
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
