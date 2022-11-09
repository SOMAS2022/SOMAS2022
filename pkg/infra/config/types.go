package config

type GameConfig struct {
	NumLevels              uint          `json:"numLevels"`
	StartingNumAgents      uint          `json:"startingNumAgents"`
	StartingHealthPoints   uint          `json:"healthPointStart"`
	StartingAttackStrength uint          `json:"attackStrengthStart"`
	StartingShieldStrength uint          `json:"shieldStrengthStart"`
	ThresholdPercentage    float32       `json:"thresholdPercentage"`
	AgentConfig            []AgentConfig `json:"agentConfig"`
}

type AgentConfig struct {
	Strategy string `json:"strategy"`
	Quantity uint   `json:"quantity"`
}

var Config GameConfig
