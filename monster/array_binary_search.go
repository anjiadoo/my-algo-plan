package main

import (
	"fmt"
	"math"
)

// 二分搜索-基础框架
func binarySearch(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		}
	}
	return -1
}

// 二分搜索-寻找左侧边界
func leftBound(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			right = mid - 1
		} else if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		}
	}

	// 如果越界target肯定不存在，返回-1
	if left >= len(nums) || left < 0 {
		return -1
	}

	// 结束时，left=right+1
	if nums[left] != target {
		return -1
	}

	// 去掉上面判断，target不存在时返回的是「大于target的最小索引」
	return left
}

// 二分搜索-寻找右侧边界
func rightBound(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		}
	}

	// 如果越界target肯定不存在，返回-1
	if right < 0 || right >= len(nums) {
		return -1
	}

	// 结束时，left=right+1 => right=left-1
	if nums[right] != target {
		return -1
	}

	// 去掉上面判断，target不存在时返回的是「小于target的最大索引」
	return right
}

// 在排序数组中查找元素的第一个和最后一个位置
// https://leetcode.cn/problems/find-first-and-last-position-of-element-in-sorted-array/
func searchRange(nums []int, target int) []int {
	leftIdx := leftBound(nums, target)
	if leftIdx == -1 {
		return []int{-1, -1}
	}
	rightIdx := rightBound(nums, target)
	return []int{leftIdx, rightIdx}
}

/* 🌟什么时候可以运用二分搜索算法技巧❓
 *
 *   首先要能抽象出一个自变量 x，一个关于 x 的函数 f(x)，以及一个目标值 target。
 *   同时 x, f(x), target 还要满足以下条件：
 *   1、f(x) 必须是在 x 上的单调函数（单调增单调减都可以）。
 *   2、题目是计算满足约束条件 f(x) == target 时的 x 的最值。
 *
 * 🌟具体操作步骤如下👇：
 *
 *   1、确定 x, f(x), target 分别是什么，并写出函数 f 的代码，即：
 *    - x 是什么？（我能控制的变量，注意：通常x本身就是一个'最值'）
 *    - f(x) 是什么？（它是关于x单调函数，x变化时(穷举所有可能)，结果如何变化）
 *    - target 是什么？求左边界还是右边界？（优化方向，也就是穷举的过程）
 *   2、找到 x 的取值范围作为二分搜索的搜索区间，初始化 left 和 right 变量。
 *   3、根据题目的要求，确定应该使用「搜索左侧」还是「搜索右侧」的二分搜索算法，写出解法代码。
 *
 * 🌟如何训练这种抽象能力？
 *   🔑 第一步：忘掉故事，提取"操作"
 *   不要被"船"、"包裹"、"天数"这些词干扰。问自己：
 *    - 我在操作什么？ → 一个有序序列
 *    - 我能做什么操作？ → 按顺序切分成连续段
 *    - 约束是什么？ → 每段有个上限 / 段数有个上限
 *    - 优化什么？ → 在满足一个约束的前提下，最小化另一个
 *
 *   🔑 第二步：识别"对偶关系"
 *   这类题目的核心洞察是——两个变量之间存在单调的对偶关系：
 *    - 每段上限 x ↑  ⟹  需要的段数 f(x) ↓
 *    - 每段上限 x ↓  ⟹  需要的段数 f(x) ↑
 *   这就是你代码注释里写的 "f(x) 必须是在 x 上的单调函数"。
 *   一旦识别出这种单调对偶关系，就可以用二分搜索。
 *
 *   🔑 第三步：统一建模模板
 *    "连续分段"类二分问题的统一模型：
 *     - 输入：数组 arr，限制值 target
 *     - x：每段的容量上限
 *     - f(x)：按容量 x 分段，最少需要几段
 *     - 单调性：x 越大 → f(x) 越小
 *     - 求：满足 f(x) ≤ target 的最小 x（左边界）
 *    x 的范围：[max(arr), sum(arr)]
 */

