/*
 * ============================================================================
 *                   📘 BFS广度优先搜索全集 · 核心记忆框架
 * ============================================================================
 * 【一句话理解各题型】
 *
 *   BFS框架：  用队列逐层扩散，外层 for 控制层数（步数），内层 for 遍历当前层所有节点，
 *              天然求无权图最短路径。
 *   滑动谜题： 将二维棋盘「序列化」为字符串作为节点状态，getNeighbors 生成每次合法移动后
 *              的新状态，BFS 求从初始状态到 "123450" 的最少交换次数。
 *   密码锁：   将四位拨轮序列化为字符串，将 deadends 预加入 visited（相当于删除死亡节点），
 *              getSelectList 枚举 8 个邻居状态，BFS 求最少拨动次数。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【BFS 核心框架三要素 + 标准模板（背熟！）】
 *
 *   ① 队列 (Queue)  ：存储待访问节点，先进先出，保证逐层扩展。
 *   ② visited 集合  ：记录已访问节点，「入队时」立即标记，防止重复入队。
 *   ③ 层级计数 depth：内层 for 处理完整一层的所有节点后 depth++，表示走了一步。
 *
 *     q := []T{start}
 *     visited[start] = true          // ← 起点入队时立即标记
 *     depth := 0
 *     for len(q) > 0 {
 *         sz := len(q)               // ← 快照当前层大小，内层必须用 sz 而非 len(q)
 *         for i := 0; i < sz; i++ {
 *             cur := q[0]; q = q[1:] // ← 出队（Go 用 slice 模拟队列）
 *             if cur == target { return depth }  // ← 出队时检查目标
 *             for _, next := range neighbors(cur) {
 *                 if !visited[next] {
 *                     q = append(q, next)
 *                     visited[next] = true  // ← 入队时标记！
 *                 }
 *             }
 *         }
 *         depth++                    // ← 整层处理完后才 ++
 *     }
 *     return -1
 *
 *   注意：树结构做 BFS（层序遍历）不需要 visited，因为树无环，不会重复访问。
 *         图结构必须加 visited，否则环会导致无限循环。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【BFS 四大经典应用场景】
 *
 *   场景①：无权图最短路径
 *     代表题：二叉树最小深度、迷宫最短路、密码锁最少步数。
 *     特征：每条边权重相同（都是"走一步"），BFS 天然保证最短。
 *     ⚠️ 有权图不能用 BFS，要用 Dijkstra（堆优化）或 Bellman-Ford。
 *
 *   场景②：多源 BFS（Multi-source BFS）
 *     代表题：LeetCode 542(01矩阵)、994(腐烂的橘子)。
 *     特征：有多个起点，求每个格子到「最近起点」的距离。
 *     做法：将所有起点同时入队，visited 同时标记，整体跑一遍 BFS 即可。
 *     模板：for each source: q.add(source); visited[source]=true  → 正常跑 BFS
 *     优势：比从每个格子分别跑 BFS 快得多（O(m*n) vs O((m*n)²)）。
 *
 *   场景③：状态空间 BFS（State-space BFS）
 *     代表题：滑动谜题、密码锁、单词接龙(WordLadder)。
 *     特征：节点不是坐标/整数，而是复杂状态（棋盘布局、字符串）。
 *     做法：将状态序列化为字符串/哈希，作为 visited 的 key。
 *     状态数量决定时间复杂度：密码锁 10^4 = 10000 状态，滑动谜题 6! = 720 状态。
 *
 *   场景④：拓扑排序 BFS（Kahn 算法）
 *     代表题：课程表(CourseSchedule)、任务调度。
 *     特征：有向无环图(DAG)，求拓扑序或判断是否有环。
 *     做法：统计每个节点的入度，将入度为 0 的节点入队，出队时将邻居入度 -1，
 *           入度变为 0 的邻居再入队，最终处理节点数 == 总节点数则无环。
 *     模板：
 *       inDegree := countInDegree(graph)
 *       for node where inDegree[node]==0: q.add(node)
 *       for len(q) > 0:
 *           cur = q.pop()
 *           result.add(cur)
 *           for next in neighbors(cur):
 *               inDegree[next]--
 *               if inDegree[next] == 0: q.add(next)
 *       return len(result) == numNodes  // true=无环，false=有环
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【网格 BFS 专项：四方向模板 + 边界处理两种写法】
 *
 *   网格问题（island、matrix）是 BFS 最高频的考场景，有两种常用写法：
 *
 *   写法A：方向数组（推荐，代码最简洁）
 *     dirs := [][2]int{{0,1},{0,-1},{1,0},{-1,0}}  // 右左下上
 *     for _, d := range dirs {
 *         nr, nc := r+d[0], c+d[1]
 *         if nr>=0 && nr<m && nc>=0 && nc<n && !visited[nr][nc] {
 *             q = append(q, [2]int{nr, nc})
 *             visited[nr][nc] = true
 *         }
 *     }
 *
 *   写法B：一维索引（序列化棋盘时用，见 slidingPuzzle）
 *     一维索引 i 的四向邻居边界检查：
 *       左：i-1，条件 i%n != 0      （i%n==0 表示在最左列）
 *       右：i+1，条件 i%n != n-1    （i%n==n-1 表示在最右列）
 *       上：i-n，条件 i-n >= 0
 *       下：i+n，条件 i+n < m*n
 *     ⚠️ 左右边界必须用 i%n 判断，不能用 i>0/i<m*n-1
 *        i=n（第二行第一列）满足 i>0 但没有左邻居，i%n==0 才是精确判断。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【BFS vs DFS 选择指南】
 *
 *   维度           BFS                            DFS
 *   ──────────────────────────────────────────────────────────────────────
 *   数据结构       队列（先进先出）                栈 / 递归调用栈
 *   遍历顺序       逐层扩展（层序）                沿一条路走到底
 *   最短路径       ✅ 天然保证（无权图）            ❌ 不保证最短
 *   空间复杂度     O(w)，w = 最大宽度             O(h)，h = 最大深度
 *   树 vs 图       树不需要 visited               图必须 visited 防环
 *   适用场景       最少步数/层数/拓扑排序          路径枚举/连通性/回溯
 *
 *   选 BFS 的信号词：「最少」「最短」「最小步数」「最近」「层数」
 *   选 DFS 的信号词：「所有路径」「是否存在」「连通分量」「全排列」
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【双向 BFS 优化原理与模板】
 *
 *   普通 BFS：从 start 扩展，最坏扩展 O(k^d) 个节点（k=分支数，d=路径长度）
 *   双向 BFS：start 和 target 同时扩展，各 O(k^(d/2))，复杂度指数级下降。
 *   前提：必须知道终点（target 确定），否则无法从终点反向扩展。
 *
 *   关键技巧：每次选择「节点数较少」的那侧扩展，保持两侧平衡：
 *     if len(q1) > len(q2): swap(q1,q2); swap(visited1,visited2)
 *
 *   相遇判断：扩展 next 时检查 visited2[next]，若为 true 则两侧已相遇，返回 depth。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【BFS 不适用的场景（常见误区）】
 *
 *   ❌ 有权图最短路径：边权不同时，BFS 不保证最短 → 用 Dijkstra（非负权）
 *   ❌ 负权图最短路径：有负权边 → 用 Bellman-Ford 或 SPFA
 *   ❌ 求所有路径：BFS 只找最短，找所有路径要用 DFS + 回溯
 *   ❌ 状态空间极大：状态数 > 10^7 级别，BFS 内存爆炸 → 考虑 A* 或双向 BFS
 *   ❌ 树的最大深度/路径总和：需要遍历所有路径，DFS 更自然
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【状态抽象解题法（复杂场景 BFS 通用思路）】
 *
 *   ① 定义「节点」：问题的一个完整状态（棋盘布局、锁读数、人物位置+钥匙状态）
 *   ② 定义「边」  ：一次合法操作后的状态转移
 *   ③ 序列化状态  ：多维状态 → 字符串/整数，用作 map key 或 visited 下标
 *   ④ 枚举邻居    ：从当前状态生成所有合法的下一状态
 *
 *   序列化方式选择：
 *     字符串拼接：通用，适合任意状态，但哈希稍慢（如棋盘、锁）
 *     整数编码：状态量小时更快（如 4 位 0~9 的锁：直接用 4 位十进制数 0000~9999）
 *     位掩码：状态是布尔集合时（如「已收集哪些钥匙」用 bitmask 压缩）
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次手写 BFS 后对照检查
 *
 *     ✅ visited 是在「入队时」标记，而不是「出队时」？
 *     ✅ 起点入队的同时有没有 visited[start] = true？
 *     ✅ 内层 for 前用 sz := len(q) 快照当前层大小，内层用 sz 而非 len(q)？
 *     ✅ depth++ 放在「内层 for 之后」，不是内层 for 之中？
 *     ✅ 图结构有没有加 visited 防环？（树结构不需要）
 *     ✅ 多源 BFS：所有起点是否同时入队并同时标记 visited？
 *     ✅ 网格 BFS：用方向数组写法，边界统一用 nr>=0&&nr<m&&nc>=0&&nc<n 判断？
 *     ✅ 序列化棋盘：左右边界用 i%n 判断，而不是 i>0 / i<m*n-1？
 *     ✅ 密码锁：deadends 有没有提前加入 visited，起点在 deadends 中特判返回-1？
 *     ✅ 环形拨轮：'9'+1→'0'，'0'-1→'9'，有没有手动处理字符环绕？
 *     ✅ 拓扑排序：最终处理节点数是否等于总节点数（判断有无环）？
 *     ✅ 双向BFS：每轮扩展较小的那一侧，两个 visited 集合分别维护？
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *
 *     0. bfs(graph Graph, s int, target int) int          // BFS遍历图框架，返回步数
 *     1. slidingPuzzle(board [][]int) int                 // 滑动谜题，状态BFS
 *     2. openLock(deadends []string, target string) int   // 密码锁，状态BFS+deadends预处理
 *     辅助：getNeighbors(m, n int, src string) []string   // 滑动谜题邻居（二维→一维索引边界判断）
 *     辅助：getSelectList(src string) []string            // 密码锁邻居（4拨轮×2方向=8个状态）
 * ============================================================================
 */

