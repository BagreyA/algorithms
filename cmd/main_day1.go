package main

import (
	"fmt"
	"myproject/graph"
	"myproject/dist"
	"myproject/mapreduce"
	"math/rand"
	"time"
)

func RandomGraph(numVertices, numEdges int) *graph.Graph {
	g := graph.NewGraph()
	rand.Seed(time.Now().UnixNano())

	// Подсчет добавленных рёбер
	edgesAdded := 0
	for edgesAdded < numEdges {
		u := rand.Intn(numVertices)
		v := rand.Intn(numVertices)

		// Проверка на дубликаты и наличие самоциклов
		if u != v && !graph.HasEdge(g, u, v) {
			g.AddEdge(u, v)
			edgesAdded++
		}
	}

	return g
}

func main() {
	numVertices := 7 // Количество вершин
	numEdges := 5    // Количество рёбер

	g := RandomGraph(numVertices, numEdges)

	fmt.Println("Структура графа:", g.adj)

	fmt.Println("Есть ли связь между 0 и 1?", graph.HasEdge(g, 0, 1))
	fmt.Println("Есть ли связь между 0 и 3?", graph.HasEdge(g, 0, 3))

	fmt.Println("Порядок обхода BFS:", graph.BFS(g, 0))
	fmt.Println("Порядок обхода DFS:", graph.DFS(g, 0))

	count, comp := graph.ConnectedComponents(g)
	fmt.Println("Количество компонентов:", count)
	for v, c := range comp {
		fmt.Printf("Вершина %d принадлежит компоненте %d\n", v, c)
	}
}
