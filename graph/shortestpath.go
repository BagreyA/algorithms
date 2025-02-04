package graph

import (
	"container/heap"
	"math"
)

// вершина в приоритетной очереди
type Item struct {
	vertex, dist int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }

func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Item))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// Dijkstra
func Dijkstra(g *Graph, start int) ([]int, []int) {
	n := len(g.adj)
	dist := make([]int, n)
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = int(^uint(0) >> 1)
		parent[i] = -1
	}
	dist[start] = 0

	pq := &PriorityQueue{}
	heap.Push(pq, &Item{vertex: start, dist: 0})

	for pq.Len() > 0 {
		u := heap.Pop(pq).(*Item)
		for _, neighbor := range g.adj[u.vertex] {
			if dist[u.vertex]+neighbor.weight < dist[neighbor.to] {
				dist[neighbor.to] = dist[u.vertex] + neighbor.weight
				parent[neighbor.to] = u.vertex
				heap.Push(pq, &Item{vertex: neighbor.to, dist: dist[neighbor.to]})
			}
		}
	}
	return dist, parent
}

// BellmanFord
func BellmanFord(g *Graph, start int) ([]int, bool) {
	n := len(g.adj)
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = math.MaxInt32
	}
	dist[start] = 0

	for i := 0; i < n-1; i++ {
		for _, edge := range g.GetAllEdges() {
			u, v, w := edge.u, edge.v, edge.w
			if dist[u]+w < dist[v] {
				dist[v] = dist[u] + w
			}
		}
	}

	for _, edge := range g.GetAllEdges() {
		u, v, w := edge.u, edge.v, edge.w
		if dist[u]+w < dist[v] {
			return dist, true
		}
	}

	return dist, false
}