package main

import (
	"fmt"
	"strings"
)

// BFS（广度优先搜索）算法实现：
// 🌟技巧1：层序遍历框架技巧 - 外层for控制层数（步数），内层for用 sz:=len(q) 快照遍历当前层所有节点，depth++放在内层for后，保证步数计算正确
// 🌟技巧2：入队时标记visited技巧 - 入队时立即 visited[next]=true，而非出队时标记；出队时标记会导致同一节点多次入队，引发重复计算甚至超时
// 🌟技巧3：状态序列化技巧 - 多维状态（棋盘/拨轮）转字符串作为 map 的 key；二维棋盘逐行拼接，多位数字直接拼接，统一节点表示
// 🌟技巧4：二维坐标一维化技巧 - 棋盘 m×n 中，一维索引 i 的四向邻居：左 i-1(i%n!=0)，右 i+1(i%n!=n-1)，上 i-n(i-n>=0)，下 i+n(i+n<m*n)，用 i%n 判断左右边界而非 i>0
// 🌟技巧5：deadends预处理技巧 - 将所有死亡数字提前加入 visited，BFS扩展时自然跳过，代码逻辑统一；但起点 "0000" 在 deadends 中需在初始化阶段特判返回 -1
// 🌟技巧6：环形拨轮技巧 - 拨轮上下各转一格共 4×2=8 个邻居；'9'+1→'0'，'0'-1→'9'，不能直接 ±1，必须手动处理环绕（rune运算无法自动回绕）
// 🌟技巧7：双向BFS优化技巧 - 知道起点和终点时，两端同时BFS，每次扩展「较小」的一侧，两侧 visited 相交即停止，时间复杂度从O(k^d)→O(k^(d/2))

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
