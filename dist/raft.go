package dist

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// LogEntry - структура записи лога
type LogEntry struct {
	Term    int
	Command string
}

// Message - структура для передачи сообщений между узлами
type Message struct {
	kind         string // "VoteRequest", "VoteResponse", "AppendEntries", "Heartbeat"
	term         int
	from, to     int
	entry        LogEntry
	leaderID     int
	leaderCommit int
}

// RaftNode - структура для представления узла в Raft
type RaftNode struct {
	id            int
	term          int
	state         string // "Follower", "Candidate", "Leader"
	log           []LogEntry
	votedFor      int
	leaderID      int
	commitIndex   int
	inbox         chan Message
	mu            sync.Mutex
	alive         bool
	votesReceived int // количество полученных голосов
	clientCommands chan string  // Канал для команд от клиента
}

// NewRaftNode создает новый узел Raft
func NewRaftNode(id int) *RaftNode {
	node := &RaftNode{
		id:            id,
		term:          0,
		state:         "Follower",
		votedFor:      -1,
		commitIndex:   0,
		inbox:         make(chan Message, 10),
		log:           []LogEntry{},
		alive:         true,
		votesReceived: 0,
		clientCommands: make(chan string, 1), // Инициализация канала для команд клиента
	}
	go node.run()
	return node
}

// run - основной цикл работы узла
func (rn *RaftNode) run() {
	timeout := time.Duration(rand.Intn(150)+150) * time.Millisecond
	heartbeatChan := time.After(timeout)

	for rn.alive {
		select {
		case msg := <-rn.inbox:
			rn.handleMessage(msg)
		case <-heartbeatChan:
			if rn.state == "Follower" {
				fmt.Printf("Узел %d: Временной лимит истек, становлюсь кандидатом\n", rn.id)
				rn.state = "Candidate"
				rn.term++
				rn.votedFor = rn.id
				rn.votesReceived = 1 // Голосуем за себя
				rn.requestVotes()
			}
		case command := <-rn.clientCommands: // Реагируем на команды клиента
			if rn.state == "Leader" {
				logEntry := LogEntry{Term: rn.term, Command: command}
				rn.log = append(rn.log, logEntry)
				rn.sendHeartbeat() // Отправляем heartbeat с новым логом
				fmt.Println("Лидер добавил команду в лог:", command)
			}
		}
	}
}

// handleMessage обрабатывает полученные сообщения
func (rn *RaftNode) handleMessage(msg Message) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	switch msg.kind {
	case "VoteRequest":
		fmt.Printf("Узел %d получил запрос на голосование от узла %d\n", rn.id, msg.from)
		if msg.term >= rn.term && (rn.votedFor == -1 || rn.votedFor == msg.from) {
			rn.votedFor = msg.from
			rn.term = msg.term
			rn.state = "Follower"
			rn.inbox <- Message{kind: "VoteResponse", term: rn.term, from: rn.id, to: msg.from}
		} else {
			rn.inbox <- Message{kind: "VoteResponse", term: rn.term, from: rn.id, to: msg.from}
		}
	case "VoteResponse":
		fmt.Printf("Узел %d получил ответ на голосование от узла %d\n", rn.id, msg.from)
		if rn.state == "Candidate" && msg.term == rn.term {
			// Увеличиваем количество полученных голосов
			rn.votesReceived++
			// Если мы получили достаточно голосов, станем лидером
			if rn.votesReceived > len(nodes)/2 {
				rn.state = "Leader"
				fmt.Printf("Узел %d стал ЛИДЕРОМ\n", rn.id)
				rn.sendHeartbeat()
			}
		} else {
			// Если узел не получил достаточное количество голосов или у него не хватает времени
			fmt.Printf("Узел %d не получил достаточное количество голосов\n", rn.id)
		}
	case "AppendEntries":
		fmt.Printf("Узел %d получил AppendEntries от узла %d\n", rn.id, msg.from)
		if msg.term >= rn.term {
			rn.term = msg.term
			rn.leaderID = msg.leaderID
			rn.log = append(rn.log, msg.entry)
			rn.commitIndex = msg.leaderCommit
			rn.inbox <- Message{kind: "AppendReply", term: rn.term, to: msg.leaderID, from: rn.id}
		}
	case "AppendReply":
		fmt.Printf("Узел %d получил AppendReply от узла %d\n", rn.id, msg.from)
	}
}

// requestVotes отправляет запросы на голосование другим узлам
func (rn *RaftNode) requestVotes() {
	fmt.Printf("Узел %d запрашивает голоса на термин %d\n", rn.id, rn.term)
	// Здесь необходимо добавить логику отправки сообщений запросов голосования
	for _, node := range nodes {
		if node.id != rn.id {
			rn.sendVoteRequest(node)
		}
	}
}

// sendVoteRequest отправляет запрос на голосование конкретному узлу
func (rn *RaftNode) sendVoteRequest(node *RaftNode) {
	// Имитация потери сообщения с вероятностью 30%
	if rand.Float32() < 0.1 {
		fmt.Printf("Узел %d: Сообщение от %d к %d было потеряно\n", rn.id, rn.id, node.id)
		return
	}
	node.inbox <- Message{kind: "VoteRequest", term: rn.term, from: rn.id}
}

// sendHeartbeat отправляет сообщения о сердцебиении всем фолловерам
func (rn *RaftNode) sendHeartbeat() {
	for _, node := range nodes {
		if node.id != rn.id {
			// Имитация потери сообщения с вероятностью 20%
			if rand.Float32() < 0.1 {
				fmt.Printf("Узел %d: Сообщение от лидера %d к %d было потеряно\n", rn.id, rn.id, node.id)
				continue
			}
			node.inbox <- Message{
				kind:        "AppendEntries",
				term:        rn.term,
				leaderID:    rn.id,
				leaderCommit: rn.commitIndex,
			}
		}
	}
}

// simulateLeaderCrash имитирует падение лидера
func (rn *RaftNode) simulateLeaderCrash() {
	// Падение лидера после того, как команда была записана, но еще не зафиксирована
	if rn.state == "Leader" && len(rn.log) > 0 {
		rn.alive = false
		fmt.Printf("Лидер %d упал. Лог был не зафиксирован и потерян.\n", rn.id)
	}
}

var node RaftNode
rn.sendVoteRequest(&node) // Передаем указатель