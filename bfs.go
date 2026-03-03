package main

import (
	"fmt"
	"strings"
)

// BFS（广度优先搜索）算法实现：
// 🌟技巧1：层序遍历框架 - BFS核心是用队列逐层扩散，外层for控制层数（步数），内层for遍历当前层所有节点，每次出队一个节点并将其未访问邻居入队
// 🌟技巧2：visited防重技巧 - 用visited集合记录已访问节点，入队时立即标记（而非出队时），避免同一节点被重复加入队列导致超时
// 🌟技巧3：状态抽象技巧 - BFS不仅适用于图/树，任何问题只要能抽象成「状态+状态转移」就能用BFS求最短路径（如滑动谜题用字符串表示棋盘状态、密码锁用字符串表示锁的状态）
// 🌟技巧4：双向BFS优化 - 当知道起点和终点时，可从两端同时BFS，在中间相遇时停止，将时间复杂度从O(k^d)降低到O(2*k^(d/2))

// ⚠️易错点1：visited标记时机 - 必须在入队时就标记visited，而不是出队时标记；否则同一节点会被多次入队，导致重复计算甚至死循环
// ⚠️易错点2：depth递增位置 - depth++必须放在内层for循环之后（即当前层全部处理完后），放错位置会导致步数计算错误
// ⚠️易错点3：起点就是终点 - 别忘了处理start==target的情况，BFS框架中判断目标是在出队时进行的，起点入队后第一次出队就会被检查到
// ⚠️易错点4：状态空间爆炸 - 将问题抽象为BFS时要注意状态数量，如果状态空间太大（如棋盘状态），需要考虑剪枝或双向BFS优化

// 何时运用BFS算法：
// ❓1、问题是否要求「最短路径」或「最少步数」？BFS天然适合求无权图的最短路径
// ❓2、能否将问题抽象为图的遍历？即定义清楚「节点」和「边」（状态和状态转移）
// ❓3、状态空间是否可控？如果状态数量太大，BFS的队列和visited集合会占用大量内存

// 0、func bfs(graph Graph, s int, target int) int             // BFS遍历图框架
// 1、func slidingPuzzle(board [][]int) int                    // 滑动谜题
// 2、func openLock(deadends []string, target string) int      // 解开密码锁的最少次数

// 从 s 开始 BFS 遍历图的所有节点，且记录遍历的步数，当走到目标节点 target 时，返回步数
func bfs(graph Graph, s int, target int) int {
	visited := map[int]bool{}
	q := []int{s}
	visited[s] = true

	// 记录从 s 开始走到当前节点的步数
	depth := 0

	for len(q) > 0 {
		sz := len(q)
		for i := 0; i < sz; i++ {
			cur := q[0]
			q = q[1:]

			// 判断是否到达终点
			if cur == target {
				return depth
			}

			// 将邻居节点加入队列，向四周扩散搜索
			for _, e := range graph.Neighbors(cur) {
				if !visited[e.to] {
					q = append(q, e.to)
					visited[e.to] = true
				}
			}
		}
		depth++
	}
	// 如果走到这里，说明在图中没有找到目标节点
	return -1
}

// 滑动谜题
func slidingPuzzle(board [][]int) int {
	// 在一个2x3的板上（board）有5块砖瓦，用数字1~5来表示, 以及一块空缺用0来表示
	// 一次移动定义为选择0与一个相邻的数字（上下左右）进行交换，
	// 最终当板 board 的结果是 "123450" 时谜板被解开。
	// 给出一个谜板的初始状态board，返回「最少」可以通过多少次移动解开谜板，如果不能解开谜板，则返回-1。

	m, n := len(board), len(board[0])

	target := "123450"
	start := ""
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			start += string('0' + rune(board[i][j]))
		}
	}

	visited := make(map[string]bool)
	q := []string{start}
	visited[start] = true

	depth := 0

	for len(q) > 0 {
		sz := len(q)
		for i := 0; i < sz; i++ {
			cur := q[0]
			q = q[1:]

			if cur == target {
				return depth
			}

			for _, neighbor := range getNeighbors(m, n, cur) {
				if visited[neighbor] {
					continue
				}
				q = append(q, neighbor)
				visited[neighbor] = true
			}
		}
		depth++
	}

	return -1
}

func swap(str string, i, j int) string {
	chs := []rune(str)
	chs[i], chs[j] = chs[j], chs[i]
	return string(chs)
}

