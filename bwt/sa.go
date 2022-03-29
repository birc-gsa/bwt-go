package bwt

/*

   If you are here, you must want to learn about prefix doubling.
   Good for you! You make me proud.

   This is a suffix array construction algorithm that, like many others,
   are based on bucket sorting, but it doesn't do divide-and-conquer.
   Instead, it sorts the suffixes according to progressively longer
   prefixes, doubling their length in each iteration (thus the name).

   Take the string 'mississippi$' (it is as good an example as any).
   The suffixes are shown below, with their length-1 prefixes in square
   brackets:

   [m]ississippi$
   [i]ssissippi$
   [s]sissippi$
   [s]issippi$
   [i]ssippi$
   [s]sippi$
   [s]ippi$
   [i]ppi$
   [p]pi$
   [p]i$
   [i]$
   [$]

   If we sort them with respect to these length-1 prefixes we get

   [$]
   [i]ssissippi$
   [i]ssippi$
   [i]ppi$
   [i]$
   [m]ississippi$
   [p]pi$
   [p]i$
   [s]sissippi$
   [s]issippi$
   [s]sippi$
   [s]ippi$

   The suffixes are not completely sorted, but they are sorted with
   respect to their first letter. That's all we get with one sweep of
   bucket sort.

   Then we sort them with respect to their length-2 prefixes:

   [$ ]
   [i$]
   [ip]pi$
   [is]sissippi$
   [is]sippi$
   [mi]ssissippi$
   [pi]$
   [pp]i$
   [si]ssippi$
   [si]ppi$
   [ss]issippi$
   [ss]ippi$

   Again, they are not quite sorted. We can recognise that there might
   be more work to do by the prefixes not being unique, similar to
   how it works with skew or SAIS; so if the prefixes we have sorted
   with respect to aren't unique we need to do more, and we will double
   the length of prefixes we look at, now to length 4, and sort again.

   [$   ]
   [i$  ]
   [ippi]$
   [issi]ssippi$
   [issi]ppi$
   [miss]issippi$
   [pi$ ]
   [ppi$]
   [sipp]i$
   [siss]ippi$
   [ssip]pi$
   [ssis]sippi$

   Once again the prefixes aren't unique (we have 'issi' twice) and thus
   we need another doubling of prefix length, now to length 8.

   [$       ]
   [i$      ]
   [ippi$   ]
   [issippi$]
   [ississip]pi$
   [mississi]ppi$
   [pi$     ]
   [ppi$    ]
   [sippi$  ]
   [sissippi]$
   [ssippi$ ]
   [ssissipp]i$

   Now, the prefixes are unique, so we know that we have sorted all
   the strings.

   What is the worst-case number of times that we need to sort? Well,
   if we double the length of the prefixes every time, we will cover
   the entire string in O(log n) iterations, and then the prefix must
   be unique, so we need at most O(log n) iterations.

   Usually, though, prefixes are unique much before. If we had random
   strings, the chance that two strings are identical drops off
   exponetially, and in that case we only expect a constant number of
   iterations. So, in practice, the approach can be quite efficient.

   Now, what about the sorting? We can obviously sort the length-1 prefix
   in O(n) using a bucket sort, but how do we handle the longer prefixes?

   Here, we use a trick similar to what we have seen in skew and sais.
   We replace strings with numbers, such that the numbers preserve the
   lexicographical ordering of the strings. For length-1 prefixes, we can
   just use the letters we have, or we can map them if you want to. Let's
   do the mapping:

   $ => 0
   i => 1
   m => 2
   p => 3
   s => 4

   then the length-1 prefix are

    0: [2]ississippi$
    1: [1]ssissippi$
    2: [4]sissippi$
    3: [4]issippi$
    4: [1]ssippi$
    5: [4]sippi$
    6: [4]ippi$
    7: [1]ppi$
    8: [3]pi$
    9: [3]i$
   10: [1]$
   11: [0]

   There is nothing new here. However, when we want to sort with length-2
   prefixes we can use that the key for suffix i should be the number for i
   followed by the number for i+1:

    0: [2]ississippi$ is followed by 1: [1]ssissippi$ so the pair is [2,1]
    1: [1]ssissippi$ is followed by 2: [4]sissippi$ so the pair is [1,4]

    and so on:

    0: [2,1]ssissippi$
    1: [1,4]sissippi$
    2: [4,4]issippi$
    3: [4,1]ssippi$
    4: [1,4]sippi$
    5: [4,4]ippi$
    6: [4,1]ppi$
    7: [1,3]pi$
    8: [3,3]i$
    9: [3,1]$
   10: [1,0]
   11: [0,0]  // put 0 if there isn't a i+1

   Again, there isn't much new here, I've just renamed the letters. However,
   I have two numbers of magnitude <= n and I can radix sort these in time
   O(n).

   11: [0,0]
   10: [1,0]
    7: [1,3]pi$
    1: [1,4]sissippi$
    4: [1,4]sippi$
    0: [2,1]ssissippi$
    9: [3,1]$
    8: [2,2]i$
    3: [4,1]ssippi$
    6: [4,1]ppi$
    2: [4,4]issippi$
    5: [4,4]ippi$

   and once I have the pairs sorted, I can assign new numbers to the pairs
   to get a single number for each of the prefixes.

   [0,0] => 0
   [1,0] => 1
   [1,3] => 2
   [1,4] => 3
   [1,4] => 3
   [2,1] => 4
   [3,1] => 5
   [2,2] => 6
   [4,1] => 7
   [4,1] => 7
   [4,4] => 8
   [4,4] => 8

   and then you have the 2-prefix encoding

   11: [0]
   10: [1]
    7: [2]pi$
    1: [3]sissippi$
    4: [3]sippi$
    0: [4]ssissippi$
    9: [5]$
    8: [6]i$
    3: [7]ssippi$
    6: [7]ppi$
    2: [8]issippi$
    5: [8]ippi$

   To get pairs of numbers for the length-4 prefixes, you use the same
   approach of pairing up numbers, but now you don't want to look at i
   and i+1, but at i and i+2.

   When you looked at i and i+1 you concatenated length-1 prefixes

        i
   ....[x] y  z...
   .... x [y] z...
          i+1

   but if you look at i and i+2 you can concatenate length-2 prefixes

        i
   ....[xy]  zw ...
   .... xy  [zw] ...
            i+2

   In general, when you want to concatenate length-k prefixes, you will
   look at i and i+k and take the numbers from these.

   For the length-2 prefixes, with the numbers they have are

   11: [0]
   10: [1]
    7: [2]pi$
    1: [3]sissippi$
    4: [3]sippi$
    0: [4]ssissippi$
    9: [5]$
    8: [6]i$
    3: [7]ssippi$
    6: [7]ppi$
    2: [8]issippi$
    5: [8]ippi$

   but let us rearrange them so the suffixes come in the order they have
   in x, so it is easier to see which suffixes are two letters apart

    0: [4]ssissippi$
    1: [3]sissippi$
    2: [8]issippi$
    3: [7]ssippi$
    4: [3]sippi$
    5: [8]ippi$
    6: [7]ppi$
    7: [2]pi$
    8: [6]i$
    9: [5]$
   10: [1]
   11: [0]

   To get the pair for the first suffix, look at 0 and 2:

    0: [4]ssissippi$
    2: [8]issippi$

   which gives you (4,8). For the second suffix, look at 1 and 3:

    1: [3]sissippi$
    3: [7]ssippi$

   so it should get the key (3,7). Continue like that to get pairs for
   all the suffixes.

    0: [4,8]issippi$
    1: [3,7]ssippi$
    2: [8,3]sippi$
    3: [7,8]ippi$
    4: [3,7]ppi$
    5: [8,2]pi$
    6: [7,6]i$
    7: [2,5]$
    8: [6,1]
    9: [5,0]
   10: [1,0]
   11: [0,0]

   Sort these

   11: [0,0]
   10: [1,0]
    7: [2,5]$
    1: [3,7]ssippi$
    4: [3,7]ppi$
    0: [4,8]issippi$
    9: [5,0]
    8: [6,1]
    6: [7,6]i$
    3: [7,8]ippi$
    5: [8,2]pi$
    2: [8,3]sippi$

   and assign the suffixes new numbers

   11: [0]                [0,0] => 0
   10: [1]                [1,0] => 1
    7: [2]$               [2,5] => 2
    1: [3]ssippi$         [3,7] => 3
    4: [3]ppi$            [3,7] => 3
    0: [4]issippi$        [4,8] => 4
    9: [5]                [5,0] => 5
    8: [6]                [6,1] => 6
    6: [7]i$              [7,6] => 7
    3: [8]ippi$           [7,8] => 8
    5: [9]pi$             [8,2] => 9
    2: [10]sippi$         [8,3] => 10

   and then do the whole thing again (but now with i and i+4 when
   building pairs).

   Since you can radix sort pairs (i,j) where i,j <= n in O(n), each
   sort iteration runs in O(n), and you can trivially assign numbers
   to the suffixes in O(n) as well.

   The only remaining caveat is that while, true, you can do a radix sort
   of pairs of numbers i,j <= n in O(n), you probably don't have the memory
   for it (nor the desire to spend the full time on it either, frankly).

   If you need bucket tables of size n, and n is in the hundreds of millions,
   you might very well run out of memory. But that isn't a problem. When
   I said that you could do radix sort, I didn't mean just two iterations of
   bucket sort on the two integers. You can use radix sort on the integers
   as well. If, say, the integers are 32 bit, we can split them into 4 bytes
   each, so sorting the pair becomes 8 iterations of bucket sort. And it ends
   up even better than that, because it would be 8 iterations if we had to sort
   the pairs from scratch, but we don't have to. After one iteration, the numbers
   are already sorted with respect to the first of the two integers in the pair.
   When we construct the pairs (i,j), we can exploit that we already have the
   i-component in order. Within each i-block, we need to sort the j's, but that
   is just one integer, so four bucket-sorts of bytes.

   There are different ways of doing this, but they are pretty much all
   highly efficient. The code below is just one approach, and should you
   feel adventurous, you can try other approaches. In any case, I suggest
   you read the code to get a feeling for this approach to suffix array
   construction.

*/

