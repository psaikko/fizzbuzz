package ints

import "testing"

func TestLog10(t *testing.T) {
	for _, tc := range []struct{ input, output int }{
		{1, 0},
		{9, 1},
		{999, 3},
		{1001, 4},
	} {
		if Log10(tc.input) != tc.output {
			t.Errorf("Expected Log10(%d) == %d, was %d", tc.input, tc.output, Log10(tc.input))
		}
	}
}

func TestPow(t *testing.T) {
	for _, tc := range []struct{ n, k, output int }{
		{0, 1, 0},
		{1, 0, 1},
		{10, 4, 10000},
		{2, 2, 4},
		{3, 3, 27},
	} {
		if Pow(tc.n, tc.k) != tc.output {
			t.Errorf("Expected Pow(%d,%d) == %d, was %d", tc.n, tc.k, tc.output, Pow(tc.n, tc.k))
		}
	}
}
