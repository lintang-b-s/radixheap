package radixheap

// implementation of radix heap (rahuja et al. (1990) [1])
// ref:
// [1] Ravindra K. Ahuja, Kurt Mehlhorn, James Orlin, and Robert E. Tarjan. 1990.
// Faster algorithms for the shortest path problem.
// J. ACM 37, 2 (April 1990), 213–223. https://doi.org/10.1145/77600.77615
// adapted from takuya akiba radix heap implementation: https://github.com/iwiwi/radix-heap
// this implementation follows section 5 variant of section 2 one-level radix heap described in section 5 & section 2 ref [1]

import (
	"math/bits"
)

type UnsignedKeyType interface {
	uint32 | uint64
}

type pair[T UnsignedKeyType, G comparable] struct {
	key   T
	value G
}

func newPair[T UnsignedKeyType, G comparable](k T, v G) pair[T, G] {
	return pair[T, G]{k, v}
}

type RadixHeap[T UnsignedKeyType, G comparable] struct {
	// the size (number of pairs inserted) of the radix heap
	size uint64
	// last extracted key.
	// in dijkstra algorithm, each iteration we call  exctract-min to get the unexplored vertex with minimum distance from s.
	// this lastExtractedKey store the distance from s (or key of the priority queue node) of the latest extracted vertex in current iteration of dijkstra.
	lastExtractedKey T
	// each bucket contains (key, value) pairs, with lower buckets store pairs with smaller key.
	// this one-level radix heap contains ceil(log(C+1))+1 buckets, where C=maximum key of any priority queue node (in this implementation  C=MaxUint32/MaxUint64).
	buckets [][]pair[T, G]
	// i-th bucketMin contains the minimum-key pair of the i-th bucket
	bucketsMin []T
	// is the key type uint32
	isUint32 bool
}

func infinity[T UnsignedKeyType]() T {
	return ^T(0)
}

func NewRadixHeap[T UnsignedKeyType, G comparable]() *RadixHeap[T, G] {
	keyBits := bits.OnesCount64(uint64(^T(0)))
	bucketsCount := 0
	isUint32 := false
	switch keyBits {
	case 32:
		bucketsCount = 33
		isUint32 = true
	case 64:
		bucketsCount = 65
		isUint32 = false
	}
	bucketsMin := make([]T, bucketsCount)
	for i := 0; i < bucketsCount; i++ {
		bucketsMin[i] = infinity[T]()
	}
	return &RadixHeap[T, G]{
		size:             0,
		lastExtractedKey: 0,
		buckets:          make([][]pair[T, G], bucketsCount),
		bucketsMin:       bucketsMin,
		isUint32:         isUint32,
	}
}

// findBucket find bucket-i that will contain key x
// by computing i such that the position i-1 is the largest position (0-indexed) in which
// the binary representation of x and r.lastExtractedKey differ. O(1).
// this implementation follows section 5 variant of section 2 one-level radix heap described in section 5 & section 2 ref [1]
func (r *RadixHeap[T, G]) findBucket(x T) uint64 {
	if r.isUint32 {
		i := 32 - bits.LeadingZeros32(uint32(x^r.lastExtractedKey))
		return uint64(i)
	}
	i := 64 - bits.LeadingZeros64(uint64(x^r.lastExtractedKey))
	return uint64(i)
}

// Insert insert (key, val) pair into the radix heap. O(1) amortized
func (r *RadixHeap[T, G]) Insert(key T, val G) {
	r.size++
	bucket := r.findBucket(key)                                      // find bucket to insert pair (key,val). O(1).
	r.buckets[bucket] = append(r.buckets[bucket], newPair(key, val)) // O(1) amortized dynamic array doubling,
	r.bucketsMin[bucket] = min(r.bucketsMin[bucket], key)
}

// ExtractMin returns minimum-key pair (key,val) and delete this pair from the radix heap.
// extract-min performed as follows. if bucket 0 is nonempty, return any pair in bucket 0.
// otherwise, find the nonempty bucket of smallest index, say bucket j, returns minimum-key pair in this bucket j,
// then distribute the remaining pairs in bucket j among buckets of smaller index by using findBucket().
// ilustration of this insert & extract-min operations of radix heap (the insert() and extractmin() in this implementation differs slightly): https://ocw.mit.edu/courses/15-082j-network-optimization-fall-2010/b9a6748a2d111ab5831efab32f6c47ff_MIT15_082JF10_av09.pdf
// number of times extract-min can be called in nonempty radix heap, is the same as number of insert operations which at most n times.
// any sequence of n times insert and extract-min operations on initially empty radix-heap has a total cost at most O(n*logC), where C=maximum key of any priority queue node (in this implementation C=MaxUint32/MaxUint64)
// why? after all n extract-min operations, for each pair, the distribute operations for this pair from old position to the latest position (until this pair extracted from the pq) has total cost at most O(logC),
// this is because the number of buckets is ceil(log(C+1))+1. summing over all n pairs inserted in this n sequence of insert and extract-min operations, we got total cost O(n*logC)
// the amortized (aggregate analysis) cost of an extract-min operation is O(n*logc)/n=O(logC)
// amortized O(logC), where C=maximum key of any priority queue node (in this implementation  C=MaxUint32/MaxUint64)
func (r *RadixHeap[T, G]) ExtractMin() (T, G) {
	if len(r.buckets[0]) == 0 {
		// if first bucket is empty
		// find first non-empty bucket
		j := uint64(1)
		for ; len(r.buckets[j]) == 0; j++ { // O(logC), because the number of buckets is ceil(log(C+1))+1
		}
		r.lastExtractedKey = r.bucketsMin[j] // update last extracted key

		// distribute all pairs (key,val)'s in j-th bucket to lower buckets
		for k := 0; k < len(r.buckets[j]); k++ { // amortized O(logC)
			pair := r.buckets[j][k]
			lowerBucket := r.findBucket(pair.key)
			r.buckets[lowerBucket] = append(r.buckets[lowerBucket], pair) // amortized O(1) dynamic array doubling
			r.bucketsMin[lowerBucket] = min(r.bucketsMin[lowerBucket], pair.key)
		}

		// empty the j-th bucket
		r.buckets[j] = r.buckets[j][:0]
		r.bucketsMin[j] = infinity[T]()
	}

	// return minimum-key pair
	// return last pair in first bucket
	last := len(r.buckets[0]) - 1
	minPair := r.buckets[0][last]
	r.lastExtractedKey = minPair.key
	r.buckets[0] = r.buckets[0][:last] // remove the last pair in the first bucket
	if len(r.buckets[0]) == 0 {
		r.bucketsMin[0] = infinity[T]()
	}
	r.size--
	return minPair.key, minPair.value
}

func (r *RadixHeap[T, G]) Empty() bool {
	return r.size == 0
}