// 二分搜索题目的思维框架：在 f(x)==target 的约束下求 x 的最值
func binarySearchSolution(nums []int, target int) int {
	// 函数 f 是关于自变量 x 的单调函数
	f := func(x int) int {
		// 这是内层for循环，x 已经固定了
		for i := 0; i < len(nums); {
			// ...
		}
		return -1
	}

	left := -1  // 自变量 x 的最小值是多少❓
	right := -1 // 自变量 x 的最大值是多少？❓

	// 外侧for循环，目的是穷举 x 所有可能
	for left <= right {

		mid := left + (right-left)/2
		fv := f(mid)

		if fv == target {
			// 要求解的最值在左边界❓
			left = mid + 1
			// 还是在右边界❓
			right = mid - 1

		} else if fv > target {
			// 怎么让 f(x) 小一点❓

		} else if fv < target {
			// 怎么让 f(x) 大一点❓

		}
	}
	return left
}

// 爱吃香蕉的珂珂 https://leetcode.cn/problems/koko-eating-bananas/
func minEatingSpeed(piles []int, h int) int {
	// x: 吃香蕉的速度 x
	// f(x): 吃完全部香蕉需要的时间 f(x)
	// target: 给定的时间限制 h
	f := func(x int) (hours int) {
		for _, pile := range piles {
			hours += pile / x
			if pile%x > 0 {
				hours++
			}
		}
		return hours
	}

	// 吃香蕉速度最小:1根/h
	left := 1
	// 吃香蕉速度最大:max(piles)根/h
	right := int(math.Pow10(9) + 1)

	// 穷举自变量 x 的所有可能
	for left <= right {
		mid := left + (right-left)/2
		fv := f(mid)
		if fv == h {
			right = mid - 1
		} else if fv > h {
			left = mid + 1
		} else if fv < h {
			right = mid - 1
		}
	}
	return left
}

