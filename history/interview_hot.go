package history

import (
	"container/heap"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

/**********************************************************************************************************************/
// 高频面试题
/**********************************************************************************************************************/
// https://leetcode.cn/problems/reverse-linked-list/description/
// 考查递归函数的巧妙，当然也可以用迭代实现
// reverseList 反转链表，递归解法
func reverseList(head *ListNode) *ListNode {
	//base case，空节点或者只有一个节点就不需要反转了
	if head == nil || head.Next == nil {
		return head
	}
	last := reverseList(head.Next)
	head.Next.Next = head
	head.Next = nil // 断掉旧链
	return last
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/reverse-linked-list-ii/description/
// reverseBetween 反转链表 II，递归解法(不容易出错)
func reverseBetween(head *ListNode, left int, right int) *ListNode {
	//base case
	if left == 1 {
		return reverseN(head, right)
	}
	//前进到反转的起始触发base case
	head.Next = reverseBetween(head.Next, left-1, right-1)
	return head
}

// 递归反转前N个节点
func reverseN(head *ListNode, n int) *ListNode {
	var reverse func(head *ListNode, n int) *ListNode
	var tail *ListNode

	reverse = func(head *ListNode, n int) *ListNode {
		if n == 1 {
			tail = head.Next //保存剩下的节点
			return head
		}
		last := reverse(head.Next, n-1)
		head.Next.Next = head //反向指
		head.Next = tail      //指向剩下的节点
		return last
	}
	return reverse(head, n)
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/advantage-shuffle/
// 思路：比得过就比，比不过就用最差的跟对方比保存己方实力
// 输入：numbers1 = [2,7,11,15], numbers2 = [1,10,4,11]
// 输出：[2,11,7,15]
// advantageCount 优势洗牌/田忌赛马
func advantageCount(numbers1 []int, numbers2 []int) []int {
	type Pair struct {
		index int
		value int
	}
	_numbers2 := make([]Pair, 0, len(numbers2))
	for i, value := range numbers2 {
		_numbers2 = append(_numbers2, Pair{index: i, value: value})
	}

	//升序
	sort.Ints(numbers1)
	sort.Slice(_numbers2, func(i, j int) bool { return _numbers2[i].value < _numbers2[j].value })

	left, right := 0, len(numbers1)-1
	res := make([]int, len(numbers1))

	for i := len(_numbers2) - 1; i >= 0; i-- {
		pair := _numbers2[i]
		if numbers1[right] > pair.value {
			res[pair.index] = numbers1[right]
			right--
		} else {
			res[pair.index] = numbers1[left]
			left++
		}
	}
	return res
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/remove-duplicate-letters/description/
// 输入：s = "bcabc"
// 输出："abc"
// 思路：因为要去重，所以需要map，因为不能改变相对顺序，所以只能顺序遍历，不能排序；
// 如何达到字典序最小？可以用单调栈来实现，把从前往后压栈，遇到栈顶字符大于当前字符的就先出栈再入栈(前提：出栈的字符后面还存在)
// removeDuplicateLetters 去除重复字母/不同字符的最小子序列
func removeDuplicateLetters(s string) string {
	inStack := [256]bool{}
	counter := [256]int{}
	for i := 0; i < len(s); i++ {
		counter[s[i]]++
	}

	stack := []uint8{}
	for i := 0; i < len(s); i++ {
		c := s[i]
		counter[c]--
		if inStack[c] {
			//已经存在就跳过（去重）
			continue
		}
		for len(stack) > 0 && stack[len(stack)-1] > c {
			//如果后面没有栈顶元素了就不能再pop了
			if counter[stack[len(stack)-1]] == 0 {
				break
			}
			//后面还有就直接pop，反正后面还会加进来
			inStack[stack[len(stack)-1]] = false
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, c)
		inStack[c] = true
	}
	return string(stack)
}

/**********************************************************************************************************************/
// 二分搜索题型分析案例
// 1、确定x, f(x), target分别是什么，并写出函数f的代码；
// 2、找到x的取值范围作为二分搜索的搜索区间，初始化left和right变量；
// 3、根据题目要求，确实应该使用搜索左侧还是搜索右侧的二分搜索算法，写出解法代码。
/**********************************************************************************************************************/
// https://leetcode.cn/problems/koko-eating-bananas/
// 思路：二分搜索的运用
// 输入：piles = [3,6,7,11], h = 8
// 输出：4
// minEatingSpeed 爱吃香蕉的珂珂
func minEatingSpeed(piles []int, h int) int {
	// x => 吃香蕉的速度，也就是待求解的自变量x
	// f(x) => 若吃香蕉的速度为x根/小时，则需要f(x)个小时吃完所有香蕉
	// target => 吃香蕉的时间限制h就是target
	f := func(x int) int {
		res := 0
		for _, pile := range piles {
			res += pile / x
			if pile%x > 0 {
				res++
			}
		}
		return res
	}

	left, right := 1, 1
	for _, val := range piles {
		if val > right {
			right = val
		}
	}

	for left < right {
		mid := left + (right-left)/2
		if f(mid) == h {
			right = mid
		} else if f(mid) > h {
			left = mid + 1
		} else if f(mid) < h {
			right = mid
		}
	}
	return left
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/capacity-to-ship-packages-within-d-days/description/
// 思路：二分搜索的运用
// 输入：weights = [3,2,2,4,1,4], days = 3
// 输出：6
// 解释：船舶最低载重6就能够在3天内送达所有包裹，如下所示：第1天：3,2、第2天：2,4、第3天：1,4，数组元素顺序不能变
// shipWithinDays 在 D 天内送达包裹的能力
func shipWithinDays(weights []int, days int) int {
	// x => 船的最低运载能力
	// f(x) => 当运载能力为x时，需要f(x)天运送完所有货物
	// target => 运送完所有货物的时间限制days就是target

	f := func(weights []int, x int) int {
		days := 0
		sum := 0
		for i := 0; i < len(weights); i++ {
			sum += weights[i]
			if sum > x {
				days++
				i-- //已经放进来的weights[i]要重新放回去
				sum = 0
			}
		}
		days++
		return days
	}

	//[left,right]是x的取值范围，left是weights中的最大值,right按道理是weights的总和
	var left, right int
	for _, w := range weights {
		right += w
		if w > left {
			left = w
		}
	}

	for left < right {
		mid := left + (right-left)/2
		if f(weights, mid) == days {
			right = mid
		} else if f(weights, mid) > days {
			left = mid + 1
		} else if f(weights, mid) < days {
			right = mid
		}
	}
	return left
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/trapping-rain-water/description/
//思路：站在某一个位置上，看看左边最大/右边最大是多少，然后取最小值减去当前位置高度就是能接的雨水
// 输入：height = [0,1,0,2,1,0,1,3,2,1,2,1]
// 输出：6
// trap 接雨水，双指针求解思路
func _trap(height []int) int {
	if len(height) == 0 {
		return 0
	}

	left, right := 0, len(height)-1
	l_max, r_max, res := 0, 0, 0

	for left < right {
		l_max = max(l_max, height[left])
		r_max = max(r_max, height[right])

		if l_max > r_max {
			res += r_max - height[right]
			right--
		} else {
			res += l_max - height[left]
			left++
		}
	}
	return res
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/container-with-most-water/description/
// 思路：双指针，跟上面👆一题思路差不多的
// 输入：[1,8,6,2,5,4,8,3,7]
// 输出：49
// maxArea 盛最多水的容器
func maxArea(height []int) int {
	if len(height) == 0 {
		return 0
	}

	left, right := 0, len(height)-1
	res := 0

	for left < right {
		res = max(res, min(height[left], height[right])*(right-left))
		if height[left] > height[right] {
			right--
		} else {
			left++
		}
	}
	return res
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/two-sum/
// twoSumTarget 两数之和 III
func twoSumTarget(numbers []int, pos, target int) [][]int {
	result := make([][]int, 0, len(numbers))
	sort.Ints(numbers)

	lo, hi := pos, len(numbers)-1
	for lo < hi {
		left, right := numbers[lo], numbers[hi]
		sum := left + right
		if sum > target {
			for lo < hi && numbers[hi] == right {
				hi--
			}
		} else if sum < target {
			for lo < hi && numbers[lo] == left {
				lo++
			}
		} else if sum == target {
			result = append(result, []int{left, right})
			for lo < hi && numbers[lo] == left {
				lo++
			}
			for lo < hi && numbers[hi] == right {
				hi--
			}
		}
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/3sum/description/
// threeSum 三数之和
func threeSum(numbers []int) [][]int {
	result := make([][]int, 0, len(numbers))
	sort.Ints(numbers)

	for i := 0; i < len(numbers); {
		res := twoSumTarget(numbers, i+1, 0-numbers[i])
		num := numbers[i]

		for _, item := range res {
			item = append(item, num)
			result = append(result, item)
		}
		for i < len(numbers) && numbers[i] == num {
			i++
		}
	}
	return result
}

/**********************************************************************************************************************/
// nSumTarget N数之和
func nSumTarget(numbers []int, N, target int) [][]int {
	var nSum func(numbers []int, n, start, target int) [][]int
	sort.Ints(numbers)

	nSum = func(numbers []int, n, start, target int) [][]int {
		res := [][]int{}
		if n < 2 || len(numbers) < n {
			return res
		}
		if n == 2 { //等于2时，转换为两数之和问题
			lo, hi := start, len(numbers)-1
			for lo < hi {
				left, right := numbers[lo], numbers[hi]
				towSum := left + right
				if towSum == target {
					res = append(res, []int{left, right})
					for lo < hi && numbers[lo] == left {
						lo++
					}
					for lo < hi && numbers[hi] == right {
						hi--
					}
				} else if towSum > target {
					for lo < hi && numbers[hi] == right {
						hi--
					}
				} else {
					for lo < hi && numbers[lo] == left {
						lo++
					}
				}
			}
		} else { //大于2时，递归计算(n-1)Sum的结果
			for i := start; i < len(numbers); {
				curr := numbers[i]
				sum := nSum(numbers, n-1, i+1, target-curr)
				for _, arr := range sum {
					arr = append(arr, curr)
					res = append(res, arr)
				}
				for i < len(numbers) && numbers[i] == curr {
					i++
				}
			}
		}
		return res
	}
	return nSum(numbers, N, 0, target)
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/lowest-common-ancestor-of-a-binary-tree/description/
// lowestCommonAncestor 二叉树的最近公共祖先
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	//不能用val判断，因为值可能重复出现
	if root == p || root == q {
		return root
	}
	left := lowestCommonAncestor(root.Left, p, q)
	right := lowestCommonAncestor(root.Right, p, q)
	if left != nil && right != nil {
		return root
	}
	if left != nil {
		return left
	}
	if right != nil {
		return right
	}
	return nil
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/lowest-common-ancestor-of-a-binary-search-tree/
// 二叉搜索树特点：左小右大
// _lowestCommonAncestor 二叉搜索树的最近公共祖先
func _lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	var find func(root, p, q *TreeNode) *TreeNode

	val1 := min(p.Val, q.Val)
	val2 := max(p.Val, q.Val)

	find = func(root, p, q *TreeNode) *TreeNode {
		if root == nil {
			return nil
		}
		if root.Val > val2 {
			return find(root.Left, p, q)
		}
		if root.Val < val1 {
			return find(root.Right, p, q)
		}

		// val1 < root.Val && root.Val < val2
		return root
	}
	return find(root, p, q)
}

/**********************************************************************************************************************/
/*
二进制基本操作
1、零 异或 任意数 = 任意数；a=0^a=a^0
2、相同整数 异或  = 零；0=a^a
由上面两个推导出：a = a^b^b

3、异或 满足交换定理；
4、移除最后一个1：a=n&(n-1)
5、获取最后一个1：diff=(n&(n-1))^n
6、异或=>相同为0，不相同为1

交换两个数: a=a^b; b=a^b; a=a^b

移除最后一个 1: a=n&(n-1)
获取最后一个 1: diff=(n&(n-1))^n

典型栗子：只出现一次的数字，其余的都出现了K次
*/

/**********************************************************************************************************************/
// https://leetcode.cn/problems/single-number-ii/description/
// 将每个数想象成64位的二进制，对于每一位的二进制的1和0累加起来必然是3N或者3N+1，
// 为3N代表目标值在这一位没贡献，3N+1代表目标值在这一位有贡献(=1)，
// 然后将所有有贡献的位|起来就是结果。这样做的好处是如果题目改成K个一样，只需要把代码改成cnt%k，很通用
// 比如：
// [000001011101100011000011101] nums[0]
// [110001001010100110101010101] nums[1]
// [000001011101100011000011101] nums[2]
// [000001011101100011000011101] nums[3]
// singleNumberII 只出现一次的数字 II
func singleNumberII(numbers []int) int {
	ans := 0
	for i := 0; i < 64; i++ {
		total := 0
		for _, num := range numbers {
			total += num >> i & 1
		}
		if total%3 > 0 {
			ans |= 1 << i
		}
	}
	return ans
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/single-number-iii/submissions/403898275/
// 输入：numbers = [1,2,1,3,2,5]
// 输出：[3,5]
// 解释：[5, 3] 也是有效的答案。
// singleNumberIII 只出现一次的数字 III
func singleNumberIII(numbers []int) []int {
	diff := 0
	for i := 0; i < len(numbers); i++ {
		diff ^= numbers[i] //最后结果diff = 目标值1 ^ 目标值2，因为其他的都异或抵消了
	}
	result := []int{diff, diff}

	diff = (diff & (diff - 1)) ^ diff // 去掉末尾的1后，异或diff就得到最后一个1的位置

	//异或=>相同为0，不相同为1，那么结果diff的最后一位“1”，&上numbers数组，就可以吧数组分为两类
	//一类包含了目标值1，一类包含了目标值2，然后对两个集合全部做异或操作，就可以得到答案。

	for i := 0; i < len(numbers); i++ {
		if diff&numbers[i] == 0 {
			result[0] ^= numbers[i]
		} else {
			result[1] ^= numbers[i]
		}
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/spiral-matrix/description/
// 题目：给你一个 m 行 n 列的矩阵 matrix ，请按照 顺时针螺旋顺序 ，返回矩阵中的所有元素。
// spiralOrder 螺旋矩阵
func spiralOrder(matrix [][]int) []int {
	row, col := len(matrix), len(matrix[0])
	upper_bound, lower_bound := 0, row-1
	left_bound, right_bound := 0, col-1
	res := make([]int, 0, row*col)
	// len(res) == row * col 则遍历完整个数组
	for len(res) < row*col {
		if upper_bound <= lower_bound {
			// 在顶部从左向右遍历
			for j := left_bound; j <= right_bound; j++ {
				res = append(res, matrix[upper_bound][j])
			}
			// 上边界下移
			upper_bound++
		}

		if left_bound <= right_bound {
			// 在右侧从上向下遍历
			for i := upper_bound; i <= lower_bound; i++ {
				res = append(res, matrix[i][right_bound])
			}
			// 右边界左移
			right_bound--
		}

		if upper_bound <= lower_bound {
			// 在底部从右向左遍历
			for j := right_bound; j >= left_bound; j-- {
				res = append(res, matrix[lower_bound][j])
			}
			// 下边界上移
			lower_bound--
		}

		if left_bound <= right_bound {
			// 在左侧从下向上遍历
			for i := lower_bound; i >= upper_bound; i-- {
				res = append(res, matrix[i][left_bound])
			}
			// 左边界右移
			left_bound++
		}
	}
	return res
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/find-minimum-in-rotated-sorted-array/description/
// 思路：最后一个值作为target，然后往左移动，最后比较left、right的值
// findMin 寻找旋转排序数组中的最小值
func findMin(numbers []int) int {
	left, right := 0, len(numbers)-1

	for left < right {
		mid := (right + left) / 2
		target := numbers[right]

		if numbers[mid] < target {
			right = mid
		} else if numbers[mid] >= target {
			left = mid + 1
		}
	}
	return numbers[left]
}

/**********************************************************************************************************************/
// factorialSum 求 30! 以内的阶乘
func factorialSum(n int64) string {
	// import "math/big"
	var calc func(n *big.Int) *big.Int
	b := big.NewInt(n)

	calc = func(n *big.Int) *big.Int {
		if n.Int64() == 1 {
			return n
		} else {
			return n.Mul(n, calc(big.NewInt(n.Int64()-1)))
		}
	}
	return calc(b).String()
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/multiply-strings/description/
// 思路：
// multiply 字符串相乘
func multiply(number1 string, number2 string) string {
	if number1 == "0" || number2 == "0" {
		return "0"
	}
	m, n := len(number1), len(number2)
	array := make([]int, m+n)

	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			array[i+j+1] += int(number2[i]-'0') * int(number1[j]-'0')
		}
	}

	for i := len(array) - 1; i >= 1; i-- {
		array[i-1] += array[i] / 10
		array[i] = array[i] % 10
	}

	var res string
	if array[0] == 0 {
		array = array[1:]
	}
	for i := 0; i < len(array); i++ {
		res += strconv.Itoa(array[i])
	}
	return res
}

/**********************************************************************************************************************/
// strcasecmp 比较两个字符串 忽略大小写 是否相等
func strcasecmp(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	toLower := func(b byte) byte {
		if b >= 'A' && b <= 'Z' {
			b += 32
		}
		return b
	}

	for i := 0; i < len(s1); i++ {
		if toLower(s1[i]) != toLower(s2[i]) {
			return false
		}
	}
	return true
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/restore-ip-addresses/description/
// 思路：考察回溯算法
// restoreIpAddresses 复原 IP 地址
func restoreIpAddresses(s string) []string {
	var backtrack func(track []string, pos int, result *[]string)

	track := []string{}
	result := []string{}

	isValid := func(strNum string) bool {
		if len(strNum) > 1 && strNum[0] == '0' {
			return false
		}
		num, err := strconv.Atoi(strNum)
		if err != nil {
			return false
		}
		if num < 0 || num > 255 {
			return false
		}
		return true
	}

	backtrack = func(track []string, pos int, result *[]string) {
		if len(track) == 4 && pos == len(s) {
			res := strings.Join(track, ".")
			*result = append(*result, res)
			return
		}
		if pos > len(s) || len(track) > 4 {
			return
		}

		for i := 0; i < 3 && i+pos < len(s); i++ {
			strNum := s[pos : i+pos+1]
			if !isValid(strNum) {
				continue
			}
			track = append(track, strNum)
			backtrack(track, i+pos+1, result)
			track = track[:len(track)-1]
		}
	}
	backtrack(track, 0, &result)
	return result
}

/**********************************************************************************************************************/
// 1+2*5-6/2 求结果
func calculate(text string) int {
	getNum := func(i int) (string, int) {
		temp := []byte{}
		for i < len(text) && text[i] >= '0' && text[i] <= '9' {
			temp = append(temp, text[i])
			i++
		}
		i--
		return string(temp), i
	}

	stack := []string{}
	for i := 0; i < len(text); {
		switch text[i] {
		case '*':
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			i++
			right, index := getNum(i)
			i = index

			leftNum, _ := strconv.Atoi(left)
			rightNum, _ := strconv.Atoi(right)
			stack = append(stack, strconv.Itoa(leftNum*rightNum))
		case '/':
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			i++
			right, index := getNum(i)
			i = index

			leftNum, _ := strconv.Atoi(left)
			rightNum, _ := strconv.Atoi(right)
			stack = append(stack, strconv.Itoa(leftNum/rightNum))
		case '+', '-':
			stack = append(stack, string(text[i]))
		default:
			num, index := getNum(i)
			i = index
			stack = append(stack, num)
		}
		i++
	}

	for len(stack) > 1 {
		right := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		op := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		left := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		leftNum, _ := strconv.Atoi(left)
		rightNum, _ := strconv.Atoi(right)

		switch op {
		case "+":
			stack = append(stack, strconv.Itoa(leftNum+rightNum))
		case "-":
			stack = append(stack, strconv.Itoa(leftNum-rightNum))
		}
	}
	if len(stack) > 0 {
		res, _ := strconv.Atoi(stack[0])
		return res
	} else {
		return 0
	}
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/evaluate-reverse-polish-notation/description/
// 输入：tokens = ["2","1","+","3","*"]
// 输出：9
// 解释：该算式转化为常见的中缀算术表达式为：((2 + 1) * 3) = 9
//
// 输入：tokens = ["10","6","9","3","+","-11","*","/","*","17","+","5","+"]
// 输出：22
// 解释：该算式转化为常见的中缀算术表达式为：
//   ((10 * (6 / ((9 + 3) * -11))) + 17) + 5
// = ((10 * (6 / (12 * -11))) + 17) + 5
// = ((10 * (6 / -132)) + 17) + 5
// = ((10 * 0) + 17) + 5
// = (0 + 17) + 5
// = 17 + 5
// = 22
// evalRPN 逆波兰表达式求值
func evalRPN(tokens []string) int {
	if len(tokens) == 0 {
		return 0
	}
	stack := []int{}
	for i := 0; i < len(tokens); i++ {
		switch tokens[i] {
		case "*", "/", "+", "-":
			if len(stack) < 2 {
				return -1
			}
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var result int
			switch tokens[i] {
			case "*":
				result = b * a
			case "/":
				result = b / a
			case "+":
				result = b + a
			case "-":
				result = b - a
			}
			stack = append(stack, result)
		default:
			num, _ := strconv.Atoi(tokens[i])
			stack = append(stack, num)
		}
	}
	return stack[0]
}

/**********************************************************************************************************************/
// 存在数组[1, 5, 3, 6, 9, 7]，求出满足如下条件的元素：
// a. 前面所有数字都比它小；
// b. 后面所有数字都比它大；
func findNumbers(numbers []int) []int {
	type pair struct {
		val, index int
	}
	array := make([]pair, len(numbers))
	for i, val := range numbers {
		array[i] = pair{val: val, index: i}
	}

	sort.Slice(array, func(i, j int) bool {
		return array[i].val < array[j].val
	})

	result := []int{}
	for index, pair := range array {
		if index == pair.index {
			result = append(result, pair.val)
		}
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/g5c51o/description/
// 求前k个高频的元素，按频次从高到低输出
// 输入: [1,6,5,5,5,6,3,4,4,8] k=3
// 输出: [5,6,4]
// topKFrequent 前 K 个高频元素
func topKFrequent(numbers []int, k int) []int {
	mapping := map[int]int{}
	for _, val := range numbers {
		mapping[val]++
	}
	pq := heapPair{}
	heap.Init(&pq)

	for val, cnt := range mapping {
		heap.Push(&pq, [2]int{val, cnt})
	}

	res := []int{}
	for ; k > 0; k-- {
		pair := heap.Pop(&pq).([2]int)
		res = append(res, pair[0])
	}
	return res
}

type heapPair [][2]int //[{val,cnt},{val,cnt}]

func (h heapPair) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h heapPair) Len() int           { return len(h) }
func (h heapPair) Less(i, j int) bool { return h[i][1] > h[j][1] }
func (h *heapPair) Pop() any          { a := *h; pair := a[len(a)-1]; *h = a[:len(a)-1]; return pair }
func (h *heapPair) Push(val any)      { *h = append(*h, val.([2]int)) }

/**********************************************************************************************************************/
// https://leetcode.cn/problems/sort-list/description/
// sortList 排序链表
func sortList(head *ListNode) *ListNode {
	return mergeSortList(head)
}

// 归并排序
func mergeSortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	middle := findMiddleNode(head)
	tail := middle.Next
	middle.Next = nil

	left := mergeSortList(head)
	right := mergeSortList(tail)

	return mergeSubListNode(left, right)
}

// 快慢指针
func findMiddleNode(head *ListNode) *ListNode {
	slow := head
	fast := head.Next
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}
	return slow
}

func mergeSubListNode(left, right *ListNode) *ListNode {
	dummyHead := &ListNode{Val: -1}
	temp := dummyHead
	for left != nil && right != nil {
		if left.Val < right.Val {
			temp.Next = left
			left = left.Next
		} else {
			temp.Next = right
			right = right.Next
		}
		temp = temp.Next
	}
	if left != nil {
		temp.Next = left
	}
	if right != nil {
		temp.Next = right
	}
	return dummyHead.Next
}

/**********************************************************************************************************************/
/*
将一个数组的所有元素向右移动若干单位，并把数组右侧溢出的元素填补
在数组左侧的空缺中，这种操作称为数组的循环平移。
给你一个不少于3 个元素的数组 a， 已知 a 是从一个有序且不包含重复元素的数组循环平移k (k 大于等于 0 且小于数组长度) 个单位而来。
请写一个函数，输 入 int 类型数组 a，返回 k 的值。
例如，对于数组 a = []int{5, 1, 2, 3, 4}，它由有序数组{1, 2, 3, 4, 5}循环平移 1个单位 而来，因此 k = 1。
*/
func findMinNumIndex(numbers []int) int {
	left, right := 0, len(numbers)-1
	for left < right {
		mid := left + (right-left)/2
		target := numbers[right]

		if numbers[mid] < target {
			right = mid
		} else if numbers[mid] >= target {
			left = mid + 1
		}
	}
	return left
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/SsGoHC/description/
// 先排序，再合并
// 输入：intervals = [[1,3],[2,6],[8,10],[15,18]]
// 输出：[[1,6],[8,10],[15,18]]
// 解释：区间 [1,3] 和 [2,6] 重叠, 将它们合并为 [1,6].
// mergeRange 合并区间
func mergeRange(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	result := [][]int{}
	result = append(result, intervals[0])

	for i := 1; i < len(intervals); i++ {
		start := intervals[i][0]
		end := intervals[i][1]
		if start > result[len(result)-1][1] {
			result = append(result, intervals[i])
		} else {
			maxNum := max(result[len(result)-1][1], end)
			result[len(result)-1][1] = maxNum
		}
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/maximum-width-of-binary-tree/description/
// 思路：对节点进行编号。一个编号为index的左子节点的编号记为 2×index2，右子节点的编号记为2×index+1计算每层宽度时，用每层节点的最大编号减去最小编号再加1即为宽度。
// widthOfBinaryTree 二叉树最大宽度
func widthOfBinaryTree(root *TreeNode) int {
	type pair struct {
		node  *TreeNode
		index int
	}

	maxWidth := 0
	queue := []pair{}
	queue = append(queue, pair{node: root, index: 1})

	for len(queue) > 0 {
		maxWidth = max(maxWidth, queue[len(queue)-1].index-queue[0].index+1)
		sz := len(queue)
		for i := 0; i < sz; i++ {
			pa := queue[0]
			queue = queue[1:]
			if pa.node.Left != nil {
				queue = append(queue, pair{pa.node.Left, pa.index * 2})
			}
			if pa.node.Right != nil {
				queue = append(queue, pair{pa.node.Right, pa.index*2 + 1})
			}
		}
	}
	return maxWidth
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/he-wei-sde-lian-xu-zheng-shu-xu-lie-lcof/description/
// 输入：numbers[1~target-1] ,target = 12
// 输出：[[3, 4, 5]]
// 解释：在上述示例中，存在一个连续正整数序列的和为 12，为 [3, 4, 5]。
// fileCombination 文件组合
func fileCombination(target int) [][]int {
	var backtrack func(track []int, pos int, result *[][]int)

	result := [][]int{}
	track := []int{}

	backtrack = func(track []int, pos int, result *[][]int) {
		totalSum := sum(track)
		if totalSum == target {
			temp := make([]int, len(track))
			copy(temp, track)
			*result = append(*result, temp)
		}
		if totalSum > target {
			return
		}

		for i := pos; i < target-1; i++ {
			if len(track) > 0 && track[len(track)-1]+1 != i {
				break
			}
			track = append(track, i)
			backtrack(track, i+1, result)
			track = track[:len(track)-1]
		}
	}
	backtrack(track, 1, &result)
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/find-largest-value-in-each-tree-row/description/
// largestValues 在每个树行中找最大值
func largestValues(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	queue := []*TreeNode{root}
	result := []int{}
	for len(queue) > 0 {
		sz := len(queue)
		maxVal := math.MinInt
		for i := 0; i < sz; i++ {
			node := queue[0]
			queue = queue[1:]
			maxVal = max(maxVal, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		result = append(result, maxVal)
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/QTMn0o/description/
// 给定一个整数数组和一个整数 k ，请找到该数组中和为 k 的连续子数组的个数。
// subArraySum 和为 K 的子数组
func subArraySum(numbers []int, k int) int {
	var backtrack func(track []int, pos int)

	subArrayCnt := 0
	track := []int{}

	sum := func(indexs []int) int {
		total := 0
		for _, index := range indexs {
			total += numbers[index]
		}
		return total
	}

	backtrack = func(track []int, pos int) {
		totalSum := sum(track)
		if len(track) > 0 && totalSum == k {
			subArrayCnt++
			return
		}
		if totalSum > k {
			return
		}
		for i := pos; i < len(numbers); i++ {
			if len(track) > 0 && track[len(track)-1]+1 != i {
				break
			}
			track = append(track, i)
			backtrack(track, i+1)
			track = track[0 : len(track)-1]
		}
	}
	backtrack(track, 0)
	return subArrayCnt
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/subarray-sum-equals-k/description/
// 连续子数组求和 => 前缀和数组
// subarraySum 和为 K 的子数组(数组元素可能为负数)
func subarraySum(numbers []int, k int) int {
	preSum := make([]int, len(numbers)+1)
	for i := 1; i <= len(numbers); i++ {
		preSum[i] = preSum[i-1] + numbers[i-1]

	}

	subArrayCnt := 0
	for i := 0; i < len(preSum); i++ {
		for j := i + 1; j < len(preSum); j++ {
			if preSum[j]-preSum[i] == k {
				subArrayCnt++
			}
		}
	}
	return subArrayCnt
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock/description/
// 思路：1、股票下跌的时候记录最小值
// 2、股票上市的时候计算收益值
// maxProfit 买卖股票的最佳时机
func maxProfit(prices []int) int {
	maxPrice := 0
	minNum := math.MaxInt

	for _, price := range prices {
		if price < minNum { //股票下跌的时候记录最小值
			minNum = price
		}
		if price > minNum { //股票上市的时候计算收益值
			maxPrice = max(maxPrice, price-minNum)
		}
	}
	return maxPrice
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/word-break/description/
// 思路：动态规划，dp[i]表示字符串的前i个位置可以由wordDict拼接
// 输入: s = "leetcode", wordDict = ["leet", "code"]
// 输出: true
// 解释: 返回 true 因为 "leetcode" 可以由 "leet" 和 "code" 拼接成。
// wordBreak 单词拆分
func wordBreak(s string, wordDict []string) bool {
	mapping := map[string]bool{}
	for _, key := range wordDict {
		mapping[key] = true
	}

	dp := make(map[int]bool, len(s)+1)
	dp[0] = true

	for i := 1; i <= len(s); i++ {
		for j := 0; j < i; j++ {
			if dp[j] && mapping[s[j:i]] {
				dp[i] = true
			}
		}
	}
	return dp[len(s)]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/majority-element/description/
// 多数元素是指在数组中出现次数 大于 ⌊ n/2 ⌋ 的元素，那么排序之后的中间位置即为该数
// majorityElement 多数元素
func majorityElement(numbers []int) int {
	sort.Ints(numbers)
	return numbers[len(numbers)/2]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/sort-colors/description/
// 思路：双指针，遇到0放前面，遇到2放后面
// 输入：numbers = [2,0,2,1,1,0]
// 输出：[0,0,1,1,2,2]
// sortColors 颜色分类
func sortColors(numbers []int) {
	p0, p2 := 0, len(numbers)-1
	for i := 0; i <= p2; i++ {
		for ; i <= p2 && numbers[i] == 2; p2-- {
			numbers[i], numbers[p2] = numbers[p2], numbers[i]
		}
		if numbers[i] == 0 {
			numbers[i], numbers[p0] = numbers[p0], numbers[i]
			p0++
		}
	}
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/next-permutation/description/
// 输入：numbers = [4, 5, 2, 6, 3, 1]
// 输出：[4, 5, 3, 1, 2, 6]
// nextPermutation 下一个排列（字典序升序）
func nextPermutation(numbers []int) {
	// 从后往前找到第一个升序对[i,j]，然后在[j,end]找到第一个大于arr[i]的数
	minIndex, maxIndex := 0, 0
	for i := len(numbers) - 1; i > 0; i-- {
		if numbers[i-1] < numbers[i] {
			minIndex = i - 1
			for j := len(numbers) - 1; j > minIndex; j-- {
				if numbers[j] > numbers[minIndex] {
					maxIndex = j
					break
				}
			}
			break
		}
	}

	if minIndex == 0 && maxIndex == 0 { //最大字典序
		for i, j := 0, len(numbers)-1; i < j; {
			numbers[i], numbers[j] = numbers[j], numbers[i]
			i++
			j--
		}
	} else {
		numbers[minIndex], numbers[maxIndex] = numbers[maxIndex], numbers[minIndex]
		//然后对[j,end]进行升序排列
		for i, j := minIndex+1, len(numbers)-1; i < j; {
			numbers[i], numbers[j] = numbers[j], numbers[i]
			i++
			j--
		}
	}
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/find-the-duplicate-number/
// 输入：numbers = [3,1,3,4,2]
// 输出：3
// findDuplicate 寻找重复数
func findDuplicate(numbers []int) int {
	sort.Ints(numbers)
	for i := 0; i < len(numbers)-1; i++ {
		if numbers[i] == numbers[i+1] {
			return numbers[i]
		}
	}
	return -1
}

// 位操作
func _findDuplicate(numbers []int) int {
	ans, bit_max := 0, 31
	for bit := 0; bit <= bit_max; bit++ {
		x, y := 0, 0
		for i := 0; i < len(numbers); i++ {
			if (numbers[i] & (1 << bit)) > 0 {
				x++
			}
			if i&(1<<bit) > 0 {
				y++
			}
		}
		if x > y {
			ans |= 1 << bit
		}
	}
	return ans
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/group-anagrams/description/
// 输入: strs = ["eat", "tea", "tan", "ate", "nat", "bat"]
// 输出: [["bat"],["nat","tan"],["ate","eat","tea"]]
// groupAnagrams 字母异位词分组
func groupAnagrams(strs []string) [][]string {
	mapping := map[string][]string{}
	for i := 0; i < len(strs); i++ {
		byteArr := []byte(strs[i])
		sort.Slice(byteArr, func(i, j int) bool { return byteArr[i] > byteArr[j] })
		str := string(byteArr)
		mapping[str] = append(mapping[str], strs[i])
	}

	result := [][]string{}
	for _, group := range mapping {
		result = append(result, group)
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/binary-tree-level-order-traversal/
// levelOrder 二叉树的层序遍历
func levelOrder(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}
	result := [][]int{}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		sz := len(queue)
		levelArr := []int{}
		for i := 0; i < sz; i++ {
			node := queue[0]
			queue = queue[1:]

			levelArr = append(levelArr, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		result = append(result, levelArr)
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/binary-tree-zigzag-level-order-traversal/
// 二叉树的锯齿形层序遍历
func zigzagLevelOrder(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}
	result := [][]int{}
	queue := []*TreeNode{root}
	isSwap := false
	for len(queue) > 0 {
		sz := len(queue)
		levelArr := []int{}
		for i := 0; i < sz; i++ {
			node := queue[0]
			queue = queue[1:]
			levelArr = append(levelArr, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		if isSwap {
			for left, right := 0, len(levelArr)-1; left < right; {
				levelArr[left], levelArr[right] = levelArr[right], levelArr[left]
				left++
				right--
			}
		}
		isSwap = !isSwap
		result = append(result, levelArr)
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/rotate-image/description/
// 思路：二维矩阵先上下翻转，然后在左对角线翻转
// rotate 旋转图像
func rotate(matrix [][]int) {
	upper_row, down_row := 0, len(matrix)-1
	for upper_row < down_row {
		for col := 0; col < len(matrix[0]); col++ {
			matrix[upper_row][col], matrix[down_row][col] = matrix[down_row][col], matrix[upper_row][col]
		}
		upper_row++
		down_row--
	}

	for row := 0; row < len(matrix); row++ {
		for col := row + 1; col < len(matrix[0]); col++ {
			matrix[row][col], matrix[col][row] = matrix[col][row], matrix[row][col]
		}
	}
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/unique-paths-ii/
// 思路：动态规划
// uniquePathsWithObstacles 不同路径 II
func uniquePathsWithObstacles(obstacleGrid [][]int) int {
	if obstacleGrid[0][0] == 1 {
		return 0
	}
	dp := make([][]int, len(obstacleGrid))
	for i := 0; i < len(obstacleGrid); i++ {
		dp[i] = make([]int, len(obstacleGrid[0]))
	}

	//base case
	dp[0][0] = 1
	for i := 1; i < len(obstacleGrid); i++ {
		if obstacleGrid[i][0] == 0 && dp[i-1][0] == 1 {
			dp[i][0] = 1
		}
	}
	for i := 1; i < len(obstacleGrid[0]); i++ {
		if obstacleGrid[0][i] == 0 && dp[0][i-1] == 1 {
			dp[0][i] = 1
		}
	}

	// dp[i][j] = dp[i-1][j] + dp[i][j-1]
	for i := 1; i < len(obstacleGrid); i++ {
		for j := 1; j < len(obstacleGrid[0]); j++ {
			if obstacleGrid[i][j] == 0 {
				dp[i][j] = dp[i-1][j] + dp[i][j-1]
			}
		}
	}

	return dp[len(obstacleGrid)-1][len(obstacleGrid[0])-1]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/remove-duplicates-from-sorted-list-ii/
// deleteDuplicatesII 删除排序链表中的重复元素 II
func deleteDuplicatesII(head *ListNode) *ListNode {
	dummy := &ListNode{}
	dummy.Next = head
	head = dummy
	for head.Next != nil && head.Next.Next != nil {
		if head.Next.Val == head.Next.Next.Val {
			//把要删除的val暂存一下
			remVal := head.Next.Val
			for head.Next != nil && head.Next.Val == remVal {
				head.Next = head.Next.Next
			}
		} else {
			head = head.Next
		}
	}
	return dummy.Next
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/partition-list/
// 思路：将大于x的节点，放到另外一个链表，最后连接这两个链表
// partition_ 分隔链表
func partition_(head *ListNode, x int) *ListNode {
	headDummy := &ListNode{Val: 0}
	headDummy.Next = head
	head = headDummy

	tailDummy := &ListNode{Val: 0}
	tail := tailDummy

	for head.Next != nil {
		if head.Next.Val < x {
			head = head.Next
		} else {
			// 移除大于等于x的节点
			del := head.Next
			head.Next = head.Next.Next
			// 放到另外一个链表
			tail.Next = del
			tail = tail.Next
		}
	}
	// 拼接两个链表
	tail.Next = nil
	head.Next = tailDummy.Next
	return headDummy.Next
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/clone-graph/description/
// cloneGraph 克隆图
type GraphNode struct {
	Val       int
	Neighbors []*GraphNode
}

func cloneGraph(node *GraphNode) *GraphNode {
	visited := make(map[*GraphNode]*GraphNode)
	return clone(node, visited)
}

// 递归克隆，传入已经访问过的元素作为过滤条件
func clone(node *GraphNode, visited map[*GraphNode]*GraphNode) *GraphNode {
	if node == nil {
		return nil
	}
	// 已经访问过直接返回
	if v, ok := visited[node]; ok {
		return v
	}
	newNode := &GraphNode{
		Val:       node.Val,
		Neighbors: make([]*GraphNode, len(node.Neighbors)),
	}
	visited[node] = newNode
	for i := 0; i < len(node.Neighbors); i++ {
		newNode.Neighbors[i] = clone(node.Neighbors[i], visited)
	}
	return newNode
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/decode-string/description/
// 输入：s = "3[a]2[bc]"
// 输出："aaabcbc"
// 输入：s = "2[abc]3[cd]ef"
// 输出："abcabccdcdcdef"
// decodeString 字符串解码
func decodeString(s string) string {
	stack := []string{}
	for i := 0; i < len(s); i++ {
		ss := string(s[i])
		if ss == "]" {
			begin := len(stack) - 1
			str := ""
			for ; stack[begin] != "["; begin-- {
				str = stack[begin] + str
			}
			begin--
			numStr := ""
			for ; begin >= 0 && stack[begin] >= "0" && stack[begin] <= "9"; begin-- {
				numStr = stack[begin] + numStr
			}
			begin++
			stack = stack[:begin]
			n, _ := strconv.Atoi(numStr)
			outStr := ""
			for i := 0; i < n; i++ {
				outStr += str
			}
			stack = append(stack, outStr)

		} else {
			stack = append(stack, ss)
		}
	}
	out := ""
	for i := 0; i < len(stack); i++ {
		out += stack[i]
	}
	return out
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/reorder-list/description/
// 思路：找到中间节点，断链，然后把后面半截链表反转，然后再把两个链表交替合成一个链表
// reorderList 重排链表
func reorderList(head *ListNode) {
	if head == nil {
		return
	}
	//找到中间节点
	mid := findMiddleNode(head)
	//反转后面半截
	tail := reverseNode(mid.Next)
	//断链
	mid.Next = nil
	//交替合并链表
	head = mergeNode(head, tail)
}

// 反转链表
func reverseNode(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	last := reverseNode(head.Next)
	head.Next.Next = head
	head.Next = nil
	return last
}

func mergeNode(list1, list2 *ListNode) *ListNode {
	dummy := &ListNode{Val: -1}
	head := dummy
	toggle := true
	for list1 != nil && list2 != nil {
		if toggle {
			head.Next = list1
			list1 = list1.Next
		} else {
			head.Next = list2
			list2 = list2.Next
		}
		toggle = !toggle
		head = head.Next
	}
	if list1 != nil {
		head.Next = list1
	}
	if list2 != nil {
		head.Next = list2
	}
	return dummy.Next
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/01-matrix/
// 输入：mat = [[0,0,0],[0,1,0],[1,1,1]]
// 输出：[[0,0,0],[0,1,0],[1,2,1]]
// updateMatrix 01矩阵
func updateMatrix(matrix [][]int) [][]int {
	q := make([][]int, 0)
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			if matrix[i][j] == 0 {
				// 进队列
				point := []int{i, j}
				q = append(q, point)
			} else {
				matrix[i][j] = -1
			}
		}
	}
	//     [-1,0]
	//[0,-1][0,0][0,1]
	//      [1,0]
	directions := [][]int{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}
	for len(q) != 0 {
		// 出队列
		point := q[0]
		q = q[1:]
		for _, v := range directions {
			x := point[0] + v[0]
			y := point[1] + v[1]
			if x >= 0 && x < len(matrix) && y >= 0 && y < len(matrix[0]) && matrix[x][y] == -1 {
				matrix[x][y] = matrix[point[0]][point[1]] + 1
				// 将当前的元素进队列，进行一次BFS
				q = append(q, []int{x, y})
			}
		}
	}
	return matrix
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/bitwise-and-of-numbers-range/description/
// 输入：left = 5, right = 7
// 输出：4
// rangeBitwiseAnd 数字范围按位与
func rangeBitwiseAnd(m int, n int) int {
	// m 5 1 0 1
	//   6 1 1 0
	// n 7 1 1 1
	// 把可能包含0的全部右移变成
	// m 5 1 0 0
	//   6 1 0 0
	// n 7 1 0 0
	// 所以最后结果就是m<<count
	var count int
	for m != n {
		m >>= 1
		n >>= 1
		count++
	}
	return m << count
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/insert-into-a-binary-search-tree/
// insertIntoBST 二叉搜索树中的插入操作
func insertIntoBST(root *TreeNode, val int) *TreeNode {
	if root == nil {
		root = &TreeNode{Val: val}
		return root
	}
	if root.Val > val {
		root.Left = insertIntoBST(root.Left, val)
	} else {
		root.Right = insertIntoBST(root.Right, val)
	}
	return root
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/swap-nodes-in-pairs/
// swapPairs 两两交换链表中的节点
func swapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	nextHead := head.Next.Next

	head.Next.Next = head
	head = head.Next
	head.Next.Next = swapPairs(nextHead)
	return head
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/copy-list-with-random-pointer/description/
type _Node struct {
	Val    int
	Random *_Node
	Next   *_Node
}

// copyRandomList 随机链表的复制
func copyRandomList(head *_Node) *_Node {
	var build func(head *_Node) *_Node

	//解题关键在这个map的定义
	visited := map[*_Node]*_Node{}

	build = func(head *_Node) *_Node {
		if head == nil {
			return nil
		}
		if node, ok := visited[head]; ok {
			return node
		}
		root := &_Node{Val: head.Val}
		visited[head] = root

		root.Next = build(head.Next)
		root.Random = build(head.Random)
		return root
	}
	return build(head)
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/min-stack/description/
// 最小栈；设计一个支持 push ，pop ，top 操作，并能在常数时间内检索到最小元素的栈。
type MinStack struct {
	rawStack []int //原始栈
	minStack []int //对应索引处的最小值，没有则不放入
}

// 初始化堆栈对象。
func Constructor10() MinStack {
	return MinStack{
		rawStack: make([]int, 0),
		minStack: make([]int, 0),
	}
}

// 将元素val推入堆栈。
func (this *MinStack) Push(val int) {
	this.rawStack = append(this.rawStack, val)
	if len(this.minStack) == 0 {
		this.minStack = append(this.minStack, val)
	} else {
		if this.minStack[len(this.minStack)-1] >= val {
			this.minStack = append(this.minStack, val)
		}
	}
}

// 删除堆栈顶部的元素。
func (this *MinStack) Pop() {
	if len(this.rawStack) == 0 {
		return
	}
	if this.minStack[len(this.minStack)-1] == this.Top() {
		this.minStack = this.minStack[:len(this.minStack)-1]
	}
	this.rawStack = this.rawStack[:len(this.rawStack)-1]
}

// 获取堆栈顶部的元素。
func (this *MinStack) Top() int {
	if len(this.rawStack) == 0 {
		panic("empty stack")
	}
	return this.rawStack[len(this.rawStack)-1]
}

// 获取堆栈中的最小元素。
func (this *MinStack) GetMin() int {
	return this.minStack[len(this.minStack)-1]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/implement-queue-using-stacks/description/
// 用栈实现队列
type MyQueue struct {
	inStack, outStack []int
}

func Constructor11() MyQueue {
	return MyQueue{
		inStack:  []int{},
		outStack: []int{},
	}
}

// 将元素 x 推到队列的末尾
func (this *MyQueue) Push(x int) {
	this.inStack = append(this.inStack, x)
}

// 从队列的开头移除并返回元素
func (this *MyQueue) Pop() int {
	if len(this.outStack) == 0 {
		this.in2out()
	}
	val := this.outStack[len(this.outStack)-1]
	this.outStack = this.outStack[:len(this.outStack)-1]
	return val
}

// 返回队列开头的元素
func (this *MyQueue) Peek() int {
	if len(this.outStack) == 0 {
		this.in2out()
	}
	return this.outStack[len(this.outStack)-1]
}

// 如果队列为空，返回 true ；否则，返回 false
func (this *MyQueue) Empty() bool {
	return len(this.inStack)+len(this.outStack) == 0
}

func (this *MyQueue) in2out() {
	for len(this.inStack) > 0 {
		this.outStack = append(this.outStack, this.inStack[len(this.inStack)-1])
		this.inStack = this.inStack[:len(this.inStack)-1]
	}
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/search-a-2d-matrix/description/
// 给你一个满足下述两条属性的 m x n 整数矩阵：
// 每行中的整数从左到右按非严格递增顺序排列。
// 每行的第一个整数大于前一行的最后一个整数。
// 给你一个整数 target ，如果 target 在矩阵中，返回 true ；否则，返回 false 。
//
// 输入：matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 3
// 输出：true
// searchMatrix 搜索二维矩阵
func searchMatrix(matrix [][]int, target int) bool {
	M, N := len(matrix), len(matrix[0])
	left, right := 0, M*N-1

	for left <= right {
		mid := left + (right-left)/2
		row, col := mid/N, mid%N

		if matrix[row][col] == target {
			return true
		} else if matrix[row][col] > target {
			right = mid - 1
		} else if matrix[row][col] < target {
			left = mid + 1
		}
	}
	return false
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/find-minimum-in-rotated-sorted-array-ii/description/
// 考察：存在重复元素的二分查找，重复的元素++/--跳过
// 输入：numbers = [2,2,2,0,1]
// 输出：0
// findMinII 寻找旋转排序数组中的最小值 II
func findMinII(numbers []int) int {
	left, right := 0, len(numbers)-1
	for left < right {
		mid := left + (right-left)/2
		target := numbers[right]

		if numbers[mid] > target {
			left = mid + 1
		} else if numbers[mid] < target {
			right = mid
		} else if numbers[mid] == target {
			right--
		}
	}
	return numbers[right]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/search-in-rotated-sorted-array-ii/description/
// 思路：两条上升直线，四种情况判断，并且处理重复数字
// 输入：numbers = [2,5,6,0,0,1,2], target = 0
// 输出：true
// searchII 搜索旋转排序数组 II
func searchII(nums []int, target int) bool {
	if len(nums) == 0 {
		return false
	}
	left, right := 0, len(nums)-1

	for left <= right {
		// 处理重复数字
		for left < right && nums[left] == nums[left+1] {
			left++
		}
		for left < right && nums[right] == nums[right-1] {
			right--
		}

		mid := (left + right) / 2
		if nums[mid] == target {
			return true
		}

		// 先判断在那个区间
		if nums[left] <= nums[mid] {
			if nums[left] <= target && target < nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else if nums[left] > nums[mid] {
			if nums[mid] < target && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}
	return false
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/unique-binary-search-trees-ii/description/
// 思路：二叉搜索树关键的性质是根节点的值大于左子树所有节点的值，小于右子树所有节点的值，且左子树和右子树也同样为二叉搜索树。
// 因此在生成所有可行的二叉搜索树的时候，假设当前序列长度为 n，如果我们枚举根节点的值为 i，
// 那么根据二叉搜索树的性质我们可以知道左子树的节点值的集合为[1…i−1]，右子树的节点值的集合为 [i+1…n]。
// 而左子树和右子树的生成相较于原问题是一个序列长度缩小的子问题，因此我们可以想到用回溯的方法来解决这道题目。
// generateTrees 不同的二叉搜索树 II
func generateTrees(n int) []*TreeNode {
	return buildBinTree(1, n)
}

func buildBinTree(left, right int) []*TreeNode {
	if left > right { //注意区间范围，这里“==”不能返回
		return []*TreeNode{nil}
	}

	allTree := []*TreeNode{}
	for i := left; i <= right; i++ {
		leftTree := buildBinTree(left, i-1)
		rightTree := buildBinTree(i+1, right)

		for _, lTree := range leftTree {
			for _, rTree := range rightTree {
				root := &TreeNode{Val: i}
				root.Left = lTree
				root.Right = rTree
				allTree = append(allTree, root)
			}
		}
	}
	return allTree
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/delete-node-in-a-bst/
// 思路：
// 站在删除节点的父节点角度看，分为三种情况：
// 1、删除节点的左子树为空，根节点替换为右孩子
// 2、删除节点的右子树为空，根节点替换为左孩子
// 3、删除节点同时存在左右子树，左子树连接到右边最左节点即可
// deleteNode 删除二叉搜索树中的节点
func deleteNode(root *TreeNode, key int) *TreeNode {
	if root == nil {
		return root
	}
	switch {
	case key > root.Val:
		root.Right = deleteNode(root.Right, key)
	case key < root.Val:
		root.Left = deleteNode(root.Left, key)
	case root.Val == key:
		if root.Left == nil {
			return root.Right
		} else if root.Right == nil {
			return root.Left
		} else {
			minLeftNode := root.Right
			// 一直向左找到最后一个左节点即可
			for minLeftNode.Left != nil {
				minLeftNode = minLeftNode.Left
			}
			minLeftNode.Left = root.Left
			return root.Right
		}
	}
	return root
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/triangle/
// 2
// 3 4
// 6 5 7
// 4 1 8 3
// 思路：动态规划
// minimumTotal 三角形最小路径和
func minimumTotal(triangle [][]int) int {
	if len(triangle) == 0 || len(triangle[0]) == 0 {
		return 0
	}
	dp := make([][]int, len(triangle))
	for i := 0; i < len(triangle); i++ {
		for j := 0; j < len(triangle[i]); j++ {
			if dp[i] == nil {
				dp[i] = make([]int, len(triangle[i]))
			}
			dp[i][j] = triangle[i][j]
		}
	}

	for i := len(triangle) - 2; i >= 0; i-- {
		for j := 0; j < len(triangle[i]); j++ {
			dp[i][j] = min(dp[i+1][j], dp[i+1][j+1]) + triangle[i][j]
		}
	}
	return dp[0][0]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/longest-consecutive-sequence/
// 给定一个未排序的整数数组 numbers ，找出数字连续的最长序列（不要求序列元素在原数组中连续）的长度。
// 请你设计并实现时间复杂度为 O(n) 的算法解决此问题。
// 输入：numbers = [100,4,200,1,3,2]
// 输出：4
// 解释：最长数字连续序列是 [1, 2, 3, 4]。它的长度为 4。
// longestConsecutive 最长连续序列
func longestConsecutive(numbers []int) int {
	// 放到map里面，就可以从头到尾遍历，然后map判断是否连续
	// 连续则for++，不连续，则保存序列长度，然后跳过已经遍历过的元素
	numSet := map[int]bool{}
	for _, num := range numbers {
		numSet[num] = true
	}
	maxLen := 0
	for num := range numSet {
		if !numSet[num-1] {
			curr := num
			long := 1
			for numSet[curr+1] {
				long++
				curr++
			}
			if long > maxLen {
				maxLen = long
			}
		}
	}
	return maxLen
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/unique-paths/
// 思路：动态规划
// uniquePaths 不同路径
func uniquePaths(m int, n int) int {
	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dp[i][0] = 1
			dp[0][j] = 1
		}
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = dp[i-1][j] + dp[i][j-1]
		}
	}
	return dp[m-1][n-1]
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/jump-game-ii/
// 思路：这道题是典型的贪心算法，通过局部最优解得到全局最优解。
// 输入: numbers = [2,3,1,1,4]
// 输出: 2
// 解释: 跳到最后一个位置的最小跳跃数是 2。从下标为 0 跳到下标为 1 的位置，跳 1 步，然后跳 3 步到达数组的最后一个位置。
// jump 跳跃游戏 II
func jump(numbers []int) int {
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	jumps := 0    //记录跳跃的次数
	end := 0      //记录一次起跳的最远距离下标
	longJump := 0 //记录每走一步可到达的最远距离下标

	for i := 0; i < len(numbers)-1; i++ {
		longJump = max(numbers[i]+i, longJump)
		//超过一次起跳的最远距离后，就需要再次起跳
		if end == i {
			jumps++
			end = longJump
		}
	}
	return jumps
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/repeated-dna-sequences/
// 输入：s = "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"
// 输出：["AAAAACCCCC","CCCCCAAAAA"]
// findRepeatedDnaSequences 重复的DNA序列
func findRepeatedDnaSequences(s string) []string {
	mapExist := map[string]int{}
	result := []string{}

	for i := 0; i+10 <= len(s); i++ {
		sub := s[i : i+10]
		mapExist[sub]++
		if mapExist[sub] == 2 {
			result = append(result, sub)
		}
	}
	return result
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/reverse-nodes-in-k-group/
// 思路：K个一组翻转链表，然后递归执行
// 输入：head = [1,2,3,4,5], k = 2
// 输出：[2,1,4,3,5]
// reverseKGroup K个一组翻转链表
func reverseKGroup(head *ListNode, k int) *ListNode {
	if head == nil {
		return head
	}
	var reverse func(a, b *ListNode) *ListNode

	// 反转区间 [a, b) 的元素，注意是左闭右开
	reverse = func(a, b *ListNode) *ListNode {
		var tail *ListNode
		if a.Next == b {
			tail = a.Next
			return a
		}
		last := reverse(a.Next, b)
		a.Next.Next = a
		a.Next = tail
		return last
	}

	// 区间 [a, b) 包含 k 个待反转元素
	a, b := head, head
	for i := 0; i < k; i++ {
		// 不足 k 个，不需要反转，base case
		if b == nil {
			return head
		}
		b = b.Next
	}

	newHead := reverse(a, b)     // 反转前 k 个元素
	a.Next = reverseKGroup(b, k) // 递归反转后续链表并连接起来
	return newHead
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/spiral-matrix-ii/description/
// generateMatrix 螺旋矩阵 II
func generateMatrix(n int) [][]int {
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
	}

	upper, lower := 0, n-1
	left, right := 0, n-1

	number := 1
	for number <= n*n {
		// 在顶部从左向右遍历
		if upper <= lower {
			for i := left; i <= right; i++ {
				matrix[upper][i] = number
				number++
			}
			upper++ // 上边界下移
		}
		// 在右侧从上向下遍历
		if left <= right {
			for i := upper; i <= lower; i++ {
				matrix[i][right] = number
				number++
			}
			right-- // 右边界左移
		}
		// 在底部从右向左遍历
		if upper <= lower {
			for i := right; i >= left; i-- {
				matrix[lower][i] = number
				number++
			}
			lower-- // 下边界上移
		}
		// 在左侧从下向上遍历
		if left <= right {
			for i := lower; i >= upper; i-- {
				matrix[i][left] = number
				number++
			}
			left++ // 左边界右移
		}
	}
	return matrix
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/jump-game/description/
// 输入：numbers = [3,2,1,0,4]
// 输出：false
// 解释：无论怎样，总会到达下标为 3 的位置。但该下标的最大跳跃长度是 0 ， 所以永远不可能到达最后一个下标。
// canJump 跳跃游戏
func canJump(numbers []int) bool {
	longIndex := 0
	for i := 0; i < len(numbers); i++ {
		longIndex = max(longIndex, numbers[i]+i)
		if longIndex >= len(numbers)-1 {
			return true
		}
		if longIndex <= i {
			return false
		}
	}
	return false
}

/**********************************************************************************************************************/
// https://leetcode.cn/problems/palindrome-linked-list/description/
// 思路：找到中间节点，然后翻转，然后对比
// isPalindrome 回文链表
func isPalindrome(head *ListNode) bool {
	var reverse func(head *ListNode) *ListNode
	reverse = func(head *ListNode) *ListNode {
		if head == nil {
			return head
		}
		last := reverse(head.Next)
		head.Next.Next = head
		head.Next = nil
		return last
	}

	slow := head
	fast := head
	for fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	tail := slow.Next
	slow.Next = nil
	tail = reverse(tail)

	for head != nil && tail != nil {
		if head.Val != tail.Val {
			return false
		}
		head = head.Next
		tail = tail.Next
	}
	return true
}
