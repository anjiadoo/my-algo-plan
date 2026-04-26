package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// 简化路径 https://leetcode.cn/problems/simplify-path/description/
func simplifyPath(path string) string {
	// 一定要明确栈中存储的是啥
	// 这里存的是文件夹组成路径
	var stack []string

	for _, str := range strings.Split(path, "/") {
		switch str {
		case "", ".":
			continue
		case "..":
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		default:
			stack = append(stack, str)
		}
	}
	if len(stack) == 0 {
		return "/"
	}
	return "/" + strings.Join(stack, "/")
}

// 有效的括号 https://leetcode.cn/problems/valid-parentheses/description/
func isValid(s string) bool {
	var stack []byte
	for _, ch := range s {
		switch ch {
		case '(', '{', '[':
			stack = append(stack, byte(ch))
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false
			}
			stack = stack[:len(stack)-1]
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				return false
			}
			stack = stack[:len(stack)-1]
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	if len(stack) > 0 {
		return false
	}
	return true
}

// 逆波兰表达式求值 https://leetcode.cn/problems/evaluate-reverse-polish-notation/
func evalRPN(tokens []string) int {
	var stack []int
	for _, token := range tokens {
		if strings.Contains("+-*/", token) {

			a := stack[len(stack)-2]
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "*":
				stack = append(stack, a*b)
			case "-":
				stack = append(stack, a-b)
			case "/":
				stack = append(stack, a/b)
			}
		} else {
			num, _ := strconv.Atoi(token)
			stack = append(stack, num)
		}
	}
	return stack[0]
}

// 文件的最长绝对路径 https://leetcode.cn/problems/longest-absolute-file-path/
func lengthLongestPath(input string) int {
	var stack []string
	maxLength := 0

	for _, part := range strings.Split(input, "\n") {
		level := strings.LastIndexByte(part, '\t') + 1

		// 平级需要出栈
		for len(stack) > level {
			stack = stack[:len(stack)-1]
		}

		// 子级需要入栈
		stack = append(stack, part[level:])

		// 如果是文件，更新最长路径
		if strings.Contains(part, ".") {
			fullPath := strings.Join(stack, "/")
			maxLength = max(maxLength, len(fullPath))
		}
	}
	return maxLength
}

// 字符串解码 https://leetcode.cn/problems/decode-string/description/
func decodeString(s string) string {
	var stack []byte
	for _, ch := range s {
		if ch == ']' {
			str := ""
			for len(stack) > 0 && stack[len(stack)-1] != '[' {
				str = string(stack[len(stack)-1]) + str
				stack = stack[:len(stack)-1]
			}
			// 去掉'['
			stack = stack[:len(stack)-1]

			sNum := ""
			for len(stack) > 0 && stack[len(stack)-1] >= '0' && stack[len(stack)-1] <= '9' {
				sNum = string(stack[len(stack)-1]) + sNum
				stack = stack[:len(stack)-1]
			}

			num, _ := strconv.Atoi(sNum)
			subStr := strings.Repeat(str, num)
			stack = append(stack, []byte(subStr)...)
		} else {
			stack = append(stack, byte(ch))
		}
	}
	return string(stack)
}

// MinStack 最小栈 https://leetcode.cn/problems/min-stack/
type MinStack struct {
	stack    []int // 原始栈
	minStack []int // 存最小元素栈长度与stack相等
}

func (s *MinStack) Push(val int) {
	s.stack = append(s.stack, val)
	if len(s.minStack) == 0 || val < s.minStack[len(s.minStack)-1] {
		s.minStack = append(s.minStack, val)
	} else {
		s.minStack = append(s.minStack, s.minStack[len(s.minStack)-1])
	}
}

func (s *MinStack) Pop() {
	s.stack = s.stack[:len(s.stack)-1]
	s.minStack = s.minStack[:len(s.minStack)-1]
}

func (s *MinStack) Top() int {
	return s.stack[len(s.stack)-1]
}

func (s *MinStack) GetMin() int {
	return s.minStack[len(s.minStack)-1]
}

// FreqStack 最大频率栈 https://leetcode.cn/problems/maximum-frequency-stack/description/
type FreqStack struct {
	maxFreq   int           // 记录FreqStack中元素的最大频率
	val2Freq  map[int]int   // 记录FreqStack中每个val对应的出现频率
	freq2Vals map[int][]int // 记录频率freq对应的val列表
}

