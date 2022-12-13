package team6

import "golang.org/x/exp/constraints"

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
