package main

import (
	"fmt"
	"strings"
)

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

func getNeighbors(m, n int, str string) []string {
	var neighbors []string

	i := strings.Index(str, "0")

	// 如果不是第一列，有左侧邻居
	if i%n != 0 {
		s := swap(str, i-1, i)
		neighbors = append(neighbors, s)
	}

	// 如果不是最后一列，有右侧邻居
	if i%n != n-1 {
		s := swap(str, i+1, i)
		neighbors = append(neighbors, s)
	}

	// 如果不是第一行，有上方邻居
	if i-n >= 0 {
		s := swap(str, i-n, i)
		neighbors = append(neighbors, s)
	}

	// 如果不是最后一行，有下方邻居
	if i+n < m*n {
		s := swap(str, i+n, i)
		neighbors = append(neighbors, s)
	}
	return neighbors
}

// 解开密码锁的最少次数
func openLock(deadends []string, target string) int {
	return -1
}
func main() {

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
	//
	//fmt.Println(bfs(graph, 0, 9))

	board := [][]int{{4, 1, 2}, {5, 0, 3}}
	fmt.Println(slidingPuzzle(board))

	board = [][]int{{1, 2, 3}, {5, 4, 0}}
	fmt.Println(slidingPuzzle(board))

	board = [][]int{{1, 2, 3}, {4, 0, 5}}
	fmt.Println(slidingPuzzle(board))
}