func Constructor() FreqStack {
	return FreqStack{
		val2Freq:  make(map[int]int),
		freq2Vals: make(map[int][]int),
	}
}

func (s *FreqStack) Push(val int) {
	freq := s.val2Freq[val] + 1
	s.val2Freq[val] = freq

	s.freq2Vals[freq] = append(s.freq2Vals[freq], val)
	if freq > s.maxFreq {
		s.maxFreq = freq
	}
}

func (s *FreqStack) Pop() int {
	vals := s.freq2Vals[s.maxFreq]
	res := vals[len(vals)-1]
	vals = vals[:len(vals)-1]
	s.freq2Vals[s.maxFreq] = vals

	s.val2Freq[res]--
	if len(s.freq2Vals[s.maxFreq]) == 0 {
		delete(s.freq2Vals, s.maxFreq)
		s.maxFreq--
	}
	return res
}

// 使括号有效的最少添加 https://leetcode.cn/problems/minimum-add-to-make-parentheses-valid/
func minAddToMakeValid(s string) int {
	needRight := 0 // 对右括号的需求
	needLeft := 0  // 对左括号的需求

	for _, ch := range s {
		if ch == '(' {
			needRight++
		}
		if ch == ')' {
			needRight--
			if needRight < 0 {
				needRight = 0
				needLeft++
			}
		}
	}
	return needLeft + needRight
}

// 平衡括号字符串的最少插入次数 https://leetcode.cn/problems/minimum-insertions-to-balance-a-parentheses-string/description/
func minInsertions(s string) int {
	needRight := 0 // 对右括号的需求
	needLeft := 0  // 对左括号的需求

	for _, ch := range s {
		if ch == '(' {
			needRight += 2

			// 难点：当遇到左括号时
			// 若对右括号的需求量为奇数
			// 需要插入 1 个右括号
			// 一个左括号对应两个右括号
			if needRight%2 == 1 {
				needLeft++
				needRight--
			}
		}
		if ch == ')' {
			needRight--
			if needRight == -1 {
				needLeft++
				needRight = 1
			}
		}
	}
	return needLeft + needRight
}

// RecentCounter 最近的请求次数 https://leetcode.cn/problems/number-of-recent-calls/
type RecentCounter struct {
	queue []int
}

func (q *RecentCounter) Ping(t int) int {
	q.queue = append(q.queue, t)
	for len(q.queue) > 0 && q.queue[0] < t-3000 {
		q.queue = q.queue[1:]
	}
	return len(q.queue)
}

// MyCircularQueue 设计循环队列 https://leetcode.cn/problems/design-circular-queue/
type MyCircularQueue struct {
	first int
	last  int
	size  int
	cap   int
	data  []int
}

func NewMyCircularQueue(k int) MyCircularQueue {
	return MyCircularQueue{
		cap:  k,
		data: make([]int, k),
	}
}

func (m *MyCircularQueue) EnQueue(value int) bool {
	if m.IsFull() {
		return false
	}
	m.data[m.last] = value
	m.last = (m.last + 1 + m.cap) % m.cap
	m.size++
	return true
}

func (m *MyCircularQueue) DeQueue() bool {
	if m.IsEmpty() {
		return false
	}
	m.data[m.first] = -1
	m.first = (m.first + 1 + m.cap) % m.cap
	m.size--
	return true
}

func (m *MyCircularQueue) Front() int {
	if m.IsEmpty() {
		return -1
	}
	return m.data[m.first]
}

func (m *MyCircularQueue) Rear() int {
	if m.IsEmpty() {
		return -1
	}
	return m.data[(m.last-1+m.cap)%m.cap]
}

func (m *MyCircularQueue) IsEmpty() bool {
	return m.size == 0
}

func (m *MyCircularQueue) IsFull() bool {
	return m.size == m.cap
}

// MyCircularDeque 设计循环双端队列 https://leetcode.cn/problems/design-circular-deque/
type MyCircularDeque struct {
	start int
	end   int
	size  int
	cap   int
	data  []int
}

