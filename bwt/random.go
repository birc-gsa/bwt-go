package bwt

import (
	"math/rand"
	"testing"
	"time"
)

// newRandomSeed creates a new random number generator
func newRandomSeed(tb testing.TB) *rand.Rand {
	tb.Helper()

	seed := time.Now().UTC().UnixNano()
	return rand.New(rand.NewSource(seed))
}

// randomStringN constructs a random string of length in n, over the alphabet alpha.
func randomStringN(n int, alpha string, rng *rand.Rand) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = alpha[rng.Intn(len(alpha))]
	}

	return string(bytes)
}
