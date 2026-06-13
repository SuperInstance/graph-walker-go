# Graph Walker Go

A **comprehensive graph algorithm library in Go** providing traversal (BFS, DFS), shortest paths (Dijkstra, A*), connectivity (Tarjan's SCC), spanning trees (Kruskal's MST), and coloring (greedy) — all on a unified directed/undirected weighted graph type.

## Why It Matters

Most graph libraries implement one or two algorithms. This library provides seven production-quality algorithms on a single `Graph` type, enabling composition: find SCCs, then run Dijkstra within each component; compute an MST, then BFS for diameter; use A* for pathfinding with graph coloring for register allocation. The Go implementation leverages the `container/heap` interface for type-safe priority queues, `sort.Slice` for edge sorting (Kruskal's), and Go's map-based adjacency list for sparse-graph efficiency. It's the reference implementation for the SuperInstance graph algorithm suite, ported to Go for environments where Go's garbage collection and concurrency primitives (goroutines, channels) are preferred over Rust's ownership model.

## How It Works

**Graph representation**: Adjacency list using `map[int][]Edge`, with a separate `map[int]bool` for node existence. This supports sparse graphs efficiently: O(1) node lookup, O(degree(v)) neighbor iteration. For undirected graphs, edges are stored in both adjacency lists.

**Algorithm inventory**:

| Algorithm | Time | Space | Key Insight |
|-----------|------|-------|-------------|
| BFS | O(V+E) | O(V) | Queue (FIFO) |
| DFS | O(V+E) | O(V) | Stack (LIFO) / recursion |
| Dijkstra | O((V+E) log V) | O(V) | Min-heap priority queue |
| A* | O(E) best | O(V) | Heuristic-guided Dijkstra |
| SCC (Tarjan) | O(V+E) | O(V) | Single DFS + lowlink |
| MST (Kruskal) | O(E log E) | O(V) | Sort edges + union-find |
| Coloring | O(V+E) | O(V) | Greedy: first available color |

**Dijkstra implementation**: Uses Go's `container/heap` with a custom `priorityQueue` type implementing `heap.Interface`. Nodes are pushed as `(node, dist)` pairs; stale entries are skipped via the standard `if cur.dist > dist[cur.node]` check.

**A* search**: Extends Dijkstra with a heuristic function `h(n)` estimating distance to goal. The priority becomes `f(n) = g(n) + h(n)` where g is actual cost and h is the heuristic. With an admissible heuristic (never overestimates), A* is guaranteed optimal.

**Tarjan's SCC**: Single-pass DFS maintaining `index` (discovery time) and `lowlink` (lowest index reachable). When `lowlink[v] == index[v]`, v is the root of an SCC. Nodes on the stack above v form the component. This is O(V+E) — optimal.

**Kruskal's MST**: Sort all edges by weight ascending. Use union-find (disjoint set with path compression and union by rank) to detect cycles. Each edge connecting two different sets is added to the MST. Near-constant amortized per union-find operation (inverse Ackermann).

**Graph coloring**: Process nodes in order. For each, check colors of neighbors and pick the smallest unused color. This greedy approach uses at most Δ+1 colors (Δ = max degree) but may not be optimal (chromatic number can be much less).

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/SuperInstance/graph-walker-go"
)

func main() {
    g := graph.New(false) // undirected
    g.AddEdge(0, 1, 4.0)
    g.AddEdge(0, 2, 1.0)
    g.AddEdge(1, 2, 2.0)
    g.AddEdge(1, 3, 5.0)
    g.AddEdge(2, 3, 3.0)

    // Shortest paths
    dist := g.Dijkstra(0)
    fmt.Println("Distances:", dist) // map[0:0 1:3 2:1 3:4]

    // MST
    mst := g.MST()
    fmt.Println("MST edges:", len(mst)) // 3 edges

    // Traversal
    g.BFS(0, func(n int) { fmt.Print(n, " ") })

    // A* with zero heuristic = Dijkstra
    cost, ok := g.AStar(0, 3, func(n int) float64 { return 0 })
    fmt.Printf("\nA* cost: %f, found: %v\n", cost, ok)
}
```

## API

| Method | Description |
|--------|-------------|
| `New(directed bool) *Graph` | Create a graph |
| `.AddNode(n)` / `.AddEdge(from, to, w)` | Add vertices/edges |
| `.BFS(start, visit)` / `.DFS(start, visit)` | Traverse with callback |
| `.Dijkstra(source) map[int]float64` | Single-source shortest paths |
| `.AStar(start, goal, h) (float64, bool)` | Heuristic-guided search |
| `.SCC() [][]int` | Strongly connected components (Tarjan) |
| `.MST() []Edge` | Minimum spanning tree (Kruskal) |
| `.GraphColoring() map[int]int` | Greedy node coloring |

## Architecture Notes

Graph Walker Go is the Go-language reference for SuperInstance's graph suite, paralleling the Rust crates (graph-bfs, graph-dfs, graph-dijkstra, graph-bellman-ford). The Go implementation prioritizes composability — all algorithms share one `Graph` type, enabling multi-step analyses (SCC → Dijkstra per component → MST). In **γ + η = C**, having one library instead of seven reduces γ (integration cost). See [Architecture](https://github.com/SuperInstance/SuperInstance/blob/main/ARCHITECTURE.md).

## References

- Tarjan, R. "Depth-First Search and Linear Graph Algorithms," SIAM J. Comput. (1972).
- Skiena, S. *The Algorithm Design Manual*, 3rd ed., Springer (2020).
- Sedgewick, R. & Wayne, K. *Algorithms*, 4th ed., Addison-Wesley (2011).

## License

MIT
