package commons

type ImmutableList[A any] struct {
	internalList []A
}

func (l ImmutableList[A]) Len() int {
	return len(l.internalList)
}

func NewImmutableList[A any](internalList []A) *ImmutableList[A] {
	return &ImmutableList[A]{internalList: internalList}
}

type Iterator[A any] struct {
	internal ImmutableList[A]
	index    int
}

func NewIterator[A any](internal ImmutableList[A]) *Iterator[A] {
	return &Iterator[A]{internal: internal}
}

func (l ImmutableList[A]) Iterator() Iterator[A] {
	return *NewIterator(l)
}

func (i *Iterator[A]) Done() bool {
	return i.index == len(i.internal.internalList)
}

func (i *Iterator[A]) Next() (A, bool) {
	defer func() {
		i.index++
	}()

	if i.index >= len(i.internal.internalList) {
		return i.internal.internalList[0], false
	}
	return i.internal.internalList[i.index], true
}
