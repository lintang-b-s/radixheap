package radixheap

import (
	"bufio"
	"io"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// Vertex constants for Undirected Graph (A, B, C, D, E)
// figure 3 https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-fall-2011/6277a1f06100c26a7ff21031af6757b5_MIT6_006F11_lec16.pdf
const (
	uA uint32 = iota
	uB
	uC
	uD
	uE
)

// Vertex constants for Directed Graph G (s, a, b, c, d)
// from https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-spring-2020/d819e7f4568aced8d5b59e03db6c7b67_MIT6_006S20_lec13.pdf
const (
	dS uint32 = iota
	dA
	dB
	dC
	dD
)

// Vertex constants for directed graph G (1, 2, 3, 4, 5, 6)
// from https://ocw.mit.edu/courses/15-082j-network-optimization-fall-2010/b9a6748a2d111ab5831efab32f6c47ff_MIT15_082JF10_av09.pdf
const (
	r1 uint32 = iota
	r2
	r3
	r4
	r5
	r6
)

func TestDijkstraRadixHeap(t *testing.T) {

	tests := []struct {
		name    string
		start   uint32
		adjList [][]edge[uint32]
		want    []uint32
	}{
		{
			// figure 3 https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-fall-2011/6277a1f06100c26a7ff21031af6757b5_MIT6_006F11_lec16.pdf
			name:  "Undirected graph from vertex A",
			start: uA,
			adjList: [][]edge[uint32]{
				uA: {{v: uB, weight: 19}, {v: uC, weight: 7}},
				uB: {{v: uA, weight: 19}, {v: uC, weight: 11}, {v: uD, weight: 4}},
				uC: {{v: uA, weight: 7}, {v: uB, weight: 11}, {v: uD, weight: 15}, {v: uE, weight: 5}},
				uD: {{v: uB, weight: 4}, {v: uC, weight: 15}, {v: uE, weight: 13}},
				uE: {{v: uC, weight: 5}, {v: uD, weight: 13}},
			},
			// want shortest path cost: [A: 0, B: 18, C: 7, D: 22, E: 12]
			want: []uint32{0, 18, 7, 22, 12},
		},
		{
			// https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-spring-2020/d819e7f4568aced8d5b59e03db6c7b67_MIT6_006S20_lec13.pdf
			name:  "Directed graph G from vertex s",
			start: dS,
			adjList: [][]edge[uint32]{
				dS: {{v: dA, weight: 10}, {v: dC, weight: 3}},
				dA: {{v: dB, weight: 2}, {v: dC, weight: 1}},
				dB: {{v: dD, weight: 7}},
				dC: {{v: dA, weight: 4}, {v: dB, weight: 8}, {v: dD, weight: 2}},
				dD: {{v: dB, weight: 5}},
			},
			// want shortest path cost: [s: 0, a: 7, b: 9, c: 3, d: 5]
			want: []uint32{0, 7, 9, 3, 5},
		},
		{
			// https://ocw.mit.edu/courses/15-082j-network-optimization-fall-2010/b9a6748a2d111ab5831efab32f6c47ff_MIT15_082JF10_av09.pdf
			name:  "Directed graph G from vertex 1",
			start: r1,
			adjList: [][]edge[uint32]{
				r1: {{v: r2, weight: 13}, {v: r3, weight: 0}, {v: r4, weight: 15}, {v: r5, weight: 20}},
				r2: {{v: r4, weight: 5}},
				r3: {{v: r5, weight: 9}},
				r4: {{v: r6, weight: 2}},
				r5: {{v: r6, weight: 8}},
				r6: {},
			},
			// want shortest path cost: [1: 0, 2: 13, 3: 0, 4: 15, 5: 9, 6: 17]
			want: []uint32{0, 13, 0, 15, 9, 17},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dijkstra(tt.start, tt.adjList)
			if !slices.Equal(got, tt.want) {
				t.Errorf("dijkstra(%d) = %v, want %v", tt.start, got, tt.want)
			}
		})
	}
}

