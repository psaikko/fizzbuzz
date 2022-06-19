package fixedwidth

import (
	"bufio"
	"fizzbuzz/baseline"
	"os"
	"testing"
)

func BenchmarkFizzBuzz1M(b *testing.B) {
	devnull, _ := os.Open("/dev/null")
	os.Stdout = devnull
	for i := 0; i < b.N; i++ {
		FizzBuzz(1, 1000000)
	}
}

func BenchmarkFizzBuzz10M(b *testing.B) {
	devnull, _ := os.Open("/dev/null")
	os.Stdout = devnull
	for i := 0; i < b.N; i++ {
		FizzBuzz(1000000, 10000000)
	}
}

func BenchmarkFizzBuzz1G(b *testing.B) {
	devnull, _ := os.Open("/dev/null")
	os.Stdout = devnull
	for i := 0; i < b.N; i++ {
		FizzBuzz(100000000, 1000000000)
	}
}

func BenchmarkFizzBuzz10G(b *testing.B) {
	devnull, _ := os.Open("/dev/null")
	os.Stdout = devnull
	for i := 0; i < b.N; i++ {
		FizzBuzz(1000000000, 10000000000)
	}
}

func TestFixedWidthSmall(t *testing.T) {
	const testLines = 10101

	r1, w1, _ := os.Pipe()
	os.Stdout = w1
	FizzBuzz(1, testLines)

	r2, w2, _ := os.Pipe()
	os.Stdout = w2
	baseline.FizzBuzz(1, testLines)

	br1 := bufio.NewReader(r1)
	br2 := bufio.NewReader(r2)

	for i := 1; i < testLines; i++ {
		s1, _ := br1.ReadString('\n')
		s2, _ := br2.ReadString('\n')
		if s1 != s2 {
			t.Errorf("Output '%s' did not match reference implementation '%s'", s1, s2)
		}
	}
}

func TestParallelSmall(t *testing.T) {
	const testLines = 10101

	r1, w1, _ := os.Pipe()
	os.Stdout = w1
	ParallelFizzBuzz(1, testLines)

	r2, w2, _ := os.Pipe()
	os.Stdout = w2
	baseline.FizzBuzz(1, testLines)

	br1 := bufio.NewReader(r1)
	br2 := bufio.NewReader(r2)

	for i := 1; i < testLines; i++ {
		s1, _ := br1.ReadString('\n')
		s2, _ := br2.ReadString('\n')
		if s1 != s2 {
			t.Errorf("Output '%s' did not match reference implementation '%s'", s1, s2)
		}
	}
}

func TestWidthRanges(t *testing.T) {
	cases := []struct {
		from, to int
		expected []widthRange
	}{
		{1, 5, []widthRange{{1, 5, 1}}},
		{1, 10, []widthRange{{1, 9, 1}, {10, 10, 2}}},
		{57, 1010, []widthRange{{57, 99, 2}, {100, 999, 3}, {1000, 1010, 4}}},
		{9, 10, []widthRange{{9, 9, 1}, {10, 10, 2}}},
	}

	for _, tc := range cases {
		res := getWidthRanges(tc.from, tc.to)

		if len(res) != len(tc.expected) {
			t.Fatalf("Incorrect number of ranges %d != expected %d", len(res), len(tc.expected))
		}

		for i := range res {
			if res[i].from != tc.expected[i].from || res[i].to != tc.expected[i].to {
				t.Fatalf("Incorrect range %v != expected %v", res[i], tc.expected[i])
			}
		}
	}
}
