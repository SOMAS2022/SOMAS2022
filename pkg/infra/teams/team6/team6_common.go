package team6

import (
	"infra/game/commons"

	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](a T, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min[T constraints.Ordered](a T, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

func SCSaturatingAdd(a uint, b uint, max uint) uint {
	if a+b > max {
		return max
	} else {
		return a + b
	}
}

func SafeMapReadOrDefault(in map[commons.ID]uint, idx commons.ID, def uint) uint {
	val, ok := in[idx]
	if ok {
		return val
	} else {
		return def
	}
}

func FindMaxAgentInMap[K comparable, V constraints.Ordered](in map[K]V) (K, V) {
	var smallest K
	for k, v := range in {
		if v < in[smallest] {
			smallest = k
		}
	}

	return smallest, in[smallest]
}
