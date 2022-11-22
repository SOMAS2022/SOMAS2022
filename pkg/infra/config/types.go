package config

import "infra/logging"

type GameConfig struct {
	Initialised            bool
	NumLevels              uint
	StartingHealthPoints   uint
	StartingAttackStrength uint
	StartingShieldStrength uint
	ThresholdPercentage    float32
	InitialNumAgents       uint
	Stamina                uint
}

var config = GameConfig{}

func InitConfig(values GameConfig) {
	if config.Initialised {
		logging.Log(logging.Error, logging.LogField{
			"levels": values.NumLevels,
		}, "Trying to reinitialise GameConfig")
		return
	}

	config = GameConfig{
		Initialised:            true,
		NumLevels:              values.NumLevels,
		StartingHealthPoints:   values.StartingHealthPoints,
		StartingAttackStrength: values.StartingAttackStrength,
		StartingShieldStrength: values.StartingShieldStrength,
		ThresholdPercentage:    values.ThresholdPercentage,
		InitialNumAgents:       values.InitialNumAgents,
		Stamina:                values.Stamina,
	}
}

func ViewConfig() GameConfig {
	return config
}
