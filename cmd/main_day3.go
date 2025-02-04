package main

import (
	"fmt"
	"myproject/graph"
	"myproject/dist"
	"myproject/mapreduce"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Создание узлов
	numNodes := 5
	nodes := make([]dist.Node, numNodes)
	for i := range nodes {
		nodes[i] = *dist.NewNode(i, (i+1)%numNodes)
		nodes[i].StartListening()
	}

	// Демонстрация выбора лидера
	startingNode := rand.Intn(numNodes)
	nodes[startingNode].StartElection(nodes)

	// Завершаем работу через некоторое время
	time.Sleep(10 * time.Second)
}
