// mapreduce/mapreduce.go
package mapreduce

import (
    "strings"
)

type Pair struct {
    Key   string
    Value int
}

type MapTask struct {
    ChunkID int
    Data    string
}

type ReduceTask struct {
    Key    string
    Values []int
}

// MapFunc выполняет операцию Map, разбивая данные на ключ-значение пары
func mapFunc(data string) []Pair {
    words := strings.Fields(data)
    result := []Pair{}
    for _, word := range words {
        result = append(result, Pair{Key: word, Value: 1})
    }
    return result
}

// ReduceFunc выполняет операцию Reduce, суммируя все значения по ключу
func reduceFunc(key string, values []int) (string, int) {
    total := 0
    for _, value := range values {
        total += value
    }
    return key, total
}

// Master обрабатывает распределение задач Map и Reduce
type Master struct {
    MapWorkers    []chan MapTask
    ReduceWorkers []chan ReduceTask
    Chunks        map[int]string
}

func NewMaster() *Master {
    return &Master{
        MapWorkers:    make([]chan MapTask, 3),
        ReduceWorkers: make([]chan ReduceTask, 2),
        Chunks:        make(map[int]string),
    }
}

// RunMapReduce запускает процесс MapReduce
func (m *Master) RunMapReduce() {
    // Распределение MapTask
    m.distributeMapTasks()
    // Shuffle и распределение ReduceTask
    m.shuffleAndReduce()
}

// distributeMapTasks распределяет задачи Map среди работников
func (m *Master) distributeMapTasks() {
    for chunkID, data := range m.Chunks {
        task := MapTask{ChunkID: chunkID, Data: data}
        workerID := chunkID % len(m.MapWorkers)
        m.MapWorkers[workerID] <- task
    }
}

// shuffleAndReduce группирует результаты и передает ReduceWorkers
func (m *Master) shuffleAndReduce() {
    // Группируем по ключу
    grouped := make(map[string][]int)
    for _, data := range m.Chunks {
        result := mapFunc(data)
        for _, pair := range result {
            grouped[pair.Key] = append(grouped[pair.Key], pair.Value)
        }
    }

    // Распределяем по ReduceWorkers
    for key, values := range grouped {
        workerID := hashKey(key) % len(m.ReduceWorkers)
        task := ReduceTask{Key: key, Values: values}
        m.ReduceWorkers[workerID] <- task
    }
}

func hashKey(key string) int {
    // Простейшая хеш-функция для распределения по ReduceWorker
    sum := 0
    for i := 0; i < len(key); i++ {
        sum += int(key[i])
    }
    return sum
}

func getMapResults(chunkID int) []Pair {
	// Получаем результаты Map
	return []Pair{}
}