package main

import "fmt"

// 贪心算法实现：
// 🌟技巧1：局部最优推全局最优 - 贪心的核心是每一步都选择当前看起来最好的选择，无需回溯，通过数学归纳法证明局部最优能推导出全局最优
// 🌟技巧2：维护"能到达的最远位置" - 跳跃类问题的贪心关键是维护一个farthest变量，遍历每个位置时更新可达的最远索引，不需要枚举所有跳法
// 🌟技巧3：区间贪心框架 - 跳跃游戏II的精髓是用[i, end]表示"当前步数可达的索引区间"，当i==end时说明当前步数的区间已遍历完，必须再跳一步
// 🌟技巧4：贪心 vs 动态规划 - 贪心是DP的特例：当DP的状态转移中，每次都只需要取最优子问题而非所有子问题时，可以用贪心替代DP，将O(n²)优化为O(n)

// ⚠️易错点1：循环边界是len(nums)-1 - 遍历时不需要处理最后一个元素，因为目标是「到达」最后位置，不需要从最后位置出发再跳
// ⚠️易错点2：farthest <= i 而非 farthest < i - canJump中判断卡住的条件是farthest<=i（最远只能到当前位置，无法继续前进），写成<会漏掉恰好卡住的情况
// ⚠️易错点3：step++的时机 - jump中必须在i==end时才递增step，不能在farthest更新时递增；end更新为farthest后，要检查是否已经能到达终点，避免多跳一步
// ⚠️易错点4：贪心不适用的场景 - 贪心要求问题具有"无后效性"，即当前选择不会影响之前已做选择的最优性；若每步选择会影响全局（如带权最短路），则必须用DP而非贪心

// 何时运用贪心算法：
// ❓1、能否证明局部最优可以推出全局最优？（数学归纳法或反证法）
// ❓2、问题是否具有无后效性？即当前步骤的选择不会影响之前步骤的结果
// ❓3、是否存在"显然最优"的局部选择？如每次选可达最远、每次选最小代价等

// 0、func canJumpDp(nums []int) bool     // 跳跃游戏 - 动态规划解法（展示用）
// 1、func canJump(nums []int) bool       // 跳跃游戏 - 贪心算法
// 2、func jumpDp(nums []int) int         // 跳跃游戏II - 动态规划解法（展示用）
// 3、func jump(nums []int) int           // 跳跃游戏II - 贪心算法

// 跳跃游戏 - 能否到达最后位置（无法通过时间限制，仅做展示用）
func canJumpDp(nums []int) bool {
	var dp func(nums []int, start int) bool
	// 动态规划解法：
	// 明确状态：当前位置 start
	// 明确选择列表：从start位置可以跳跃的步数（1到nums[start]步）
	// 明确dp函数定义：dp(start) 表示从位置start能否到达最后一个位置
	// 明确base case：start >= len(nums)-1（到达终点）或 nums[start] == 0（无法移动）

	dp = func(nums []int, start int) bool {
		if start >= len(nums)-1 {
			return true
		}
		if nums[start] == 0 {
			return false
		}

		for i := 1; i <= nums[start]; i++ {
			res := dp(nums, start+i)
			if res {
				return true
			}
		}
		return false
	}
	return dp(nums, 0)
}

// 跳跃游戏 - 贪心算法思路
func canJump(nums []int) bool {
	// 注意i边界是：i < len(nums)-1
	var farthest = 0

	// 注意下标从0开始
	for i := 0; i < len(nums)-1; i++ {
		farthest = max(farthest, i+nums[i])

		// 碰到了0，卡住跳不动了
		if farthest <= i {
			return false
		}
	}
	return farthest >= len(nums)-1
}

// 跳跃游戏II - 最少跳跃次数（无法通过时间限制，仅做展示用）
func jumpDp(nums []int) int {
	//定义：从索引p跳到最后一格，至少需要dp(nums, p)步
	var dp func(nums []int, p int, memo []int) int

	// 备忘录
	memo := make([]int, len(nums))
	for i := range memo {
		memo[i] = len(nums)
	}

	dp = func(nums []int, p int, memo []int) int {
		if p >= len(nums)-1 {
			return 0
		}

		// 子问题已经计算过
		if memo[p] != len(nums) {
			return memo[p]
		}

		step := nums[p]
		for i := 1; i <= step; i++ {
			// 穷举每一个选择
			// 计算每一个子问题的结果
			subProb := dp(nums, p+i, memo)
			// 取其中最小的作为最终结果
			memo[p] = min(memo[p], subProb+1)
		}
		return memo[p]
	}

	return dp(nums, 0, memo)
}

// 跳跃游戏II
func jump(nums []int) int {
	if len(nums) <= 1 {
		return 0
	}

	// jumps的含义：跳到索引区间[i, end]需step步，i是变化量
	end, step := 0, 0

	// farthest的含义：从索引i可以跳到的最远索引
	farthest := 0

	// 注意下标从0开始
	for i := 0; i < len(nums)-1; i++ {
		if nums[i]+i > farthest {
			farthest = nums[i] + i
		}

		// [i, end]区间是step步可达的索引范围
		// 现在已经遍历完[i, end]，所以需要再跳一步
		if i == end {
			end = farthest
			step++
			if farthest >= len(nums)-1 {
				return step
			}
		}
	}
	return -1
}

func main() {
	//fmt.Println(canJump([]int{2, 3, 1, 1, 4}))
	//fmt.Println(canJump([]int{3, 2, 1, 0}))
	//fmt.Println(canJump([]int{3, 2, 1, 0, 4}))

	fmt.Println(jump([]int{2, 3, 1, 1, 4}))
	fmt.Println(jump([]int{2, 3, 0, 1, 4}))
}
