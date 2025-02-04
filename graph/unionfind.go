package graph

// UnionFind - структура для работы с системой непересекающихся множеств
type UnionFind struct {
	parent []int
	rank   []int
}

// NewUnionFind - создает новую структуру UnionFind
func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.rank[i] = 0
	}
	return uf
}

// Find - находит родительский элемент
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

// Union - объединяет два множества
func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)
	if rootX == rootY {
		return false
	}

	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}
	return true
}
