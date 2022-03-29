package bwt

import (
	"testing"
)

func TestBwt(t *testing.T) {
	rng := newRandomSeed(t)
	for i := 0; i < 10; i++ {
		x := randomStringN(10, "acgt", rng)
		y := Bwt(x)
		z := Rbwt(y)
		if x != z {
			t.Errorf("%s != %s\n", x, z)
		}
	}
}
