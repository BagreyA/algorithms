package graph

import "sort"

// DisjointSet
type DisjointSet struct {
	parent []int
	rank   []int
}

func NewDisjointSet(n int) *DisjointSet {
	ds := &DisjointSet{
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		ds.parent[i] = i
		ds.rank[i] = 0
	}
	return ds
}

// возвращает корень элемента x с сжатием пути
func (ds *DisjointSet) Find(x int) int {
	if ds.parent[x] != x {
		ds.parent[x] = ds.Find(ds.parent[x])
	}
	return ds.parent[x]
}

// объединяет два множества
func (ds *DisjointSet) Union(x, y int) bool {
	rx := ds.Find(x)
	ry := ds.Find(y)
	if rx == ry {
		return false
	}
	if ds.rank[rx] < ds.rank[ry] {
		ds.parent[rx] = ry
	} else if ds.rank[rx] > ds.rank[ry] {
		ds.parent[ry] = rx
	} else {
		ds.parent[ry] = rx
		ds.rank[rx]++
	}
	return true
}

type Edge struct {
	From   int
	To     int
	Weight int
}

// MST
func MST(n int, edges []Edge) (mst []Edge, totalWeight int) {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})

	ds := NewDisjointSet(n)
	for _, edge := range edges {
		if ds.Find(edge.u) != ds.Find(edge.v) {
			ds.Union(edge.u, edge.v)
			mst = append(mst, edge)
			totalWeight += edge.w
		}
	}
	return mst, totalWeight
}