// Compute the rank each suffix has if we only look at the first character
func calcRank0(x []byte) (rank []int32, sigma int32) {
	alpha := [256]int32{}
	rank = make([]int32, len(x)+1)

	// run through x and tag occurring letters
	for _, a := range x {
		alpha[a] = 1
	}

	// assign numbers to each occurring letter
	sigma = 1 // start at 1, 0 is the sentinel
	for a := 0; a < 256; a++ {
		if alpha[a] == 1 {
			alpha[a] = sigma
			sigma++
		}
	}

	// map each letter from x to its number and place them in mapped
	for i := 0; i < len(x); i++ {
		rank[i] = alpha[x[i]]
	}
	// rank[len(x)] is already the sentinel (0) because make zeros.

	return rank, sigma
}

// Give us the first "suffix array"; just the indicies from 0 to n.
func sa0(n int) (sa []int32) {
	sa = make([]int32, n)
	for i := 0; i < n; i++ {
		sa[i] = int32(i)
	}
	return sa
}

// Get the rank for index i with padded zeros after the end
func getRank(rank []int32, i int32) int32 {
	if int(i) < len(rank) {
		return rank[i]
	}

	return 0
}

// Radix sort sa with respect to rank. k is the offset to use when
// accessing the second integer in the prefix-pair. buf is just a
// buffer we use for the sort.
func radixSortBuckets(rank, sa, buf []int32, k int32) {
	sa_p, buf_p := &sa, &buf

	for shift := 0; shift < 32; shift += 8 {
		buckets := [256]int32{}
		for i := 0; i < len(sa); i++ {
			b := getRank(rank, (*sa_p)[i]+k) >> shift
			buckets[b]++
		}
		for acc, i := int32(0), 0; i < 256; i++ {
			b := buckets[i]
			buckets[i] = acc
			acc += b
		}
		// then place sa[i] in buckets
		for i := 0; i < len(sa); i++ {
			b := getRank(rank, (*sa_p)[i]+k) >> shift
			(*buf_p)[buckets[b]] = (*sa_p)[i]
			buckets[b]++
		}

		// flip sa and buf for next iteration...
		sa_p, buf_p = buf_p, sa_p
	}

	// We run for an even number of iterations (four) so at the end, the
	// result is back in (*sa_p) == sa.
}

