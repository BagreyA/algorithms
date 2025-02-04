package day1

import (
	"fmt"
	"math/rand"
	"time"
)

// Стек
type Stack struct {
	data []int
}

func (s *Stack) Push(x int) {
	s.data = append(s.data, x)
}

func (s *Stack) Pop() (int, bool) {
	if len(s.data) == 0 {
		return 0, false
	}
	x := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return x, true
}

func (s *Stack) IsEmpty() bool {
	return len(s.data) == 0
}

// Очередь
type Queue struct {
	data []int
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

// BFS
func BFS(g *Graph, start int) []int {
	visited := make(map[int]bool)
	order := []int{}
	queue := Queue{}

	visited[start] = true
	queue.Enqueue(start)

	for !queue.IsEmpty() {
		u, _ := queue.Dequeue()
		order = append(order, u)

		for _, v := range g.adj[u] {
			if !visited[v] {
				visited[v] = true
				queue.Enqueue(v)
			}
		}
	}
	return order
}

// DFS
func (g *Graph) dfsUtil(v int, visited map[int]bool, order *[]int) {
	visited[v] = true
	*order = append(*order, v)

	for _, u := range g.adj[v] {
		if !visited[u] {
			g.dfsUtil(u, visited, order)
		}
	}
}

func DFS(g *Graph, start int) []int {
	visited := make(map[int]bool)
	order := []int{}
	g.dfsUtil(start, visited, &order)
	return order
}

func ConnectedComponents(g *Graph) (count int, comp map[int]int) {
	visited := make(map[int]bool)
	comp = make(map[int]int)
	count = 0

	for v := range g.adj {
		if !visited[v] {
			count++
			order := []int{}
			g.dfsUtil(v, visited, &order)
			for _, u := range order {
				comp[u] = count
			}
		}
	}
	return count, comp
}

func RandomGraph(numVertices, numEdges int) *Graph {
	g := NewGraph()
	rand.Seed(time.Now().UnixNano())

	// Подсчет добавленных рёбер
	edgesAdded := 0
	for edgesAdded < numEdges {
		u := rand.Intn(numVertices)
		v := rand.Intn(numVertices)

		// Проверка на дубликаты и наличие самоциклов
		if u != v && !HasEdge(g, u, v) {
			g.AddEdge(u, v)
			edgesAdded++
		}
	}

	return g
}

func main1() {
	numVertices := 7 // Количество вершин
	numEdges := 5    // Количество рёбер

	g := RandomGraph(numVertices, numEdges)

	fmt.Println("Структура графа:", g.adj)

	fmt.Println("Есть ли связь между 0 и 1?", HasEdge(g, 0, 1))
	fmt.Println("Есть ли связь между 0 и 3?", HasEdge(g, 0, 3))

	fmt.Println("Порядок обхода BFS:", BFS(g, 0))
	fmt.Println("Порядок обхода DFS:", DFS(g, 0))

	count, comp := ConnectedComponents(g)
	fmt.Println("Количество компонентов:", count)
	for v, c := range comp {
		fmt.Printf("Вершина %d принадлежит компоненте %d\n", v, c)
	}
}
