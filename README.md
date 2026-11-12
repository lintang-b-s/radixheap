# radixheap

[![GoDoc](https://godoc.org/github.com/lintang-b-s/radixheap?status.svg)](https://godoc.org/github.com/lintang-b-s/radixheap)
[![codecov](https://codecov.io/github/lintang-b-s/radixheap/graph/badge.svg?token=AVXOSD6FQH)](https://codecov.io/github/lintang-b-s/radixheap)

an radix heap implementation for Go. Radix heap [[1]](#ref1) is a monotone priority queue using binary representation of unsigned integers. A monotone priority is a priority queue that the keys returned by successive extract-min calls are monotonically increasing, like in dijkstra algorithm. The cost of extract-min operation is amortized O(logC), where C=maximum key of any priority queue node (in this implementation C=MaxUint32/MaxUint64), and insert is amortized O(1). A radix heap can be implemented to speed up Dijkstra's algorithm when edge weights are unsigned integers, as described in the paper [[1]](#ref1).

## Usage

### Installing

To start using radixheap, install Go and run `go get`:

```sh
$ go get -u github.com/lintang-b-s/radixheap
```

### Basic operations

```go
// create a radix heap
rdx := NewRadixHeap[uint32, uint32]()

// insert (key, val) pairs
val := 20
rdx.Insert(2, val)

val2 := 50
rdx.Insert(90, val2)

// extract-min
minKey, minVal := rdx.ExtractMin()
fmt.Print(minKey, minVal) // print 2, 20

```

## Benchmark

### Simple insert & extract-min operations

Simple insert & extract-min operations, preinsert 32 millions pairs.

```
BenchmarkRadixHeapSimple-12     18651518                77.19 ns/op           45 B/op          0 allocs/op

BenchmarkStdLibHeapSimple-12            33725971               573.8 ns/op             8 B/op          1 allocs/op
```

### Dijkstra SSSP USA road network DIMACS

Dijkstra Single-source shortest path (SSSP) on DIMACS 9th implementation challenge USA Road Network Graph ([download](https://web.archive.org/web/20261020102152/https://www.diag.uniroma1.it/~challenge9/download.shtml)), with 23million vertices & 58 million edges.

```
BenchmarkDijkstraRadixHeap-12    	      15	         4.462 s/op	199481824 B/op	     368 allocs/op

BenchmarkDijkstraStdLibHeap-12    	       3	        10.17 s/op	805779026 B/op	25575264 allocs/op
```

## References and Acknowledgment

<a id="ref1"></a>[1] Ravindra K. Ahuja, Kurt Mehlhorn, James Orlin, and Robert E. Tarjan. 1990. Faster algorithms for the shortest path problem. J. ACM 37, 2 (April 1990), 213–223. https://doi.org/10.1145/77600.77615

this radix heap implementation is adapted from the radix heap implementation by takuya akiba: https://github.com/iwiwi/radix-heap

## License

radixheap source code is available under the MIT License.