func TestDijkstraStdLibPq(t *testing.T) {

	tests := []struct {
		name    string
		start   uint32
		adjList [][]edge[uint32]
		want    []uint32
	}{
		{
			// figure 3 https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-fall-2011/6277a1f06100c26a7ff21031af6757b5_MIT6_006F11_lec16.pdf
			name:  "Undirected graph from vertex A",
			start: uA,
			adjList: [][]edge[uint32]{
				uA: {{v: uB, weight: 19}, {v: uC, weight: 7}},
				uB: {{v: uA, weight: 19}, {v: uC, weight: 11}, {v: uD, weight: 4}},
				uC: {{v: uA, weight: 7}, {v: uB, weight: 11}, {v: uD, weight: 15}, {v: uE, weight: 5}},
				uD: {{v: uB, weight: 4}, {v: uC, weight: 15}, {v: uE, weight: 13}},
				uE: {{v: uC, weight: 5}, {v: uD, weight: 13}},
			},
			// want shortest path cost: [A: 0, B: 18, C: 7, D: 22, E: 12]
			want: []uint32{0, 18, 7, 22, 12},
		},
		{
			// https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-spring-2020/d819e7f4568aced8d5b59e03db6c7b67_MIT6_006S20_lec13.pdf
			name:  "Directed graph G from vertex s",
			start: dS,
			adjList: [][]edge[uint32]{
				dS: {{v: dA, weight: 10}, {v: dC, weight: 3}},
				dA: {{v: dB, weight: 2}, {v: dC, weight: 1}},
				dB: {{v: dD, weight: 7}},
				dC: {{v: dA, weight: 4}, {v: dB, weight: 8}, {v: dD, weight: 2}},
				dD: {{v: dB, weight: 5}},
			},
			// want shortest path cost: [s: 0, a: 7, b: 9, c: 3, d: 5]
			want: []uint32{0, 7, 9, 3, 5},
		},
		{
			// https://ocw.mit.edu/courses/15-082j-network-optimization-fall-2010/b9a6748a2d111ab5831efab32f6c47ff_MIT15_082JF10_av09.pdf
			name:  "Directed graph G from vertex 1",
			start: r1,
			adjList: [][]edge[uint32]{
				r1: {{v: r2, weight: 13}, {v: r3, weight: 0}, {v: r4, weight: 15}, {v: r5, weight: 20}},
				r2: {{v: r4, weight: 5}},
				r3: {{v: r5, weight: 9}},
				r4: {{v: r6, weight: 2}},
				r5: {{v: r6, weight: 8}},
				r6: {},
			},
			// want shortest path cost: [1: 0, 2: 13, 3: 0, 4: 15, 5: 9, 6: 17]
			want: []uint32{0, 13, 0, 15, 9, 17},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dijkstraStdLibPq(tt.start, tt.adjList)
			if !slices.Equal(got, tt.want) {
				t.Errorf("dijkstraStdLibPq(%d) = %v, want %v", tt.start, got, tt.want)
			}
		})
	}
}

