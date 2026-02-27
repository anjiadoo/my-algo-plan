package main

import "fmt"

// 图算法实现：
// 🌟技巧1：先检查后标记技巧 - DFS遍历节点前必须先检查visited防止死循环，然后立即标记visited再访问（节点遍历）
// 🌟技巧2：二维visited技巧 - 边遍历用二维数组visited[from][to]，同一方向只遍历一次，用判断后return而非continue（边遍历）
// 🌟技巧3：onPath回溯技巧 - 找所有路径用onPath而非visited，前序位置add到path，后序位置pop出来（路径遍历）
// 🌟技巧4：层级计数技巧 - BFS用sz记录当前层的大小，每层循环结束后step++才能正确计算步数（BFS层次遍历）
// 🌟技巧5：入队即标记技巧 - BFS入队时立即标记visited，而非出队时，防止同一节点被多次入队（BFS防重复）
// 🌟技巧6：带权状态技巧 - BFS带权重用state结构体存(node, weight)，入队时用累计weight（带权重BFS）

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
func traverseGraph(graph Graph, v int, visited []bool) {
	// 1、base case：检查节点边界
	// 2、visited判断：防止重复访问
	// 3、标记visited：记录当前节点已访问
	// 4、访问节点：执行节点相关操作
	// 5、for 循环 neighbors：遍历所有邻居节点
	// 6、递归调用：对未访问的邻居进行DFS

	if v < 0 || v >= graph.Size() {
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
		traverseGraph(graph, edge.to, visited)
	}
}

// DFS 遍历图的所有边
func traverseEdge(graph Graph, v int, visited [][]bool) {
	// 1、base case：检查节点边界
	// 2、for 循环 neighbors：遍历所有出边
	//   2.1 边是否已访问：检查visited[from][to]
	//   2.2 标记边：visited[from][to] = true
	//   2.3 访问边：执行边相关操作
	//   2.4 递归调用：对目标节点继续遍历

	if v < 0 || v >= graph.Size() {
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

		traverseEdge(graph, edge.to, visited)
	}
}

// DFS 遍历图的所有路径
func traversePath(graph Graph, src, dest int, onPath []bool, path *[]int, res *[]string) {
	// 1、base case：检查节点边界
	// 2、是否成环：检查onPath[src]防止循环
	// 3、是否目标节点：src == dest时记录路径
	// 4、前序位置-标记：onPath[src] = true, path追加src
	// 5、for 循环 neighbors：遍历所有邻居
	// 6、递归调用：继续寻找路径
	// 7、后序位置-撤销：onPath[src] = false, path弹出src

	if src < 0 || src >= graph.Size() || dest < 0 || dest >= graph.Size() {
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
		traversePath(graph, edge.to, dest, onPath, path, res)
	}

	// 后序位置-撤销标记
	onPath[src] = false
	*path = (*path)[:len(*path)-1]
}

// BFS 从s开始遍历图的所有节点，且记录遍历的步数
func levelOrderTraverseGraph2(graph Graph, s int) {
	// 1、初始化：visited数组，队列，步数step=0
	// 2、起点入队并标记：q=[s], visited[s]=true
	// 3、while队列不为空：
	//   3.1 记录当前层大小sz
	//   3.2 遍历当前层所有节点：
	//       3.2.1 出队当前节点
	//       3.2.2 访问节点（可记录步数）
	//       3.2.3 遍历邻居：
	//           3.2.3.1 检查是否已访问
	//           3.2.3.2 入队并标记visited
	//   3.3 当前层处理完毕，step++

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
	type state struct {
		node   int // 当前节点 ID
		weight int // 从起点 s 到当前节点的遍历步数
	}

	visited := make([]bool, graph.Size())
	q := []*state{{node: s, weight: 0}}
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
				visited[e.to] = true
				q = append(q, &state{node: e.to, weight: cur.weight + e.weight})
			}
		}
	}
}

func main2() {
	num := 10
	graph := NewMyMatrixGraph(num)
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
	traversePath(graph, 5, 9, onPath, &path, &res)
	for i, itr := range res {
		fmt.Printf("路径%d: %s\n", i+1, itr)
	}

	levelOrderTraverseGraph3(graph, 4)

	//// 遍历节点
	//visited := make([]bool, num)
	//traverseGraph(graph, 0, visited)
	//
	//// 遍历边
	//visitedEdge := make([][]bool, num)
	//for i := 0; i < num; i++ {
	//	visitedEdge[i] = make([]bool, num)
	//}
	//traverseEdge(graph, 0, visitedEdge)
}
