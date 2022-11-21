package game_math

import (
	"math"
	"math/rand"
)

func CalculateMonsterHealth(N uint, ST uint, L uint, CL uint) uint {
	//todo : fix this to be something correct
	delta := ((float64(rand.Intn(40))) + float64(80)) / float64(100)
	return uint(math.Ceil((float64(N) * float64(ST) / float64(L)) * delta * (2.0*float64(CL)/float64(L) + float64(0.5))))
}

func CalculateMonsterDamage(N uint, HP uint, ST uint, TH float32, L uint, CL uint) uint {
	//todo : fix this to be something correct
	delta := ((float64(rand.Intn(40))) + float64(80)) / float64(100)
	M := math.Ceil(float64(N) * float64(TH))
	return uint(math.Ceil(delta * (((float64(N) * float64(HP)) + (float64(N) * float64(ST))) / float64(L)) * (2.0*float64(CL)/float64(L) + float64(0.5)) * (1 - (M / float64(N)))))
}