// Sort the elements sa according to the rank[sa[i]+k]
// (with padded zero sentinels) using a radix sort over
// 8-bit sub-integers. The result is left in sa; buf
// is a scratch buffer.
func radixSort(k int32, rank, sa, buf []int32) {
	// sa is already sorted, so we need to sort sa+k for each bucket.
	b_start, b_end := 0, 0
	for b_start < len(sa) {
		for b_end < len(sa) && rank[sa[b_start]] == rank[sa[b_end]] {
			b_end++
		}

		// Sort the bucket if it is more than one element large
		if (b_end - b_start) > 1 {
			radixSortBuckets(rank, sa[b_start:b_end], buf[b_start:b_end], k)
		}

		b_start = b_end
	}
}

// For each element in sa, assumed sorted according to
// rank[sa[i]],rank[sa[i]+k], work out what rank
// (order of rank[sa[i]],rank[sa[i]+k]) each element has
// and put the result in out.
func updateRank(sa, rank, out []int32, k int32) (sigma int32) {

	// We have 32-bit integers. To get pairs that we can
	// readily compare, we pack them in 64-bit integers
	pair := func(i, k int32) int64 {
		return int64(rank[sa[i]])<<32 | int64(getRank(rank, sa[i]+k))
	}

	a := int32(0)
	out[sa[0]] = a

	prev_pair := pair(0, k)
	for i := 1; i < len(sa); i++ {
		cur_pair := pair(int32(i), k)
		if prev_pair != cur_pair {
			a++
		}
		prev_pair = cur_pair
		out[sa[i]] = a
	}

	sigma = a + 1 // alphabet size is one past the largest letter
	return sigma
}

func PrefixDoubling(x string) (sa []int32) {
	sa = sa0(len(x) + 1)
	buf := make([]int32, len(sa))
	rank, sigma := calcRank0([]byte(x))
	radixSortBuckets(rank, sa, buf, 0)

	buf_p, rank_p := &buf, &rank
	for k := int32(1); int(sigma) < len(rank); k *= 2 {
		radixSort(k, *rank_p, sa, *buf_p)
		sigma = updateRank(sa, *rank_p, *buf_p, k)
		buf_p, rank_p = rank_p, buf_p
	}

	return sa
}
