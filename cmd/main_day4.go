package main

import (
    "fmt"
    "time"
    "math/rand"
    "myproject/graph"
    "myproject/dist"
    "myproject/mapreduce"
)

var nodes []*RaftNode

func main() {
	rand.Seed(time.Now().UnixNano())
	numNodes := 5

	// Инициализация узлов
	for i := 0; i < numNodes; i++ {
		nodes = append(nodes, NewRaftNode(i))
	}

	// Симуляция команды клиента, которая отправляется только лидеру
	clientCommand := "commandX"
	nodes[0].clientCommands <- clientCommand
	time.Sleep(3 * time.Second) // Даем время для выполнения команд

	// Имитируем падение лидера
	nodes[1].simulateLeaderCrash()

	time.Sleep(30 * time.Second) // Ждем, чтобы увидеть результаты работы
}