package initialise

import (
	"github.com/google/uuid"
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/state"
)

func InstantiateAgent[S agent.Strategy](gameConfig config.GameConfig,
	agentMap map[commons.ID]agent.Agent,
	agentStateMap map[commons.ID]state.AgentState,
	quantity uint,
	strategy S,
	agentName string,
) {
	for i := uint(0); i < quantity; i++ {
		// TODO: add peer channels
		agentId := uuid.New().String()
		agentMap[agentId] = agent.Agent{
			BaseAgent: agent.BaseAgent{AgentName: agentName},
			Strategy:  strategy,
		}

		agentStateMap[agentId] = state.AgentState{
			Hp:           gameConfig.StartingHealthPoints,
			Stamina:      gameConfig.Stamina,
			Attack:       gameConfig.StartingAttackStrength,
			Defense:      gameConfig.StartingShieldStrength,
			BonusAttack:  0,
			BonusDefense: 0,
		}
	}
}

func InitGameConfig() config.GameConfig {
	gameConfig := config.GameConfig{
		NumLevels:              config.EnvToUint("LEVELS", 60),
		StartingHealthPoints:   config.EnvToUint("STARTING_HP", 1000),
		StartingAttackStrength: config.EnvToUint("STARTING_ATTACK", 100),
		StartingShieldStrength: config.EnvToUint("STARTING_SHIELD", 100),
		ThresholdPercentage:    config.EnvToFloat("THRESHOLD_PCT", 0.6),
		InitialNumAgents:       uint(0),
		Stamina:                config.EnvToUint("BASE_STAMINA", 2000),
	}

	return gameConfig
}

func InitAgents(defaultStratergyMap map[commons.ID]agent.Strategy, gameConfig config.GameConfig) (numAgents uint, agentMap map[commons.ID]agent.Agent, agentStateMap map[commons.ID]state.AgentState) {
	agentMap = make(map[commons.ID]agent.Agent)
	agentStateMap = make(map[commons.ID]state.AgentState)

	numAgents = 0

	for agentName, strategy := range defaultStratergyMap {
		expectedEnvName := "AGENT_" + agentName + "_QUANTITY"
		quantity := config.EnvToUint(expectedEnvName, 100)

		numAgents += quantity
		InstantiateAgent(gameConfig, agentMap, agentStateMap, quantity, strategy, agentName)
	}

	return

}
