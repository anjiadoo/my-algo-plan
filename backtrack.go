package main

import (
	"fmt"
	"sort"
)

// 回溯算法实现：
// 🌟技巧1：回溯框架技巧 - 回溯本质是N叉树的遍历，在前序位置做选择（加入路径），在后序位置撤销选择（从路径移除），递归终止时收集结果
// 🌟技巧2：路径/选择列表/结束条件三要素 - 解题时先明确：①路径（已做的选择track）②选择列表（当前可做的选择）③结束条件（何时收集结果），三者明确后套框架即可
// 🌟技巧3：used数组/visited剪枝技巧 - 用used数组（或map）标记已选择的元素，避免重复选择同一元素；在排列问题中标记索引，在组合问题中用start参数控制起点
// 🌟技巧4：排列vs组合vs子集 - 排列问题用used数组从头遍历（元素顺序有关），组合/子集问题用start参数从当前位置往后遍历（元素顺序无关），这是三类问题的核心区别
//
// ⚠️易错点1：结果必须深拷贝 - 收集结果时必须copy(temp, track)深拷贝路径，不能直接append(res, track)，因为track是引用类型，后续修改会影响已收集的结果
// ⚠️易错点2：撤销选择必须对称 - 做选择和撤销选择的操作必须完全对称（如append对应切片截断，used[i]=true对应used[i]=false），漏掉撤销操作会导致结果错误
// ⚠️易错点3：去重排序前提 - 含重复元素的排列/组合问题，去重的前提是先对数组排序，然后用nums[i]==nums[i-1]跳过同层重复选择，不排序直接去重会遗漏或出错
// ⚠️易错点4：指针传递track和res - Go中切片是引用类型但append可能改变底层数组，建议用指针传递track和res（如*[]int），确保递归过程中修改对所有层可见

// 何时运用回溯算法：
// ❓1、问题是否要求「穷举所有方案」？如全排列、所有组合、所有子集
// ❓2、问题是否可以抽象为「在决策树上做选择」？每一步有多个选项，需要遍历所有分支
// ❓3、问题是否满足「选择-探索-撤销」的模式？如数独、N皇后、括号生成

// 1、func permute(nums []int) [][]int                         // 全排列

// 全排列
// 给定一个不含重复数字的数组nums，返回其「所有可能的全排列」
func permute(nums []int) [][]int {
	// 1、路径：已经做出的选择
	// 2、选择列表：还可以选择的列表
	// 3、结束条件：到达决策底层，无法再做选择的条件

	track := make([]int, 0, len(nums))
	res := make([][]int, 0)
	used := make(map[int]bool)

	backtrack(nums, used, &track, &res)
	return res
}

func backtrack(nums []int, used map[int]bool, track *[]int, res *[][]int) {
	if len(*track) == len(nums) {
		temp := make([]int, len(*track))
		copy(temp, *track)
		*res = append(*res, temp)
		return
	}

	for i := 0; i < len(nums); i++ {
		if used[nums[i]] {
			continue
		}
		*track = append(*track, nums[i])
		used[nums[i]] = true

		backtrack(nums, used, track, res)

		used[nums[i]] = false
		*track = (*track)[:len(*track)-1]
	}
}

// 组合
// 给定两个整数n和k，返回范围[1, n]中所有可能的k个数的组合
func combine(n, k int) [][]int {
	var track []int
	var res [][]int

	backtrack1(n, k, 1, &track, &res)
	return res
}

func backtrack1(n, k, start int, track *[]int, res *[][]int) {
	if len(*track) == k {
		temp := make([]int, len(*track))
		copy(temp, *track)
		*res = append(*res, temp)
		return
	}

	for i := start; i <= n; i++ {
		*track = append(*track, i)
		backtrack1(n, k, i+1, track, res)
		*track = (*track)[:len(*track)-1]
	}
}

// 子集II
// 给你一个整数数组nums ，其中可能包含重复元素，请你返回该数组所有可能的子集（幂集），解集「不能」包含重复的子集。
func subsetsWithDup(nums []int) [][]int {
	var track []int
	var res [][]int

	sort.Ints(nums)

	backtrack2(nums, 0, &track, &res)
	return res
}

func backtrack2(nums []int, start int, track *[]int, res *[][]int) {
	temp := make([]int, len(*track))
	copy(temp, *track)
	*res = append(*res, temp)

	for i := start; i < len(nums); i++ {
		// 剪枝逻辑，值相同的相邻树枝，只遍历第一条
		if i > start && nums[i] == nums[i-1] {
			continue
		}

		*track = append(*track, nums[i])
		backtrack2(nums, i+1, track, res)
		*track = (*track)[:len(*track)-1]
	}
}

// 组合总和II
// 给定一个候选人编号的集合candidates和一个目标数target，找出candidates中所有可以使数字和为target的组合
// candidates中的每个数字在每个组合中只能使用「一次」
func combinationSum2(candidates []int, target int) [][]int {
	var track []int
	var res [][]int
	trackSum := 0

	sort.Ints(candidates)

	backtrack3(candidates, 0, target, trackSum, &track, &res)
	return res
}

func backtrack3(nums []int, start, target, trackSum int, track *[]int, res *[][]int) {
	if trackSum == target {
		temp := make([]int, len(*track))
		copy(temp, *track)
		*res = append(*res, temp)
		return
	}
	if trackSum > target { // 提前剪枝
		return
	}

	for i := start; i < len(nums); i++ {
		if i > start && nums[i] == nums[i-1] {
			continue
		}

		*track = append(*track, nums[i])
		trackSum += nums[i]

		backtrack3(nums, i+1, target, trackSum, track, res)

		trackSum -= nums[i]
		*track = (*track)[:len(*track)-1]
	}
}

// 全排列II
// 给定一个可包含重复数字的序列nums，按任意顺序返回所有不重复的全排列。
func permuteUnique(nums []int) [][]int {
	return nil
}

func main() {
	fmt.Println("全排列:", permute([]int{1, 2, 3}))
	fmt.Println("组合:", combine(4, 2))

	fmt.Println("字集2:", subsetsWithDup([]int{1, 1, 2, 2}))
	fmt.Println("组合总和II:", combinationSum2([]int{10, 1, 2, 7, 6, 1, 5}, 8))
}
