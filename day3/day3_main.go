package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Message struct {
	kind   string
	fromID int
	weight int
	data   interface{}
}

type Node struct {
	id         int
	weight     int
	alive      bool
	leaderID   int
	localCount int
	inbox      chan Message
	mu         sync.Mutex
}

func NewNode(id, weight int) *Node {
	return &Node{
		id:         id,
		weight:     weight,
		alive:      true,
		leaderID:   -1,
		localCount: rand.Intn(100),         // Пример данных
		inbox:      make(chan Message, 10), // Используем буферизованный канал
	}
}

func (n *Node) StartListening() {
	go func() {
		for msg := range n.inbox {
			n.ProcessMessage(msg)
		}
	}()
}

func (n *Node) StartElection(nodes []Node) {
	fmt.Printf("Узел %d начинает выборы\n", n.id)

	// Отправляем сообщение ELECTION всем другим узлам
	for _, node := range nodes {
		if node.id != n.id && node.alive {
			msg := Message{kind: "ELECTION", fromID: n.id, weight: n.weight}
			node.inbox <- msg
		}
	}

	// Таймер на ожидание ответов
	timeout := time.After(2 * time.Second)

	go func() {
		<-timeout
		n.mu.Lock()
		if n.leaderID == -1 { // Если лидер не установлен
			fmt.Printf("Узел %d не услышал о лидере, объявляет себя лидером\n", n.id)
			n.leaderID = n.id
			for _, node := range nodes {
				if node.id != n.id && node.alive {
					node.inbox <- Message{kind: "COORDINATOR", fromID: n.leaderID}
				}
			}
		}
		n.mu.Unlock()
	}()
}

func (n *Node) ProcessMessage(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()

	switch msg.kind {
	case "ELECTION":
		fmt.Printf("Узел %d получил ELECTION от %d (вес %d)\n", n.id, msg.fromID, msg.weight)

		if msg.fromID == n.id {
			return
		}

		if msg.weight < n.weight || (msg.weight == n.weight && msg.fromID < n.id) {
			fmt.Printf("Узел %d отвечает OK на выборы узла %d\n", n.id, msg.fromID)
			n.inbox <- Message{kind: "OK", fromID: n.id, weight: n.weight}
		} else {
			fmt.Printf("Узел %d проигрывает выборы узлу %d (вес %d > вес %d)\n", n.id, msg.fromID, n.weight, msg.weight)
		}

		if n.leaderID == -1 && msg.weight > n.weight {
			n.leaderID = msg.fromID
			fmt.Printf("Узел %d объявляет узел %d лидером\n", n.id, n.leaderID)
		}

	case "OK":
		fmt.Printf("Узел %d получил OK от %d (вес %d)\n", n.id, msg.fromID, msg.weight)

	case "COORDINATOR":
		n.leaderID = msg.fromID
		fmt.Printf("Узел %d получил COORDINATOR от %d, устанавливает лидера на %d\n", n.id, msg.fromID, msg.fromID)

	case "COLLECT":
		if n.alive {
			fmt.Printf("Узел %d получил COLLECT от лидера %d, отвечает своим локальным счетом %d\n", n.id, msg.fromID, n.localCount)
			n.inbox <- Message{kind: "COLLECT_REPLY", fromID: n.id, data: n.localCount}
		} else {
			fmt.Printf("Узел %d не может ответить, так как он не работает\n", n.id)
		}
	}
}

func (n *Node) StartGlobalCollection(nodes []Node) {
	if n.leaderID == n.id {
		fmt.Printf("Лидер %d начинает глобальный сбор данных\n", n.id)
		sum := 0
		received := 0
		numNodes := 0

		// Отправляем запрос на сбор данных всем узлам
		for _, node := range nodes {
			if node.id != n.id && node.alive {
				numNodes++
				node.inbox <- Message{kind: "COLLECT", fromID: n.id}
			}
		}

		go func() {
			for {
				select {
				case msg := <-n.inbox:
					if msg.kind == "COLLECT_REPLY" {
						// Увеличиваем сумму только от узлов, которые не являются лидером
						data, ok := msg.data.(int)
						if ok {
							sum += data
						}
						received++
					}

					// Проверяем, получили ли мы ответы от всех узлов
					if received == numNodes {
						fmt.Printf("Лидер %d собрал сумму: %d\n", n.id, sum)
						return
					}

				case <-time.After(10 * time.Second):
					fmt.Println("Таймаут сбора данных, некоторые узлы не ответили.")
					return
				}
			}
		}()
	}
}

func PrintNodesInfo(nodes []Node) {
	fmt.Println("Перечень всех узлов и их весов:")
	for _, node := range nodes {
		fmt.Printf("Узел %d, Вес: %d, Данные: %d\n", node.id, node.weight, node.localCount)
	}
}

var nodes []Node

func main() {
	rand.Seed(time.Now().UnixNano())

	numNodes := 5
	nodes = make([]Node, numNodes)
	for i := range nodes {
		nodes[i] = *NewNode(i, (i+1)%numNodes)
		nodes[i].StartListening()
	}

	PrintNodesInfo(nodes)

	for i := range nodes {
		nodes[i].leaderID = -1
	}

	// Случайный узел начинает выборы
	startingNode := rand.Intn(numNodes)
	nodes[startingNode].StartElection(nodes)

	time.Sleep(3 * time.Second)

	// Поиск лидера
	for _, node := range nodes {
		if node.leaderID != -1 {
			node.StartGlobalCollection(nodes)
			break
		}
	}

	// Решаем, произойдет ли краш узла
	if rand.Float32() < 0.3 {
		crashedNode := rand.Intn(numNodes)
		nodes[crashedNode].alive = false
		fmt.Printf("Узел %d теперь не работает\n", nodes[crashedNode].id)
	}

	time.Sleep(2 * time.Second)
}
