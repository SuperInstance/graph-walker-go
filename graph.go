package graph

import (
	"container/heap"
	"sort"
)

// Graph represents a directed weighted graph
type Graph struct {
	adj  map[int][]Edge
	nodes map[int]bool
	directed bool
}

// Edge represents a weighted edge
type Edge struct {
	To     int
	Weight float64
}

// New creates a new graph
func New(directed bool) *Graph {
	return &Graph{
		adj:      make(map[int][]Edge),
		nodes:    make(map[int]bool),
		directed: directed,
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(node int) {
	g.nodes[node] = true
}

// AddEdge adds an edge with given weight
func (g *Graph) AddEdge(from, to int, weight float64) {
	g.AddNode(from)
	g.AddNode(to)
	g.adj[from] = append(g.adj[from], Edge{To: to, Weight: weight})
	if !g.directed {
		g.adj[to] = append(g.adj[to], Edge{To: from, Weight: weight})
	}
}

// AddUnweightedEdge adds an edge with weight 1
func (g *Graph) AddUnweightedEdge(from, to int) {
	g.AddEdge(from, to, 1.0)
}

// HasNode checks if a node exists
func (g *Graph) HasNode(node int) bool {
	return g.nodes[node]
}

// Neighbors returns the neighbors of a node
func (g *Graph) Neighbors(node int) []Edge {
	return g.adj[node]
}

// Nodes returns all nodes
func (g *Graph) Nodes() []int {
	nodes := make([]int, 0, len(g.nodes))
	for n := range g.nodes {
		nodes = append(nodes, n)
	}
	sort.Ints(nodes)
	return nodes
}

// NodeCount returns the number of nodes
func (g *Graph) NodeCount() int {
	return len(g.nodes)
}

// EdgeCount returns the number of edges
func (g *Graph) EdgeCount() int {
	count := 0
	for _, edges := range g.adj {
		count += len(edges)
	}
	if !g.directed {
		return count / 2
	}
	return count
}

// BFS traverses the graph breadth-first from start, calling visit for each node
func (g *Graph) BFS(start int, visit func(node int)) {
	visited := make(map[int]bool)
	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		visit(node)
		for _, edge := range g.adj[node] {
			if !visited[edge.To] {
				visited[edge.To] = true
				queue = append(queue, edge.To)
			}
		}
	}
}

// DFS traverses the graph depth-first from start, calling visit for each node
func (g *Graph) DFS(start int, visit func(node int)) {
	visited := make(map[int]bool)
	g.dfsHelper(start, visited, visit)
}

func (g *Graph) dfsHelper(node int, visited map[int]bool, visit func(int)) {
	visited[node] = true
	visit(node)
	for _, edge := range g.adj[node] {
		if !visited[edge.To] {
			g.dfsHelper(edge.To, visited, visit)
		}
	}
}

// Dijkstra finds shortest paths from source to all reachable nodes
func (g *Graph) Dijkstra(source int) map[int]float64 {
	dist := make(map[int]float64)
	for _, n := range g.Nodes() {
		dist[n] = float64(1<<63 - 1) // infinity
	}
	dist[source] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &item{node: source, dist: 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*item)
		if cur.dist > dist[cur.node] {
			continue
		}
		for _, edge := range g.adj[cur.node] {
			newDist := dist[cur.node] + edge.Weight
			if newDist < dist[edge.To] {
				dist[edge.To] = newDist
				heap.Push(pq, &item{node: edge.To, dist: newDist})
			}
		}
	}
	return dist
}

// AStar finds shortest path from start to goal using heuristic
func (g *Graph) AStar(start, goal int, heuristic func(int) float64) (float64, bool) {
	dist := make(map[int]float64)
	for _, n := range g.Nodes() {
		dist[n] = float64(1<<63 - 1)
	}
	dist[start] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &item{node: start, dist: heuristic(start)})

	cameFrom := make(map[int]int)

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*item)
		if cur.node == goal {
			return dist[goal], true
		}
		for _, edge := range g.adj[cur.node] {
			newDist := dist[cur.node] + edge.Weight
			if newDist < dist[edge.To] {
				dist[edge.To] = newDist
				cameFrom[edge.To] = cur.node
				heap.Push(pq, &item{node: edge.To, dist: newDist + heuristic(edge.To)})
			}
		}
	}
	return 0, false
}

// SCC finds strongly connected components using Tarjan's algorithm
func (g *Graph) SCC() [][]int {
	index := 0
	stack := []int{}
	onStack := make(map[int]bool)
	indices := make(map[int]int)
	lowlink := make(map[int]int)
	var result [][]int

	var strongConnect func(v int)
	strongConnect = func(v int) {
		indices[v] = index
		lowlink[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		for _, edge := range g.adj[v] {
			w := edge.To
			if _, ok := indices[w]; !ok {
				strongConnect(w)
				if lowlink[w] < lowlink[v] {
					lowlink[v] = lowlink[w]
				}
			} else if onStack[w] {
				if indices[w] < lowlink[v] {
					lowlink[v] = indices[w]
				}
			}
		}

		if lowlink[v] == indices[v] {
			var component []int
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				component = append(component, w)
				if w == v {
					break
				}
			}
			result = append(result, component)
		}
	}

	for _, v := range g.Nodes() {
		if _, ok := indices[v]; !ok {
			strongConnect(v)
		}
	}
	return result
}

// MST finds minimum spanning tree using Kruskal's algorithm
func (g *Graph) MST() []Edge {
	if g.directed {
		return nil // MST only for undirected
	}
	
	// Collect all unique edges
	type edge struct{ from, to int; weight float64 }
	var edges []edge
	seen := make(map[[2]int]bool)
	for from, adj := range g.adj {
		for _, e := range adj {
			key := [2]int{min(from, e.To), max(from, e.To)}
			if !seen[key] {
				seen[key] = true
				edges = append(edges, edge{from, e.To, e.Weight})
			}
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].weight < edges[j].weight
	})

	parent := make(map[int]int)
	rank := make(map[int]int)
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(x, y int) bool {
		px, py := find(x), find(y)
		if px == py { return false }
		if rank[px] < rank[py] { px, py = py, px }
		parent[py] = px
		if rank[px] == rank[py] { rank[px]++ }
		return true
	}

	for _, n := range g.Nodes() {
		parent[n] = n
	}

	var mst []Edge
	for _, e := range edges {
		if union(e.from, e.to) {
			mst = append(mst, Edge{To: e.to, Weight: e.weight})
			if len(mst) == g.NodeCount()-1 {
				break
			}
		}
	}
	return mst
}

// GraphColoring assigns colors using greedy algorithm, returns map[node]color
func (g *Graph) GraphColoring() map[int]int {
	colors := make(map[int]int)
	for _, node := range g.Nodes() {
		used := make(map[int]bool)
		for _, edge := range g.adj[node] {
			if c, ok := colors[edge.To]; ok {
				used[c] = true
			}
		}
		color := 0
		for used[color] { color++ }
		colors[node] = color
	}
	return colors
}

// Priority queue implementation for Dijkstra/A*
type item struct {
	node int
	dist float64
	index int
}

type priorityQueue []*item

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[:n-1]
	return item
}
