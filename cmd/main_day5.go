package main

import (
	"fmt"
	"math/rand"
	"myproject/graph"
	"myproject/dist"
	"myproject/mapreduce"
	"time"
)

func main() {
	// Инициализируем Master для MapReduce
	master := mapreduce.NewMaster()

	// Добавляем несколько кусков текста
	master.Chunks[0] = "hello world"
	master.Chunks[1] = "hello go"
	master.Chunks[2] = "go is great"

	// Запускаем процесс MapReduce
	master.RunMapReduce()

	// Выводим результаты
	fmt.Println("MapReduce completed.")
}