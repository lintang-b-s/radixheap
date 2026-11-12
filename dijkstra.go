package radixheap

import "container/heap"

type edge[T UnsignedKeyType] struct {
	v      uint32
	weight T
}

// dijkstra dijkstra Single-Source Shortest path. lazy-deletion version + radix heap.
func dijkstra[T UnsignedKeyType](s uint32, adjList [][]edge[T]) []T {
	n := len(adjList)
	rdx := NewRadixHeap[T, uint32]()
	dist := make([]T, n)
	for i := 0; i < n; i++ {
		dist[i] = infinity[T]()
	}
	dist[s] = 0
	rdx.Insert(0, s)
	for !rdx.Empty() {
		uCost, u := rdx.ExtractMin()
		if uCost > dist[u] {
			continue
		}
		for _, e := range adjList[u] {
			v, eWeight := e.v, e.weight
			nvCost := uCost + eWeight
			if nvCost < dist[v] {
				dist[v] = nvCost
				rdx.Insert(nvCost, v)
			}
		}
	}
	return dist
}

// ------- stdlib container/heap version ---------
// priority queue code adapted from: https://pkg.go.dev/container/heap

// An item is something we manage in a priority queue.
type item[T UnsignedKeyType] struct {
	vertex   uint32
	priority T
	index    int
}

type priorityQueue[T UnsignedKeyType] []*item[T]

func (pq priorityQueue[T]) Len() int { return len(pq) }

func (pq priorityQueue[T]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(*item[T])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func dijkstraStdLibPq[T UnsignedKeyType](s uint32, adjList [][]edge[T]) []T {
	n := len(adjList)
	dist := make([]T, n)
	for i := range dist {
		dist[i] = infinity[T]()
	}
	dist[s] = 0

	pq := make(priorityQueue[T], 0)
	heap.Push(&pq, &item[T]{vertex: s, priority: 0})

	for pq.Len() != 0 {
		uu := heap.Pop(&pq).(*item[T])
		uCost, u := uu.priority, uu.vertex
		if uCost > dist[u] {
			continue
		}
		for _, e := range adjList[u] {
			v, eWeight := e.v, e.weight
			nvCost := uCost + eWeight
			if nvCost < dist[v] {
				dist[v] = nvCost
				heap.Push(&pq, &item[T]{vertex: v, priority: nvCost})
			}
		}
	}
	return dist
}
