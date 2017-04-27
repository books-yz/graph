package graph

import (
	"container/heap"
)

// ShortestPath computes a shortest path from v to w.
// Only edges with non-negative costs are included.
// The number dist is the length of the path, or -1 if w cannot be reached.
//
// The time complexity is O((|E| + |V|)⋅log|V|), where |E| is the number of edges
// and |V| the number of vertices in the graph.
func ShortestPath(g Iterator, v, w int) (path []int, dist int64) {
	parent, distances := ShortestPaths(g, v)
	path, dist = []int{}, distances[w]
	if dist == -1 {
		return
	}
	for v := w; v != -1; v = parent[v] {
		path = append(path, v)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return
}

// ShortestPaths computes the shortest paths from v to all other vertices.
// Only edges with non-negative costs are included.
// The number parent[w] is the predecessor of w on a shortest path from v to w,
// or -1 if none exists.
// The number dist[w] equals the length of a shortest path from v to w,
// or is -1 if w cannot be reached.
//
// The time complexity is O((|E| + |V|)⋅log|V|), where |E| is the number of edges
// and |V| the number of vertices in the graph.
func ShortestPaths(g Iterator, v int) (parent []int, dist []int64) {
	n := g.Order()
	dist = make([]int64, n)
	parent = make([]int, n)
	for i := range dist {
		dist[i], parent[i] = -1, -1
	}
	dist[v] = 0

	// Dijkstra's algorithm
	Q := emptySpQueue(dist)
	heap.Push(Q, v)
	for Q.Len() > 0 {
		v = heap.Pop(Q).(int)
		g.Visit(v, func(w int, d int64) (skip bool) {
			if d < 0 {
				return
			}
			alt := dist[v] + d
			switch {
			case dist[w] == -1:
				dist[w], parent[w] = alt, v
				heap.Push(Q, w)
			case alt < dist[w]:
				dist[w], parent[w] = alt, v
				Q.Update(w)
			}
			return
		})
	}
	return
}

type spQueue struct {
	heap  []int // vertices in heap order
	index []int // index of each vertex in the heap
	dist  []int64
}

func emptySpQueue(dist []int64) *spQueue {
	return &spQueue{dist: dist, index: make([]int, len(dist))}
}

func (pq *spQueue) Len() int { return len(pq.heap) }

func (pq *spQueue) Less(i, j int) bool {
	return pq.dist[pq.heap[i]] < pq.dist[pq.heap[j]]
}

func (pq *spQueue) Swap(i, j int) {
	pq.heap[i], pq.heap[j] = pq.heap[j], pq.heap[i]
	pq.index[pq.heap[i]] = i
	pq.index[pq.heap[j]] = j
}

func (pq *spQueue) Push(x interface{}) {
	n := len(pq.heap)
	v := x.(int)
	pq.heap = append(pq.heap, v)
	pq.index[v] = n
}

func (pq *spQueue) Pop() interface{} {
	n := len(pq.heap) - 1
	v := pq.heap[n]
	pq.heap = pq.heap[:n]
	return v
}

func (pq *spQueue) Update(v int) {
	heap.Fix(pq, pq.index[v])
}
