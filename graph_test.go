package graph

import (
	"math"
	"testing"
)

func TestAddNode(t *testing.T) {
	g := New(true)
	g.AddNode(1)
	g.AddNode(2)
	if g.NodeCount() != 2 {
		t.Errorf("expected 2 nodes, got %d", g.NodeCount())
	}
	if !g.HasNode(1) || !g.HasNode(2) {
		t.Error("nodes should exist")
	}
}

func TestAddEdge(t *testing.T) {
	g := New(true)
	g.AddUnweightedEdge(1, 2)
	if g.EdgeCount() != 1 {
		t.Errorf("expected 1 edge, got %d", g.EdgeCount())
	}
}

func TestUndirectedEdgeCount(t *testing.T) {
	g := New(false)
	g.AddUnweightedEdge(1, 2)
	if g.EdgeCount() != 1 {
		t.Errorf("expected 1 edge, got %d", g.EdgeCount())
	}
}

func TestBFS(t *testing.T) {
	g := New(true)
	g.AddUnweightedEdge(1, 2)
	g.AddUnweightedEdge(1, 3)
	g.AddUnweightedEdge(2, 4)
	g.AddUnweightedEdge(3, 4)

	var visited []int
	g.BFS(1, func(n int) { visited = append(visited, n) })

	if len(visited) != 4 {
		t.Errorf("BFS should visit 4 nodes, got %d", len(visited))
	}
	if visited[0] != 1 {
		t.Errorf("BFS should start at 1, got %d", visited[0])
	}
}

func TestDFS(t *testing.T) {
	g := New(true)
	g.AddUnweightedEdge(1, 2)
	g.AddUnweightedEdge(2, 3)
	g.AddUnweightedEdge(3, 4)

	var visited []int
	g.DFS(1, func(n int) { visited = append(visited, n) })

	if len(visited) != 4 {
		t.Errorf("DFS should visit 4 nodes, got %d", len(visited))
	}
	if visited[0] != 1 {
		t.Errorf("DFS should start at 1, got %d", visited[0])
	}
}

func TestDijkstra(t *testing.T) {
	g := New(true)
	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 3, 2)
	g.AddEdge(1, 3, 10)

	dist := g.Dijkstra(1)
	if dist[2] != 1 {
		t.Errorf("dist[2] should be 1, got %f", dist[2])
	}
	if dist[3] != 3 { // 1->2->3 = 3, not 10
		t.Errorf("dist[3] should be 3, got %f", dist[3])
	}
}

func TestAStar(t *testing.T) {
	g := New(true)
	g.AddEdge(0, 1, 1)
	g.AddEdge(1, 2, 1)
	g.AddEdge(0, 2, 10)

	// Manhattan heuristic (perfect for this grid-like graph)
	h := func(n int) float64 { return math.Abs(float64(2 - n)) }

	cost, found := g.AStar(0, 2, h)
	if !found {
		t.Error("A* should find path")
	}
	if cost != 2 { // 0->1->2 = 2
		t.Errorf("A* cost should be 2, got %f", cost)
	}
}

func TestSCC(t *testing.T) {
	g := New(true)
	g.AddUnweightedEdge(1, 2)
	g.AddUnweightedEdge(2, 3)
	g.AddUnweightedEdge(3, 1)
	g.AddUnweightedEdge(3, 4)

	sccs := g.SCC()
	if len(sccs) < 1 {
		t.Error("should find at least 1 SCC")
	}
	// Find the SCC containing {1,2,3}
	var found123 bool
	for _, scc := range sccs {
		if len(scc) == 3 {
			m := map[int]bool{}
			for _, n := range scc { m[n] = true }
			if m[1] && m[2] && m[3] { found123 = true }
		}
	}
	if !found123 {
		t.Error("should find SCC {1,2,3}")
	}
}

func TestMST(t *testing.T) {
	g := New(false)
	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 3, 2)
	g.AddEdge(1, 3, 10)

	mst := g.MST()
	if len(mst) != 2 {
		t.Errorf("MST should have 2 edges, got %d", len(mst))
	}
	totalWeight := 0.0
	for _, e := range mst {
		totalWeight += e.Weight
	}
	if totalWeight != 3 { // edges 1-2 (1) and 2-3 (2)
		t.Errorf("MST weight should be 3, got %f", totalWeight)
	}
}

func TestGraphColoring(t *testing.T) {
	g := New(false)
	// Triangle: needs 3 colors
	g.AddUnweightedEdge(1, 2)
	g.AddUnweightedEdge(2, 3)
	g.AddUnweightedEdge(3, 1)

	colors := g.GraphColoring()
	// No two adjacent nodes should share a color
	if colors[1] == colors[2] || colors[2] == colors[3] || colors[1] == colors[3] {
		t.Error("adjacent nodes should have different colors")
	}
}

func TestGraphColoringBipartite(t *testing.T) {
	g := New(false)
	g.AddUnweightedEdge(1, 2)
	g.AddUnweightedEdge(2, 3)
	g.AddUnweightedEdge(3, 4)
	
	colors := g.GraphColoring()
	maxColor := 0
	for _, c := range colors {
		if c > maxColor { maxColor = c }
	}
	if maxColor > 1 {
		t.Errorf("bipartite graph should need only 2 colors, used %d", maxColor+1)
	}
}

func TestDijkstraUnreachable(t *testing.T) {
	g := New(true)
	g.AddNode(1)
	g.AddNode(2)
	dist := g.Dijkstra(1)
	if dist[2] != float64(1<<63-1) {
		t.Errorf("unreachable node should have infinite distance")
	}
}

func TestSCCSingle(t *testing.T) {
	g := New(true)
	g.AddNode(1)
	sccs := g.SCC()
	if len(sccs) != 1 {
		t.Errorf("single node should be 1 SCC, got %d", len(sccs))
	}
}

func TestMSTSingleNode(t *testing.T) {
	g := New(false)
	g.AddNode(1)
	mst := g.MST()
	if len(mst) != 0 {
		t.Errorf("single node MST should have 0 edges, got %d", len(mst))
	}
}
