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

	agentMap, globalState := initialise()
	gameLoop(globalState, agentMap)
}

func gameLoop(globalState state.State, agentMap map[uint]agent.Agent) {
	var prevDecisions map[uint]decision.FightDecision
	for globalState.CurrentLevel = 0; globalState.CurrentLevel < numLevels; globalState.CurrentLevel++ {
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		for globalState.MonsterHealth = monsterHealth; globalState.MonsterHealth != 0; {
			coweringAgents, attackSum, shieldSum, decisions := fight.HandleFightRound(&globalState, agentMap, prevDecisions)
			prevDecisions = decisions
			logging.Log.WithFields(logging.LogField{
				"currLevel": globalState.CurrentLevel,
				"numCoward": coweringAgents,
				"attackSum": attackSum,
				"shieldSum": shieldSum,
				"numAgents": len(agentMap),
			}).Info(fmt.Sprintf("Battle summary"))
			if coweringAgents == uint(len(agentMap)) {
				attack := globalState.MonsterAttack
				fight.DealDamage(attack, agentMap, globalState)
			} else {
				globalState.MonsterHealth = commons.SaturatingSub(globalState.MonsterHealth, attackSum)
				if globalState.MonsterHealth > 0 {
					damageTaken := globalState.MonsterAttack - shieldSum
					fight.DealDamage(damageTaken, agentMap, globalState)
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

		for _, agent := range agentMap {
			allocatedWeapon := rand.Intn(len(weaponLoot))
			allocatedShield := rand.Intn(len(shieldLoot))

			agent.State.BonusAttack = weaponLoot[allocatedWeapon]
			agent.State.BonusDefense = shieldLoot[allocatedShield]
			weaponLoot, _ = commons.DeleteElFromSlice(weaponLoot, allocatedWeapon)
			shieldLoot, _ = commons.DeleteElFromSlice(shieldLoot, allocatedShield)
		}
	}
}

func initialise() (map[uint]agent.Agent, state.State) {
	agentMap := make(map[uint]agent.Agent)

	for i := uint(0); i < numAgents; i++ {
		//todo: add peer channels

		agentMap[i] = agent.Agent{
			Strategy:             agent.RandomAgent{},
			StateChannel:         make(chan state.State),
			FightDecisionChannel: make(chan decision.FightDecision),
			State: state.AgentState{
				Hp:            5,
				Attack:        5,
				Defense:       3,
				AbilityPoints: 10,
				BonusAttack:   0,
				BonusDefense:  0,
			},
		}

	}

	globalState := state.State{
		MonsterHealth: monsterHealth,
		MonsterAttack: monsterAttack,
		// AgentState:    agentStateMap,
	}
	return agentMap, globalState
}