func NewMyCircularDeque(k int) MyCircularDeque {
	return MyCircularDeque{
		cap:  k,
		data: make([]int, k),
	}
}

func (m *MyCircularDeque) InsertFront(value int) bool {
	if m.IsFull() {
		return false
	}
	m.start = (m.start - 1 + m.cap) % m.cap
	m.data[m.start] = value
	m.size++
	return true
}

func (m *MyCircularDeque) InsertLast(value int) bool {
	if m.IsFull() {
		return false
	}
	m.data[m.end] = value
	m.end = (m.end + 1 + m.cap) % m.cap
	m.size++
	return true
}

func (m *MyCircularDeque) DeleteFront() bool {
	if m.IsEmpty() {
		return false
	}
	m.data[m.start] = -1
	m.start = (m.start + 1 + m.cap) % m.cap
	m.size--
	return true
}

func (m *MyCircularDeque) DeleteLast() bool {
	if m.IsEmpty() {
		return false
	}
	m.end = (m.end - 1 + m.cap) % m.cap
	m.size--
	return true
}

func (m *MyCircularDeque) GetFront() int {
	if m.IsEmpty() {
		return -1
	}
	return m.data[m.start]
}

func (m *MyCircularDeque) GetRear() int {
	if m.IsEmpty() {
		return -1
	}
	return m.data[(m.end-1+m.cap)%m.cap]
}

func (m *MyCircularDeque) IsEmpty() bool {
	return m.size == 0
}

func (m *MyCircularDeque) IsFull() bool {
	return m.size == m.cap
}

// 下一个更大元素I https://leetcode.cn/problems/next-greater-element-i/description/
func nextGreaterElement(nums1 []int, nums2 []int) []int {
	mapNextGreater := map[int]int{}

	// 从后往前遍历，栈里存的是当前元素的下一个更大元素
	var stack []int

	for i := len(nums2) - 1; i >= 0; i-- {
		for len(stack) > 0 && nums2[i] >= stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			mapNextGreater[nums2[i]] = -1
		} else {
			mapNextGreater[nums2[i]] = stack[len(stack)-1]
		}
		stack = append(stack, nums2[i])
	}

	var res []int
	for _, num := range nums1 {
		res = append(res, mapNextGreater[num])
	}
	return res
}

// 每日温度 https://leetcode.cn/problems/daily-temperatures/
func dailyTemperatures(temperatures []int) []int {
	res := make([]int, len(temperatures))
	var stack []int

	for i := len(temperatures) - 1; i >= 0; i-- {
		for len(stack) > 0 && temperatures[i] >= temperatures[stack[len(stack)-1]] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			res[i] = 0
		} else {
			res[i] = stack[len(stack)-1] - i
		}
		stack = append(stack, i)
	}
	return res
}

// 下一个更大元素II https://leetcode.cn/problems/next-greater-element-ii/
func nextGreaterElements(nums []int) []int {
	result := make([]int, len(nums))
	n := len(nums)
	var stack []int

	// 数组长度加倍模拟环形数组
	for i := 2*n - 1; i >= 0; i-- {
		for len(stack) > 0 && nums[i%n] >= stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			result[i%n] = -1
		} else {
			result[i%n] = stack[len(stack)-1]
		}
		stack = append(stack, nums[i%n])
	}

	return result
}

type ListNode1 struct {
	Val  int
	Next *ListNode1
}

func NewMyListNode_(nums []int) *ListNode1 {
	dummy := &ListNode1{}
	p := dummy
	for i := 0; i < len(nums); i++ {
		p.Next = &ListNode1{Val: nums[i]}
		p = p.Next
	}
	return dummy.Next
}

// 链表中的下一个更大节点 https://leetcode.cn/problems/next-greater-node-in-linked-list/description/
func nextLargerNodes(head *ListNode1) []int {
	var nums []int
	for p := head; p != nil; p = p.Next {
		nums = append(nums, p.Val)
	}

	res := make([]int, len(nums))
	var stack []int

	for i := len(nums) - 1; i >= 0; i-- {
		for len(stack) > 0 && stack[len(stack)-1] <= nums[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			res[i] = 0
		} else {
			res[i] = stack[len(stack)-1]
		}
		stack = append(stack, nums[i])
	}
	return res
}

