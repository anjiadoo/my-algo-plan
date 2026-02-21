package main

import "fmt"

type Graph interface {
	AddEdge(from, to, weight int)
	RemoveEdge(from, to int)
	HasEdge(from, to int) bool
	Weight(from, to int) int
	Size() int
	Neighbors(v int) []*Edge
	Display()
}

type Edge struct {
	from   int
	to     int
	weight int
}

// DFS 遍历图的所有节点
func traverseGraph(num int, graph Graph, v int, visited []bool) {
	if v < 0 || v > num {
		return
	}
	if visited[v] {
		// 防止死循环
		return
	}

	// 前序位置
	visited[v] = true
	fmt.Println("visit=", v)

	for _, edge := range graph.Neighbors(v) {
		traverseGraph(num, graph, edge.to, visited)
	}
}

// DFS 遍历图的所有边
func traverseEdge(num int, graph Graph, v int, visited [][]bool) {
	if v < 0 || v > num {
		return
	}
	for _, edge := range graph.Neighbors(v) {
		// 如果边已经被遍历过，则跳过
		if visited[edge.from][edge.to] {
			return
		}

		// 标记并访问边
		visited[edge.from][edge.to] = true
		fmt.Printf("edge:%d->%d\n", edge.from, edge.to)

		traverseEdge(num, graph, edge.to, visited)
	}
}

// DFS 遍历图的所有路径
func traversePath(num int, graph Graph, src, dest int, onPath []bool, path *[]int, res *[]string) {
	if src < 0 || src > num || dest < 0 || dest > num {
		return
	}

	// 防止死循环（成环）
	if onPath[src] {
		return
	}

	// 找到目标节点
	if src == dest {
		sPath := fmt.Sprintf("长度%d, path=", len(*path)+1)
		for i := 0; i < len(*path); i++ {
			sPath += fmt.Sprintf("%d->", (*path)[i])
		}
		*res = append(*res, sPath+fmt.Sprintf("%d", dest))
		return
	}

	// 前序位置-标记
	onPath[src] = true
	*path = append(*path, src)

	for _, edge := range graph.Neighbors(src) {
		traversePath(num, graph, edge.to, dest, onPath, path, res)
	}

	// 后序位置-撤销标记
	onPath[src] = false
	*path = (*path)[:len(*path)-1]
}

// BFS 从s开始遍历图的所有节点，且记录遍历的步数
func levelOrderTraverseGraph2(graph Graph, s int) {
	visited := make([]bool, graph.Size())
	q := []int{s}
	visited[s] = true
	// 记录从 s 开始走到当前节点的步数
	step := 0
	for len(q) > 0 {
		sz := len(q)
		for i := 0; i < sz; i++ {
			cur := q[0]
			q = q[1:]
			fmt.Printf("bfs: start [%d] visit [%d] at step %d\n", s, cur, step)
			for _, e := range graph.Neighbors(cur) {
				if visited[e.to] {
					continue
				}
				q = append(q, e.to)
				visited[e.to] = true
			}
		}
		step++
	}
}

// BFS 从s开始遍历图的所有节点，适配不同权重边的写法。
func levelOrderTraverseGraph3(graph Graph, s int) {
	type NodeState struct {
		node   int // 当前节点 ID
		weight int // 从起点 s 到当前节点的遍历步数
	}

	visited := make([]bool, graph.Size())
	q := []*NodeState{{node: s, weight: 0}}
	visited[s] = true

	for len(q) > 0 {
		sz := len(q)
		for i := 0; i < sz; i++ {
			cur := q[0]
			q = q[1:]

			fmt.Printf("bfs: start [%d] visit [%d] at weight %d\n", s, cur.node, cur.weight)

			for _, e := range graph.Neighbors(cur.node) {
				if visited[e.to] {
					continue
				}
				q = append(q, &NodeState{node: e.to, weight: cur.weight + e.weight})
				visited[e.to] = true
			}
		}
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

	levelOrderTraverseGraph3(graph, 0)

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
}
