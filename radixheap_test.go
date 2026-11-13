package radixheap

import (
	"container/heap"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestInsertExtractUint32(t *testing.T) {
	testCases := []struct {
		name            string
		input           []uint32
		wantSortedInput []uint32
	}{
		{
			name:            "input [90, 20, 60, 100, 70, 30, 40, 80, 10, 50], expected sorted input=[10, 20, 30, 40, 50, 60, 70, 80, 90, 100]",
			input:           []uint32{90, 20, 60, 100, 70, 30, 40, 80, 10, 50},
			wantSortedInput: []uint32{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		},
		{
			name:            "duplicate keys",
			input:           []uint32{15, 10, 15, 0, 10, 0, 10},
			wantSortedInput: []uint32{0, 0, 10, 10, 10, 15, 15},
		},
		{
			name:            "large uint32 values",
			input:           []uint32{4294967294, 1000000000, 500, 4294967290, 0},
			wantSortedInput: []uint32{0, 500, 1000000000, 4294967290, 4294967294},
		},
		{
			name:            "exactly one item per bucket",
			input:           []uint32{1 << 31, 1 << 16, 1 << 8, 1 << 4, 1 << 2, 1 << 1, 1 << 0, 0},
			wantSortedInput: []uint32{0, 1, 2, 4, 16, 256, 65536, 1 << 31},
		},

		{
			name: "exactly one item per bucket, each bucket filled",
			input: []uint32{
				1 << 17, 1 << 3, 1 << 29, 0, 1 << 11, 1 << 24, 1 << 0, 1 << 19,
				1 << 7, 1 << 30, 1 << 14, 1 << 2, 1 << 26, 1 << 9, 1 << 21, 1 << 5,
				1 << 28, 1 << 12, 1 << 1, 1 << 23, 1 << 16, 1 << 4, 1 << 27, 1 << 10,
				1 << 20, 1 << 6, 1 << 25, 1 << 15, 1 << 8, 1 << 22, 1 << 13, 1 << 18, 1 << 31,
			},
			wantSortedInput: []uint32{
				0,
				1 << 0, 1 << 1, 1 << 2, 1 << 3, 1 << 4, 1 << 5, 1 << 6, 1 << 7,
				1 << 8, 1 << 9, 1 << 10, 1 << 11, 1 << 12, 1 << 13, 1 << 14, 1 << 15,
				1 << 16, 1 << 17, 1 << 18, 1 << 19, 1 << 20, 1 << 21, 1 << 22, 1 << 23,
				1 << 24, 1 << 25, 1 << 26, 1 << 27, 1 << 28, 1 << 29, 1 << 30, 1 << 31,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rdx := NewRadixHeap[uint32, uint32]()
			for _, key := range tc.input {
				rdx.Insert(key, key)
			}

			gotSortedInput := make([]uint32, 0, len(tc.input))
			for !rdx.Empty() {
				skey, _ := rdx.ExtractMin()
				gotSortedInput = append(gotSortedInput, skey)
			}

			if len(gotSortedInput) != len(tc.wantSortedInput) {
				t.Errorf("length of gotSortedInput: %v not the same as length of  wantSortedInput: %v", len(gotSortedInput),
					len(tc.wantSortedInput))
			}

			for i := 0; i < len(gotSortedInput); i++ {
				if gotSortedInput[i] != tc.wantSortedInput[i] {
					t.Errorf("%v-th element of gotSortedInput:%v not the same as %v-th element of wantSortedInput: %v", i,
						gotSortedInput[i], i, tc.wantSortedInput[i])
				}
			}
		})
	}
}

func TestInsertExtractUint64(t *testing.T) {
	testCases := []struct {
		name            string
		input           []uint64
		wantSortedInput []uint64
	}{
		{
			name:            "input [90, 20, 60, 100, 70, 30, 40, 80, 10, 50], expected sorted input=[10, 20, 30, 40, 50, 60, 70, 80, 90, 100]",
			input:           []uint64{90, 20, 60, 100, 70, 30, 40, 80, 10, 50},
			wantSortedInput: []uint64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		},
		{
			name:            "duplicate keys",
			input:           []uint64{15, 10, 15, 0, 10, 0, 10},
			wantSortedInput: []uint64{0, 0, 10, 10, 10, 15, 15},
		},
		{
			name:            "large uint64 values",
			input:           []uint64{4294967294, 1000000000, 500, 4294967290, 0},
			wantSortedInput: []uint64{0, 500, 1000000000, 4294967290, 4294967294},
		},
		{
			name:            "exactly one item per bucket",
			input:           []uint64{1 << 31, 1 << 16, 1 << 8, 1 << 4, 1 << 2, 1 << 1, 1 << 0, 0},
			wantSortedInput: []uint64{0, 1, 2, 4, 16, 256, 65536, 1 << 31},
		},

		{
			name: "exactly one item per bucket",
			input: []uint64{
				1 << 17, 1 << 3, 1 << 29, 0, 1 << 11, 1 << 24, 1 << 0, 1 << 19,
				1 << 7, 1 << 30, 1 << 14, 1 << 2, 1 << 26, 1 << 9, 1 << 21, 1 << 5,
				1 << 28, 1 << 12, 1 << 1, 1 << 23, 1 << 16, 1 << 4, 1 << 27, 1 << 10,
				1 << 20, 1 << 6, 1 << 25, 1 << 15, 1 << 8, 1 << 22, 1 << 13, 1 << 18,
			},
			wantSortedInput: []uint64{
				0,
				1 << 0, 1 << 1, 1 << 2, 1 << 3, 1 << 4, 1 << 5, 1 << 6, 1 << 7,
				1 << 8, 1 << 9, 1 << 10, 1 << 11, 1 << 12, 1 << 13, 1 << 14, 1 << 15,
				1 << 16, 1 << 17, 1 << 18, 1 << 19, 1 << 20, 1 << 21, 1 << 22, 1 << 23,
				1 << 24, 1 << 25, 1 << 26, 1 << 27, 1 << 28, 1 << 29, 1 << 30,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rdx := NewRadixHeap[uint64, uint64]()
			for _, key := range tc.input {
				rdx.Insert(key, key)
			}

			gotSortedInput := make([]uint64, 0, len(tc.input))
			for !rdx.Empty() {
				skey, _ := rdx.ExtractMin()
				gotSortedInput = append(gotSortedInput, skey)
			}

			if len(gotSortedInput) != len(tc.wantSortedInput) {
				t.Errorf("length of gotSortedInput: %v not the same as length of  wantSortedInput: %v", len(gotSortedInput),
					len(tc.wantSortedInput))
			}

			for i := 0; i < len(gotSortedInput); i++ {
				if gotSortedInput[i] != tc.wantSortedInput[i] {
					t.Errorf("%v-th element of gotSortedInput:%v not the same as %v-th element of wantSortedInput: %v", i,
						gotSortedInput[i], i, tc.wantSortedInput[i])
				}
			}
		})
	}
}

func TestInsertExtractInterleave(t *testing.T) {
	testCases := []struct {
		name         string
		input        [][]uint32
		extractCount []uint32
		want         []uint32
	}{
		{
			name: "insert [8, 2, 9, 5, 7],  extract 2 pairs, insert [10, 19, 13, 23, 24], extract 3 pairs, insert [101, 100], extract 1 pair, extract all remaining pairs",
			input: [][]uint32{
				{8, 2, 9, 5, 7},
				{10, 19, 13, 23, 24},
				{101, 100},
			},
			extractCount: []uint32{2, 3, 1},
			want:         []uint32{13, 19, 23, 24, 100, 101},
		},
		{
			name: "insert [5000, 3000, 4000, 9000], extract 1 pair, insert [6000, 12000], extract 2 pairs, extract all remaining pairs",
			input: [][]uint32{
				{5000, 3000, 4000, 9000},
				{6000, 12000},
			},
			extractCount: []uint32{1, 2},
			want:         []uint32{6000, 9000, 12000},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rdx := NewRadixHeap[uint32, uint32]()
			for i, batch := range tc.input {
				for _, key := range batch {
					rdx.Insert(key, key)
				}
				for range tc.extractCount[i] {
					rdx.ExtractMin()
				}
			}

			got := make([]uint32, 0)
			for !rdx.Empty() {
				key, _ := rdx.ExtractMin()
				got = append(got, key)
			}

			if !slices.Equal(got, tc.want) {
				t.Errorf("got= %v, want %v", got, tc.want)
			}
		})
	}
}
func TestHeapSortLike(t *testing.T) {
	n := uint32(5_000_000)
	rdx := NewRadixHeap[uint32, uint32]()
	wantSortedArr := make([]uint32, n)
	seed1, seed2 := rand.Uint64(), rand.Uint64()
	r := rand.New(rand.NewPCG(seed1, seed2))
	for i := uint32(0); i < n; i++ {
		key := r.Uint32N(math.MaxUint32)
		rdx.Insert(key, key)
		wantSortedArr[i] = key
	}
	slices.Sort(wantSortedArr)

	gotSortedArr := make([]uint32, n)
	for i := uint32(0); i < n; i++ {
		gotKey, _ := rdx.ExtractMin()
		gotSortedArr[i] = gotKey
	}

	if !slices.Equal(gotSortedArr, wantSortedArr) {
		t.Errorf("heapsortLike= %v, want %v", gotSortedArr, wantSortedArr)
	}
}

func TestInsertClearInsertExtract(t *testing.T) {
	testCases := []struct {
		name      string
		input     []uint32
		inputLast []uint32

		wantSortedInput []uint32
	}{

		{
			name:  "insert 99, 88, 22, 666, 888, 333, clears, insert tc.inputLast",
			input: []uint32{99, 88, 22, 666, 888, 333},
			inputLast: []uint32{
				1 << 17, 1 << 3, 1 << 29, 0, 1 << 11, 1 << 24, 1 << 0, 1 << 19,
				1 << 7, 1 << 30, 1 << 14, 1 << 2, 1 << 26, 1 << 9, 1 << 21, 1 << 5,
				1 << 28, 1 << 12, 1 << 1, 1 << 23, 1 << 16, 1 << 4, 1 << 27, 1 << 10,
				1 << 20, 1 << 6, 1 << 25, 1 << 15, 1 << 8, 1 << 22, 1 << 13, 1 << 18, 1 << 31,
			},
			wantSortedInput: []uint32{
				0,
				1 << 0, 1 << 1, 1 << 2, 1 << 3, 1 << 4, 1 << 5, 1 << 6, 1 << 7,
				1 << 8, 1 << 9, 1 << 10, 1 << 11, 1 << 12, 1 << 13, 1 << 14, 1 << 15,
				1 << 16, 1 << 17, 1 << 18, 1 << 19, 1 << 20, 1 << 21, 1 << 22, 1 << 23,
				1 << 24, 1 << 25, 1 << 26, 1 << 27, 1 << 28, 1 << 29, 1 << 30, 1 << 31,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rdx := NewRadixHeap[uint32, uint32]()
			for _, key := range tc.input {
				rdx.Insert(key, key)
			}

			rdx.Clear()
			for _, key := range tc.inputLast {
				rdx.Insert(key, key)
			}

			gotSortedInput := make([]uint32, 0, len(tc.input))
			for !rdx.Empty() {
				skey, _ := rdx.ExtractMin()
				gotSortedInput = append(gotSortedInput, skey)
			}

			if !slices.Equal(gotSortedInput, tc.wantSortedInput) {
				t.Errorf("heapsortLike= %v, want %v", gotSortedInput, tc.wantSortedInput)
			}
		})
	}
}

// go test -run=^$ -bench=^BenchmarkRadixHeapSimple$  -benchtime=15s  -benchmem
// BenchmarkRadixHeapSimple-12     18651518                77.19 ns/op           45 B/op          0 allocs/op
func BenchmarkRadixHeapSimple(b *testing.B) {
	n := uint32(32_000_000)
	maxi := uint32(math.MaxUint32)
	rdx := NewRadixHeap[uint32, uint32]()
	seed1, seed2 := rand.Uint64(), rand.Uint64()
	r := rand.New(rand.NewPCG(seed1, seed2))
	for i := uint32(0); i < n; i++ {
		key := r.Uint32N(maxi)
		rdx.Insert(uint32(key), i)
	}

	for b.Loop() {
		mkey, _ := rdx.ExtractMin()
		rkey := mkey + r.Uint32N(maxi-mkey)
		rdx.Insert(rkey, rkey)
	}
}

// ------ stdlib container/heap benchmark ------
// adapted from: https://pkg.go.dev/container/heap

// An intHeap is a min-heap of ints.
type intHeap []uint32

func (h intHeap) Len() int           { return len(h) }
func (h intHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *intHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(uint32))
}

func (h *intHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// go test -run=^$ -bench=^BenchmarkStdLibHeapSimple$  -benchtime=15s  -benchmem
// BenchmarkStdLibHeapSimple-12            33725971               573.8 ns/op             8 B/op          1 allocs/op
func BenchmarkStdLibHeapSimple(b *testing.B) {
	n := uint32(32_000_000)
	maxi := uint32(math.MaxUint32)
	h := &intHeap{}
	heap.Init(h)
	seed1, seed2 := rand.Uint64(), rand.Uint64()
	r := rand.New(rand.NewPCG(seed1, seed2))
	for i := uint32(0); i < n; i++ {
		key := r.Uint32N(maxi)
		heap.Push(h, key)
	}

	for b.Loop() {
		mkey := heap.Pop(h).(uint32)
		rkey := mkey + r.Uint32N(maxi-mkey)
		heap.Push(h, rkey)
	}
}
