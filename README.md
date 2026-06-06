# graph-walker-go

Graph traversal and algorithm library for Go with visitor pattern support.

## Features

- **BFS** — Breadth-first traversal with visitor callback
- **DFS** — Depth-first traversal with visitor callback
- **Dijkstra** — Shortest paths from a single source
- **A\*** — Heuristic-guided shortest path search
- **Tarjan's SCC** — Strongly connected components
- **Kruskal's MST** — Minimum spanning tree for undirected graphs
- **Graph Coloring** — Greedy vertex coloring

## Installation

```bash
go get github.com/SuperInstance/graph-walker-go
```

## Usage

```go
package main

import (
    "fmt"
    graph "github.com/SuperInstance/graph-walker-go"
)

func main() {
    g := graph.New(true) // directed graph

    // Add weighted edges
    g.AddEdge(0, 1, 1.0)
    g.AddEdge(1, 2, 2.0)
    g.AddEdge(0, 2, 10.0)

    // Dijkstra shortest paths
    dist := g.Dijkstra(0)
    fmt.Printf("Distance to node 2: %f\n", dist[2]) // 3.0 (via node 1)

    // A* with heuristic
    heuristic := func(n int) float64 { return math.Abs(float64(2 - n)) }
    cost, found := g.AStar(0, 2, heuristic)
    fmt.Printf("A* cost: %f, found: %v\n", cost, found)

    // BFS traversal
    g.BFS(0, func(n int) { fmt.Println("Visited:", n) })

    // Graph coloring
    colors := g.GraphColoring()
    fmt.Printf("Node 0: color %d, Node 1: color %d\n", colors[0], colors[1])
}
```

## Testing

```bash
go test -v ./...    # 14 tests
```

## License

MIT
