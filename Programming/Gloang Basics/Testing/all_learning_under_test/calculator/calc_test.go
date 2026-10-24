package calculator

import (
	"testing"
)

// What is Table Driven Test ???

func TestAdd(t *testing.T) {
	t.Parallel() // Runs this code parallel with other
	var want float64 = 10
	got := Add(4, 6)
	if want != got {
		t.Errorf("got %f, want %f", got, want)
	}

	testCases := []struct {
		a, b float64
		want float64
	}{
		{4, 5, 9},
		{1.3, 4.5, 5.8},
		{-1, 5, 4},
	}

	for _, tc := range testCases {
		got := Add(tc.a, tc.b)
		if tc.want != got {
			t.Errorf("got %g, want %g", got, tc.want)
		}
	}
}
