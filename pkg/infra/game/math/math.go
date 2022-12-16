package math

import (
	"math"
	"math/rand"
	"time"

	"infra/config"
)

// Enemy Resilience Modifier
func CalculateDelta() float64 {
	min := 0.8
	max := 1.2
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

// X, the monster’s resilience
func CalculateMonsterHealth(nAgent uint, stamina uint, nLevel uint, currentLevel uint) uint {
	delta := CalculateDelta()
	return uint(math.Ceil((float64(nAgent) * float64(stamina) / float64(nLevel)) * delta * 10 * (float64(currentLevel)/float64(nLevel) + (1 / 10))))
}

// Y, monster’s damage rating
func CalculateMonsterDamage(nAgent uint, HP uint, stamina uint, thresholdPercentage float32, nLevel uint, currentLevel uint) uint {
	delta := CalculateDelta()
	// Agent Survival Threshold
	M := math.Ceil(float64(nAgent) * float64(thresholdPercentage))
	NFp := float64(nAgent)
	LFp := float64(nLevel)
	return uint(delta * (NFp / LFp) * (float64(HP) + float64(stamina)) * (2.0*float64(currentLevel)/LFp + 0.5) * (1.0 - M/NFp))
}

func GetNextLevelMonsterValues(gameConfig config.GameConfig, currentLevel uint) (uint, uint) {
	return CalculateMonsterHealth(gameConfig.InitialNumAgents, gameConfig.Stamina, gameConfig.NumLevels, currentLevel+1), CalculateMonsterDamage(gameConfig.InitialNumAgents, gameConfig.StartingHealthPoints, gameConfig.Stamina, gameConfig.ThresholdPercentage, gameConfig.NumLevels, currentLevel+1)
}

func NumberPotionDropped(P float64, nAgent uint) uint {
	delta := CalculateDelta()
	return uint(delta * P * float64(nAgent))
}

func NumberEquipmentDropped(E float64, nAgent uint) uint {
	delta := CalculateDelta()
	return uint(delta * E * float64(nAgent))
}

// function encapsulated to use same random val
func GetPotionDistribution(nAgent uint) (uint, uint) {
	tau := rand.Float64()
	// hardcoded by design
	P := 0.2
	NumberHealthPotionDropped := uint((tau) * float64(NumberPotionDropped(P, nAgent)))
	NumberStaminaPotionDropped := uint((1 - tau) * float64(NumberPotionDropped(P, nAgent)))
	return NumberHealthPotionDropped, NumberStaminaPotionDropped
}

// tau recalculated for equipment  and potions
func GetEquipmentDistribution(nAgent uint) (uint, uint) {
	tau := rand.Float64()
	// hardcoded by design
	E := 0.15
	NumberWeaponDropped := uint((tau) * float64(NumberEquipmentDropped(E, nAgent)))
	NumberShieldDropped := uint((1 - tau) * float64(NumberEquipmentDropped(E, nAgent)))
	return NumberWeaponDropped, NumberShieldDropped
}

func GetWeaponDamage(X uint, nAgent uint) uint {
	delta := CalculateDelta()
	return uint(math.Ceil((delta * float64(X)) / (4.0 * float64(nAgent) * 0.8)))
}

func GetShieldProtection(Y uint, nAgent uint) uint {
	delta := CalculateDelta()
	return uint(math.Ceil((delta * float64(Y) * 0.5) / (float64(nAgent) * 0.8)))
}

func GetHealthPotionValue(Y uint, nAgent uint) uint {
	delta := CalculateDelta()
	return uint(math.Ceil((delta * float64(Y) * 5.0) / (float64(nAgent) * 0.8)))
}

func GetStaminaPotionValue(X uint, nAgent uint) uint {
	delta := CalculateDelta()
	return uint(math.Ceil((delta * float64(X)) / (float64(nAgent) * 0.8)))
}