// readLine reads a line from the buffered reader, trimming trailing CRLF.
func readLine(br *bufio.Reader) (string, error) {
	line, err := br.ReadString('\n')
	if err != nil && (err != io.EOF || len(line) == 0) {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func readDimacs(b *testing.B, mapCoordpath, mapEdgesPath string) [][]edge[uint32] {
	b.Helper()

	f, err := os.OpenFile(filepath.Clean(mapCoordpath), os.O_RDONLY, 0600) //nolint:gosec // benchmark input file
	if err != nil {
		b.Fatalf("could not open test file: %v", mapCoordpath)
	}
	defer f.Close()

	br := bufio.NewReader(f)

	// read comments gak penting
	for i := 0; i < 4; i++ {
		_, err = readLine(br)
		if err != nil {
			b.Fatal(err.Error())
		}
	}

	line, err := readLine(br)
	if err != nil {
		b.Fatal(err.Error())
	}

	ff := strings.Fields(line)
	n, err := strconv.Atoi(ff[4]) // number of vertices
	if err != nil {
		b.Fatal(err.Error())
	}

	// read comments gak penting lagi
	for i := 0; i < 2; i++ {
		_, err = readLine(br)
		if err != nil {
			b.Fatal(err.Error())
		}
	}

	const rounder = 1e8

	for v := 0; v < n; v++ { // vertex id 0 dummy vertex.. id vertex dari file dimacs mulai dari 1
		line, err = readLine(br)
		if err != nil {
			b.Fatal(err.Error())
		}
		ff := strings.Fields(line)
		_, err := strconv.Atoi(ff[1])
		if err != nil {
			b.Fatal(err.Error())
		}
		_, err = strconv.Atoi(ff[2])
		if err != nil {
			b.Fatal(err.Error())
		}
		_, err = strconv.Atoi(ff[3])
		if err != nil {
			b.Fatal(err.Error())
		}

	}

	fInputEdges, err := os.OpenFile(mapEdgesPath, os.O_RDONLY, 0600) //nolint:gosec // benchmark input file
	if err != nil {
		b.Fatalf("could not open test file: %v", mapEdgesPath)
	}
	defer fInputEdges.Close()

	br = bufio.NewReader(fInputEdges)

	// read comments gak penting
	for i := 0; i < 4; i++ {
		_, err = readLine(br)
		if err != nil {
			b.Fatal(err.Error())
		}
	}

	line, err = readLine(br)
	if err != nil {
		b.Fatal(err.Error())
	}

	ff = strings.Fields(line)
	m, err := strconv.Atoi(ff[3]) // number of edges
	if err != nil {
		b.Fatal(err.Error())
	}

	// read comments gak penting lagi
	for i := 0; i < 2; i++ {
		_, err = readLine(br)
		if err != nil {
			b.Fatal(err.Error())
		}
	}

	adjList := make([][]edge[uint32], n)
	minWeight := math.MaxInt64
	maxWeight := math.MinInt64
	for i := 0; i < m; i++ {
		line, err = readLine(br)
		if err != nil {
			b.Fatal(err.Error())
		}
		ff := strings.Fields(line)
		u, err := strconv.Atoi(ff[1])
		if err != nil {
			b.Fatal(err.Error())
		}
		v, err := strconv.Atoi(ff[2])
		if err != nil {
			b.Fatal(err.Error())
		}
		u--
		v--
		weight, err := strconv.Atoi(ff[3])
		if err != nil {
			b.Fatal(err.Error())
		}
		if weight < minWeight {
			minWeight = weight
		}
		if weight > maxWeight {
			maxWeight = weight
		}

		adjList[u] = append(adjList[u], edge[uint32]{uint32(v), uint32(weight)})
	}

	return adjList
}

// prequisite: download the USA road network data https://web.archive.org/web/20261020102152/https://www.diag.uniroma1.it/~challenge9/download.shtml
// or download via: https://drive.google.com/file/d/1PGu7qRvp11jFr1mIg7QpLJb3eEj6L5Tr/view?usp=sharing and https://drive.google.com/file/d/1shyuOvMPLDfFDS3r9REx2sR7GTKNXkwL/view?usp=sharing
// go test -run=^$ -bench=^BenchmarkDijkstraRadixHeap$  -benchtime=60s  -benchmem
// cpu: AMD Ryzen 5 7540U w/ Radeon(TM) 740M Graphics
// BenchmarkDijkstraRadixHeap-12                 20                 3.506 s/op     99193557 B/op        358 allocs/op
func BenchmarkDijkstraRadixHeap(b *testing.B) {
	graph := readDimacs(b, "USA-road-d.USA.co", "USA-road-t.USA.gr")
	n := len(graph)
	for b.Loop() {
		s := rand.Uint32N(uint32(n))
		dijkstra(s, graph)
	}
	b.ReportMetric(b.Elapsed().Seconds()/float64(b.N), "s/op")
	b.ReportMetric(0, "ns/op")
}

// go test -run=^$ -bench=^BenchmarkDijkstraStdLibHeap$  -benchtime=60s  -benchmem
// cpu: AMD Ryzen 5 7540U w/ Radeon(TM) 740M Graphics
// BenchmarkDijkstraStdLibHeap-12                 8                 7.909 s/op     505327374 B/op  25575164 allocs/op
func BenchmarkDijkstraStdLibHeap(b *testing.B) {
	graph := readDimacs(b, "USA-road-d.USA.co", "USA-road-t.USA.gr")
	n := len(graph)

	for b.Loop() {
		s := rand.Uint32N(uint32(n))
		dijkstraStdLibPq(s, graph)
	}
	b.ReportMetric(b.Elapsed().Seconds()/float64(b.N), "s/op")
	b.ReportMetric(0, "ns/op")
}