// 队列中可以看到的人数 https://leetcode.cn/problems/number-of-visible-people-in-a-queue/description/
func canSeePersonsCount(heights []int) []int {
	res := make([]int, len(heights))
	var stack []int

	for i := len(heights) - 1; i >= 0; i-- {
		count := 0
		for len(stack) > 0 && heights[stack[len(stack)-1]] <= heights[i] {
			stack = stack[:len(stack)-1]
			count++ //⚠️记录被挤掉的人
		}
		if len(stack) == 0 {
			res[i] = count
		} else {
			res[i] = count + 1 //⚠️除了可以看到“被挤掉的人”外，比自己高的那个人也能看到
		}
		stack = append(stack, i)
	}

	return res
}

// 商品折扣后的最终价格 https://leetcode.cn/problems/final-prices-with-a-special-discount-in-a-shop/description/
func finalPrices(prices []int) []int {
	res := make([]int, len(prices))
	var stack []int

	for i := len(prices) - 1; i >= 0; i-- {
		for len(stack) > 0 && stack[len(stack)-1] > prices[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			res[i] = prices[i]
		} else {
			res[i] = prices[i] - stack[len(stack)-1]
		}
		stack = append(stack, prices[i])
	}
	return res
}

// 移掉K位数字 https://leetcode.cn/problems/remove-k-digits/description/
func removeKdigits(num string, k int) string {
	var stack []rune
	for _, ch := range num {
		// 上一个更小或相等
		for len(stack) > 0 && stack[len(stack)-1] > (ch) && k > 0 {
			stack = stack[:len(stack)-1]
			k--
		}

		// 防止 0 作为数字的开头
		if len(stack) == 0 && ch == '0' {
			continue
		}
		stack = append(stack, ch)
	}

	// ⚠️k没用完，栈一定是单调递增的，上面的for循环
	for k > 0 && len(stack) > 0 {
		stack = stack[:len(stack)-1]
		k--
	}

	if len(stack) == 0 {
		return "0"
	}
	return string(stack)
}

// 车队 https://leetcode.cn/problems/car-fleet/
func carFleet(target int, position []int, speed []int) int {
	// 关键：按起始位置排序后，到达时间快的车会被后面到达时间慢的车卡住
	// 所以应该是计算单调递减序列

	type Pair struct {
		pos   int
		speed int
	}

	var pairs []Pair
	for i := 0; i < len(position); i++ {
		pairs = append(pairs, Pair{pos: position[i], speed: speed[i]})
	}

	// 按起始位置排序
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].pos < pairs[j].pos
	})

	// 计算到达时间
	var times []float64
	for _, pair := range pairs {
		times = append(times, float64(target-pair.pos)/float64(pair.speed))
	}

	// 单调递减栈大小就是答案
	//var stack []int
	//for i := len(times) - 1; i >= 0; i-- {
	//	for len(stack) > 0 && times[i] <= stack[len(stack)-1] {
	//		stack = stack[:len(stack)-1]
	//	}
	//	stack = append(stack, times[i])
	//}

	// 避免使用栈模拟，倒序遍历取递增序列就是答案
	var maxTime float64
	res := 0
	for i := len(times) - 1; i >= 0; i-- {
		if times[i] > maxTime {
			res++
			maxTime = times[i]
		}
	}
	return res
}

// 最短无序连续子数组 https://leetcode.cn/problems/shortest-unsorted-continuous-subarray/
func findUnsortedSubarray(nums []int) int {
	temp := make([]int, len(nums))
	copy(temp, nums)
	sort.Ints(temp)

	left := 0
	for i := 0; i < len(nums); i++ {
		if temp[left] != nums[i] {
			break
		}
		left++
	}

	right := len(nums) - 1
	for i := len(nums) - 1; i >= 0; i-- {
		if temp[right] != nums[i] {
			break
		}
		right--
	}

	if left == len(nums) && right == -1 {
		return 0
	}
	return right - left + 1
}

