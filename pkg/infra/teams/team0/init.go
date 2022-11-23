package team0

import (
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/stage/initialise"
	"infra/game/state"
)

var InitAgentMap = map[commons.ID]agent.Strategy{
	"RANDOM":          agent.NewRandomAgent(),
	"AggressiveAgent": NewProbabilisticAgent(0.1, 0.8, 0.1),
	"CowardlyAgent":   NewProbabilisticAgent(0.9, 0.05, 0.05),
}

func InitAgents(defaultStratergyMap map[commons.ID]agent.Strategy, gameConfig config.GameConfig) (numAgents uint, agentMap map[commons.ID]agent.Agent, agentStateMap map[commons.ID]state.AgentState) {
	agentMap = make(map[commons.ID]agent.Agent)
	agentStateMap = make(map[commons.ID]state.AgentState)

	numAgents = 0

	for agentName, strategy := range defaultStratergyMap {
		expectedEnvName := "AGENT_" + agentName + "_QUANTITY"
		quantity := config.EnvToUint(expectedEnvName, 100)

		numAgents += quantity
		initialise.InstantiateAgent(gameConfig, agentMap, agentStateMap, quantity, strategy, agentName)
	}

	return
}
