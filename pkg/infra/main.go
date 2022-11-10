package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/stage/fight"
	"infra/game/state"
	"infra/logging"
	"log"
	"math"
	"math/rand"
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

	agentMap, stateChannels, decisionChannels, globalState, gameConfig := initialise()
	gameLoop(globalState, agentMap, decisionChannels, stateChannels, gameConfig)
}

func gameLoop(globalState state.State, agentMap map[uint]agent.Agent, decisionChannels map[uint]<-chan decision.Decision, stateChannels map[uint]chan<- state.State, gameConfig config.GameConfig) {
	for globalState.CurrentLevel = 0; globalState.CurrentLevel < gameConfig.NumLevels; globalState.CurrentLevel++ {
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		for globalState.MonsterHealth != 0 {
			coweringAgents, attackSum, shieldSum := fight.HandleFightRound(&globalState, agentMap, decisionChannels, gameConfig.StartingHealthPoints)
			logging.Log.WithFields(logging.LogField{
				"currLevel": globalState.CurrentLevel,
				"numCoward": coweringAgents,
				"attackSum": attackSum,
				"shieldSum": shieldSum,
				"numAgents": len(agentMap),
			}).Info(fmt.Sprintf("Battle summary"))
			if coweringAgents == uint(len(agentMap)) {
				attack := globalState.MonsterAttack
				fight.DealDamage(attack, agentMap, globalState, stateChannels, decisionChannels)
			} else {
				globalState.MonsterHealth = commons.SaturatingSub(globalState.MonsterHealth, attackSum)
				if globalState.MonsterHealth > 0 {
					damageTaken := globalState.MonsterAttack - shieldSum
					fight.DealDamage(damageTaken, agentMap, globalState, stateChannels, decisionChannels)
					// TODO: Monster disruptive ability
				}
			}
			if float64(len(agentMap)) < math.Ceil(float64(gameConfig.ThresholdPercentage)*float64(gameConfig.InitialNumAgents)) {
				// FIXME: This comparison is wrong. Presumably, we need to compare to number of agents that have started this
				// particular round as opposed to the game, but we don't have a way to store this afaik
				log.Printf("Lost on level %d  with %d remaining", globalState.CurrentLevel, len(agentMap))
				return
			}
		}
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

func initialise() (map[uint]agent.Agent, map[uint]chan<- state.State, map[uint]<-chan decision.Decision, state.State, config.GameConfig) {
	agentMap := make(map[uint]agent.Agent)

	agentStateMap := make(map[uint]state.AgentState)

	stateChannels := make(map[uint]chan<- state.State)
	decisionChannels := make(map[uint]<-chan decision.Decision)

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("No .env file located, using defaults")
	}

	gameConfig := config.GameConfig{
		NumLevels:              config.EnvToUint("LEVELS", 60),
		StartingHealthPoints:   config.EnvToUint("STARTING_HP", 1000),
		StartingAttackStrength: config.EnvToUint("STARTING_ATTACK", 1000),
		StartingShieldStrength: config.EnvToUint("STARTING_SHIELD", 1000),
		ThresholdPercentage:    config.EnvToFloat("THRESHOLD_PCT", 0.6),
		AgentRandomQty:         config.EnvToUint("AGENT_RANDOM_QUANTITY", 100),
		InitialNumAgents:       uint(0),
	}

	instantiateAgent(gameConfig, stateChannels, decisionChannels, agentMap, agentStateMap, agent.RandomAgent{})

	gameConfig.InitialNumAgents = gameConfig.AgentRandomQty
	delta := 0.8
	N := float64(gameConfig.InitialNumAgents)
	M := math.Ceil(N * float64(gameConfig.ThresholdPercentage))
	AT := float64(gameConfig.StartingAttackStrength)
	SH := float64(gameConfig.StartingShieldStrength)
	HP := float64(gameConfig.StartingHealthPoints)
	L := float64(gameConfig.NumLevels)

	monsterResilienceStart := math.Ceil(N * AT * delta / L)
	monsterDamageStart := math.Ceil(((N*HP + N*SH) / (5 * L)) * delta * (1 - (M / N)))

	globalState := state.State{
		MonsterHealth: uint(monsterResilienceStart),
		MonsterAttack: uint(monsterDamageStart),
		AgentState:    agentStateMap,
	}
	return agentMap, stateChannels, decisionChannels, globalState, gameConfig
}

func instantiateAgent[S agent.Strategy](gameConfig config.GameConfig,
	stateChannels map[uint]chan<- state.State,
	decisionChannels map[uint]<-chan decision.Decision,
	agentMap map[uint]agent.Agent,
	agentStateMap map[uint]state.AgentState,
	strategy S) {
	for i := uint(0); i < gameConfig.AgentRandomQty; i++ {
		// TODO: add peer channels

		stateChan := make(chan state.State)
		decisionChan := make(chan decision.Decision)

		stateChannels[i] = stateChan
		decisionChannels[i] = decisionChan

		agentMap[i] = agent.Agent{
			BaseAgent: agent.BaseAgent{
				Communication: commons.Communication{
					Peer:     nil,
					Receiver: stateChan,
					Sender:   decisionChan,
				},
				Id: i,
			},
			Strategy: strategy,
		}
		agentStateMap[i] = state.AgentState{
			Hp:           gameConfig.StartingHealthPoints,
			Attack:       gameConfig.StartingAttackStrength,
			Defense:      gameConfig.StartingShieldStrength,
			BonusAttack:  0,
			BonusDefense: 0,
		}
	}
}
