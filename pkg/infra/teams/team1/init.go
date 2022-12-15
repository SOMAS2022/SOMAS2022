/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/

package team1

import (
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/stage/initialise"
	"infra/game/state"
	"infra/logging"
	"math/rand"
	"sync"
	"time"

	"github.com/benbjohnson/immutable"
)

var InitAgentMap = map[commons.ID]func() agent.Strategy{
	"SocialAgent": NewSocialAgent,
}

func InitAgents(defaultStrategyMap map[commons.ID]func() agent.Strategy, gameConfig config.GameConfig, ptr *state.View) (numAgents uint, agentMap map[commons.ID]agent.Agent, agentStateMap map[commons.ID]state.AgentState, inventoryMap state.InventoryMap) {
	// Initialise a random seed
	rand.Seed(time.Now().UnixNano())
	agentMap = make(map[commons.ID]agent.Agent)
	agentStateMap = make(map[commons.ID]state.AgentState)
	inventoryMap = state.InventoryMap{
		Weapons: make(map[commons.ItemID]uint),
		Shields: make(map[commons.ItemID]uint),
	}

	numAgents = 0

	for agentName, strategy := range defaultStrategyMap {
		expectedEnvName := "AGENT_" + agentName + "_QUANTITY"
		quantity := config.EnvToUint(expectedEnvName, 100)

		numAgents += quantity
		initialise.InstantiateAgent(gameConfig, agentMap, agentStateMap, quantity, strategy, agentName, ptr)
	}

	// intit agent's social cap info
	allAgents := make([]string, 0, len(agentMap))
	for k := range agentMap {
		allAgents = append(allAgents, k)
	}
	for _, a := range agentMap {
		socialStrategy := a.Strategy.(*SocialAgent)
		socialStrategy.initSocialCapital(allAgents)
	}
	connectAgents(agentMap)

	return
}

func UpdateInternalStates(agentMap map[commons.ID]agent.Agent, globalState *state.State, immutableFightRounds *commons.ImmutableList[decision.ImmutableFightResult], votesResult *immutable.Map[decision.Intent, uint], logChan chan<- logging.AgentLog) {
	var wg sync.WaitGroup
	for id, a := range agentMap {
		id := id
		a := a
		wg.Add(1)
		go func(wait *sync.WaitGroup) {
			a.HandleUpdateInternalState(globalState.AgentState[id], immutableFightRounds, votesResult, logChan)
			wait.Done()
		}(&wg)
	}
	wg.Wait()

	printGraph(agentMap, globalState)
}
