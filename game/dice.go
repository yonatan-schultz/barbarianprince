package game

import "math/rand"

// Roll1d6 returns a random number 1-6
func Roll1d6() int {
	return rand.Intn(6) + 1
}

// Roll2d6 returns the sum of two 1d6 rolls (2-12)
func Roll2d6() int {
	return Roll1d6() + Roll1d6()
}

// Roll1d3 returns a random number 1-3
func Roll1d3() int {
	return rand.Intn(3) + 1
}