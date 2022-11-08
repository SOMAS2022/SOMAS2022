package commons

import (
	"testing"
)

func TestSaturatingSub(t *testing.T) {
	got := SaturatingSub(0, 1)
	if got != 0 {
		t.Errorf("SaturatingSub(0, 1) = %d; want 0", 0)
	}
	got = SaturatingSub(1, 1)
	if got != 0 {
		t.Errorf("SaturatingSub(1, 1) = %d; want 0", 0)
	}
}

func TestDeleteElFromSlice(t *testing.T) {
	uints := []uint{1, 2, 3, 4, 5}
	slice, err := DeleteElFromSlice(uints, 1)
	sliceExp := []uint{1, 5, 3, 4}
	if err != nil {
		t.Errorf("DeleteElFromSlice({1, 2, 3, 4, 5}, 1) threw error: %v", err)
	}
	if !testEq(slice, sliceExp) {
		t.Errorf("DeleteElFromSlice({1, 2, 3, 4, 5}, 1) got: %v, expected %v", slice, sliceExp)
	}

	_, err = DeleteElFromSlice(uints, -1)
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