func getNeighbors(m, n int, src string) []string {
	var neighbors []string

	i := strings.Index(src, "0")

	// 如果不是第一列，有左侧邻居
	if i%n != 0 {
		s := swap(src, i-1, i)
		neighbors = append(neighbors, s)
	}

	// 如果不是最后一列，有右侧邻居
	if i%n != n-1 {
		s := swap(src, i+1, i)
		neighbors = append(neighbors, s)
	}

	// 如果不是第一行，有上方邻居
	if i-n >= 0 {
		s := swap(src, i-n, i)
		neighbors = append(neighbors, s)
	}

	// 如果不是最后一行，有下方邻居
	if i+n < m*n {
		s := swap(src, i+n, i)
		neighbors = append(neighbors, s)
	}
	return neighbors
}

// 解开密码锁的最少次数
func openLock(deadends []string, target string) int {
	// 有一个带有四个圆形拨轮的转盘锁。每个拨轮都有10个数字： '0', '1', '2', '3', '4', '5', '6', '7', '8', '9' 。
	// 每个拨轮可以自由旋转：例如把 '9' 变为 '0'，'0' 变为 '9' 。每次旋转都只能旋转一个拨轮的一位数字。
	// 锁的初始数字为 '0000' ，一个代表四个拨轮的数字的字符串。
	// 列表 deadends 包含了一组死亡数字，一旦拨轮的数字和列表里的任何一个元素相同，这个锁将会被永久锁定，无法再被旋转。
	// 字符串 target 代表可以解锁的数字，你需要给出解锁需要的最小旋转次数，如果无论如何不能解锁，返回 -1 。

	visited := make(map[string]bool)
	q := []string{"0000"}
	visited["0000"] = true

	for _, dead := range deadends {
		if dead == "0000" {
			return -1
		}
		visited[dead] = true
	}

	depth := 0

	for len(q) > 0 {
		sz := len(q)
		for i := 0; i < sz; i++ {
			cur := q[0]
			q = q[1:]

			if cur == target {
				return depth
			}

			for _, neighbor := range getSelectList(cur) {
				if visited[neighbor] {
					continue
				}
				q = append(q, neighbor)
				visited[neighbor] = true
			}
		}
		depth++
	}

	return -1
}

func getSelectList(src string) []string {
	var neighbors []string

	for i := 0; i < 4; i++ {
		// 向上拨动
		chs := []rune(src)
		if chs[i] == '9' {
			chs[i] = '0'
		} else {
			chs[i] += 1
		}
		neighbors = append(neighbors, string(chs))

		// 向下拨动
		chs = []rune(src)
		if chs[i] == '0' {
			chs[i] = '9'
		} else {
			chs[i] -= 1
		}
		neighbors = append(neighbors, string(chs))
	}
	return neighbors
}

func main() {

	deadends := []string{"0201", "0101", "0102", "1212", "2002"}
	target := "0202"
	fmt.Println(openLock(deadends, target))

	deadends = []string{"0000"}
	target = "0009"
	fmt.Println(openLock(deadends, target))

	deadends = []string{"8887", "8889", "8878", "8898", "8788", "8988", "7888", "9888"}
	target = "8888"
	fmt.Println(openLock(deadends, target))

	//num := 6
	//graph := NewMyMatrixGraph(num)
	//graph.AddEdge(0, 1, 1)
	//graph.AddEdge(0, 2, 2)
	//graph.AddEdge(0, 3, 3)
	//graph.AddEdge(1, 2, 12)
	//graph.AddEdge(1, 3, 13)
	//graph.AddEdge(1, 4, 14)
	//graph.AddEdge(2, 3, 23)
	//graph.AddEdge(2, 4, 24)
	//graph.AddEdge(2, 5, 25)
	//graph.AddEdge(3, 4, 34)
	//graph.AddEdge(3, 5, 35)
	//graph.AddEdge(4, 5, 45)
	//graph.Display()

	//fmt.Println(bfs(graph, 0, 9))

	//board := [][]int{{4, 1, 2}, {5, 0, 3}}
	//fmt.Println(slidingPuzzle(board))
	//board = [][]int{{1, 2, 3}, {5, 4, 0}}
	//fmt.Println(slidingPuzzle(board))
	//board = [][]int{{1, 2, 3}, {4, 0, 5}}
	//fmt.Println(slidingPuzzle(board))
}