func findUnsortedSubarrayII(nums []int) int {
	left, right := len(nums), -1

	// 递增栈，弹出的元素都是乱序元素
	var incrStack []int
	for i := 0; i < len(nums); i++ {
		for len(incrStack) > 0 && nums[incrStack[len(incrStack)-1]] > nums[i] {
			left = min(left, incrStack[len(incrStack)-1])
			incrStack = incrStack[:len(incrStack)-1]
		}
		incrStack = append(incrStack, i)
	}

	// 递减栈，弹出的元素都是乱序元素
	var decrStack []int
	for i := len(nums) - 1; i >= 0; i-- {
		for len(decrStack) > 0 && nums[decrStack[len(decrStack)-1]] < nums[i] {
			right = max(right, decrStack[len(decrStack)-1])
			decrStack = decrStack[:len(decrStack)-1]
		}
		decrStack = append(decrStack, i)
	}

	// 单调栈没有弹出任何元素，说明nums本来就是有序的
	if left == len(nums) && right == -1 {
		return 0
	}
	return right - left + 1
}

// 柱状图中最大的矩形 https://leetcode.cn/problems/largest-rectangle-in-histogram/
func largestRectangleArea(heights []int) int {
	nums := make([]int, len(heights)+2)
	for i := 0; i < len(heights); i++ {
		nums[i+1] = heights[i]
	}

	// 求最大，那么就需要穷举，思路是：单调递增栈
	// 找 heights[i] 左侧第一个更小元素的索引 left
	// 找 heights[i] 右侧第一个更小元素的索引 right
	// 单栈需要 pop 元素的时候，就找到了右侧第一个更小

	var stack []int
	maxArea := 0

	for i := 0; i < len(nums); i++ {
		for len(stack) > 0 && nums[stack[len(stack)-1]] > nums[i] {
			// 先取高度，先 pop
			height := nums[stack[len(stack)-1]]
			stack = stack[:len(stack)-1]

			// 再算宽度：用pop后的新栈顶作为左边界（开区间）
			width := i - stack[len(stack)-1] - 1
			maxArea = max(maxArea, width*height)
		}
		stack = append(stack, i)
	}
	return maxArea
}

// 滑动窗口最大值 https://leetcode.cn/problems/sliding-window-maximum/
func maxSlidingWindow(nums []int, k int) []int {
	mq := NewMonotonicQueue()
	var res []int

	for i := 0; i < len(nums); i++ {
		if i < k-1 {
			mq.push(nums[i])
		} else {
			mq.push(nums[i])
			res = append(res, mq.getMax())
			mq.pop()
		}
	}
	return res
}

type MonotonicQueue struct {
	data []int //原始队列
	maxQ []int //单调递减，维护队列最大值
	minQ []int //单调递增，维护队列最小值
}

func NewMonotonicQueue() *MonotonicQueue {
	return &MonotonicQueue{}
}

func (mq *MonotonicQueue) getMax() int {
	if len(mq.maxQ) == 0 {
		return -1
	}
	return mq.maxQ[0]
}

func (mq *MonotonicQueue) getMin() int {
	if len(mq.minQ) == 0 {
		return -1
	}
	return mq.minQ[0]
}

func (mq *MonotonicQueue) pop() int {
	if len(mq.data) == 0 {
		return -1
	}
	elem := mq.data[0]
	mq.data = mq.data[1:]
	if elem == mq.getMax() {
		mq.maxQ = mq.maxQ[1:]
	}
	if elem == mq.getMin() {
		mq.minQ = mq.minQ[1:]
	}
	return elem
}

func (mq *MonotonicQueue) push(elem int) {
	mq.data = append(mq.data, elem)
	for len(mq.maxQ) > 0 && mq.maxQ[len(mq.maxQ)-1] < elem {
		mq.maxQ = mq.maxQ[:len(mq.maxQ)-1]
	}
	mq.maxQ = append(mq.maxQ, elem)

	for len(mq.minQ) > 0 && mq.minQ[len(mq.minQ)-1] > elem {
		mq.minQ = mq.minQ[:len(mq.minQ)-1]
	}
	mq.minQ = append(mq.minQ, elem)
}

// 绝对差不超过限制的最长连续子数组 https://leetcode.cn/problems/longest-continuous-subarray-with-absolute-diff-less-than-or-equal-to-limit/description/
func longestSubarray(nums []int, limit int) int {
	mq := NewMonotonicQueue()
	left, right := 0, 0
	res := 0

	for right < len(nums) {
		mq.push(nums[right])
		right++
		for mq.getMax()-mq.getMin() > limit {
			mq.pop()
			left++
		}
		// 左闭右开，因为right++了
		res = max(res, right-left)
	}
	return res
}

