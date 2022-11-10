package game_math

import (
	"math"
)

func CalculateMonsterHealth(N uint, AT uint, delta float64, L uint, CL uint) uint {
	//todo : fix this to be something correct
	return uint(math.Ceil(float64(N) * float64(AT) * delta * 14 * float64(1.0+CL/L) / float64(L)))
}

func CalculateMonsterDamage(N uint, HP uint, SH uint, delta float64, TH float32, L uint, CL uint) uint {
	//todo : fix this to be something correct
	M := math.Ceil(float64(N) * float64(TH))
	return uint(math.Ceil(((float64(N)*float64(HP) + float64(N)*float64(SH)) * 14 * float64(1.0+CL/L) / (5 * float64(L))) * delta * (1 - (M / float64(N)))))
}
