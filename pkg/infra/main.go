package main

import (
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
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

	for globalState.CurrentLevel = 0; globalState.CurrentLevel < numLevels; globalState.CurrentLevel++ {
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		for globalState.MonsterHealth = monsterHealth; globalState.MonsterHealth != 0; {
			coweringAgents, attackSum, shieldSum := handleFightRound(&globalState, agentMap, decisionChannels)
			fmt.Printf("%d cowards, %d att, %d def, %d agents\n", coweringAgents, attackSum, shieldSum, len(agentMap))
			if coweringAgents == uint(len(agentMap)) {
				attack := globalState.MonsterAttack
				dealDamage(attack, agentMap, globalState, stateChannels, decisionChannels)
			} else {
				globalState.MonsterHealth = commons.SaturatingSub(globalState.MonsterHealth, attackSum)
				if globalState.MonsterHealth > 0 {
					damageTaken := globalState.MonsterAttack - shieldSum
					dealDamage(damageTaken, agentMap, globalState, stateChannels, decisionChannels)
					// TODO: Monster disruptive ability
				}
			}
			if len(agentMap) < numAgentsRequired {
				fmt.Printf("Lost on level %d  with %d remaining", globalState.CurrentLevel, len(agentMap))
				return
			}
		}
		// TODO: End of Level looting and trading
	}
}

func dealDamage(attack uint, agentMap map[uint]agent.Agent, globalState state.State, stateChannels map[uint]chan<- state.State, decisionChannels map[uint]<-chan decision.Decision) {
	splitDamage := attack / uint(len(agentMap))
	for id, agentState := range globalState.AgentState {
		newHp := commons.SaturatingSub(agentState.Hp, splitDamage)
		if newHp == 0 {
			// kill agent
			// todo: prune peer channels somehow...
			delete(globalState.AgentState, id)
			delete(agentMap, id)
			delete(stateChannels, id)
			delete(decisionChannels, id)
		} else {
			globalState.AgentState[id] = state.AgentState{
				Hp:            newHp,
				Attack:        agentState.Attack,
				Defense:       agentState.Defense,
				AbilityPoints: agentState.AbilityPoints,
				BonusAttack:   agentState.BonusAttack,
				BonusDefense:  agentState.BonusDefense,
			}
		}
	}
}

func handleFightRound(state *state.State, agents map[uint]agent.Agent, decisionChannels map[uint]<-chan decision.Decision) (uint, uint, uint) {
	for _, a := range agents {
		go a.Strategy.HandleFight(*state, a.BaseAgent)
	}
	decisions := make(map[uint]decision.FightAction)

	for agentID, decisionC := range decisionChannels {
		handleFightDecision(decisionC, decisions, agentID, state.AgentState[agentID])
	}

	var coweringAgents uint
	var attackSum uint
	var shieldSum uint

	for agentID, d := range decisions {
		switch v := d.(type) {
		case decision.Fight:
			attackSum += v.Attack
			shieldSum += v.Defend
		case decision.Cower:
			coweringAgents++
			if entry, ok := state.AgentState[agentID]; ok {
				entry.Hp += 1
				state.AgentState[agentID] = entry
			}
		}
	}

	return coweringAgents, attackSum, shieldSum
}

func handleFightDecision(decisionC <-chan decision.Decision, decisions map[uint]decision.FightAction, agentID uint, s state.AgentState) {
	for {
		received := <-decisionC
		switch d := received.(type) {
		case decision.FightDecision:
			switch f := d.Action.(type) {
			case decision.Fight:
				if f.Attack <= s.TotalAttack() && f.Defend <= s.TotalDefense() && f.Attack+f.Defend <= s.AbilityPoints {
					decisions[agentID] = f
				} else {
					decisions[agentID] = decision.Cower{}
				}
			case decision.Cower:
				decisions[agentID] = f
			}
			return
		default:
			continue
		}
	}
}
