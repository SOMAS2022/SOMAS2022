package commons

import (
	"fmt"

	"github.com/benbjohnson/immutable"
	"golang.org/x/exp/constraints"
)

func SaturatingSub(x uint, y uint) uint {
	res := x - y
	var val uint
	if res <= x {
		val = 1
	}
	res &= -val
	return res
}

func DeleteElFromSlice(s []uint, i int) ([]uint, error) {
	if i < cap(s) && i >= 0 {
		s[i] = s[len(s)-1]
		return s[:len(s)-1], nil
	} else {
		return s, fmt.Errorf("Out of bounds error, attempted to access index %d in slice %v\n", i, s)
	}
}

func ImmutableMapKeys[K constraints.Ordered, V any](p immutable.Map[K, V]) []K {
	keys := make([]K, p.Len())
	iterator := p.Iterator()
	idx := 0
	for !iterator.Done() {
		key, _, _ := iterator.Next()
		keys[idx] = key
		idx++
	}

	return keys
}

func MapToImmutable[K constraints.Ordered, V any](m map[K]V) immutable.Map[K, V] {
	builder := immutable.NewMapBuilder[K, V](nil)
	for k, v := range m {
		builder.Set(k, v)
	}
	return *builder.Map()
}

type ID = string

type ProposalID = string

type ItemID = string
