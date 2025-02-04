package main

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
)

// DisjointSet
type DisjointSet struct {
	parent []int
	rank   []int
}

func NewDisjointSet(n int) *DisjointSet {
	ds := &DisjointSet{
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		ds.parent[i] = i
		ds.rank[i] = 0
	}
	return ds
}

// возвращает корень элемента x с сжатием пути
func (ds *DisjointSet) Find(x int) int {
	if ds.parent[x] != x {
		ds.parent[x] = ds.Find(ds.parent[x])
	}
	return ds.parent[x]
}

// объединяет два множества
func (ds *DisjointSet) Union(x, y int) bool {
	rx := ds.Find(x)
	ry := ds.Find(y)
	if rx == ry {
		return false
	}
	if ds.rank[rx] < ds.rank[ry] {
		ds.parent[rx] = ry
	} else if ds.rank[rx] > ds.rank[ry] {
		ds.parent[ry] = rx
	} else {
		ds.parent[ry] = rx
		ds.rank[rx]++
	}
	return true
}

type Edge struct {
	u, v, w int
}

// Graph
type Graph struct {
	adj   map[int][]struct{ to, weight int }
	edges []Edge
}

func NewGraph() *Graph {
	return &Graph{
		adj:   make(map[int][]struct{ to, weight int }),
		edges: []Edge{},
	}
}

// добавляет
func (g *Graph) AddEdge(u, v, w int) {
	g.adj[u] = append(g.adj[u], struct{ to, weight int }{v, w})
	g.adj[v] = append(g.adj[v], struct{ to, weight int }{u, w})
	g.edges = append(g.edges, Edge{u, v, w})
}

// возвращает все
func (g *Graph) GetAllEdges() []Edge {
	var edges []Edge
	for u, neighbors := range g.adj {
		for _, neighbor := range neighbors {
			if u < neighbor.to {
				edges = append(edges, Edge{u, neighbor.to, neighbor.weight})
			}
		}
	}
	return edges
}

// MST
func MST(n int, edges []Edge) (mst []Edge, totalWeight int) {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})

	ds := NewDisjointSet(n)
	for _, edge := range edges {
		if ds.Find(edge.u) != ds.Find(edge.v) {
			ds.Union(edge.u, edge.v)
			mst = append(mst, edge)
			totalWeight += edge.w
		}
	}
	return mst, totalWeight
}

// вершина в оп
type Item struct {
	vertex, dist int
}

type PriorityQueue []*Item

// возвращает длину
func (pq PriorityQueue) Len() int { return len(pq) }

// true, если элемент i меньший приоритет
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }

// меняет местами элементы
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Item))
}

// удаляет и возвращает элемент
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

func main() {
	ds := NewDisjointSet(6)

	ds.Union(0, 1)
	ds.Union(1, 2)
	ds.Union(3, 4)

	fmt.Println("Найдены результаты:")
	for i := 0; i < 6; i++ {
		fmt.Printf("Найдено(%d): %d\n", i, ds.Find(i))
	}

	fmt.Println("Компоненты: 0, 1, 2 - одна, 3, 4 - другая, 5 - сам по себе")

	ds.Union(2, 3)

	fmt.Printf("родитель: %v\n", ds.parent)
	fmt.Printf("ранг: %v\n", ds.rank)

	g := NewGraph()
	g.AddEdge(0, 1, 4)
	g.AddEdge(1, 2, 2)
	g.AddEdge(2, 3, 1)
	g.AddEdge(3, 4, 3)
	g.AddEdge(0, 3, 5)

	edges := g.GetAllEdges()

	mst, totalWeight := MST(len(g.adj), edges)
	fmt.Println("MST:")
	for _, edge := range mst {
		fmt.Printf("%d -- %d (вес: %d)\n", edge.u, edge.v, edge.w)
	}
	fmt.Printf("Общий вес MST: %d\n", totalWeight)

	fmt.Println("\nДейкстра:")
	start := 0
	dist, _ := Dijkstra(g, start)
	fmt.Println("Расстояния от начальной вершины:", dist)

	bfDist, hasNegativeCycle := BellmanFord(g, start)
	if hasNegativeCycle {
		fmt.Println("Обнаружен отрицательный цикл")
	} else {
		fmt.Println("Расстояния Беллмана-Форда от начальной вершины:", bfDist)
	}
}
