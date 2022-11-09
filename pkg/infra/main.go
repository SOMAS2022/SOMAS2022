package main

import (
	"flag"
	"fmt"
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/stage/fight"
	"infra/game/state"
	"infra/logging"
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

	agentMap, stateChannels, decisionChannels, globalState := initialise()
	gameLoop(globalState, agentMap, decisionChannels, stateChannels)
}

func gameLoop(globalState state.State, agentMap map[uint]agent.Agent, decisionChannels map[uint]<-chan decision.Decision, stateChannels map[uint]chan<- state.State) {
	for globalState.CurrentLevel = 0; globalState.CurrentLevel < config.Config.NumLevels; globalState.CurrentLevel++ {
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		for globalState.MonsterHealth != 0 {
			coweringAgents, attackSum, shieldSum := fight.HandleFightRound(&globalState, agentMap, decisionChannels)
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
			if float64(len(agentMap)) < math.Ceil(float64(config.Config.ThresholdPercentage)*float64(config.Config.StartingNumAgents)) {
				// FIXME: This comparison is wrong. Presumably, we need to compare to number of agents that have started this
				// particular round as opposed to the game, but we don't have a way to store this afaik
				fmt.Printf("Lost on level %d  with %d remaining", globalState.CurrentLevel, len(agentMap))
				return
			}
		}
		// TODO: End of Level looting and trading
		// FIXME: This loot allocation should not stay for long!
		var weaponLoot = make([]uint, len(agentMap))
		var shieldLoot = make([]uint, len(agentMap))

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

func initialise() (map[uint]agent.Agent, map[uint]chan<- state.State, map[uint]<-chan decision.Decision, state.State) {
	agentMap := make(map[uint]agent.Agent)

	agentStateMap := make(map[uint]state.AgentState)

	stateChannels := make(map[uint]chan<- state.State)
	decisionChannels := make(map[uint]<-chan decision.Decision)

	for a := range config.Config.AgentConfig {
		for i := uint(0); i < config.Config.AgentConfig[a].Quantity; i++ {
			// TODO: add peer channels

			stateChan := make(chan state.State)
			decisionChan := make(chan decision.Decision)

			stateChannels[i] = stateChan
			decisionChannels[i] = decisionChan

			// TODO: dynamically assign agent strategies when more are implemented
			var strategy agent.Strategy
			if config.Config.AgentConfig[a].Strategy == "random" {
				strategy = agent.RandomAgent{}
			}

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
				Hp:           config.Config.StartingHealthPoints,
				Attack:       config.Config.StartingAttackStrength,
				Defense:      config.Config.StartingShieldStrength,
				BonusAttack:  0,
				BonusDefense: 0,
			}
		}
	}

	var delta = float64(0.8)
	var M = math.Ceil(float64(config.Config.StartingNumAgents) * float64(config.Config.ThresholdPercentage))
	var N = float64(config.Config.StartingNumAgents)
	var AT = float64(config.Config.StartingAttackStrength)
	var SH = float64(config.Config.StartingShieldStrength)
	var HP = float64(config.Config.StartingHealthPoints)
	var L = float64(config.Config.NumLevels)

	var monsterResilienceStart = math.Ceil(N * AT * delta / L)
	var monsterDamageStart = math.Ceil(((N*HP + N*SH) / (5 * L)) * delta * (1 - (M / N)))

	globalState := state.State{
		MonsterHealth: uint(monsterResilienceStart),
		MonsterAttack: uint(monsterDamageStart),
		AgentState:    agentStateMap,
	}
	return agentMap, stateChannels, decisionChannels, globalState
}
