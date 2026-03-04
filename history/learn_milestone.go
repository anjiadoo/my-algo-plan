package history

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 记录刷题过程
// 每道题需要注明力扣题目ID/链接
// 结构体定义放在在上面，题目必须写注释
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// TreeNode 二叉树数据结构定义
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ListNode 单链表数据结构定义
type ListNode struct {
	Val  int
	Next *ListNode
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func sum(numbers []int) int {
	res, left, right := 0, 0, len(numbers)-1
	for left < right {
		res += numbers[left] + numbers[right]
		left++
		right--
	}
	if len(numbers)%2 != 0 {
		res += numbers[left]
	}
	return res
}

type maxIntHeap []int // 大根堆

func (h maxIntHeap) Len() int           { return len(h) }
func (h maxIntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h maxIntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h *maxIntHeap) Push(v any)        { *h = append(*h, v.(int)) }
func (h *maxIntHeap) Pop() any          { a := *h; v := a[len(a)-1]; *h = a[:len(a)-1]; return v }

/**********************************************************************************************************************/
// 核心算法框架/算法思维
/**********************************************************************************************************************/
// https://leetcode.cn/problems/binary-tree-maximum-path-sum/description/
// 二叉树的遍历框架：前中后序
// 通常二叉树的问题都可以通过遍历解决
// maxPathSum 二叉树中的最大路径和
func maxPathSum(root *TreeNode) int {
	var maxSum = math.MinInt
	var maxFunc func(root *TreeNode) int

	maxFunc = func(root *TreeNode) int {
		if root == nil {
			return 0
		}

		//递归计算左右子节点的最大贡献值
		//只有在最大贡献值>0时，才选取该节点
		left := max(0, maxFunc(root.Left))
		right := max(0, maxFunc(root.Right))

		//更新答案
		maxSum = max(maxSum, left+right+root.Val)

		//返回当前节点的最大贡献值
		return max(left, right) + root.Val
	}
	maxFunc(root)
	return maxSum
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/construct-binary-tree-from-preorder-and-inorder-traversal/description/
// 输入: preorder = [3,9,20,15,7], inorder = [9,3,15,20,7]
// 输出: [3,9,20,null,null,15,7]
// 其实也是在考察二叉树的遍历框架
// buildTree 从前序与中序遍历序列构造二叉树
func buildTree(preorder []int, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}
	mapping := map[int]int{}
	for i, val := range inorder {
		mapping[val] = i
	}
	//注意边界问题，len()-1
	return build(mapping, preorder, 0, len(preorder)-1, inorder, 0, len(inorder)-1)
}

func build(mapping map[int]int, preorder []int, preStart, preEnd int, inorder []int, inStart, inEnd int) *TreeNode {
	if preStart > preEnd {
		return nil
	}

	rootVal := preorder[preStart]

	index := mapping[rootVal]
	leftSize := index - inStart

	node := &TreeNode{Val: rootVal}
	node.Left = build(mapping, preorder, preStart+1, preStart+leftSize, inorder, inStart, index-1)
	node.Right = build(mapping, preorder, preStart+leftSize+1, preEnd, inorder, index+1, inEnd)
	return node
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/kth-smallest-element-in-a-bst/
// 输入：root = [3,1,4,null,2], k = 1
// 输出：1
// 思路：考察二叉树的遍历框架，在中序位置加逻辑
// kthSmallest 二叉搜索树中第K小的元素，力扣题目230
func kthSmallest(root *TreeNode, k int) int {
	var traverse func(root *TreeNode)
	var res, rank int

	traverse = func(root *TreeNode) {
		if root == nil {
			return
		}
		traverse(root.Left)
		//注意：这里先++，再判断
		rank++
		if rank == k {
			res = root.Val
			return
		}
		traverse(root.Right)
	}
	traverse(root)
	return res
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/search-in-rotated-sorted-array/description/
// 考查二分查找的变体
// 输入：numbers = [4,5,6,7,0,1,2], target = 0
// 输出：4
// search 搜索旋转排序数组
func search(numbers []int, target int) int {
	// 思路：两条上升直线，四种情况判断
	if len(numbers) == 0 {
		return -1
	}
	start, end := 0, len(numbers)-1
	for start+1 < end {
		mid := start + (end-start)/2
		// 相等直接返回
		if numbers[mid] == target {
			return mid
		}
		// 判断在那个区间，可能分为四种情况
		if numbers[start] < numbers[mid] {
			if numbers[start] < target && target < numbers[mid] {
				end = mid
			} else {
				start = mid
			}
		} else if numbers[mid] < numbers[end] {
			if numbers[mid] < target && target < numbers[end] {
				start = mid
			} else {
				end = mid
			}
		}
	}
	if numbers[start] == target {
		return start
	} else if numbers[end] == target {
		return end
	}
	return -1
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/fibonacci-number/description/
// F(0) = 0，F(1) = 1
// F(n) = F(n - 1) + F(n - 2)，其中 n > 1
// 思维：明确base case -> 明确“状态” -> 明确“选择” -> 定义dp数组/函数
// 自顶向下“递归”的动态规划
// fib 斐波那契数，力扣题目509
func fib(n int) int {
	var helper func(memo []int, n int) int

	helper = func(memo []int, n int) int {
		if n == 0 || n == 1 {
			return n
		}
		if memo[n] != 0 {
			return memo[n]
		}
		memo[n] = helper(memo, n-1) + helper(memo, n-2)
		return memo[n]
	}
	return helper(make([]int, n+1), n)
}

// 自底向上"递推"的动态规划
func fib2(n int) int {
	if n == 0 {
		return 0
	}
	dp := make([]int, n+1)

	// base case
	dp[0], dp[1] = 0, 1

	// 状态转移方程
	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}
	return dp[n]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/coin-change/
// 思维：明确base case -> 明确“状态” -> 明确“选择” -> 定义dp数组/函数
// 输入：coins = [1, 2, 5], amount = 11
// 输出：3
// 解释：11 = 5 + 5 + 1
// 自顶向下“递归”的动态规划
// coinChange 零钱兑换，力扣题目322
func coinChange(coins []int, amount int) int {
	var dp func(coins []int, amount int) int
	memo := make([]int, amount+1) //备忘录

	dp = func(coins []int, amount int) int {
		if amount == 0 { //base case
			return 0
		}
		if amount < 0 {
			return -1
		}
		if memo[amount] != 0 {
			return memo[amount]
		}

		var res = math.MaxInt
		for _, coin := range coins {
			subProblem := dp(coins, amount-coin) //计算子问题的结果
			if subProblem == -1 {
				continue
			}
			res = min(res, subProblem+1)
		}
		if res == math.MaxInt {
			memo[amount] = -1
			return memo[amount]
		}
		memo[amount] = res
		return memo[amount]
	}
	return dp(coins, amount)
}

// 自底向上"递推"的动态规划
func coinChange2(coins []int, amount int) int {
	dp := make([]int, amount+1)
	for i := 0; i < len(dp); i++ {
		dp[i] = amount + 1
	}

	dp[0] = 0
	for i := 0; i <= amount; i++ {
		for _, coin := range coins {
			if i-coin < 0 {
				continue
			}
			dp[i] = min(dp[i], dp[i-coin]+1)
		}
	}
	if dp[amount] == amount+1 {
		return -1
	}
	return dp[amount]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/permutations/
// 考察回溯算法：本质就是N叉数的前后序遍历问题
// 输入：numbers = [1,2,3]
// 输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
// permute 全排列
func permute(numbers []int) [][]int {
	var backtrack func(numbers []int, track []int, result *[][]int)

	visited := map[int]bool{} //访问数组
	result := [][]int{}       //结果列表
	track := []int{}          //已选择路径

	backtrack = func(numbers []int, track []int, result *[][]int) {
		if len(numbers) == len(track) {
			temp := make([]int, len(track))
			copy(temp, track)
			*result = append(*result, temp)
			return
		}
		for i := 0; i < len(numbers); i++ {
			if visited[numbers[i]] {
				continue
			}
			visited[numbers[i]] = true
			track = append(track, numbers[i])
			backtrack(numbers, track, result)
			track = track[:len(track)-1]
			delete(visited, numbers[i])
		}
	}

	backtrack(numbers, track, &result)
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/n-queens/
// 考察回溯算法：本质就是N叉数的前后序遍历问题
// 输入：n = 4
// 输出：[[".Q..","...Q","Q...","..Q."],["..Q.","Q...","...Q",".Q.."]]
// 解释：如上图所示，4 皇后问题存在两个不同的解法。
// solveNQueens N皇后
func solveNQueens(n int) [][]string {
	result := [][]string{} //结果数组
	board := []string{}    //初始化棋盘
	for i := 0; i < n; i++ {
		row := ""
		for k := 0; k < n; k++ {
			row += "."
		}
		board = append(board, row)
	}

	_backtrack(board, 0, &result)
	return result
}

func _backtrack(board []string, row int, result *[][]string) {
	if row == len(board) {
		temp := make([]string, len(board))
		copy(temp, board)
		*result = append(*result, temp)
		return
	}
	for col := 0; col < len(board); col++ {
		if isVisited(board, row, col) {
			continue
		}
		str := []byte(board[row])
		str[col] = 'Q'
		board[row] = string(str)

		_backtrack(board, row+1, result)

		str = []byte(board[row])
		str[col] = '.'
		board[row] = string(str)
	}
}

func isVisited(board []string, row, col int) bool {
	//检查列是否有皇后冲突
	for i := 0; i <= row; i++ {
		if board[i][col] == 'Q' {
			return true
		}
	}
	//检查左上方是否有冲突
	for i, j := row-1, col-1; i >= 0 && j >= 0; {
		if board[i][j] == 'Q' {
			return true
		}
		i--
		j--
	}
	//检查右上方是否冲突
	for i, j := row-1, col+1; i >= 0 && j < len(board); {
		if board[i][j] == 'Q' {
			return true
		}
		i--
		j++
	}
	return false
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/minimum-depth-of-binary-tree/
// 本质就是：在一幅“图”中找到从“起点”到“终点”的最短距离，从中心向四周扩散（面）
// BFS 广度优先算法
// minDepth 二叉树的最小深度
func minDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	queue := []*TreeNode{}
	queue = append(queue, root)
	step := 1
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]
			if node.Left == nil && node.Right == nil {
				return step
			}
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		step++
	}
	return step
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/open-the-lock/description/
// 本质就是：在一幅“图”中找到从“起点”到“终点”的最短距离，从中心向四周扩散（面），考查：BFS 广度优先算法
// 输入：deadends = ["0201","0101","0102","1212","2002"], target = "0202"
// 输出：6
// 解释：
// 可能的移动序列为 "0000" -> "1000" -> "1100" -> "1200" -> "1201" -> "1202" -> "0202"。
// 注意 "0000" -> "0001" -> "0002" -> "0102" -> "0202" 这样的序列是不能解锁的，
// 因为当拨动到 "0102" 时这个锁就会被锁定。
// openLock 打开转盘锁，本质考察BFS算法
func openLock(deadends []string, target string) int {
	visited := map[string]bool{}
	for _, dead := range deadends {
		visited[dead] = true
	}
	if visited["0000"] || visited[target] {
		return -1
	}

	queue := []string{"0000"}
	visited["0000"] = true

	pushOne := func(str string, idx int) string {
		byteStr := []byte(str)
		if byteStr[idx] == '9' {
			byteStr[idx] = '0'
		} else {
			byteStr[idx] += 1
		}
		return string(byteStr)
	}
	downOne := func(str string, idx int) string {
		byteStr := []byte(str)
		if byteStr[idx] == '0' {
			byteStr[idx] = '9'
		} else {
			byteStr[idx] -= 1
		}
		return string(byteStr)
	}

	step := 0
	for len(queue) > 0 {
		size := len(queue)
		for k := 0; k < size; k++ {
			curr := queue[0]
			queue = queue[1:]
			if curr == target {
				return step
			}
			for i := 0; i < 4; i++ {
				if up := pushOne(curr, i); !visited[up] {
					queue = append(queue, up)
					visited[up] = true
				}
				if down := downOne(curr, i); !visited[down] {
					queue = append(queue, down)
					visited[down] = true
				}
			}
		}
		step++
	}
	return -1
}

/**********************************************************************************************************************/
//https://leetcode.cn/problems/maximum-depth-of-binary-tree/description/
// 通过遍历的思维模式计算答案，其实也就是在考察二叉树的遍历框架，只是需要在特殊时间点添加逻辑（前/后序位置）
// maxDepth 二叉树的最大深度
func maxDepth(root *TreeNode) int {
	var traverse func(root *TreeNode)
	var depth, result int
	traverse = func(root *TreeNode) {
		if root == nil {
			result = max(result, depth)
			return
		}
		depth++
		traverse(root.Left)
		traverse(root.Right)
		depth--
	}
	traverse(root)
	return result
}

// 通过分解子问题的思维模式计算答案
func _maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return max(_maxDepth(root.Left), _maxDepth(root.Right)) + 1
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/diameter-of-binary-tree/description/
// 后序位置插入操作的栗子，后序位置比较特殊，可以从子节点获取数据
// diameterOfBinaryTree 二叉树的直径
func diameterOfBinaryTree(root *TreeNode) int {
	var maxDepth func(root *TreeNode) int
	var maxDiameter = 0

	maxDepth = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		leftDepth := maxDepth(root.Left)
		rightDepth := maxDepth(root.Right)

		//后序位置计算出最大直径
		maxDiameter = max(maxDiameter, leftDepth+rightDepth)
		return max(leftDepth, rightDepth) + 1
	}
	maxDepth(root)
	return maxDiameter
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/find-first-and-last-position-of-element-in-sorted-array/description/
// 考察二分查找：通常出现寻找一个数/寻找左侧边界/右侧边界，那么基本都是用二分查找框架解决
// 输入：numbers = [5,7,7,8,8,10]
// target = 8，输出：[3,4]
// target = 6，输出：[-1,-1]
// target = 10，输出：[5,5]
// 代码复用 => 可以转化题目为：查找第一个小于target的坐标
// searchRange 在排序数组中查找元素的第一个和最后一个位置
func searchRange(numbers []int, target int) []int {
	leftBound := findLeftBound(numbers, target)
	if leftBound == -1 {
		return []int{-1, -1}
	}
	rightBound := findRightBound(numbers, target)
	return []int{leftBound, rightBound}
}

func findLeftBound(numbers []int, target int) int {
	left, right := 0, len(numbers)-1

	// [left, right] => left = right+1
	for left <= right {
		mid := left + (right-left)/2
		if numbers[mid] == target {
			right = mid - 1
		} else if numbers[mid] > target {
			right = mid - 1
		} else if numbers[mid] < target {
			left = mid + 1
		}
	}
	//要不要判断边界，就看nums[mid] == target时索引的关系
	if left == len(numbers) {
		return -1
	}
	if numbers[left] == target {
		return left
	}
	return -1
}

func findRightBound(numbers []int, target int) int {
	left, right := 0, len(numbers)-1

	// [left, right] => left = right+1
	for left <= right {
		mid := left + (right-left)/2
		if numbers[mid] == target {
			left = mid + 1
		} else if numbers[mid] > target {
			right = mid - 1
		} else if numbers[mid] < target {
			left = mid + 1
		}
	}
	//避免数组越界，原因是left取值可能为0
	//left要不要加1/减1，就看nums[mid] == target时索引的关系（寻找有边界的特殊点）
	if left-1 >= 0 && numbers[left-1] == target {
		return left - 1
	}
	return -1
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/minimum-window-substring/description/
// 考察滑动窗口算法：一般涉及到子串问题，都可以用滑动窗口算法框架解决
// 输入：s = "ADOBECODEBANC", t = "ABC"
// 输出："BANC"
// 解释：最小覆盖子串 "BANC" 包含来自字符串 t 的 'A'、'B' 和 'C'。
// minWindow 最小覆盖子串
func minWindow(s string, t string) string {
	needs := map[byte]int{}
	for i := 0; i < len(t); i++ {
		needs[t[i]]++
	}

	left, right := 0, 0
	window := map[byte]int{}

	match := 0 //单个字符已经匹配计数
	minStartIdx, mixLength := 0, math.MaxInt

	for right < len(s) {
		in := s[right]
		right++
		if _, ok := needs[in]; ok {
			window[in]++
			if window[in] == needs[in] {
				match++
			}
		}

		for match == len(needs) {
			//最小覆盖子串的索引（几乎都是在这个位置做判断）
			if right-left < mixLength {
				minStartIdx = left
				mixLength = right - left
			}

			out := s[left]
			left++

			if _, ok := needs[out]; ok {
				window[out]--
				if window[out] < needs[out] {
					match--
				}
			}
		}
	}

	if mixLength == math.MaxInt {
		return ""
	}
	return s[minStartIdx : minStartIdx+mixLength]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/permutation-in-string/
// 考察滑动窗口算法，算法框架都一样，只是具体判断逻辑不同而已
// 输入：s1 = "ab" s2 = "eidbaooo"
// 输出：true
// 解释：s2 包含 s1 的排列之一 ("ba").
// checkInclusion 字符串的排列
func checkInclusion(s1 string, s2 string) bool {
	needs := map[byte]int{}
	for i := 0; i < len(s1); i++ {
		needs[s1[i]]++
	}

	left, right := 0, 0
	window := map[byte]int{}

	match := 0 //单个字符已经匹配计数

	for right < len(s2) {
		in := s2[right]
		right++
		if _, ok := needs[in]; ok {
			window[in]++
			if window[in] == needs[in] {
				match++
			}
		}

		for match == len(needs) {
			//是否是排列之一（几乎都是在这个位置做判断）
			if right-left == len(s1) {
				return true
			}

			out := s2[left]
			left++

			if _, ok := needs[out]; ok {
				window[out]--
				if window[out] < needs[out] {
					match--
				}
			}
		}
	}
	return false
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/find-all-anagrams-in-a-string/
// 本质也是考察滑动窗口算法，该算法的特殊逻辑都是在开始收缩窗口的时候判断
// 输入: s = "cbaebabacd", p = "abc"
// 输出: [0,6]
// 解释:
// 起始索引等于 0 的子串是 "cba", 它是 "abc" 的异位词。
// 起始索引等于 6 的子串是 "bac", 它是 "abc" 的异位词。
// findAnagrams 找到字符串中所有字母异位词
func findAnagrams(s string, p string) []int {
	needs := map[byte]int{}
	for i := 0; i < len(p); i++ {
		needs[p[i]]++
	}

	left, right := 0, 0
	window := map[byte]int{}

	match := 0
	result := []int{}

	for right < len(s) {
		in := s[right]
		right++
		if _, ok := needs[in]; ok {
			window[in]++
			if window[in] == needs[in] {
				match++
			}
		}

		for match == len(needs) {
			//判断逻辑（几乎都是在这个位置做判断）
			if right-left == len(p) {
				result = append(result, left)
			}

			out := s[left]
			left++
			if _, ok := needs[out]; ok {
				window[out]--
				if window[out] < needs[out] {
					match--
				}
			}
		}
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/longest-substring-without-repeating-characters/
// 涉及到了子串的问题，不用想滑动窗口直接往上贴
// 输入: s = "abcabcbb"
// 输出: 3
// 解释: 因为无重复字符的最长子串是 "abc"，所以其长度为 3。
// lengthOfLongestSubstring 无重复字符的最长子串
func lengthOfLongestSubstring(s string) int {
	left, right := 0, 0
	window := map[byte]int{}

	maxSubLen := 0

	for right < len(s) {
		in := s[right]
		right++
		window[in]++

		for window[in] > 1 {
			out := s[left]
			left++
			window[out]--
		}
		//因为收缩的条件是存在重复字符，所以收缩完成后就保证了一定没有重复字符
		maxSubLen = max(maxSubLen, right-left)
	}
	return maxSubLen
}

/**********************************************************************************************************************/
//数据结构题型——链表and数组
/**********************************************************************************************************************/
// https://leetcode.cn/problems/merge-two-sorted-lists/
// 输入：p1 = [1,2,4], p2 = [1,3,4]
// 输出：[1,1,2,3,4,4]
// mergeTwoLists 合并两个有序链表
func mergeTwoLists(p1 *ListNode, p2 *ListNode) *ListNode {
	dummy := &ListNode{Val: -1} //虚拟头节点

	//p是游走在合并后的链表上的指针
	//有点像拉链的拉锁
	p := dummy

	for p1 != nil && p2 != nil {
		if p1.Val > p2.Val {
			p.Next = p2
			p2 = p2.Next
		} else {
			p.Next = p1
			p1 = p1.Next
		}
		p = p.Next
	}
	if p1 != nil {
		p.Next = p1
	}
	if p2 != nil {
		p.Next = p2
	}
	return dummy.Next
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/merge-k-sorted-lists/
// 本题需要使用到优先级队列-最小堆（可以实现标准库的接口）
// 输入：lists = [[1,4,5],[1,3,4],[2,6]]
// 输出：[1,1,2,3,4,4,5,6]
// 解释：链表数组如下：
//   1->4->5,
//   1->3->4,
//   2->6
// 将它们合并到一个有序链表中得到。
// 1->1->2->3->4->4->5->6
// mergeKLists 合并 K 个升序链表
func mergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	dummy := &ListNode{Val: -1}
	p := dummy

	pq := &minHeap{}
	heap.Init(pq)

	for _, head := range lists {
		if head != nil {
			heap.Push(pq, head)
		}
	}

	for len(*pq) > 0 {
		val := heap.Pop(pq)
		node := val.(*ListNode)
		p.Next = node
		if node.Next != nil {
			heap.Push(pq, node.Next)
		}
		p = p.Next
	}
	return dummy.Next
}

type minHeap []*ListNode

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].Val < h[j].Val } // 最小堆
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(v any)        { *h = append(*h, v.(*ListNode)) }
func (h *minHeap) Pop() any          { a := *h; v := a[len(a)-1]; *h = a[:len(a)-1]; return v }

/**********************************************************************************************************************/
// https://leetcode.cn/problems/kth-node-from-end-of-list-lcci/
// 考察双指针，一般涉及到数组或者链表都可以考虑下是否能用双指针解决
// 输入： 1->2->3->4->5 和 k = 2
// 输出： 4
// kthToLast 返回倒数第 k 个节点
func kthToLast(head *ListNode, k int) int {
	if head == nil {
		return 0
	}
	p1, p2 := head, head
	for i := 0; i < k; i++ {
		if p1 != nil {
			p1 = p1.Next
		}
	}
	for p1 != nil {
		p1 = p1.Next
		p2 = p2.Next
	}
	return p2.Val
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/remove-nth-node-from-end-of-list/description/
// 考察双指针，设置虚拟头节点dummy，可以减少边界判断
// 输入：head = [1,2,3,4,5], n = 2
// 输出：[1,2,3,5]
// removeNthFromEnd 删除链表的倒数第 N 个结点
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Val: -1}
	dummy.Next = head

	p1, p2 := dummy, dummy
	for i := 0; i <= n; i++ {
		if p1 != nil {
			p1 = p1.Next
		}
	}

	for p1 != nil {
		p1 = p1.Next
		p2 = p2.Next
	}

	delNode := p2.Next
	p2.Next = delNode.Next
	return dummy.Next
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/middle-of-the-linked-list/
// 考察双指针-快慢指针
// 输入：head = [1,2,3,4,5]
// 输出：[3,4,5]
// 解释：链表只有一个中间结点，值为 3 。
// middleNode 链表的中间结点
func middleNode(head *ListNode) *ListNode {
	fast, slow := head, head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}
	return slow
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/linked-list-cycle-ii/description/
// 考察快慢指针
// 输入：head = [3,2,0,-4], pos = 1
// 输出：返回索引为 1 的链表节点
// 解释：链表中有一个环，其尾部连接到第二个节点。
// detectCycle 环形链表 II（环的起点）
func detectCycle(head *ListNode) *ListNode {
	fast, slow := head, head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
		if fast == slow {
			break
		}
	}
	//fast遇到了空指针，说明没有环
	if fast == nil || fast.Next == nil {
		return nil
	}
	fast = head
	for fast != slow {
		fast = fast.Next
		slow = slow.Next
	}
	return slow
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/intersection-of-two-linked-lists/
// 给你两个单链表的头节点 headA 和 headB ，请找出并返回两个单链表相交的起始节点。如果两个链表不存在相交节点，返回 null 。
// getIntersectionNode 两个链表是否相交
func getIntersectionNode(headA, headB *ListNode) *ListNode {
	p1, p2 := headA, headB
	for p1 != p2 {
		if p1 == nil {
			p1 = headB
		} else {
			p1 = p1.Next
		}
		if p2 == nil {
			p2 = headA
		} else {
			p2 = p2.Next
		}
	}
	return p1
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/remove-duplicates-from-sorted-array/description/
// 考察快慢指针，一般原地操作/零拷贝等问题大概率都是双指针问题
// 输入：numbers = [0,0,1,1,1,2,2,3,3,4]
// 输出：5, numbers = [0,1,2,3,4]
// 解释：函数应该返回新的长度 5 ， 并且原数组 numbers 的前五个元素被修改为 0, 1, 2, 3, 4 。不需要考虑数组中超出新长度后面的元素。
// removeDuplicates 删除有序数组中的重复项
func removeDuplicates(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}
	fast, slow := 0, 0
	for fast < len(numbers) {
		if numbers[slow] != numbers[fast] {
			slow++
			numbers[slow] = numbers[fast]
		}
		fast++
	}
	return slow + 1
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/remove-duplicates-from-sorted-list/description/
// 考察快慢指针，跟上面👆道题一样，一个是数组，一个是链表而已
// 输入：head = [1,1,2,3,3]
// 输出：[1,2,3]
// deleteDuplicates  删除排序链表中的重复元素
func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	fast, slow := head, head
	for fast != nil {
		if fast.Val != slow.Val {
			slow.Next = fast
			slow = slow.Next
		}
		fast = fast.Next
	}
	slow.Next = nil //断开与后面重复原因的连接
	return head
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/remove-element/description/
// 考察双指针，一般原地操作/零拷贝等问题大概率都是双指针问题
// 输入：numbers = [3,2,2,3], val = 3
// 输出：2, numbers = [2,2]
// 输入：numbers = [0,1,2,2,3,0,4,2], val = 2
// 输出：5, numbers = [0,1,3,0,4]
// removeElement 移除元素
func removeElement(numbers []int, val int) int {
	fast, slow := 0, 0
	for fast < len(numbers) {
		if numbers[fast] != val {
			numbers[slow] = numbers[fast]
			slow++
		}
		fast++
	}
	return slow
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/two-sum-ii-input-array-is-sorted/
// 考察双指针-左右指针
// 输入：numbers = [2,7,11,15], target = 9
// 输出：[1,2]
// twoSum 两数之和 II - 输入有序数组
func twoSum(numbers []int, target int) []int {
	start, end := 0, len(numbers)-1
	for start < end {
		if numbers[start]+numbers[end] == target {
			return []int{start + 1, end + 1}
		} else if numbers[start]+numbers[end] < target {
			start++
		} else if numbers[start]+numbers[end] > target {
			end--
		}
	}
	return nil
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/longest-palindromic-substring/description/
// 这个也是考察双指针问题，palindrome方法解决了中心奇/偶数的问题
// 输入：s = "babad"
// 输出："bab"
// 解释："aba" 同样是符合题意的答案。
// longestPalindrome 最长回文子串
func longestPalindrome(s string) string {
	palindrome := func(s string, left, right int) string {
		for left >= 0 && right < len(s) && s[left] == s[right] {
			left--
			right++
		}
		return s[left+1 : right]
	}
	maxPalindrome := ""
	for i := 0; i < len(s); i++ {
		str1 := palindrome(s, i, i)
		if len(str1) > len(maxPalindrome) {
			maxPalindrome = str1
		}

		str2 := palindrome(s, i, i+1)
		if len(str2) > len(maxPalindrome) {
			maxPalindrome = str2
		}
	}
	return maxPalindrome
}

/**********************************************************************************************************************/
//数据结构题型-前缀和数组and差分数据
/**********************************************************************************************************************/
// https://leetcode.cn/problems/range-sum-query-immutable/description/
// 典型的前缀和数组技巧
// SumRange 区域和检索 - 数组不可变
func (this *NumArray) SumRange(left int, right int) int {
	return this.preNums[right+1] - this.preNums[left]
}

type NumArray struct {
	preNums []int
}

func Constructor(numbers []int) NumArray {
	na := NumArray{preNums: make([]int, len(numbers)+1)}
	for i := 1; i <= len(numbers); i++ {
		na.preNums[i] = na.preNums[i-1] + numbers[i-1]
	}
	return na
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/range-sum-query-2d-immutable/
// 也是典型的前缀和数组技巧，通过空间换时间
// 这个一定得画个图出来，不然不好想象
// _SumRegion二维区域和检索 - 矩阵不可变
func (this *NumMatrix) _SumRegion(row1 int, col1 int, row2 int, col2 int) int {
	return this.preSum[row2+1][col2+1] - this.preSum[row2+1][col1] - this.preSum[row1][col2+1] + this.preSum[row1][col1]
}

type NumMatrix struct {
	preSum [][]int
}

func _Constructor(matrix [][]int) NumMatrix {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return NumMatrix{}
	}

	//一定要初始化多一行/多一列
	nm := NumMatrix{make([][]int, len(matrix)+1)}
	for i := 0; i <= len(matrix); i++ {
		nm.preSum[i] = make([]int, len(matrix[0])+1)
	}

	for i := 1; i <= len(matrix); i++ {
		for j := 1; j <= len(matrix[0]); j++ {
			nm.preSum[i][j] = nm.preSum[i-1][j] + nm.preSum[i][j-1] + matrix[i-1][j-1] - nm.preSum[i-1][j-1]
		}
	}
	return nm
}

/**********************************************************************************************************************/
// 考察差分数组
// getModifiedArray 区间加法
func getModifiedArray(length int, updates [][]int) []int {
	diff := make([]int, length)
	for _, update := range updates {
		leftIndex := update[0]
		rightIndex := update[1]
		value := update[2]

		diff[leftIndex] += value
		//小心数组越界
		if rightIndex+1 < length {
			diff[rightIndex+1] -= value
		}
	}

	result := make([]int, length)
	result[0] = diff[0]

	for i := 1; i < length; i++ {
		result[i] = result[i-1] + diff[i]
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/corporate-flight-bookings/description/
// 考察差分数组
// 输入：bookings = [[1,2,10],[2,3,20],[2,5,25]], n = 5
// 输出：[10,55,45,25,25]
// 解释：
// 航班编号        1   2   3   4   5
// 预订记录 1 ：   10  10
// 预订记录 2 ：       20  20
// 预订记录 3 ：       25  25  25  25
// 总座位数：      10  55  45  25  25
// 因此，answer = [10,55,45,25,25]
// corpFlightBookings 航班预订统计
func corpFlightBookings(bookings [][]int, n int) []int {
	diff := make([]int, n)
	for _, booking := range bookings {
		start := booking[0] - 1
		end := booking[1] - 1
		val := booking[2]

		diff[start] += val
		if end+1 < n {
			diff[end+1] -= val
		}
	}

	result := make([]int, n)
	result[0] = diff[0]

	for i := 1; i < n; i++ {
		result[i] = result[i-1] + diff[i]
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/car-pooling/description/
// 考察差分数组 - 关键是构造diff数组
// 输入：trips = [[2,1,5],[3,3,7]], capacity = 4
// 输出：false
// 输入：trips = [[2,1,5],[3,3,7]], capacity = 5
// 输出：true
// carPooling 拼车，没有告诉有多少个站点，但是告诉了范围，直接取最大
func carPooling(trips [][]int, capacity int) bool {
	diff := make([]int, 1000)
	for _, trip := range trips {
		diff[trip[1]] += trip[0]
		if trip[2]+1 < 1000 {
			diff[trip[2]] -= trip[0]
		}
	}

	result := make([]int, 1000)
	result[0] = diff[0]
	if result[0] > capacity {
		return false
	}

	for i := 1; i < 1000; i++ {
		result[i] = result[i-1] + diff[i]
		if result[i] > capacity {
			return false
		}
	}
	return true
}

/**********************************************************************************************************************/
//https://leetcode.cn/problems/lru-cache/description/
// 核心考察：
// O(1)时间复杂的的“查询”要求 => map
// O(1)时间复杂的的“插入/删除”要求 => 双向链表
// 写代码可以考虑分层实现，这样代码逻辑清晰易懂，先定义出方法，然后在一个个实现
// LRUCache 实现 LRU (最近最少使用) 缓存
type LRUCache struct {
	capacity int
	hMap     map[int]*Node //哈希表- 关键数据结构
	dList    *DoubleList   //双链表- 关键数据结构
}

type Node struct {
	key, val   int
	prev, next *Node
}

type DoubleList struct {
	length     int
	head, tail *Node
}

// head->n->x->y->tail
// head<-n<-x<-y<-tail
// 在双链表中移除一个节点（重新点开后台应用的场景）
func (d *DoubleList) remove(x *Node) {
	x.prev.next = x.next
	x.next.prev = x.prev
	d.length--
}

// head->n->x->y->tail
// head<-n<-x<-y<-tail
// 在双链表尾部添加一个节点（新打开一个应用的场景）
func (d *DoubleList) addLast(x *Node) *Node {
	x.prev = d.tail.prev
	x.next = d.tail

	d.tail.prev.next = x
	d.tail.prev = x
	d.length++
	return x
}

// 删除双链表的第一个节点（到达最大存储容量时，需要删除最久未使用的）
func (d *DoubleList) removeFirst() *Node {
	if d.head.next == nil {
		return nil
	}
	first := d.head.next
	d.remove(first)
	return first
}

func (d *DoubleList) Size() int { return d.length }

func Constructor_(capacity int) LRUCache {
	lru := LRUCache{}
	lru.hMap = make(map[int]*Node, capacity)

	head, tail := &Node{}, &Node{}
	head.next, tail.prev = tail, head

	lru.capacity = capacity
	lru.dList = &DoubleList{length: 0, head: head, tail: tail}
	return lru
}

func (this *LRUCache) Get(key int) int {
	x, ok := this.hMap[key]
	if !ok {
		return -1
	}
	this.makeRecently(key)
	return x.val
}

func (this *LRUCache) Put(key int, value int) {
	x, ok := this.hMap[key]
	if ok {
		this.dList.remove(x)
	}
	if this.dList.Size() >= this.capacity {
		this.removeLeastRecently()
	}
	this.addRecently(key, value)
}

// 提升节点为最近使用
func (this *LRUCache) makeRecently(key int) {
	x := this.hMap[key]
	if x == nil {
		return
	}
	this.dList.remove(x)
	this.dList.addLast(x)
}

// 添加一个最近节点到lur中
func (this *LRUCache) addRecently(key, val int) {
	x := &Node{key: key, val: val}
	x = this.dList.addLast(x)
	this.hMap[key] = x
}

// 删除最久未使用的元素
func (this *LRUCache) removeLeastRecently() {
	if this.dList.Size() == 0 {
		return
	}
	x := this.dList.removeFirst()
	delete(this.hMap, x.key)
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/insert-delete-getrandom-o1/
// 考察对数据结构优缺点的理解，合理利用多个数据结构的优点来实现解决问题
// 数组可以在O(1)时间等概率取出元素
// 如果只在数组尾部插入/删除，时间复杂的则是O(1)
// RandomizedSet O(1) 时间插入、删除和获取随机元素
type RandomizedSet struct {
	numbers []int       // O(1)时间等概率取出元素，只能是数组
	mapping map[int]int // Val->Index 要让数组在O(1)时间插入/删除，那么需要额外借助哈希表
}

func _Constructor_() RandomizedSet {
	return RandomizedSet{numbers: make([]int, 0), mapping: map[int]int{}}
}

func (this *RandomizedSet) Insert(val int) bool {
	if _, ok := this.mapping[val]; ok {
		return false
	}
	this.numbers = append(this.numbers, val)
	this.mapping[val] = len(this.numbers) - 1
	return true
}

func (this *RandomizedSet) Remove(val int) bool {
	if _, ok := this.mapping[val]; !ok {
		return false
	}
	index := this.mapping[val]
	delete(this.mapping, val)

	size := len(this.numbers) - 1
	this.numbers[index] = this.numbers[size]
	this.mapping[this.numbers[index]] = index
	this.numbers = this.numbers[:size]
	return true
}

func (this *RandomizedSet) GetRandom() int {
	return this.numbers[rand.Intn(len(this.numbers))]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/random-pick-with-blacklist/description/
// 这个主要还是考察对基本数据结构的掌握，算法不复杂但是比较巧妙
// 实现 Solution 类:
// Solution(int n, int[] blacklist) 初始化整数 n 和被加入黑名单 blacklist 的整数
// int pick() 返回一个范围为 [0, n - 1] 且不在黑名单 blacklist 中的随机整数
// Solution.Pick 黑名单中的随机数
type Solution struct {
	sz      int         //有效元素区间[0:n-len(blacklist)]
	mapping map[int]int //black->index 黑名单到非黑名单元素的索引映射
}

func Constructor1(n int, blacklist []int) Solution {
	s := Solution{}
	s.mapping = map[int]int{}
	for _, black := range blacklist {
		s.mapping[black] = -1
	}

	s.sz = n - len(blacklist)
	last := n - 1

	for _, black := range blacklist {
		//黑名单已经在【SZ:N】区间就不用管了
		if black >= s.sz {
			continue
		}
		//避免last也在黑名单里面
		for {
			if _, ok := s.mapping[last]; !ok {
				break
			}
			last--
		}
		s.mapping[black] = last
		last--
	}
	return s
}

func (this *Solution) Pick() int {
	ele := rand.Intn(math.MaxInt) % this.sz
	if index, ok := this.mapping[ele]; ok {
		return index
	}
	return ele
}

/***********************************************************************************************************************
// 单调栈-算法框架：
int [] nextGreaterElement(int[] numbers) {
	int n = numbers.length;
	// 存放答案的数组
	int[] res = new int[n]
	Stack<Integer> s = new Stack<>();
	// 倒着往栈里面放
	for (int i= n-1; i>=0; i--) {
		// 判断个子高矮
		while (!s.isEmpty() && s.peek() <= numbers[i]) {
			// 矮个子起开，反正也被挡着了
			s.pop()
		}
		// numbers[i] 身后的 next great number
		res[i] = s.isEmpty() ? -1: s.peek();
		s.push(numbers[i])
    }
	return res
}
/**********************************************************************************************************************/
// https://leetcode.cn/problems/daily-temperatures/
// 考察单调栈
// 输入: temperatures = [73,74,75,71,69,72,76,73]
// 输出: [1,1,4,2,1,1,0,0]
// 输入：temperatures = [89,62,70,58,47,47,46,76,100,70]
// 输出：[8,1,5,4,3,2,1,1,0,0]
// dailyTemperatures 每日温度(类似下一个更大的元素)
func dailyTemperatures(temperatures []int) []int {
	result := make([]int, len(temperatures)) //结果数组
	stack := []int{}                         //单调栈-记录的时索引，不是元素

	for i := len(temperatures) - 1; i >= 0; i-- {
		for len(stack) > 0 && temperatures[i] >= temperatures[stack[len(stack)-1]] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			//没有更大的元素，默认为0
		} else {
			result[i] = stack[len(stack)-1] - i
		}
		stack = append(stack, i)
	}
	return result
}

/**********************************************************************************************************************/
//https://leetcode.cn/problems/next-greater-element-i/
// 考察单调栈
// 找出满足nums1[i]==nums2[j]的下标j ，并且在nums2确定nums2[j]的下一个更大元素 。如果不存在下一个更大元素，那么本次查询的答案是-1 。
// 输入：nums1 = [2,4], nums2 = [1,2,3,4].
// 输出：[3,-1]
// 解释：nums1 中每个值的下一个更大元素如下所述：
// - 2 ，用加粗斜体标识，nums2 = [1,2,3,4]。下一个更大元素是 3 。
// - 4 ，用加粗斜体标识，nums2 = [1,2,3,4]。不存在下一个更大元素，所以答案是 -1 。
// nextGreaterElement 下一个更大元素 I
func nextGreaterElement(nums1 []int, nums2 []int) []int {
	stack := []int{}         //单调栈-记录元素
	mapping := map[int]int{} //记录nums2的下一个最大元素

	for j := len(nums2) - 1; j >= 0; j-- {
		for len(stack) > 0 && nums2[j] >= stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			mapping[nums2[j]] = -1
		} else {
			mapping[nums2[j]] = stack[len(stack)-1]
		}
		stack = append(stack, nums2[j])
	}

	result := make([]int, len(nums1)) //结果数组，在nums2中查表
	for i, num := range nums1 {
		result[i] = mapping[num]
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/next-greater-element-ii/description/
// 考察单调栈
// 循环数组技巧：数组长度加倍，进行下标判断时i := index % arr.len（[1,2,3,4] => [1,2,3,4,1,2,3,4]）
// 输入: numbers = [1,2,3,4,3]
// 输出: [2,3,4,-1,4]
// nextGreaterElements 下一个更大元素 II，在循环数组中查询
func nextGreaterElements(numbers []int) []int {
	result := make([]int, len(numbers))
	stack := []int{}
	n := len(numbers)
	for i := 2*n - 1; i >= 0; i-- {
		for len(stack) > 0 && numbers[i%n] >= stack[len(stack)-1] {
			// pop，小个子走开
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			result[i%n] = -1
		} else {
			result[i%n] = stack[len(stack)-1]
		}
		stack = append(stack, numbers[i%n])
	}
	return result
}

/***********************************************************************************************************************
// 单调队列-算法框架：
int[] maxSlidingWindow(int[] numbers, int k){
	MonotonicQueue window = new MonotonicQueue();
	List<Integer> res = new ArrayList<>();

	for (int i=0; i<numbers.Length(); i++) {
		if (i < k-1) {
			// 先把窗口的前k-1个位置填满
			window.push(numbers[i)
		}else {
			// 窗口开始向前移动，移入新元素
			window.push(numbers[i]);

			// 将当前窗口的最大元素记入结果
			res.add(window.max());

			// 移除最后的元素
			window.pop(numbers[i-k+1])
		}
	}
	return res
}
/**********************************************************************************************************************/
// https://leetcode.cn/problems/sliding-window-maximum/description/
// 单调队列+数据结构设计
// 输入：numbers = [1,3,-1,-3,5,3,6,7], k = 3
// 输出：[3,3,5,5,6,7]
// maxSlidingWindow 滑动窗口最大值
func maxSlidingWindow(numbers []int, k int) []int {
	queue := []int{}
	push := func(i int) {
		//这里是算法的关键，按道理是超过k，入队时才del队头，但其实保留最大的就好了
		//入队的元素比队列最后一个大，就一直踢掉最后一个（窗口滑动，把小的干掉，留下最大的）
		for len(queue) > 0 && numbers[i] >= numbers[queue[len(queue)-1]] {
			queue = queue[:len(queue)-1]
		}
		queue = append(queue, i)
	}

	result := []int{}
	for i := 0; i < len(numbers); i++ {
		if i < k-1 {
			push(i)
		} else {
			// 窗口开始向前移动，移入新元素的下标（因为保存元素的话，元素可能本身会重复）
			push(i)

			//这里是为了防止队列要超过固定大小 k
			for queue[0] <= i-k {
				queue = queue[1:]
			}
			//记入结果
			result = append(result, numbers[queue[0]])
		}
	}
	return result
}

/**********************************************************************************************************************/
// 下面开始二叉树系列问题，二叉树解题的思维模式：
//  1、是否可以通过遍历一遍二叉树得到答案？
//  2、是否可以定义一个递归函数，通过子问题(子树)的答案推导出愿问题的答案？
// 如果单独抽出一个二叉树节点，需要对它做什么事情？需要在什么时候（前/中/后序位置）做？
/**********************************************************************************************************************/
// https://leetcode.cn/problems/invert-binary-tree/
// invertTree 翻转二叉树（通过“遍历”的思维模式解决）
func _invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	invertTree(root.Left)
	invertTree(root.Right)
	root.Left, root.Right = root.Right, root.Left
	return root
}

// invertTree 翻转二叉树（通过“分解问题”的思维模式解决）
func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	left := invertTree(root.Left)
	right := invertTree(root.Right)
	root.Left, root.Right = right, left
	return root
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/populating-next-right-pointers-in-each-node/description/
// 通过“遍历”的思维模式解决，把相邻非同一个父节点的看成一个整体，二叉树=>三叉树
// connect 填充每个节点的下一个右侧节点指针
func connect(root *BinaryNode) *BinaryNode {
	if root == nil {
		return nil
	}
	var traverse func(node1, node2 *BinaryNode)
	traverse = func(node1, node2 *BinaryNode) {
		if node1 == nil || node2 == nil {
			return
		}
		node1.Next = node2
		traverse(node1.Left, node1.Right)
		traverse(node2.Left, node2.Right)
		traverse(node1.Right, node2.Left)
	}
	traverse(root.Left, root.Right)
	return root
}

type BinaryNode struct {
	Val   int
	Left  *BinaryNode
	Right *BinaryNode
	Next  *BinaryNode
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/flatten-binary-tree-to-linked-list/description/
// 通过“分解问题”的思维模式解决
// flatten 二叉树展开为链表
func flatten(root *TreeNode) {
	if root == nil {
		return
	}
	flatten(root.Left)
	flatten(root.Right)

	// 1、左右子树已经被拉成链表
	left := root.Left
	right := root.Right

	// 2、将左子树作为右子树
	root.Left = nil
	root.Right = left

	// 3、将原先的右子树接到当前右子树的末端
	p := root
	for p.Right != nil {
		p = p.Right
	}
	p.Right = right
}

/**********************************************************************************************************************/
//https://leetcode.cn/problems/maximum-binary-tree/description/
// 思路：通过“分解问题”的思维模式解决，只需要关注一个节点该怎么操作，其他节点交给递归函数
// 二叉树构造问题，一般通过“分解问题”的思维模式解决
// constructMaximumBinaryTree 构造最大二叉树
func constructMaximumBinaryTree(numbers []int) *TreeNode {
	var build func(numbers []int, lo, hi int) *TreeNode

	build = func(numbers []int, lo, hi int) *TreeNode {
		if lo > hi {
			return nil
		}
		maxNum, index := math.MinInt, -1
		for i := lo; i <= hi; i++ {
			if numbers[i] > maxNum {
				maxNum = numbers[i]
				index = i
			}
		}
		root := &TreeNode{Val: maxNum}
		root.Left = build(numbers, lo, index-1)
		root.Right = build(numbers, index+1, hi)
		return root
	}

	return build(numbers, 0, len(numbers)-1)
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/construct-binary-tree-from-inorder-and-postorder-traversal/
// 思路：通过“分解问题”的思维模式解决，只需要关注一个节点该怎么操作，其他节点交给递归函数
// 二叉树构造问题，一般通过“分解问题”的思维模式解决
// buildTree 从中序与后序遍历序列构造二叉树
func _buildTree(inorder []int, postorder []int) *TreeNode {
	if len(inorder) == 0 {
		return nil
	}
	mapping := map[int]int{}
	for i, val := range inorder {
		mapping[val] = i
	}
	return _build(mapping, inorder, 0, len(inorder)-1, postorder, 0, len(postorder)-1)
}

func _build(mapping map[int]int, inorder []int, inStart, inEnd int, postorder []int, postStart, postEnd int) *TreeNode {
	if inStart > inEnd {
		return nil
	}

	rootVal := postorder[postEnd]
	index := mapping[rootVal]

	leftSize := index - inStart

	root := &TreeNode{Val: rootVal}
	root.Left = _build(mapping, inorder, inStart, index-1, postorder, postStart, postStart+leftSize-1)
	root.Right = _build(mapping, inorder, index+1, inEnd, postorder, postStart+leftSize, postEnd-1)
	return root
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/construct-binary-tree-from-preorder-and-postorder-traversal/
// 二叉树构造问题，一般通过“分解问题”的思维模式解决
// constructFromPrePost 根据前序和后序遍历构造二叉树（答案不唯一）
func constructFromPrePost(preorder []int, postorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}
	mapping := map[int]int{}
	for i, val := range postorder {
		mapping[val] = i
	}
	return build_(mapping, preorder, 0, len(postorder)-1, postorder, 0, len(postorder)-1)
}

func build_(mapping map[int]int, preorder []int, preStart, preEnd int, postorder []int, postStart, postEnd int) *TreeNode {
	if preStart > preEnd {
		return nil
	}
	if preStart == preEnd {
		return &TreeNode{Val: preorder[preStart]}
	}

	rootVal := preorder[preStart]
	leftRootVal := preorder[preStart+1] //假设前序遍历第二个节点为root的左子树的根节点

	index := mapping[leftRootVal]
	leftSize := index - postStart + 1

	root := &TreeNode{Val: rootVal}
	root.Left = build_(mapping, preorder, preStart+1, preStart+leftSize, postorder, postStart, index)
	root.Right = build_(mapping, preorder, preStart+leftSize+1, preEnd, postorder, index+1, postEnd-1)
	return root
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/serialize-and-deserialize-binary-tree/
// 输入：root = [1,2,3,null,null,4,5]
// 输出：[1,2,3,null,null,4,5]
// 二叉树的序列化与反序列化
type Codec struct{}

func Constructor_1() Codec { return Codec{} }

// 序列化-前序遍历方式
func (this *Codec) serialize(root *TreeNode) string {
	if root == nil {
		return ""
	}

	var traverse func(root *TreeNode)
	result := []string{}

	traverse = func(root *TreeNode) {
		if root == nil {
			result = append(result, "null")
			return
		}
		result = append(result, strconv.Itoa(root.Val))
		traverse(root.Left)
		traverse(root.Right)
	}
	traverse(root)

	str := ""
	for _, nodeVal := range result {
		if str == "" {
			str = nodeVal
		} else {
			str = str + "," + nodeVal
		}
	}
	return str
}

// 反序列化-前序遍历方式
func (this *Codec) deserialize(data string) *TreeNode {
	if len(data) == 0 {
		return nil
	}

	var build func() *TreeNode

	//注意nodes不要通过参数传递进去
	nodes := strings.Split(data, ",")

	build = func() *TreeNode {
		if len(nodes) == 0 {
			return nil
		}
		rootVal := nodes[0]
		nodes = nodes[1:]

		if rootVal == "null" {
			return nil
		}

		val, _ := strconv.Atoi(rootVal)
		root := &TreeNode{Val: val, Left: build(), Right: build()}
		return root
	}
	return build()
}

func (this *Codec) String(root *TreeNode) string {
	var traverse func(root *TreeNode) string

	traverse = func(root *TreeNode) string {
		if root == nil {
			return "null"
		}
		rootVal := root.Val
		leftVal := traverse(root.Left)
		rightVal := traverse(root.Right)
		return fmt.Sprintf("%v,%v,%v", rootVal, leftVal, rightVal)
	}
	return traverse(root)
}

// 序列化-层序遍历方式
func (this *Codec) serializeByLevel(root *TreeNode) string {
	if root == nil {
		return ""
	}
	queue := []*TreeNode{}
	queue = append(queue, root)

	result := []string{}
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]
			if node == nil {
				result = append(result, "null")
				continue
			}
			result = append(result, strconv.Itoa(node.Val))
			queue = append(queue, node.Left)
			queue = append(queue, node.Right)
		}

	}
	str := ""
	for _, nodeVal := range result {
		if str == "" {
			str = nodeVal
		} else {
			str = str + "," + nodeVal
		}
	}
	return str
}

// 反序列化-层序遍历方式
func (this *Codec) deserializeByLevel(data string) *TreeNode {
	if len(data) == 0 {
		return nil
	}
	nodes := strings.Split(data, ",")

	index := 0
	rootVal, _ := strconv.Atoi(nodes[index])
	index++

	queue := []*TreeNode{}
	root := &TreeNode{Val: rootVal}
	queue = append(queue, root)

	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]

			left := nodes[index]
			index++
			if left != "null" {
				leftVal, _ := strconv.Atoi(left)
				node.Left = &TreeNode{Val: leftVal}
				queue = append(queue, node.Left)
			}

			right := nodes[index]
			index++
			if right != "null" {
				rightVal, _ := strconv.Atoi(right)
				node.Right = &TreeNode{Val: rightVal}
				queue = append(queue, node.Right)
			}
		}
	}
	return root
}

// a := Constructor_1()
//	tree := a.deserialize("1,2,null,null,3,4,null,null,5,null,null")
//	fmt.Println(a.String(tree))
//	fmt.Println(a.serialize(tree))

/**********************************************************************************************************************/
// https://leetcode.cn/problems/sort-an-array/description/
// 归并排序解决--二叉树的后序遍历
// mergeSort 排序数组
func mergeSort(numbers []int) []int {
	var sort func(numbers []int, lo, hi int)
	var merge func(numbers []int, lo, mid, hi int)

	temp := make([]int, len(numbers)) //避免在递归函数里面频繁创建销毁

	merge = func(numbers []int, lo, mid, hi int) {
		for i := lo; i <= hi; i++ {
			temp[i] = numbers[i]
		}

		i, j := lo, mid+1
		for p := lo; p <= hi; p++ {
			if i == mid+1 { //左边遍历完了
				numbers[p] = temp[j]
				j++
			} else if j == hi+1 { //右边遍历完了
				numbers[p] = temp[i]
				i++
			} else if temp[i] > temp[j] {
				numbers[p] = temp[j]
				j++
			} else if temp[i] <= temp[j] {
				numbers[p] = temp[i]
				i++
			}
		}
	}
	sort = func(numbers []int, lo, hi int) {
		if lo == hi {
			return
		}
		mid := lo + (hi-lo)/2
		sort(numbers, lo, mid)
		sort(numbers, mid+1, hi)
		merge(numbers, lo, mid, hi)
	}
	sort(numbers, 0, len(numbers)-1)
	return numbers
}

// 快速排序解决--二叉树的前序遍历
// quickSort 快排方式
func quickSort(numbers []int) []int {
	var sort func(numbers []int, lo, hi int)
	sort = func(numbers []int, lo, hi int) {
		if lo >= hi {
			return
		}
		p := partition(numbers, lo, hi)
		sort(numbers, lo, p-1)
		sort(numbers, p+1, hi)
	}
	sort(numbers, 0, len(numbers)-1)
	return numbers
}

func partition(numbers []int, lo, hi int) int {
	pivot := numbers[lo]
	i, j := lo+1, hi

	for i <= j {
		for i < hi && numbers[i] <= pivot {
			i++
		}
		for j > lo && numbers[j] > pivot {
			j--
		}
		if i >= j {
			break
		}
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	numbers[lo], numbers[j] = numbers[j], numbers[lo]
	return j
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/count-of-smaller-numbers-after-self/description/
// 题目描述：返回 countArr[i] = COUNT(*) WHERE j > i AND numbers[j] < numbers[i]
// 思路：归并排序时记录答案
// 原因：我们在使用merge函数合并两个有序数组的时候，其实是可以知道一个元素nums[i]后面有多个元素比nums[i]小的；
// 即，当nums[p]=temp[i]是，[mid+1，j) 则为 temp[i]后面比它小的元素，也就是 j-(mid+1)
// countSmaller 计算右侧小于当前元素的个数
func countSmaller(numbers []int) []int {
	type pair struct {
		val   int
		index int
	}

	var sort func(numbers []pair, lo, hi int)
	var merge func(numbers []pair, lo, mid, hi int)

	list := make([]pair, len(numbers))
	for i := 0; i < len(numbers); i++ {
		list[i] = pair{val: numbers[i], index: i}
	}

	temp := make([]pair, len(numbers)) //避免在递归函数里面频繁创建销毁
	countArr := make([]int, len(numbers))

	merge = func(numbers []pair, lo, mid, hi int) {
		for i := lo; i <= hi; i++ {
			temp[i] = numbers[i]
		}
		i, j := lo, mid+1
		for p := lo; p <= hi; p++ {
			if i == mid+1 { //左边遍历完了
				numbers[p] = temp[j]
				j++
			} else if j == hi+1 { //右边遍历完了
				numbers[p] = temp[i]
				i++
				countArr[numbers[p].index] += j - mid - 1 //本题的关键在这里
			} else if temp[i].val > temp[j].val {
				numbers[p] = temp[j]
				j++
			} else if temp[i].val <= temp[j].val {
				numbers[p] = temp[i]
				i++
				countArr[numbers[p].index] += j - mid - 1 //本题的关键在这里
			}
		}
	}
	sort = func(numbers []pair, lo, hi int) {
		if lo == hi {
			return
		}
		mid := lo + (hi-lo)/2
		sort(numbers, lo, mid)
		sort(numbers, mid+1, hi)
		merge(numbers, lo, mid, hi)
	}
	sort(list, 0, len(numbers)-1)
	return countArr
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/shu-zu-zhong-de-ni-xu-dui-lcof/
// 本题跟上面👆道思路一模一样，只是返回值不同而已
// 递归思路求解
// reversePairs 交易逆序对的总数
func reversePairs(record []int) int {
	var sort func(numbers []int, lo, hi int)
	var merge func(numbers []int, lo, mid, hi int)

	temp := make([]int, len(record)) //避免在递归函数里面频繁创建销毁
	result := [][]int{}
	count := 0

	sort = func(numbers []int, lo, hi int) {
		if lo == hi {
			return
		}
		mid := lo + (hi-lo)/2
		sort(numbers, lo, mid)
		sort(numbers, mid+1, hi)
		merge(numbers, lo, mid, hi)
	}
	merge = func(numbers []int, lo, mid, hi int) {
		for i := lo; i <= hi; i++ {
			temp[i] = numbers[i]
		}

		i, j := lo, mid+1
		for p := lo; p <= hi; p++ {
			if i == mid+1 { //左边遍历完了
				numbers[p] = temp[j]
				j++
			} else if j == hi+1 { //右边遍历完了
				numbers[p] = temp[i]
				i++

				//本题的关键在这里：[mid+1，j)
				for k := mid + 1; k < j; k++ {
					result = append(result, []int{numbers[p], temp[k]})
				}
				count += j - mid - 1
			} else if temp[i] > temp[j] {
				numbers[p] = temp[j]
				j++
			} else {
				numbers[p] = temp[i]
				i++

				//本题的关键在这里：当 temp[i] < temp[j] 时，[mid+1，j)范围都是比temp[i]小的元素
				for k := mid + 1; k < j; k++ {
					result = append(result, []int{numbers[p], temp[k]})
				}
				count += j - mid - 1
			}
		}
	}
	sort(record, 0, len(record)-1)
	fmt.Println(result)
	return count
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/convert-bst-to-greater-tree/description/
// BST的特性，中序遍历是升序排列，如果把左右子树顺序换一下，再中序遍历就是倒序排列了
// 输入：[4,1,6,0,2,5,7,null,null,null,3,null,null,null,8]
// 输出：[30,36,21,36,35,26,15,null,null,null,33,null,null,null,8]
// convertBST 把二叉搜索树转换为累加树
func convertBST(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}

	var traverse func(root *TreeNode)
	sum := 0

	traverse = func(root *TreeNode) {
		if root == nil {
			return
		}
		traverse(root.Right)
		sum += root.Val
		root.Val = sum
		traverse(root.Left)
	}
	traverse(root)
	return root
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/validate-binary-search-tree/description/
// 二叉搜索树的特性：左小右大，中序遍历是升序
// 思路：使用辅助函数，增加函数参数列表，在参数中携带额外的信息
// isValidBST 验证二叉搜索树
func isValidBST(root *TreeNode) bool {
	var helper func(root *TreeNode, lower, upper int) bool

	helper = func(root *TreeNode, lower, upper int) bool {
		if root == nil {
			return true
		}
		if root.Val <= lower || root.Val >= upper {
			return false
		}
		return helper(root.Left, lower, root.Val) && helper(root.Right, root.Val, upper)
	}
	return helper(root, math.MinInt64, math.MaxInt64)
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/kth-largest-element-in-an-array/
// 快速排序的变体——快速选择，有点像二分查找，但是原始数组没有排序，所以通过快排的思路，先排个序，在折半查找
// findKthLargest 数组中的第K个最大元素
func findKthLargest(numbers []int, k int) int {
	return quickSelect(numbers, 0, len(numbers)-1, len(numbers)-k)
}

func quickSelect(a []int, l, r, target int) int {
	q := _partition(a, l, r)
	if q == target {
		return a[q]
	}

	if q < target {
		return quickSelect(a, q+1, r, target)
	} else {
		return quickSelect(a, l, q-1, target)
	}
}

func _partition(numbers []int, lo, hi int) int {
	pivot := numbers[lo]
	i, j := lo+1, hi

	for i <= j {
		for i < hi && numbers[i] <= pivot {
			i++
		}
		for j > lo && numbers[j] > pivot {
			j--
		}
		if i >= j {
			break
		}
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	numbers[lo], numbers[j] = numbers[j], numbers[lo]
	return j
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/all-paths-from-source-to-target/
// 图的遍历算法
// 图的遍历，本质就是多叉树的遍历,只是图可能有环，所以需要visited数组去重，避免死循环
// allPathsSourceTarget 所有可能的路径
func allPathsSourceTarget(graph [][]int) [][]int {
	result := make([][]int, 0, len(graph))
	path := []int{}
	traverse(graph, 0, path, &result)
	return result
}

func traverse(graph [][]int, s int, path []int, result *[][]int) {
	path = append(path, s) //加入路径

	if s == len(graph)-1 { //到达终点
		temp := make([]int, len(path))
		copy(temp, path)
		*result = append(*result, temp)

		//提前结束一定记得删除路径
		path = path[:len(path)-1]
		return
	}
	for _, v := range graph[s] {
		traverse(graph, v, path, result)
	}

	path = path[:len(path)-1] //删除路径
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/satisfiability-of-equality-equations/
// 思路：UnionFind算法（并查集算法），将equations中的算式根据“==”和“!=”分成两部分，先处理“==”的，使其构成连通量，
// 然后处理“!=”算式，查看连通量是否被破坏。
// 输入：["a==b","b==c","a==c"]
// 输出：true
// 输入：["a==b","b!=c","c==a"]
// 输出：false
// equationsPossible 等式方程的可满足性
func equationsPossible(equations []string) bool {
	uf := NewUF(26)
	for _, str := range equations {
		if str[1] == '=' {
			uf.union(int(str[0]-'a'), int(str[3]-'a'))
		}
	}

	for _, str := range equations {
		if str[1] == '!' {
			if uf.connected(int(str[0]-'a'), int(str[3]-'a')) {
				return false
			}
		}
	}
	return true
}

// UnionFind算法（并查集算法）的实际应用案例：最小生成树，思路：Kruskal算法，本质其实就是树的合法性判断 + 权重排序逻辑
// 即： 将所有边按照权重从小到大排序，从权重小的开始遍历，如果这条边在mset中不会形成环，
// 那么这条边就是最小生成树的一部分，把它加入集合，否则，这条边这条边不是最小生成树的一部分，不要把它加入mset集合。
// 输入：{{1, 2, 10}, {5, 3, 15}, {2, 3, 20}, {3, 4, 5}, {1, 5, 25}}
// 输出：50
// minimumCost 最低成本连通所有城市
func minimumCost(n int, connections [][]int) int {
	sort.Slice(connections, func(i, j int) bool {
		return connections[i][2] < connections[j][2]
	})

	uf := NewUF(n + 1)
	mset := 0

	for _, edge := range connections {
		u, v, weight := edge[0], edge[1], edge[2]
		if uf.connected(u, v) {
			continue
		}
		uf.union(u, v)
		mset += weight
	}

	if uf.Size() == 2 {
		return mset
	} else {
		return -1
	}
}

// //////////////////////////////////////////////////////////////////////
// UnionFind算法（并查集算法），如果，p/q有相同的根节点，则说明p/q是连通的
type UF struct {
	parent []int
	size   int
}

func NewUF(n int) *UF {
	uf := UF{size: n, parent: make([]int, n)}
	for i := 0; i < n; i++ {
		uf.parent[i] = i
	}
	return &uf
}

func (uf *UF) union(p, q int) {
	rootP := uf.find(p)
	rootQ := uf.find(q)

	if rootP == rootQ {
		return
	}
	uf.size--
	uf.parent[rootP] = rootQ
}

func (uf *UF) find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x]) //路径压缩
	}
	return uf.parent[x]
}
func (uf *UF) Size() int               { return uf.size }
func (uf *UF) connected(p, q int) bool { return uf.find(p) == uf.find(q) }

/**********************************************************************************************************************/
// https://leetcode.cn/problems/min-cost-to-connect-all-points/
// 也是一道最小生成树的问题，只是边和权重需要前置处理一下而已
// 输入：points = [[0,0],[2,2],[3,10],[5,2],[7,0]]
// 输出：20
// 输入：points = [[3,12],[-2,5],[-4,1]]
// 输出：18
// minCostConnectPoints 连接所有点的最小费用
func minCostConnectPoints(points [][]int) int {
	// 生成所有的边和权重
	edges := [][]int{}
	n := len(points)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			xi := float64(points[i][0])
			xj := float64(points[j][0])
			yi := float64(points[i][1])
			yj := float64(points[j][1])
			edges = append(edges, []int{i, j, int(math.Abs(xi-xj) + math.Abs(yi-yj))})
		}
	}

	// 下面几乎都是模版代码
	sort.Slice(edges, func(i, j int) bool {
		return edges[i][2] < edges[j][2]
	})

	uf := NewUF(n)
	mset := 0

	for _, edge := range edges {
		u, v, weight := edge[0], edge[1], edge[2]
		if uf.connected(u, v) {
			continue
		}
		uf.union(u, v)
		mset += weight
	}

	if uf.Size() == 1 {
		return mset
	} else {
		return -1
	}
}

/**********************************************************************************************************************/
//回溯算法解决：子集/组合、排列问题
/**********************************************************************************************************************/
// https://leetcode.cn/problems/subsets/
// 子集/组合/排列问题解决思路：回溯算法
// 输入：numbers = [1,2,3]
// 输出：[[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
// subsets 标准的子集问题
func subsets(numbers []int) [][]int {
	var backtrack func([]int, int, []int, *[][]int)
	var result = make([][]int, 0, len(numbers))
	var track = []int{}

	backtrack = func(numbers []int, start int, track []int, result *[][]int) {
		temp := make([]int, len(track))
		copy(temp, track)
		*result = append(*result, temp)

		for i := start; i < len(numbers); i++ {
			track = append(track, numbers[i]) //做选择
			backtrack(numbers, i+1, track, result)
			track = track[:len(track)-1] //撤销选择
		}
	}
	backtrack(numbers, 0, track, &result)
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/combinations/
// 输入：n = 4, k = 2
// 输出：[[2,4],[3,4],[2,3],[1,2],[1,3],[1,4],
// combine 标准的组合问题
func combine(n int, k int) [][]int {
	var backtrack func(int, []int, *[][]int)
	var result = make([][]int, 0, n)
	var track = []int{}

	backtrack = func(pos int, track []int, result *[][]int) {
		if len(track) == k {
			temp := make([]int, len(track))
			copy(temp, track)
			*result = append(*result, temp)
			return
		}
		for i := pos; i <= n; i++ {
			track = append(track, i) //做选择
			backtrack(i+1, track, result)
			track = track[:len(track)-1] //撤销选择
		}
	}
	backtrack(1, track, &result)
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/permutations/description/
// 输入：numbers = [1,2,3]
// 输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
// _permute 标准的全排列问题
func _permute(numbers []int) [][]int {
	var backtrack func([]int, []int, *[][]int)
	var result = [][]int{}
	var track = []int{}
	var visited = map[int]bool{}

	backtrack = func(numbers []int, track []int, result *[][]int) {
		if len(track) == len(numbers) {
			temp := make([]int, len(track))
			copy(temp, track)
			*result = append(*result, temp)
			return
		}
		for i := 0; i < len(numbers); i++ {
			if visited[numbers[i]] {
				continue
			}
			visited[numbers[i]] = true
			track = append(track, numbers[i])
			backtrack(numbers, track, result)
			delete(visited, numbers[i])
			track = track[:len(track)-1]
		}
	}
	backtrack(numbers, track, &result)
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/subsets-ii/description/
// 输入：nums = [1,2,2]
// 输出：[[],[1],[1,2],[1,2,2],[2],[2,2]]
// subsetsWithDup 子集II，可能包含重复元素
func subsetsWithDup(numbers []int) [][]int {
	var backtrack func([]int, int, []int, *[][]int)
	var result = make([][]int, 0, len(numbers))
	var track = []int{}

	sort.Ints(numbers)

	backtrack = func(numbers []int, start int, track []int, result *[][]int) {
		temp := make([]int, len(track))
		copy(temp, track)
		*result = append(*result, temp)

		for i := start; i < len(numbers); i++ {
			//剪枝逻辑，值相同的相邻树枝，只遍历一次
			if i > start && numbers[i] == numbers[i-1] {
				continue
			}
			track = append(track, numbers[i]) //做选择
			backtrack(numbers, i+1, track, result)
			track = track[:len(track)-1] //撤销选择
		}
	}
	backtrack(numbers, 0, track, &result)
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/combination-sum-ii/
// 输入: candidates = [10,1,2,7,6,1,5], target = 8,
// 输出: [[1,1,6],[1,2,5],[1,7],[2,6]]
// combinationSum2 组合总和II，candidates存在重复元素
func combinationSum2(candidates []int, target int) [][]int {
	var backtrack func([]int, int, []int, *[][]int)
	var result = make([][]int, 0, len(candidates))
	var track = []int{}

	sort.Ints(candidates)

	backtrack = func(candidates []int, pos int, track []int, result *[][]int) {
		if sum(track) == target {
			temp := make([]int, len(track))
			copy(temp, track)
			*result = append(*result, temp)
			return
		}
		if sum(track) > target {
			return
		}
		for i := pos; i < len(candidates); i++ {
			//剪枝逻辑，值相同的相邻树枝，只遍历一次
			if i > pos && candidates[i] == candidates[i-1] {
				continue
			}
			track = append(track, candidates[i]) //做选择
			backtrack(candidates, i+1, track, result)
			track = track[:len(track)-1] //撤销选择
		}
	}
	backtrack(candidates, 0, track, &result)
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/permutations-ii/description/
// 输入：nums = [1,1,2]
// 输出：[[1,1,2],[1,2,1],[2,1,1]]
// permuteUnique 全排列II, 可包含重复数字的序列（元素可重复不可复选）
func permuteUnique(numbers []int) [][]int {
	var backtrack func([]int, []int, *[][]int)
	var result = [][]int{}
	var track = []int{}
	var visited = map[int]bool{}

	sort.Ints(numbers)
	backtrack = func(numbers []int, track []int, result *[][]int) {
		if len(track) == len(numbers) {
			temp := make([]int, len(track))
			copy(temp, track)
			*result = append(*result, temp)
			return
		}
		for i := 0; i < len(numbers); i++ {
			if visited[i] {
				continue
			}
			//新增的剪枝逻辑，保证了相同元素在排列中的相对位置保证固定
			if i > 0 && numbers[i] == numbers[i-1] && !visited[i-1] {
				continue
			}

			visited[i] = true
			track = append(track, numbers[i])
			backtrack(numbers, track, result)
			delete(visited, i)
			track = track[:len(track)-1]
		}
	}
	backtrack(numbers, track, &result)
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/combination-sum/description/
// 输入：candidates = [2,3,6,7], target = 7
// 输出：[[2,2,3],[7]]
// 输入: candidates = [2,3,5], target = 8
// 输出: [[2,2,2,2],[2,3,3],[3,5]]
// combinationSum 组合总和(元素无重复可复选)
func combinationSum(candidates []int, target int) [][]int {
	var backtrack func([]int, int, *[][]int)

	track := make([]int, 0, len(candidates))
	result := make([][]int, 0)
	pathSum := 0

	backtrack = func(candidates []int, pos int, result *[][]int) {
		if pathSum == target {
			temp := make([]int, len(track))
			copy(temp, track)
			*result = append(*result, temp)
			return
		}
		if pathSum > target {
			return
		}
		for i := pos; i < len(candidates); i++ {
			pathSum += candidates[i]
			track = append(track, candidates[i])
			backtrack(candidates, i, result) // i不加1就等于元素可以复选了
			pathSum -= candidates[i]
			track = track[:len(track)-1]
		}
	}
	backtrack(candidates, 0, &result)
	return result
}

/**********************************************************************************************************************/
//回溯算法解决：集合划分问题
/**********************************************************************************************************************/
// https://leetcode.cn/problems/partition-to-k-equal-sum-subsets/description/
// 输入： nums = [4, 3, 2, 3, 5, 2, 1], k = 4
// 输出： True
// 说明： 有可能将其分成 4 个子集（5），（1,4），（2,3），（2,3）等于总和。
// canPartitionKSubsets 划分为k个相等的子集(会超时)
func canPartitionKSubsets(numbers []int, k int) bool {
	i, j, arrSum := 0, len(numbers)-1, sum(numbers)

	//优化剪枝逻辑，提前命中if
	sort.Ints(numbers)
	for i < j {
		numbers[i], numbers[j] = numbers[j], numbers[i]
		i++
		j--
	}
	//桶的数量大于数组长度，或者，不能总和平均放到每一个桶
	if k > len(numbers) || arrSum%k != 0 {
		return false
	}

	var backtrack func([]int, int, []int) bool
	bucket := make([]int, k)
	target := arrSum / k

	backtrack = func(numbers []int, index int, bucket []int) bool {
		if index == len(numbers) {
			for _, bSum := range bucket {
				if bSum != target {
					return false
				}
			}
			return true
		}

		for i := 0; i < len(bucket); i++ {
			if bucket[i]+numbers[index] > target {
				continue
			}
			bucket[i] += numbers[index]
			if backtrack(numbers, index+1, bucket) {
				return true
			}
			bucket[i] -= numbers[index]
		}
		return false
	}
	return backtrack(numbers, 0, bucket)
}

/**********************************************************************************************************************/
//DFS算法搞定岛屿系列
/**********************************************************************************************************************/
// https://leetcode.cn/problems/number-of-islands/description/
// 输入：grid = [
// ["1","1","1","1","0"],
// ["1","1","0","1","0"],
// ["1","1","0","0","0"],
// ["0","0","0","0","0"]
// ]
// 输出：1
// numIslands 岛屿数量
func numIslands(grid [][]byte) int {
	dirs := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	res, m, n := 0, len(grid), len(grid[0])

	var dfs func([][]byte, int, int)

	//核心算法逻辑
	dfs = func(grid [][]byte, i, j int) {
		//边界判断
		if i < 0 || j < 0 || i >= m || j >= n {
			return
		}
		//是否已经访问
		if grid[i][j] == '0' {
			return
		}
		//标记访问（淹没掉）
		grid[i][j] = '0'
		for _, dir := range dirs {
			//dfs遍历
			dfs(grid, i+dir[0], j+dir[1])
		}
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' {
				res++
				dfs(grid, i, j)
			}
		}
	}
	return res
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/number-of-closed-islands/
// closedIsland 统计封闭岛屿的数目，和上面👆题基本一样，只是需要提前处理下矩阵四周的情况
func closedIsland(grid [][]int) int {
	dirs := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	res, m, n := 0, len(grid), len(grid[0])

	var dfs func([][]int, int, int)

	//核心算法逻辑
	dfs = func(grid [][]int, i, j int) {
		//边界判断
		if i < 0 || j < 0 || i >= m || j >= n {
			return
		}
		//是否已经访问
		if grid[i][j] == 1 {
			return
		}
		//标记访问（淹没掉）
		grid[i][j] = 1
		for _, dir := range dirs {
			dfs(grid, i+dir[0], j+dir[1])
		}
	}

	//前置处理
	for i := 0; i < n; i++ {
		dfs(grid, 0, i)   //淹没上边
		dfs(grid, m-1, i) //淹没下边
	}
	for i := 0; i < m; i++ {
		dfs(grid, i, 0)   //淹没左边
		dfs(grid, i, n-1) //淹没右边
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 0 {
				res++
				dfs(grid, i, j)
			}
		}
	}
	return res
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/number-of-enclaves/description/
// 输入：grid = [[0,0,0,0],[1,0,1,0],[0,1,1,0],[0,0,0,0]]
// 输出：3
// 解释：有三个 1 被 0 包围。一个 1 没有被包围，因为它在边界上。
// 思路同上，只不过是把跨越过的陆地加起来而已
// numEnclaves 飞地的数量
func numEnclaves(grid [][]int) int {
	dirs := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	res, m, n := 0, len(grid), len(grid[0])

	var dfs func([][]int, int, int)

	//核心算法逻辑
	dfs = func(grid [][]int, i, j int) {
		//边界判断
		if i < 0 || j < 0 || i >= m || j >= n {
			return
		}
		//是否已经访问
		if grid[i][j] == 0 {
			return
		}
		//标记访问（淹没掉）
		grid[i][j] = 0
		res++
		for _, dir := range dirs {
			dfs(grid, i+dir[0], j+dir[1])
		}
	}

	//前置处理
	for i := 0; i < n; i++ {
		dfs(grid, 0, i)   //淹没上边
		dfs(grid, m-1, i) //淹没下边
	}
	for i := 0; i < m; i++ {
		dfs(grid, i, 0)   //淹没左边
		dfs(grid, i, n-1) //淹没右边
	}

	res = 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				dfs(grid, i, j)
			}
		}
	}
	return res
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/max-area-of-island/
// 输入：grid = [
// [0,0,1,0,0,0,0,1,0,0,0,0,0],
// [0,0,0,0,0,0,0,1,1,1,0,0,0],
// [0,1,1,0,1,0,0,0,0,0,0,0,0],
// [0,1,0,0,1,1,0,0,1,0,1,0,0],
// [0,1,0,0,1,1,0,0,1,1,1,0,0],
// [0,0,0,0,0,0,0,0,0,0,1,0,0],
// [0,0,0,0,0,0,0,1,1,1,0,0,0],
// [0,0,0,0,0,0,0,1,1,0,0,0,0],
// ]
// 输出：6
// 解释：答案不应该是 11 ，因为岛屿只能包含水平或垂直这四个方向上的 1 。
// maxAreaOfIsland 岛屿的最大面积
func maxAreaOfIsland(grid [][]int) int {
	dirs := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	m, n := len(grid), len(grid[0])

	var dfs func([][]int, int, int, *int)

	//核心算法逻辑
	dfs = func(grid [][]int, i, j int, land *int) {
		//边界判断
		if i < 0 || j < 0 || i >= m || j >= n {
			return
		}
		//是否已经访问
		if grid[i][j] == 0 {
			return
		}
		//标记访问（淹没掉）
		grid[i][j] = 0
		*land++
		for _, dir := range dirs {
			dfs(grid, i+dir[0], j+dir[1], land)
		}
	}

	maxLand := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				land := 0
				dfs(grid, i, j, &land)
				maxLand = max(maxLand, land)
			}
		}
	}
	return maxLand
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/count-sub-islands/
// 思路：把grid2中的可能岛屿对应在grid1中是水域的全部淹没掉，剩下的就都是子岛屿了
// countSubIslands 统计子岛屿
func countSubIslands(grid1 [][]int, grid2 [][]int) int {
	dirs := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	res, m, n := 0, len(grid2), len(grid2[0])

	var dfs func([][]int, int, int)

	//核心算法逻辑
	dfs = func(grid [][]int, i, j int) {
		//边界判断
		if i < 0 || j < 0 || i >= m || j >= n {
			return
		}
		//是否已经访问
		if grid[i][j] == 0 {
			return
		}
		//标记访问（淹没掉）
		grid[i][j] = 0
		for _, dir := range dirs {
			dfs(grid, i+dir[0], j+dir[1])
		}
	}

	//前置处理
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid2[i][j] == 1 && grid1[i][j] == 0 {
				dfs(grid2, i, j)
			}
		}
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid2[i][j] == 1 {
				res++
				dfs(grid2, i, j)
			}
		}
	}
	return res
}

/**********************************************************************************************************************/
//动态规划系列
/**********************************************************************************************************************/
// https://leetcode.cn/problems/minimum-falling-path-sum/
// 动态规划-典型例子🌰，到达matrix[i][j]的最短路径为：
// matrix[i][j] + min(matrix[i-1][j-1]，matrix[i-1][j]，matrix[i-1][j+])（这个就是状态转移方程）
// 输入：matrix = [[2,1,3],[6,5,4],[7,8,9]]
// 输出：13
// 输入：matrix = [[-19,57],[-40,-5]]
// 输出：-59
// minFallingPathSum 下降路径最小和
func minFallingPathSum(matrix [][]int) int {
	var dp func(matrix [][]int, i, j int) int
	res := math.MaxInt
	n := len(matrix)
	memo := make([][]int, n)

	//初始化备忘录
	for i := 0; i < n; i++ {
		memo[i] = make([]int, len(matrix[0]))
		for j := 0; j < len(matrix[0]); j++ {
			memo[i][j] = 666666
		}
	}

	dp = func(matrix [][]int, i, j int) int {
		//边界判断
		if i < 0 || j < 0 || i >= len(matrix) || j >= len(matrix[0]) {
			return 9999999
		}
		//base case
		if i == 0 {
			return matrix[0][j]
		}
		//备忘录
		if memo[i][j] != 666666 {
			return memo[i][j]
		}
		//状态转移
		memo[i][j] = matrix[i][j] + min(min(dp(matrix, i-1, j), dp(matrix, i-1, j-1)), dp(matrix, i-1, j+1))
		return memo[i][j]
	}

	for j := 0; j < n; j++ {
		res = min(res, dp(matrix, n-1, j))
	}
	return res
}

/**********************************************************************************************************************/
//动态规划 之 子序列问题
/**********************************************************************************************************************/
// 设计迭代的动态规划算法时，不是要定义一个dp数组嘛，我们可以假设dp[0..i-1]都已经被计算出来了，然后问问自己，怎么通过这些结果计算出dp[i]。
// 子序列问题解题模版：
// 🚫首先，一旦涉及到子序列和最值，那么几乎可以肯定，考查的就是动态规划技巧
// 思路1：定义一个一维的dp数组；(解释：在子数组arr[0..i]中，以arr[i]结尾的子序列的长度为dp[i])
// for i:=1; i<n; i++ {
//   for j:=0; j<n; j++ {
//     dp[i] = 最值(dp[i], dp[j]+...)
//  }
// }
//
// 思路2：定义一个二维的dp数组；(解释：涉及两个字符串/数组场景时，在子数组arr1[0..i]和arr2[0..j]中，我们要求的子序列长度为dp[i][j]；
// 只涉及一个字符串/数组的场景时，在子数组arr[i..j]中，我们要求的子序列的长度为dp[i][j]。）
// for i:=0; i<n; i++ {
//   for j:=0; j<n; j++ {
//     if arr[i] == arr[j] {
//       dp[i][j] = dp[i][j] + ...
//    } else {
//       dp[i][j] = 最值(...)
//   }
//  }
// }
/**********************************************************************************************************************/
// https://leetcode.cn/problems/longest-increasing-subsequence/description/
// 思路：动态规划-自底向上递推的动态规划
// 输入：nums = [10,9,2,5,3,7,101,18]
// 输出：4
// 解释：最长递增子序列是 [2,3,7,101]，因此长度为 4 。
// lengthOfLIS 最长递增子序列
func lengthOfLIS(numbers []int) int {
	if len(numbers) == 0 || len(numbers) == 1 {
		return len(numbers)
	}

	//base case
	dp := make([]int, len(numbers))
	for i := 0; i < len(dp); i++ {
		dp[i] = 1
	}

	//状态转移方程
	for i := 1; i < len(numbers); i++ {
		for j := 0; j < i; j++ {
			if numbers[i] > numbers[j] {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
	}

	res := dp[0]
	for i := 0; i < len(dp); i++ {
		res = max(res, dp[i])
	}
	return res
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/maximum-subarray/description/
// 输入：nums = [-2,1,-3,4,-1,2,1,-5,4]
// 输出：6
// 解释：连续子数组 [4,-1,2,1] 的和最大，为 6 。
// maxSubArray 最大子数组和，滑动窗口思路
func maxSubArray(numbers []int) int {
	windowSum, maxSum := 0, math.MinInt
	left, right := 0, 0

	for right < len(numbers) {
		windowSum += numbers[right]
		right++
		maxSum = max(maxSum, windowSum)

		for windowSum < 0 {
			windowSum -= numbers[left]
			left++
		}
	}
	return maxSum
}

// _maxSubArray 最大子数组和，动态规划解法
func _maxSubArray(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}

	//base case
	dp := make([]int, len(numbers))
	dp[0] = numbers[0]

	//状态转移方程
	for i := 1; i < len(numbers); i++ {
		dp[i] = max(numbers[i], dp[i-1]+numbers[i])
	}

	res := dp[0]
	for i := 0; i < len(dp); i++ {
		res = max(res, dp[i])
	}
	return res
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/longest-common-subsequence/
// 思路：对于两个字符串求子序列的问题，都是用两个指针 i 和 j 分别在两个字符串上移动，大概率是动态规划的思路。
// 算法简单理解：就是把相同字符放到同一个lcs桶🪣里，最后计算lcs的长度，如果s1[i] != s2[j] ，那么要么
// s1前移一步，要么s2前移一步，然后穷举所有可能。
// 输入：text1 = "abcde", text2 = "ace"
// 输出：3
// 解释：最长公共子序列是 "ace" ，它的长度为 3 。
// longestCommonSubsequence 最长公共子序列
func longestCommonSubsequence(text1 string, text2 string) int {
	var dp func(s1 string, i int, s2 string, j int) int

	memo := make([][]int, len(text1))
	for i := 0; i < len(text1); i++ {
		memo[i] = make([]int, len(text2))
		for j := 0; j < len(text2); j++ {
			memo[i][j] = -1
		}
	}

	dp = func(s1 string, i int, s2 string, j int) int {
		//base case
		if i == len(s1) || j == len(s2) {
			return 0
		}
		//备忘录
		if memo[i][j] != -1 {
			return memo[i][j]
		}

		//状态转移方程
		if s1[i] == s2[j] {
			//相同字符则都属于lcs桶，那么结果+1(相当于放入桶里面了)
			memo[i][j] = dp(s1, i+1, s2, j+1) + 1
		} else {
			memo[i][j] = max(dp(s1, i+1, s2, j), dp(s1, i, s2, j+1))
		}
		return memo[i][j]
	}

	return dp(text1, 0, text2, 0)
}

// 自底向上迭代的动态规划思路
func _longestCommonSubsequence(text1, text2 string) int {
	// dp[i][j] 表示text1[0:i]和text2[0:j]的最长公共子序列的长度。
	//上述表示中，text1[0:i]表示text1的长度为i的前缀，text2[0:j]表示text2的长度为j的前缀。
	m, n := len(text1), len(text2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i, c1 := range text1 {
		for j, c2 := range text2 {
			// 相等取左上元素+1，否则取左或上的较大值，原因如下：
			// maxLen := 最长公共子序列的长度
			// c1 == c2，意味着是公共字符，那么直接在原来的 $maxLen+1 就好了
			// c1 != c2，意味着本次新加入的字符不会导致 $maxLen 增加，那么选之前最大的保存就好了
			if c1 == c2 {
				dp[i+1][j+1] = dp[i][j] + 1
			} else {
				dp[i+1][j+1] = max(dp[i][j+1], dp[i+1][j])
			}
		}
	}
	return dp[m][n]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/delete-operation-for-two-strings/
// 本质跟上面一个题一样，只是返回值变化了一下
// 输入: word1 = "sea", word2 = "eat"
// 输出: 2
// 解释: 第一步将 "sea" 变为 "ea" ，第二步将 "eat "变为 "ea"
// minDistance 两个字符串的删除操作
func minDistance(word1 string, word2 string) int {
	var dp func(s1 string, i int, s2 string, j int) int

	memo := make([][]int, len(word1))
	for i := 0; i < len(word1); i++ {
		memo[i] = make([]int, len(word2))
		for j := 0; j < len(word2); j++ {
			memo[i][j] = -1
		}
	}

	dp = func(s1 string, i int, s2 string, j int) int {
		//base case
		if i == len(word1) || j == len(word2) {
			return 0
		}
		//备忘录
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		//状态转移方程
		if s1[i] == s2[j] {
			memo[i][j] = dp(word1, i+1, word2, j+1) + 1
		} else {
			memo[i][j] = max(dp(word1, i+1, word2, j), dp(word1, i, word2, j+1))
		}
		return memo[i][j]
	}

	lcs := dp(word1, 0, word2, 0)

	return len(word1) - lcs + len(word2) - lcs
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/minimum-ascii-delete-sum-for-two-strings/
// 思路同上
// 输入: s1 = "sea", s2 = "eat"
// 输出: 231
// 解释: 在 "sea" 中删除 "s" 并将 "s" 的值(115)加入总和。 在 "eat" 中删除 "t" 并将 116 加入总和。
// 结束时，两个字符串相等，115 + 116 = 231 就是符合条件的最小和。
// minimumDeleteSum 两个字符串的最小ASCII删除和
func minimumDeleteSum(s1 string, s2 string) int {
	var dp func(s1 string, i int, s2 string, j int) int

	sumToEnd := func(s string, start int) int {
		res := 0 //这里不能用uint8，会导致溢出
		for i := start; i < len(s); i++ {
			res += int(s[i])
		}
		return res
	}

	memo := make([][]int, len(s1))
	for i := 0; i < len(s1); i++ {
		memo[i] = make([]int, len(s2))
		for j := 0; j < len(s2); j++ {
			memo[i][j] = -1
		}
	}

	dp = func(s1 string, i int, s2 string, j int) int {
		//base case
		if i == len(s1) {
			return sumToEnd(s2, j)
		}
		if j == len(s2) {
			return sumToEnd(s1, i)
		}
		//备忘录
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		//状态转移方程
		if s1[i] == s2[j] {
			memo[i][j] = dp(s1, i+1, s2, j+1)
		} else {
			memo[i][j] = min(
				dp(s1, i+1, s2, j)+int(s1[i]),
				dp(s1, i, s2, j+1)+int(s2[j]),
			)
		}
		return memo[i][j]
	}

	return dp(s1, 0, s2, 0)
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/longest-palindromic-subsequence/
// 思路：动态规划，在子串s[i..j]中，最长回文子序列的长度为dp[i][j]
// 状态转移方程：dp[i+1][j-1], dp[i+1][j], dp[i][j-1] => dp[i][j]
// 输入：s = "bbbab"
// 输出：4
// 解释：一个可能的最长回文子序列为 "bbbb" 。
// longestPalindromeSubseq 最长回文子序列
func longestPalindromeSubseq(s string) int {
	// 自底向上迭代（递推）的动态规划思路
	dp := make([][]int, len(s))
	n := len(s)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
		dp[i][i] = 1
	}

	// 要求dp[i][j]，就得需要知道dp[i+1][j-1], dp[i+1][j], dp[i][j-1]
	// 把dp二维数组画出来，然后反转遍历就能提前计算出dp[i+1][j-1], dp[i+1][j], dp[i][j-1]
	for i := n - 1; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			if s[i] == s[j] {
				dp[i][j] = dp[i+1][j-1] + 2
			} else {
				dp[i][j] = max(dp[i+1][j], dp[i][j-1])
			}
		}
	}
	return dp[0][n-1]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/minimum-insertion-steps-to-make-a-string-palindrome/
// 思路与上面👆题类似，dp[i][j] 表示把字符串s[i..j]变成回文串的最少插入次数
// minInsertions 让字符串成为回文串的最少插入次数
func minInsertions(s string) int {
	dp := make([][]int, len(s))
	n := len(s)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
	}

	for i := n - 1; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			if s[i] == s[j] {
				dp[i][j] = dp[i+1][j-1]
			} else {
				dp[i][j] = min(dp[i+1][j], dp[i][j-1]) + 1
			}
		}
	}
	return dp[0][n-1]
}

/**********************************************************************************************************************/
// 动态规划 之 背包问题
// dp[i][w]表示：对于前i个物品（从1开始计数），当前背包的容量为w时，这种情况下可以装下的最大价值时dp[i][w]
// 把第i个物品装入背包：dp[i][w] = dp[i-1][w-wt[i-1]] + val[i-1]
// 不把第i个物品装入背包：dp[i][w] = dp[i-1][w]（⚠️注意：i-1是因为计数从1开始的，访问数组时需要下标-1。）
// for i:=1; i<N; i++ {
//  for w:=1; w<W; w++ {
// 		if w - wt[i-1] <0 {
// 			dp[i][w] = dp[i-1][w] //放入的物品太重，装不下了
//		}else{
// 			dp[i][w] = max(dp[i-1][w], dp[i-1][w-wt[i-1]]+val[i-1])
//		}
// 	}
//}
/**********************************************************************************************************************/
// https://leetcode.cn/problems/partition-equal-subset-sum/description/
// 0-1背包问题变体
// 可以把问题转化为: 给定一个可装载 W 的背包和 N 个物品，每个物品的重量为numbers[i]，现在问是否存在一种装法，能够恰好装满背包？
// dp[i][j]=x 的含义：对于前i个物品，当前背包容量为j时，若x为true表示可以恰好将背包装满，x为false表示不能恰好装满背包。
// 输入：numbers = [1,5,11,5]
// 输出：true
// 解释：数组可以分割成 [1, 5, 5] 和 [11] 。
// canPartition 分割等和子集
func canPartition(numbers []int) bool {
	totalSum := 0
	for _, num := range numbers {
		totalSum += num
	}
	if totalSum%2 != 0 {
		return false
	}

	W := totalSum / 2 //单个集合的承重量
	N := len(numbers)

	//base case
	dp := make([][]bool, N+1)
	for i := 0; i <= N; i++ {
		dp[i] = make([]bool, W+1)
		dp[i][0] = true //背包容量为0的时候，相当于背包满了
	}

	for i := 1; i <= N; i++ {
		for w := 1; w <= W; w++ {
			if w-numbers[i-1] < 0 {
				dp[i][w] = dp[i-1][w] //背包装不下了
			} else {
				dp[i][w] = dp[i-1][w] || dp[i-1][w-numbers[i-1]]
			}
		}
	}
	return dp[N][W]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/coin-change-ii/description/
// 输入：amount = 5, coins = [1, 2, 5]
// 输出：4
// 解释：有四种方式可以凑成总金额：
// 5=5
// 5=2+2+1
// 5=2+1+1+1
// 5=1+1+1+1+1
// 思路：dp[i][j] 表示使用coins中前i个硬币的面值(可以重复使用)，刚好凑出金额为j，有dp[i][j]中凑法
// change 零钱兑换 II
func change(amount int, coins []int) int {
	dp := make([][]int, len(coins)+1)
	for i := 0; i <= len(coins); i++ {
		dp[i] = make([]int, amount+1)
		dp[i][0] = 1 //base case，金额为0直接满足了，方式都为一种（不使用任何硬币）
	}

	for i := 1; i <= len(coins); i++ {
		for j := 1; j <= amount; j++ {
			if j-coins[i-1] >= 0 {
				//不同的方式 = 把物品装进去 + 不把物品装进去
				dp[i][j] = dp[i-1][j] + dp[i][j-coins[i-1]]
			} else {
				dp[i][j] = dp[i-1][j]
			}
		}
	}
	return dp[len(coins)][amount]
}

/**********************************************************************************************************************/
// 动态规划 之 玩游戏
/**********************************************************************************************************************/
// https://leetcode.cn/problems/minimum-path-sum/
// 典型的动态规划题目
// 输入：grid = [[1,3,1],[1,5,1],[4,2,1]]
// 输出：7
// 解释：因为路径 1→3→1→1→1 的总和最小。
// minPathSum 最小路径和，自顶向下递归的动态规划
func minPathSum(grid [][]int) int {
	var dp func(grid [][]int, i, j int) int

	memo := make([][]int, len(grid))
	for i := 0; i < len(grid); i++ {
		memo[i] = make([]int, len(grid[0]))
		for j := 0; j < len(grid[0]); j++ {
			memo[i][j] = -1
		}
	}

	dp = func(grid [][]int, i, j int) int {
		if i == 0 && j == 0 {
			return grid[0][0]
		}
		if i < 0 || j < 0 {
			return math.MaxInt
		}
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		memo[i][j] = min(dp(grid, i-1, j), dp(grid, i, j-1)) + grid[i][j]
		return memo[i][j]
	}
	return dp(grid, len(grid)-1, len(grid[0])-1)
}

// _minPathSum 最小路径和，自底向上递推（迭代）的动态规划
func _minPathSum(grid [][]int) int {
	dp := make([][]int, len(grid))
	for i := 0; i < len(grid); i++ {
		dp[i] = make([]int, len(grid[0]))
	}

	//base case
	dp[0][0] = grid[0][0]
	for i := 1; i < len(grid); i++ {
		dp[i][0] = dp[i-1][0] + grid[i][0]
	}
	for i := 1; i < len(grid[0]); i++ {
		dp[0][i] = dp[0][i-1] + grid[0][i]
	}

	//状态转移方程
	for i := 1; i < len(grid); i++ {
		for j := 1; j < len(grid[0]); j++ {
			dp[i][j] = grid[i][j] + min(dp[i-1][j], dp[i][j-1])
		}
	}
	return dp[len(grid)-1][len(grid[0])-1]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/dungeon-game/
// 思路：需要一点点技巧，反向思考，即从grid[i][j]到达终点所需的最少生命值是dp(grid,i,j)
// calculateMinimumHP 地下城游戏
func calculateMinimumHP(grid [][]int) int {
	var dp func(grid [][]int, i, j int) int
	m, n := len(grid), len(grid[0])

	memo := make([][]int, m)
	for i := 0; i < m; i++ {
		memo[i] = make([]int, n)
		for j := 0; j < n; j++ {
			memo[i][j] = -1
		}
	}

	dp = func(grid [][]int, i, j int) int {
		//base case
		if i == m-1 && j == n-1 {
			if grid[i][j] >= 0 {
				//最后一步不消耗血量，所以只需要满足最低要求1即可
				return 1
			}
			//最后一步需要消耗血量，所以必须满足-grid[i][j] + 1
			return -grid[i][j] + 1
		}
		//边界判断
		if i == m || j == n {
			return math.MaxInt
		}
		//备忘录
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		//状态转移方程
		res := min(dp(grid, i+1, j), dp(grid, i, j+1)) - grid[i][j]
		if res <= 0 {
			memo[i][j] = 1
		} else {
			memo[i][j] = res
		}
		return memo[i][j]
	}
	return dp(grid, 0, 0)
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/super-egg-drop/description/
// 输入：k = 1, n = 2
// 输出：2
// 解释：鸡蛋从 1 楼掉落。如果它碎了，肯定能得出 f = 0 。否则，鸡蛋从 2 楼掉落。
// 如果它碎了，肯定能得出 f = 1 。如果它没碎，那么肯定能得出 f = 2 。因此，在最坏的情况下我们需要移动 2 次以确定 f 是多少。
// superEggDrop 鸡蛋掉落
func superEggDrop(k int, n int) int {
	//状态（变化量）：鸡蛋🥚个数K，楼层数N
	//选择：其实就是选择哪层楼扔鸡蛋🥚
	//定义：手握K个鸡蛋🥚，面对N层楼，最少的扔鸡蛋🥚次数为dp(k,n)
	var dp func(k, n int) int

	memo := make([][]int, k+1)
	for i := 0; i <= k; i++ {
		memo[i] = make([]int, n+1)
		for j := 0; j <= n; j++ {
			memo[i][j] = -666
		}
	}

	dp = func(k, n int) int {
		//base case
		if k == 1 {
			return n //只有一个鸡蛋时，只能线性扫描
		}
		if n == 0 {
			return 0 //楼层为0时，不需要扔鸡蛋
		}
		//备忘录
		if memo[k][n] != -666 {
			return memo[k][n]
		}
		//状态转移方程
		res := math.MaxInt
		//for i := 1; i <= n; i++ {
		//	//在所有楼层进行尝试（+1就是尝试次数），取最少扔鸡蛋次数
		//	res = min(
		//		res,
		//		max(dp(k, n-i), dp(k-1, i-1))+1, //碎和没碎取最坏情况
		//	)
		//}

		//优化：二分搜索替代线性搜索
		lo, hi := 1, n
		for lo <= hi {
			mid := lo + (hi-lo)/2
			broken := dp(k-1, mid-1)   //碎了
			not_borken := dp(k, n-mid) //没碎

			//res = min(res, max(碎，没碎)+1)
			if broken > not_borken {
				hi = mid - 1
				res = min(res, broken+1)
			} else {
				lo = mid + 1
				res = min(res, not_borken+1)
			}
		}

		memo[k][n] = res
		return memo[k][n]
	}
	return dp(k, n)
}

/**********************************************************************************************************************/
// 题目转化：在一排气球points中,请你戳破气球0和气球n+1之间的所有气球（不包含0，n+1），使得最终只剩下气球0和气球n+1两个气球，最多能够得多少分？
// https://leetcode.cn/problems/burst-balloons/
// 输入：numbers = [3,1,5,8]
// 输出：167
// 解释：numbers = [3,1,5,8] --> [3,5,8] --> [3,8] --> [8] --> []; coins =  3*1*5    +   3*5*8   +  1*3*8  + 1*8*1 = 167
// maxCoins 戳气球
func maxCoins(numbers []int) int {
	//添加两侧的虚拟气球
	points := make([]int, len(numbers)+2)
	n := len(numbers)
	points[0], points[n+1] = 1, 1
	for i := 1; i <= n; i++ {
		points[i] = numbers[i-1]
	}

	// dp[i][k] = x表示，戳破气球i和气球j（开区间，不包括i和j）之间的所有气球可以获得的最高分数为x。
	dp := make([][]int, n+2)
	for i := 0; i < n+2; i++ {
		dp[i] = make([]int, n+2)
	}

	//i,j反正遍历，i,j 取值范围(0,n+2)，0和n+1是虚拟气球🎈
	for i := n; i >= 0; i-- {
		for j := i + 1; j < n+2; j++ {
			//在(i,j)区间内，把所有气球都戳破一遍，记录最大值
			for k := i + 1; k < j; k++ {
				dp[i][j] = max(dp[i][j], dp[i][k]+dp[k][j]+points[i]*points[k]*points[j])
			}
		}
	}
	return dp[0][n+1]
}
