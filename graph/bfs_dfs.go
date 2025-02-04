package graph

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