func main() {

	fmt.Println(longestSubarray([]int{8, 2, 4, 7}, 4))
	fmt.Println(longestSubarray([]int{10, 1, 2, 4, 7, 2}, 5))
	fmt.Println(longestSubarray([]int{4, 2, 2, 2, 4, 4, 2, 2}, 0))

	//fmt.Println(maxSlidingWindow([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3))
	//fmt.Println(maxSlidingWindow([]int{1, 3, 1, 2, 0, 5}, 3))

	//fmt.Println(largestRectangleArea([]int{2, 1, 5, 6, 2, 3}))
	//fmt.Println(largestRectangleArea([]int{2, 1, 2}))

	//fmt.Println(findUnsortedSubarrayII([]int{2, 6, 4, 8, 10, 9, 15}))
	//fmt.Println(findUnsortedSubarrayII([]int{1, 2, 3, 4, 5}))

	//fmt.Println(carFleet(10, []int{6, 8}, []int{3, 2}))
	//fmt.Println(carFleet(12, []int{10, 8, 0, 5, 3}, []int{2, 4, 1, 1, 3}))

	//fmt.Println(removeKdigits("1432219", 3))
	//fmt.Println(removeKdigits("10200", 1))

	//fmt.Println(finalPrices([]int{10, 1, 1, 6}))
	//fmt.Println(finalPrices([]int{8, 4, 6, 2, 3}))
	//fmt.Println(finalPrices([]int{1, 2, 3, 4, 5}))

	//fmt.Println(canSeePersonsCount([]int{10, 6, 8, 5, 11, 9}))
	//fmt.Println(canSeePersonsCount([]int{5, 1, 2, 3, 10}))

	//fmt.Println(nextLargerNodes(NewMyListNode([]int{2, 1, 5})))
	//fmt.Println(nextLargerNodes(NewMyListNode([]int{2, 7, 4, 3, 5})))

	//fmt.Println(nextGreaterElements([]int{1, 2, 1}))
	//fmt.Println(nextGreaterElements([]int{1, 2, 3, 4, 3}))

	//fmt.Println(dailyTemperatures([]int{73, 74, 75, 71, 69, 72, 76, 73}))
	//fmt.Println(dailyTemperatures([]int{30, 40, 50, 60}))

	//fmt.Println(nextGreaterElement([]int{4, 1, 2}, []int{1, 3, 4, 2}))
	//fmt.Println(nextGreaterElement([]int{2, 4}, []int{1, 2, 3, 4}))

	//fmt.Println(minInsertions("(()))(()))()())))"))
	//fmt.Println(minInsertions("))())("))
	//fmt.Println(minInsertions("(((((("))

	//fmt.Println(minAddToMakeValid("()))(("))
	//fmt.Println(minAddToMakeValid("((("))
	//fmt.Println(minAddToMakeValid("()()()"))

	//fmt.Println(decodeString("3[a2[c]]"))
	//fmt.Println(decodeString("2[abc]3[cd]ef"))

	//fmt.Println(lengthLongestPath("a\n\tb1\n\t\tf1.txt\n\taaaaa\n\t\tf2.txt"))
	//fmt.Println(lengthLongestPath("a"))
	//fmt.Println(lengthLongestPath("dir\n\tsubdir1\n\t\tfile1.ext\n\t\tsubsubdir1\n\tsubdir2\n\t\tsubsubdir2\n\t\t\tfile2.ext"))

	//fmt.Println(evalRPN([]string{"2", "1", "+", "3", "*"}))
	//fmt.Println(evalRPN([]string{"4", "13", "5", "/", "+"}))
	//fmt.Println(evalRPN([]string{"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"}))

	//fmt.Println(isValid("()"))
	//fmt.Println(isValid("()[]{}"))
	//fmt.Println(isValid("([])"))
	//fmt.Println(isValid("([)]"))

	//fmt.Println(simplifyPath("/home/"))
	//fmt.Println(simplifyPath("/home//foo/"))
	//fmt.Println(simplifyPath("/home/user/Documents/../Pictures"))
	//fmt.Println(simplifyPath("/../"))
}
