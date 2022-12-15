package commons_test

import (
	"testing"

	"github.com/benbjohnson/immutable"

	"infra/game/commons"
)

func TestSaturatingSub(t *testing.T) {
	t.Parallel()

	got := commons.SaturatingSub(0, 1)
	if got != 0 {
		t.Errorf("SaturatingSub(0, 1) = %d; want 0", 0)
	}
	got = commons.SaturatingSub(1, 1)
	if got != 0 {
		t.Errorf("SaturatingSub(1, 1) = %d; want 0", 0)
	}
}

func TestDeleteElFromSlice(t *testing.T) {
	t.Parallel()

	uints := []uint{1, 2, 3, 4, 5}
	slice, err := commons.DeleteElFromSlice(uints, 1)
	if err != nil {
		t.Errorf("DeleteElFromSlice({1, 2, 3, 4, 5}, 1) threw error: %v", err)
	}
	if !testEq(slice, []uint{1, 5, 3, 4}) {
		t.Errorf("DeleteElFromSlice({1, 2, 3, 4, 5}, 1) got: %v, expected %v", slice, []uint{1, 5, 3, 4})
	}

	_, err = commons.DeleteElFromSlice(uints, -1)
	if err == nil {
		t.Errorf("Called DeleteElFromSlice({1,2,3,4,5}, -1) got: nil, expected: error")
	}
}

func TestImmutableListEquality(t *testing.T) {
	t.Parallel()

	builder := immutable.NewListBuilder[int]()
	builder.Append(1)
	builder.Append(2)
	builder.Append(3)
	listA := builder.List()

	builder2 := immutable.NewListBuilder[int]()
	builder2.Append(1)
	builder2.Append(2)
	builder2.Append(3)
	listB := builder2.List()

	builder3 := immutable.NewListBuilder[int]()
	builder3.Append(3)
	builder3.Append(2)
	builder3.Append(1)
	listC := builder3.List()

	got := commons.ImmutableListEquality(*listA, *listB)

	if got != true {
		t.Errorf("Called ImmutableListEquality({1,2,3}, {1,2,3}) got: false, expected: true")
	}
	got = commons.ImmutableListEquality(*listA, *listC)
	if got != false {
		t.Errorf("Called ImmutableListEquality({1,2,3}, {3,2,1}) got: true, expected: false")
	}
}

func testEq[Type comparable](a, b []Type) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
