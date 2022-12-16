package config

type GameConfig struct {
	NumLevels              uint
	StartingHealthPoints   uint
	StartingAttackStrength uint
	StartingShieldStrength uint
	ThresholdPercentage    float32
	InitialNumAgents       uint
	Stamina                uint
	VotingStrategy         uint
	VotingPreferences      uint
	Defection              bool
}