// 在D天内送达包裹的能力 https://leetcode.cn/problems/capacity-to-ship-packages-within-d-days/
// 本质：给定一个数组，将其按顺序分成若干连续段，每段的"总量"不超过 x。求使得段数 ≤ target 时，x 的最小值
func shipWithinDays(weights []int, days int) int {
	// x: 船的载能力 x
	// f(x): 按运载能力x运送，需要f(x)天运完
	// target: 给定的时间限制 days

	f := func(x int) (days int) {
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

	// 船的最小载能力 max(weights)
	left := -1
	// 船的最大载能力 sum(weights)
	right := -1

	for i := range weights {
		left = max(left, weights[i])
		right += weights[i]
	}

	// 穷举自变量 x 的所有可能
	for left <= right {
		mid := left + (right-left)/2
		fv := f(mid)
		if fv == days {
			right = mid - 1
		} else if fv > days {
			left = mid + 1
		} else if fv < days {
			right = mid - 1
		}
	}
	return left
}

// 分割数组的最大值 https://leetcode.cn/problems/split-array-largest-sum/
// 本质：给定一个数组，将其按顺序分成若干连续段，每段的"总量"不超过 x。求使得段数 ≤ target 时，x 的最小值
func splitArray(nums []int, k int) int {
	// x: 子数组和的上限 x
	// f(x): 按子数组和sum为x来拆分，需要f(x)个连续子数组
	// target: 给定的子数组个数限制 k

	f := func(x int) (num int) {
		sum := 0
		for i := 0; i < len(nums); i++ {
			sum += nums[i]
			if sum > x {
				num++
				i-- //已经放进来的nums[i]要重新放回去
				sum = 0
			}
		}
		num++
		return num
	}

	// 子数组和最小 max(weights)
	left := -1
	// 子数组和最大 sum(weights)
	right := -1

	for i := range nums {
		left = max(left, nums[i])
		right += nums[i]
	}

	// 穷举自变量 x 的所有可能
	for left <= right {
		mid := left + (right-left)/2
		fv := f(mid)
		if fv == k {
			right = mid - 1
		} else if fv > k {
			left = mid + 1
		} else if fv < k {
			right = mid - 1
		}
	}
	return left
}

// 搜索二维矩阵 https://leetcode.cn/problems/search-a-2d-matrix/description/
// 每行中的整数从左到右按非严格递增顺序排列。
// 每行的第一个整数大于前一行的最后一个整数。
func searchMatrix(matrix [][]int, target int) bool {
	m, n := len(matrix), len(matrix[0])
	left, right := 0, m*n-1

	for left <= right {
		mid := left + (right-left)/2
		find := matrix[mid/n][mid%n]

		if find == target {
			return true
		} else if find > target {
			right = mid - 1
		} else if find < target {
			left = mid + 1
		}
	}
	return false
}

// 搜索二维矩阵II https://leetcode.cn/problems/search-a-2d-matrix-ii/
// 每行的元素从左到右升序排列。
// 每列的元素从上到下升序排列。
func searchMatrixII(matrix [][]int, target int) bool {
	// 从右上角开始，规定只能向左或向下移动
	m, n := len(matrix), len(matrix[0])
	i, j := 0, n-1

	for i < m && j >= 0 {
		//fmt.Printf("i=%d j=%d next=%d\n", i, j, matrix[i][j])

		if matrix[i][j] == target {
			return true
		} else if matrix[i][j] > target {
			j--
		} else if matrix[i][j] < target {
			i++
		}
	}
	return false
}

// 匹配子序列的单词数 https://leetcode.cn/problems/number-of-matching-subsequences/
func numMatchingSubseq(s string, words []string) int {
	mapIndexs := make(map[byte][]int)
	for i := range s {
		if mapIndexs[s[i]] == nil {
			mapIndexs[s[i]] = []int{}
		}
		mapIndexs[s[i]] = append(mapIndexs[s[i]], i)
	}

	leftBound := func(nums []int, target int) int {
		left, right := 0, len(nums)-1
		for left <= right {
			mid := left + (right-left)/2
			if nums[mid] == target {
				right = mid - 1
			} else if nums[mid] > target {
				right = mid - 1
			} else if nums[mid] < target {
				left = mid + 1
			}
		}
		return left
	}

	res := 0
	for _, word := range words {
		i := 0 // word的字符索引
		j := 0 // s的字符索引
		for i < len(word) {
			ch := word[i]

			indexList := mapIndexs[ch]
			if indexList == nil {
				break
			}

			// 二分搜索大于等于j的最小索引（s已经走到j-1的位置了）
			// 即在 s[j..] 中搜索等于 word[i] 的最小索引
			pos := leftBound(indexList, j)
			if pos == len(indexList) {
				break
			}

			// 找到 word[i] == s[j]，继续往后匹配
			j = indexList[pos]

			i++
			j++
		}
		if i == len(word) {
			res++
		}
	}
	return res
}

func main() {

	fmt.Println(numMatchingSubseq("abcde", []string{"a", "bb", "acd", "ace"}))
	fmt.Println(numMatchingSubseq("dsahjpjauf", []string{"ahjpjau", "ja", "ahbwzgqnuk", "tnmlanowax"}))

	//fmt.Println(searchMatrixII([][]int{
	//	{1, 4, 7, 11, 15},
	//	{2, 5, 8, 12, 19},
	//	{3, 6, 9, 16, 22},
	//	{10, 13, 14, 17, 24},
	//	{18, 21, 23, 26, 30},
	//}, 13))

	//fmt.Println(searchMatrix([][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 60))
	//fmt.Println(searchMatrix([][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 13))

	//fmt.Println(splitArray([]int{7, 2, 5, 10, 8}, 2))
	//fmt.Println(splitArray([]int{1, 2, 3, 4, 5}, 2))
	//fmt.Println(splitArray([]int{1, 4, 4}, 3))

	//fmt.Println(shipWithinDays([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5))
	//fmt.Println(shipWithinDays([]int{3, 2, 2, 4, 1, 4}, 3))
	//fmt.Println(shipWithinDays([]int{1, 2, 3, 1, 1}, 4))

	//fmt.Println(minEatingSpeed([]int{3, 6, 7, 11}, 8))
	//fmt.Println(minEatingSpeed([]int{30, 11, 23, 4, 20}, 5))
	//fmt.Println(minEatingSpeed([]int{30, 11, 23, 4, 20}, 6))

	//fmt.Println(searchRange([]int{5, 7, 7, 8, 8, 10}, 8))
	//fmt.Println(searchRange([]int{5, 7, 7, 8, 8, 10}, 6))

	//fmt.Println(rightBound([]int{-1, 1, 2, 3, 3, 3, 3, 3, 5, 6}, 3))
	//fmt.Println(leftBound([]int{-1, 1, 2, 3, 3, 3, 3, 3, 5, 6}, 3))
	//fmt.Println(binarySearch([]int{-1, 1, 2, 3, 3, 3, 3, 3, 5, 6}, 3))
}
