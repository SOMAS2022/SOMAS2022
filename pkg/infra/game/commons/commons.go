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
		return s, fmt.Errorf("out of bounds error, attempted to access index %d in slice %v", i, s)
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

func ListToImmutableList[I comparable](l []I) immutable.List[I] {
	v := immutable.NewListBuilder[I]()

	for _, x := range l {
		v.Append(x)
	}
	return *v.List()
}

func ListToImmutableSortedSet[I constraints.Ordered](l []I) immutable.SortedMap[I, struct{}] {
	builder := immutable.NewSortedMapBuilder[I, struct{}](nil)
	for _, elem := range l {
		builder.Set(elem, struct{}{})
	}
	return *builder.Map()
}

func MapToSortedImmutable[K constraints.Ordered, V any](m map[K]V) immutable.SortedMap[K, V] {
	builder := immutable.NewSortedMapBuilder[K, V](nil)
	for k, v := range m {
		builder.Set(k, v)
	}
	return *builder.Map()
}

func ImmutableListEquality[I comparable](a immutable.List[I], b immutable.List[I]) bool {
	if a.Len() != b.Len() {
		return false
	}
	iteratorA := a.Iterator()
	iteratorB := b.Iterator()

	for !iteratorA.Done() && !iteratorB.Done() {
		_, a := iteratorA.Next()
		_, b := iteratorB.Next()
		if a != b {
			return false
		}
	}
	return true
}

func ImmutableSetEquality[I constraints.Ordered](a immutable.SortedMap[I, struct{}], b immutable.SortedMap[I, struct{}]) bool {
	if a.Len() != b.Len() {
		return false
	}
	iteratorA := a.Iterator()
	iteratorB := b.Iterator()
	for !iteratorA.Done() && !iteratorB.Done() {
		vA, _, _ := iteratorA.Next()
		vB, _, _ := iteratorB.Next()
		if vA != vB {
			return false
		}
	}
	return true
}

type ItemType bool

const (
	Shield ItemType = true
	Weapon ItemType = false
)

type ID = string

type ProposalID = string

type ItemID = string

type TradeID = string

func ImmutableListToSlice[V comparable](list immutable.List[V]) []V {
	slice := make([]V, list.Len())
	for i := 0; i < list.Len(); i++ {
		slice[i] = list.Get(i)
	}
	return slice
}

func SliceToImmutableList[V comparable](slice []V) *immutable.List[V] {
	list := immutable.NewListBuilder[V]()
	for _, item := range slice {
		list.Append(item)
	}
	return list.List()
}
