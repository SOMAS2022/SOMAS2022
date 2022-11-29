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
	"infra/game/example"
	"infra/game/stage/initialise"
	"infra/game/state"
	"infra/teams/team1/utils"
	"math/rand"
	"time"
)

var InitAgentMap = map[commons.ID]agent.Strategy{
	"RANDOM":          example.NewRandomAgent(),
	"AggressiveAgent": CreateAggressiveAgent(),
	"DefensiveAgent":  CreateDefensiveAgent(),
}

func InitAgents(defaultStrategyMap map[commons.ID]func(), gameConfig config.GameConfig) (numAgents uint, agentMap map[commons.ID]agent.Agent, agentStateMap map[commons.ID]state.AgentState) {
	// Initialise a random seed
	rand.Seed(time.Now().UnixNano())
	utils.Config = gameConfig // TODO: Not needed when confg is globally accessible
	agentMap = make(map[commons.ID]agent.Agent)
	agentStateMap = make(map[commons.ID]state.AgentState)

	numAgents = 0

	for agentName, strategy := range defaultStrategyMap {
		expectedEnvName := "AGENT_" + agentName + "_QUANTITY"
		quantity := config.EnvToUint(expectedEnvName, 100)

		numAgents += quantity
		initialise.InstantiateAgent(gameConfig, agentMap, agentStateMap, quantity, strategy, agentName)
	}

	return
}
