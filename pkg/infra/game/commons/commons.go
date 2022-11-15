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

func ImmutableMapKeys[K constraints.Ordered, V any](p immutable.Map[K, V]) (keys []K) {
	keys = make([]K, p.Len())
	iterator := p.Iterator()
	for !iterator.Done() {
		key, _, _ := iterator.Next()
		keys = append(keys, key)
	}
	return
}

type ID = string
