package main

import (
    "fmt"
    "strings"
    "sync"
    "time"
)

type MapTask struct {
    chunkID int
    data    string
}

type Pair struct {
    key   string
    value int
}

var chunks = map[int]string{
    0: "This is the first chunk of the data.",
    1: "This is the second chunk of the data.",
    2: "This is the third chunk of the data.",
}

var taskStatus = make(map[int]string)
var results = make(map[string]int)
var mu sync.Mutex

func splitWords(data string) []string {
    return strings.Fields(data)
}

func mapFunc(data string) []Pair {
    wordCount := make(map[string]int)
    words := splitWords(data)

    for _, word := range words {
        wordCount[word]++
    }

    pairs := make([]Pair, 0, len(wordCount))
    for word, count := range wordCount {
        pairs = append(pairs, Pair{key: word, value: count})
    }
    return pairs
}

func sendMapTask(task MapTask) []Pair {
    time.Sleep(1 * time.Second) // Имитация времени обработки
    return mapFunc(task.data)
}

func master() {
    var wg sync.WaitGroup
    mapResults := make(chan []Pair)
    timeout := time.After(5 * time.Second)

    // Перебор частей
    for i := 0; i < len(chunks); i++ {
        taskStatus[i] = "in-progress"
        wg.Add(1)

        go func(chunkID int) {
            defer wg.Done()
            workerResult := sendMapTask(MapTask{chunkID: chunkID, data: chunks[chunkID]})
            fmt.Printf("Worker finished for chunk %d: %v\n", chunkID, workerResult) // Debug message

            if workerResult != nil {
                mapResults <- workerResult
                taskStatus[chunkID] = "done"
            }
        }(i)
    }

    go func() {
        wg.Wait()
        close(mapResults)
    }()

    combinedResults := make(map[string]int)

    for {
        select {
        case result, ok := <-mapResults:
            if !ok {
                fmt.Println("No more map results.")
                return
            }
            for _, pair := range result {
                mu.Lock()
                combinedResults[pair.key] += pair.value
                mu.Unlock()
            }

        case <-timeout:
            fmt.Println("Timeout reached, checking task statuses...")
            for i, status := range taskStatus {
                if status == "in-progress" {
                    fmt.Printf("Task %d still in progress!\n", i)
                }
            }
            return // Завершение работы в таймаут
        }
    }

    mu.Lock()
    for k, v := range combinedResults {
        results[k] = v
    }
    mu.Unlock()
}

func main() {
    master()
    fmt.Println("Final Results:")
    mu.Lock()
    for k, v := range results {
        fmt.Printf("%s: %d\n", k, v)
    }
    mu.Unlock()
}