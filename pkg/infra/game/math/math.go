package math

import (
	"math"
	"math/rand"

	"infra/config"
)

func CalculateMonsterHealth(n uint, st uint, l uint, cl uint) uint {
	delta := ((float64(rand.Intn(40))) + float64(80)) / float64(100)
	return uint(math.Ceil((float64(n) * float64(st) / float64(l)) * delta * (2.0*float64(cl)/float64(l) + 0.5)))
}

func CalculateMonsterDamage(n uint, hp uint, st uint, th float32, l uint, cl uint) uint {
	delta := ((float64(rand.Intn(40))) + float64(80)) / float64(100)
	M := math.Ceil(float64(n) * float64(th))
	NFp := float64(n)
	LFp := float64(l)
	return uint(delta * (NFp / LFp) * (float64(hp) + float64(st)) * (2.0*float64(cl)/LFp + 0.5) * (1.0 - M/NFp))
}

func GetNextLevelMonsterValues(gameConfig config.GameConfig, currentLevel uint) (uint, uint) {
	return CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.Stamina, gameConfig.NumLevels, currentLevel+1), CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.Stamina, gameConfig.ThresholdPercentage, gameConfig.NumLevels, currentLevel+1)
}
