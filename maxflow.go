package graph

// MaxFlow computes a maximum flow from s to t in a graph
// with nonnegative edge capacities.
// The time complexity is O(|V|⋅|E|²), where |V| is the number of vertices
// and |E| the number of edges in the graph.
func MaxFlow(g Iterator, s, t int) (flow int64, graph Iterator) {
	// Edmonds-Karp's algorithm
	n := g.Order()
	residual := New(n)
	for v := 0; v < n; v++ {
		g.Visit(v, func(w int, c int64) (skip bool) {
			residual.AddCost(v, w, c)
			return
		})
	}
	prev := make([]int, n)
	for residualFlow(residual, s, t, prev) && flow < Max {
		pathFlow := Max
		for v := t; v != s; {
			u := prev[v]
			if c := residual.Cost(u, v); c < pathFlow {
				pathFlow = c
			}
			v = u
		}
		flow += pathFlow
		for v := t; v != s; {
			u := prev[v]
			residual.AddCost(u, v, residual.Cost(u, v)-pathFlow)
			residual.AddCost(v, u, residual.Cost(v, u)+pathFlow)
			v = u
		}
	}
	res := New(n)
	for v := 0; v < n; v++ {
		g.Visit(v, func(w int, c int64) (skip bool) {
			if flow := c - residual.Cost(v, w); flow > 0 {
				res.AddCost(v, w, flow)
			}
			return
		})
	}
	graph = Sort(res)
	return
}

func residualFlow(g Iterator, s, t int, prev []int) bool {
	visited := make([]bool, g.Order())
	prev[s], visited[s] = -1, true
	for queue := []int{s}; len(queue) > 0; {
		v := queue[0]
		queue = queue[1:]
		g.Visit(v, func(w int, c int64) (skip bool) {
			if !visited[w] && c > 0 {
				prev[w] = v
				visited[w] = true
				queue = append(queue, w)
			}
			return
		})
	}
	return visited[t]
}
