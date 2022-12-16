package team0

import (
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/example"
	"infra/game/stage/initialise"
	"infra/game/state"
)

/**
 * This is an example of a private experiment:
 *
 * Try running this several times (set `MODE=0` in .env) and observe the final levels reached.
 * Now uncomment the "DefensiveAgent" amd comment out the "AggressiveAgent", what differences
 * do you observe?
 */
var InitAgentMap = map[commons.ID]func() agent.Strategy{
	"RANDOM": example.NewRandomAgent,
	// "AggressiveAgent": NewProbabilisticAgent(0.1, 0.8, 0.1),
	//	"DefensiveAgent": NewProbabilisticAgent(0.1, 0.8, 0.1),
	// "CowardlyAgent": NewProbabilisticAgent(0.9, 0.05, 0.05),
}

func InitAgents(defaultStrategyMap map[commons.ID]func() agent.Strategy, gameConfig config.GameConfig, ptr *state.View) (numAgents uint, agentMap map[commons.ID]agent.Agent, agentStateMap map[commons.ID]state.AgentState, inventoryMap state.InventoryMap) {
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

	return
}
