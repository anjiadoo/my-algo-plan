package main

import (
	"fmt"
)

type WeightedDigraph struct {
	// 邻接矩阵，matrix[from][to]存储从节点 from 到节点 to 的边的权重
	// 0 表示没有连接
	size   int
	matrix [][]int // 有向加权图-邻接矩阵
}

func NewWeightedDigraph(n int) *WeightedDigraph {
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
		for j := 0; j < n; j++ {
			matrix[i][j] = -1
		}
	}
	return &WeightedDigraph{
		size:   n,
		matrix: matrix,
	}
}

func (w *WeightedDigraph) AddEdge(from, to, weight int) {
	w.matrix[from][to] = weight
}

func (w *WeightedDigraph) RemoveEdge(from, to int) {
	w.matrix[from][to] = -1
}

func (w *WeightedDigraph) HasEdge(from, to int) bool {
	return w.matrix[from][to] != -1
}

func (w *WeightedDigraph) Weight(from, to int) int {
	return w.matrix[from][to]
}

func (w *WeightedDigraph) Size() int {
	return w.size
}

func (w *WeightedDigraph) Neighbors(v int) []*Edge {
	var edges []*Edge
	for to, weight := range w.matrix[v] {
		if weight != -1 {
			edges = append(edges, &Edge{from: v, to: to, weight: weight})
		}
	}
	return edges
}

func (w *WeightedDigraph) Display() {
	for from := 0; from < len(w.matrix); from++ {
		var edges []string
		for to := 0; to < len(w.matrix); to++ {
			if w.matrix[from][to] != -1 {
				edges = append(edges, fmt.Sprintf("%d->%d(%d)", from, to, w.matrix[from][to]))
			}
		}
		fmt.Printf("v=%d, edges=%+v\n", from, edges)
	}
}

func main_() {
	num := 10
	graph := NewWeightedDigraph(num)
	graph.AddEdge(0, 1, 10)
	graph.AddEdge(0, 2, 10)
	graph.AddEdge(0, 4, 10)
	graph.AddEdge(1, 3, 20)
	graph.AddEdge(1, 4, 20)
	graph.AddEdge(1, 5, 20)
	graph.AddEdge(2, 5, 30)
	graph.AddEdge(2, 3, 30)
	graph.AddEdge(2, 4, 30)
	graph.AddEdge(3, 4, 40)
	graph.AddEdge(4, 5, 50)
	graph.AddEdge(5, 1, 60)
	graph.Display()

	fmt.Println(graph.HasEdge(0, 1))
	fmt.Println(graph.HasEdge(1, 0))

	for _, edge := range graph.Neighbors(2) {
		fmt.Printf("%d -> %d, weight: %d\n", edge.from, edge.to, edge.weight)
	}

	graph.RemoveEdge(0, 1)
	graph.RemoveEdge(0, 2)
	graph.RemoveEdge(0, 4)
	graph.Display()
}
