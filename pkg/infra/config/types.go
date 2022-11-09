package config

type GameConfig struct {
	NumLevels              uint
	StartingHealthPoints   uint
	StartingAttackStrength uint
	StartingShieldStrength uint
	ThresholdPercentage    float32
	AgentRandomQty         uint
	InitialNumAgents       uint
}
