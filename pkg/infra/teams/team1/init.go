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
	"infra/game/stage/initialise"
	"infra/game/state"
	"infra/teams/team1/utils"
	"math/rand"
	"time"
)

var InitAgentMap = map[commons.ID]func() agent.Strategy{
	"CollaborativeAgent": CreateCollaborativeAgent,
	"SelfishAgent":       CreateSelfishAgent,
}

func InitAgents(defaultStrategyMap map[commons.ID]func() agent.Strategy, gameConfig config.GameConfig, ptr *state.View) (numAgents uint, agentMap map[commons.ID]agent.Agent, agentStateMap map[commons.ID]state.AgentState, inventoryMap state.InventoryMap) {
	// Initialise a random seed
	rand.Seed(time.Now().UnixNano())
	utils.Config = gameConfig // TODO: Not needed when confg is globally accessible
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
	sci := SocialCapInfo{}
	sci.arr = [4]float64{0.5, 0.5, 0.5, 0.5}
	sciMap := map[string]SocialCapInfo{}
	for k, _ := range agentMap {
		sci.ID = k
		sciMap[k] = sci
	}
	for k, a := range agentMap {
		socialStrategy := a.Strategy.(*SocialAgent)
		socialStrategy.initSocialCapital(k, sciMap)
	}
	connectAgents(agentMap)

	return
}
