package graph

// Очередь
type Queue struct {
	data []int
}

type Neighbor struct {
	To     int
	Weight int
}

func (q *Queue) Enqueue(x int) {
	q.data = append(q.data, x)
}

func (q *Queue) Dequeue() (int, bool) {
	if len(q.data) == 0 {
		return 0, false
	}
	x := q.data[0]
	q.data = q.data[1:]
	return x, true
}

func (q *Queue) IsEmpty() bool {
	return len(q.data) == 0
}

// Граф
type Graph struct {
	adj map[int][]int
}

func NewGraph() *Graph {
	return &Graph{adj: make(map[int][]int)}
}

func (g *Graph) AddEdge(u, v int) {
	if _, ok := g.adj[u]; !ok {
		g.adj[u] = []int{}
	}
	if _, ok := g.adj[v]; !ok {
		g.adj[v] = []int{}
	}

	if !HasEdge(g, u, v) {
		g.adj[u] = append(g.adj[u], v)
		g.adj[v] = append(g.adj[v], u)
	}
}

func HasEdge(g *Graph, u, v int) bool {
	for _, neighbor := range g.adj[u] {
		if neighbor == v {
			return true
		}
	}
	return false
}


func (g *Graph) GetAllEdges() []Edge {
	var edges []Edge
	for _, node := range g.Nodes {
			for _, neighbor := range node.Neighbors {
					edges = append(edges, Edge{From: node.ID, To: neighbor.To, Weight: neighbor.Weight})
			}
	}
	return edges
}
