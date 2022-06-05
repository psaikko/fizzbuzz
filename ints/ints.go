package ints

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func CeilDiv(val, divisor int) int {
	return (val + divisor - 1) / divisor
}

func FloorDiv(val, divisor int) int {
	return val / divisor
}

// Pow returns n^k
func Pow(n, k int) int {
	if k == 0 {
		return 1
	} else if k == 1 {
		return n
	} else {
		return Pow(n, k/2) * Pow(n, k-k/2)
	}
}

// Log10 returns base 10 logarithm for positive n
func Log10(n int) int {
	if n <= 0 {
		panic("bad input")
	}
	i := 0
	c := 1
	for c < n {
		c *= 10
		i++
	}
	return i
}
