package main

import (
	"fmt"
)

type MyMatrixGraph struct {
	// 邻接矩阵，matrix[from][to]存储从节点 from 到节点 to 的边的权重
	// 0 表示没有连接
	size   int
	matrix [][]int // 有向加权图-邻接矩阵
}

func NewMyMatrixGraph(n int) *MyMatrixGraph {
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
		for j := 0; j < n; j++ {
			matrix[i][j] = -1
		}
	}
	return &MyMatrixGraph{
		size:   n,
		matrix: matrix,
	}
}

func (w *MyMatrixGraph) AddEdge(from, to, weight int) {
	w.matrix[from][to] = weight
}

func (w *MyMatrixGraph) RemoveEdge(from, to int) {
	w.matrix[from][to] = -1
}

func (w *MyMatrixGraph) HasEdge(from, to int) bool {
	return w.matrix[from][to] != -1
}

func (w *MyMatrixGraph) Weight(from, to int) int {
	return w.matrix[from][to]
}

func (w *MyMatrixGraph) Size() int {
	return w.size
}

func (w *MyMatrixGraph) Neighbors(v int) []*Edge {
	var edges []*Edge
	for to, weight := range w.matrix[v] {
		if weight != -1 {
			edges = append(edges, &Edge{from: v, to: to, weight: weight})
		}
	}
	return edges
}

func (w *MyMatrixGraph) Display() {
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
	graph := NewMyMatrixGraph(num)
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
