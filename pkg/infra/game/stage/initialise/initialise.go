package initialise

import (
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/state"

	"github.com/google/uuid"
)

func InstantiateAgent[S agent.Strategy](gameConfig config.GameConfig,
	agentMap map[commons.ID]agent.Agent,
	agentStateMap map[commons.ID]state.AgentState,
	quantity uint,
	strategyConstructor func() S,
	agentName string,
) {
	for i := uint(0); i < quantity; i++ {
		// TODO: add peer channels
		agentId := uuid.New().String()
		agentMap[agentId] = agent.Agent{
			BaseAgent: agent.NewBaseAgent(nil, agentId, agentName),
			Strategy:  strategyConstructor(),
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
		VotingStrategy:         config.EnvToUint("VOTING_STRATEGY", 1),
		VotingPreferences:      config.EnvToUint("VOTING_PREFERENCES", 1),
	}

	return gameConfig
}

func InitAgents(defaultStrategyMap map[commons.ID]func() agent.Strategy, gameConfig config.GameConfig) (numAgents uint, agentMap map[commons.ID]agent.Agent, agentStateMap map[commons.ID]state.AgentState) {
	agentMap = make(map[commons.ID]agent.Agent)
	agentStateMap = make(map[commons.ID]state.AgentState)

	numAgents = 0

	for agentName, strategy := range defaultStrategyMap {
		expectedEnvName := "AGENT_" + agentName + "_QUANTITY"
		quantity := config.EnvToUint(expectedEnvName, 100)

		numAgents += quantity
		InstantiateAgent(gameConfig, agentMap, agentStateMap, quantity, strategy, agentName)
	}

	return
}
