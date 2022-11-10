package main

import (
	"flag"
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/stage/fight"
	"infra/game/state"
	"infra/logging"
	"math/rand"
)

/*
Each agent is initialised with these modes of communication
Whereby the peer channels enable bidirectional p2p communication
Enables peers to send and receive messages with broadcasting possible via non-blocking sends
*/

const numLevels = 60

const numAgents = 100

const numAgentsRequired = numAgents * 2 / 5

const monsterHealth = 100
const monsterAttack = 240

func main() {
	// define flags
	useJSONFormatter := flag.Bool("j", false, "whether to use JSONFormatter for logging")
	flag.Parse()

	logging.InitLogger(*useJSONFormatter)

	agentMap, stateChannels, decisionChannels, globalState := initialise()
	gameLoop(globalState, agentMap, decisionChannels, stateChannels)
}

func gameLoop(globalState state.State, agentMap map[uint]agent.Agent, decisionChannels map[uint]<-chan decision.Decision, stateChannels map[uint]chan<- state.State) {
	for globalState.CurrentLevel = 0; globalState.CurrentLevel < numLevels; globalState.CurrentLevel++ {
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		for globalState.MonsterHealth = monsterHealth; globalState.MonsterHealth != 0; {
			coweringAgents, attackSum, shieldSum := fight.HandleFightRound(&globalState, agentMap, decisionChannels)
			logging.Log.Info(fmt.Sprintf("%4d cowards, %4d att, %4d def, %4d agents\n", coweringAgents, attackSum, shieldSum, len(agentMap)))
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
			if len(agentMap) < numAgentsRequired {
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

	for i := uint(0); i < numAgents; i++ {
		//todo: add peer channels

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
			Strategy: agent.RandomAgent{},
		}
		agentStateMap[i] = state.AgentState{
			Hp:            5,
			Attack:        5,
			Defense:       3,
			AbilityPoints: 10,
			BonusAttack:   0,
			BonusDefense:  0,
		}
	}

	globalState := state.State{
		MonsterHealth: monsterHealth,
		MonsterAttack: monsterAttack,
		AgentState:    agentStateMap,
	}
	return agentMap, stateChannels, decisionChannels, globalState
}
