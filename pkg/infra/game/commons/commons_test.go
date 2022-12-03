package commons_test

import (
	"infra/game/commons"
	"testing"
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

func testSlice2Map(t *testing.T) {
	s := []uint{1, 2, 3, 4, 5}
	m := commons.Slice2Map(s)
	mt := make(map[int]uint)
	mt[1] = 1
	mt[2] = 2
	mt[3] = 3
	mt[4] = 4
	mt[5] = 5

	if m != nil {
		for i := 1; i <= 5; i++ {
			if m[i] != mt[i] {
				t.Errorf("m[%d] is %d, want %d", i, m[i], i)
			}
		}
	} else {
		t.Errorf("emplty map, expect 5 mapping")
	}
}
