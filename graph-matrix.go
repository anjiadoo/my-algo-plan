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

func main() {
	num := 10
	graph := NewWeightedDigraph(num)
	graph.AddEdge(0, 1, 1)
	graph.AddEdge(0, 2, 2)
	graph.AddEdge(0, 3, 3)
	graph.AddEdge(1, 2, 12)
	graph.AddEdge(1, 3, 13)
	graph.AddEdge(1, 4, 14)
	graph.AddEdge(2, 3, 23)
	graph.AddEdge(2, 4, 24)
	graph.AddEdge(2, 5, 25)
	graph.AddEdge(3, 4, 34)
	graph.AddEdge(3, 5, 35)
	graph.AddEdge(4, 5, 45)
	graph.AddEdge(4, 6, 46)
	graph.AddEdge(5, 6, 56)
	graph.AddEdge(5, 7, 57)
	graph.AddEdge(6, 7, 67)
	graph.AddEdge(6, 8, 68)
	graph.AddEdge(7, 8, 78)
	graph.AddEdge(7, 9, 79)
	graph.AddEdge(8, 9, 89)
	graph.AddEdge(8, 0, 80)
	graph.AddEdge(9, 0, 90)
	graph.AddEdge(9, 1, 91)
	graph.Display()

	// 遍历所有路径
	onPath := make([]bool, num)
	var path []int
	var res []string
	traversePath(num, graph, 5, 9, onPath, &path, &res)
	for i, itr := range res {
		fmt.Printf("路径%d: %s\n", i+1, itr)
	}

	bfs(graph, 0)

	//// 遍历节点
	//visited := make([]bool, num)
	//traverseGraph(num, graph, 0, visited)
	//
	//// 遍历边
	//visitedEdge := make([][]bool, num)
	//for i := 0; i < num; i++ {
	//	visitedEdge[i] = make([]bool, num)
	//}
	//traverseEdge(num, graph, 0, visitedEdge)
	//
	//fmt.Println(graph.HasEdge(0, 1))
	//fmt.Println(graph.HasEdge(1, 0))
	//
	//for _, edge := range graph.Neighbors(2) {
	//	fmt.Printf("%d -> %d, weight: %d\n", edge.from, edge.to, edge.weight)
	//}
	//
	//graph.RemoveEdge(0, 1)
	//graph.RemoveEdge(0, 2)
	//graph.RemoveEdge(0, 4)
	//graph.Display()
}
