package main

import "fmt"

// 🌟技巧1：把graph设置为 map[int][]Edge，可以动态添加新节点
// 🌟技巧2：HasEdge/RemoveEdge/Weight方法遍历List可以优化，比如用map[int]map[int]int存储，就可以避免遍历List，复杂度能降到O(1)。
// 🌟技巧3：参数涉及到切片的，注意扩容带来的影响

type MyWeightedDigraph struct {
	size  int
	graph [][]*Edge // 有向加权图-邻接表
}

func NewMyWeightedDigraph(n int) *MyWeightedDigraph {
	return &MyWeightedDigraph{
		size:  n,
		graph: make([][]*Edge, n),
	}
}

func (m *MyWeightedDigraph) AddEdge(from, to, weight int) {
	m.graph[from] = append(m.graph[from], &Edge{from: from, to: to, weight: weight})
}

func (m *MyWeightedDigraph) RemoveEdge(from, to int) {
	for i := 0; i < len(m.graph[from]); i++ {
		if m.graph[from][i].to == to {
			m.graph[from] = append(m.graph[from][:i], m.graph[from][i+1:]...)
			break
		}
	}
}

func (m *MyWeightedDigraph) HasEdge(from, to int) bool {
	for i := 0; i < len(m.graph[from]); i++ {
		if m.graph[from][i].to == to {
			return true
		}
	}
	return false
}

func (m *MyWeightedDigraph) Weight(from, to int) int {
	for i := 0; i < len(m.graph[from]); i++ {
		if m.graph[from][i].to == to {
			return m.graph[from][i].weight
		}
	}
	return -1
}

func (m *MyWeightedDigraph) Size() int {
	return m.size
}

func (m *MyWeightedDigraph) Neighbors(v int) []*Edge {
	return m.graph[v]
}

func (m *MyWeightedDigraph) Display() {
	for i := 0; i < len(m.graph); i++ {
		var edges []string
		for _, e := range m.graph[i] {
			edge := fmt.Sprintf("%d->%d(%d)", e.from, e.to, e.weight)
			edges = append(edges, edge)
		}
		fmt.Printf("v=%d, edges=%+v\n", i, edges)
	}
}

func main__() {
	num := 6
	graph := NewMyWeightedDigraph(num)
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
		fmt.Printf("%d -> %d, weight: %d\n", 2, edge.to, edge.weight)
	}

	graph.RemoveEdge(0, 1)
	graph.RemoveEdge(0, 2)
	graph.RemoveEdge(0, 4)
	graph.Display()
}
