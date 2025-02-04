package main

import (
	"fmt"
	"myproject/graph"
	"myproject/dist"
	"myproject/mapreduce"
)

func main() {
	ds := graph.NewDisjointSet(6)

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

	g := graph.NewGraph()
	g.AddEdge(0, 1, 4)
	g.AddEdge(1, 2, 2)
	g.AddEdge(2, 3, 1)
	g.AddEdge(3, 4, 3)
	g.AddEdge(0, 3, 5)

	edges := g.GetAllEdges()

	mst, totalWeight := graph.MST(len(g.adj), edges)
	fmt.Println("MST:")
	for _, edge := range mst {
		fmt.Printf("%d -- %d (вес: %d)\n", edge.u, edge.v, edge.w)
	}
	fmt.Printf("Общий вес MST: %d\n", totalWeight)

	fmt.Println("\nДейкстра:")
	start := 0
	dist, _ := graph.Dijkstra(g, start)
	fmt.Println("Расстояния от начальной вершины:", dist)

	bfDist, hasNegativeCycle := graph.BellmanFord(g, start)
	if hasNegativeCycle {
		fmt.Println("Обнаружен отрицательный цикл")
	} else {
		fmt.Println("Расстояния Беллмана-Форда от начальной вершины:", bfDist)
	}
}
