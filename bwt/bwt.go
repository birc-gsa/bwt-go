package bwt

func Bwt(x string) string {
	sa := PrefixDoubling(x)
	y := make([]byte, len(x)+1) // + 1 for the sentinel, not included in x

	for i, j := range sa {
		if j == 0 {
			y[i] = 0
		} else {
			y[i] = x[j-1]
		}
	}

	return string(y)
}

// CTab Structur holding the C-table for BWT search.
// This is a map from letters in the alphabet to the
// cumulative sum of how often we see letters in the
// BWT
type CTab struct {
	CumSum []int
}

// Rank How many times does the BWT hold a letter smaller
// than a? Undefined behaviour if a isn't in the table.
func (ctab *CTab) Rank(a byte) int {
	return ctab.CumSum[a]
}

// NewCTab builds the c-table from a string.
func NewCTab(bwt []byte, asize int) *CTab {
	// First, count how often we see each character
	counts := make([]int, asize)
	for _, b := range bwt {
		counts[b]++
	}
	// Then get the accumulative sum
	var n int
	for i, count := range counts {
		counts[i] = n
		n += count
	}

	return &CTab{counts}
}

// OTab Holds the o-table (rank table) from a BWT string
type OTab struct {
	nrow, ncol int
	table      []int
}

func (otab *OTab) offset(a byte, i int) int {
	// -1 to a because we don't store the sentinel
	// and -1 to i because we don't store the first
	// row (which is always zero)
	return otab.ncol*(int(a)-1) + (i - 1)
}

func (otab *OTab) get(a byte, i int) int {
	return otab.table[otab.offset(a, i)]
}

func (otab *OTab) set(a byte, i, val int) {
	otab.table[otab.offset(a, i)] = val
}

// Rank How many times do we see letter a before index i
// in the BWT string?
func (otab *OTab) Rank(a byte, i int) int {
	// We don't explicitly store the first column,
	// since it is always empty anyway.
	if i == 0 {
		return 0
	}

	return otab.get(a, i)
}

// NewOTab builds the o-table from a string. It uses
// the suffix array to get the BWT and a c-table
// to handle the alphabet.
func NewOTab(bwt []byte, asize int) *OTab {
	// We index for all characters except $, so
	// nrow is the alphabet size minus one.
	// We index all indices [0,len(sa)], but we emulate
	// row 0, since it is always zero, so we only need
	// len(sa) columns.
	nrow, ncol := asize-1, len(bwt)
	table := make([]int, nrow*ncol)
	otab := OTab{nrow, ncol, table}

	// The character at the beginning of bwt gets a count
	// of one at row one.
	otab.set(bwt[0], 1, 1)

	// The remaining entries either copies or increment from
	// the previous column. We count a from 1 to alpha size
	// to skip the sentinel, then -1 for the index
	for a := 1; a < asize; a++ {
		ba := byte(a) // get the right type for accessing otab
		for i := 2; i <= len(bwt); i++ {
			val := otab.get(ba, i-1)
			if bwt[i-1] == ba {
				val++
			}

			otab.set(ba, i, val)
		}
	}

	return &otab
}

func Rbwt(y string) string {
	z := []byte(y)
	sigma := 256
	ctab := NewCTab(z, sigma) // Lazyness makes me use 256 for the alphabet size.
	otab := NewOTab(z, sigma) // It would be better to map the string to a smaller alphabet.

	x := make([]byte, len(y)-1) // y has the sentinel; x shouldn't have it
	i := 0

	// We start at len(y) - 2 because we already
	// (implicitly) have the sentinel at len(y) - 1
	// and this way we don't need to start at the index
	// in bwt that has the sentinel (so we save a search).
	for j := len(y) - 2; j >= 0; j-- {
		a := y[i]
		x[j] = a
		i = ctab.Rank(a) + otab.Rank(a, i)
	}

	return string(x)
}
