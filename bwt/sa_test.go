package bwt

import (
	"sort"
	"testing"
)

// checkSAIndices checks that the suffix array sa has all the
// indices in x (plus one for the sentinel if len(sa) == len(x) + 1).
// Reports an error to t otherwise
func checkSAIndices(t *testing.T, x string, sa []int32) bool {
	t.Helper()

	if len(sa) != len(x) && len(sa) != len(x)+1 {
		t.Errorf("Suffix %v has an invalid length: %d. "+
			"It should be %d without sentinel or %d with.",
			sa, len(sa), len(x), len(x)+1)
	}

	indices := make([]int, len(sa))
	for i, j := range sa {
		indices[i] = int(j)
	}

	sort.Ints(indices)

	for i, j := range indices {
		if j < 0 || j > len(x) {
			t.Errorf("Index %d is not valid for a suffix array over a string of length %d.",
				j, len(x))
		}

		if i < j {
			t.Errorf("Index %d is missing from the suffix array.",
				i)
			return false
		}
	}

	return true
}

// checkSASorted checks if a suffix array sa actually
// represents the sorted suffix in the string x. Reports
// errors to t.
func checkSASorted(t *testing.T, x string, sa []int32) bool {
	t.Helper()

	result := true

	for i := 1; i < len(sa); i++ {
		if x[sa[i-1]:] >= x[sa[i]:] {
			t.Errorf("Suffix array is not sorted! %q >= %q",
				x[sa[i-1]:], x[sa[i]:])

			result = false
		}
	}

	return result
}

// checkSuffixArray runs all the consistency checks for
// suffix array sa over string x, reporting errors to t.
func checkSuffixArray(t *testing.T, x string, sa []int32) bool {
	t.Helper()

	result := true
	result = result && checkSAIndices(t, x, sa)
	result = result && checkSASorted(t, x, sa)

	return result
}

func TestSaConstruction(t *testing.T) {
	rng := newRandomSeed(t)
	for i := 0; i < 10; i++ {
		x := randomStringN(10, "acgt", rng)
		sa := PrefixDoubling(x)
		checkSuffixArray(t, x, sa)
	}
}
