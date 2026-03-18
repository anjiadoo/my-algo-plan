package main

import (
	"fmt"
	"sort"
)

// 两数之和，返回和为target的两个数，注意不能返回重复数对儿
func twoSumTarget(nums []int, target int) [][]int {
	// 解题核心是：先排序，再++和--时使用for循环跳过相同元素
	sort.Ints(nums)
	var lo, hi = 0, len(nums) - 1
	var res [][]int

	for lo < hi {
		left, right := nums[lo], nums[hi]
		sum := nums[lo] + nums[hi]

		if sum < target {
			for lo < hi && nums[lo] == left {
				lo++
			}
		} else if sum > target {
			for lo < hi && nums[hi] == right {
				hi--
			}
		} else {
			res = append(res, []int{left, right})
			for lo < hi && nums[lo] == left {
				lo++
			}
			for lo < hi && nums[hi] == right {
				hi--
			}
		}
	}
	return res
}

// 3数之和，返回三元组nums[i]+nums[j]+nums[k]=target且i!=j!=k的所有元素对儿
func threeSumTarget(nums []int, target int) [][]int {
	sort.Ints(nums)
	var result [][]int

	twoSumTarget := func(nums []int, start, target int) [][]int {
		sort.Ints(nums)
		var lo, hi = start, len(nums) - 1
		var res [][]int
		for lo < hi {
			left, right := nums[lo], nums[hi]
			sum := nums[lo] + nums[hi]

			if sum < target {
				for lo < hi && nums[lo] == left {
					lo++
				}
			} else if sum > target {
				for lo < hi && nums[hi] == right {
					hi--
				}
			} else {
				res = append(res, []int{left, right})
				for lo < hi && nums[lo] == left {
					lo++
				}
				for lo < hi && nums[hi] == right {
					hi--
				}
			}
		}
		return res
	}

	// 穷举 threeSum 的第一个数
	for i := 0; i < len(nums); i++ {
		tuples := twoSumTarget(nums, i+1, target-nums[i])
		for _, tuple := range tuples {
			tuple = append(tuple, nums[i])
			result = append(result, tuple)
		}
		// 跳过第一个数字重复的情况，否则会出现重复结果
		for i < len(nums)-1 && nums[i] == nums[i+1] {
			i++
		}
	}

	return result
}

// n数之和，返回n元组nums[i]+nums[j]...nums[n]=target且i!=j...n的所有元素对儿
func nSumTarget(nums []int, target int) [][]int {
	sort.Ints(nums)
	return _nSumTarget(nums, 4, 0, target)
}

func _nSumTarget(nums []int, n int, start int, target int) [][]int {
	// nums 必须是升序数组
	var result [][]int
	if n == 2 {
		lo, hi := start, len(nums)-1
		for lo < hi {
			left, right := nums[lo], nums[hi]
			sum := nums[lo] + nums[hi]
			if sum < target {
				for lo < hi && nums[lo] == left {
					lo++
				}
			} else if sum > target {
				for lo < hi && nums[hi] == right {
					hi--
				}
			} else {
				result = append(result, []int{left, right})
				for lo < hi && nums[lo] == left {
					lo++
				}
				for lo < hi && nums[hi] == right {
					hi--
				}
			}
		}
	} else {
		for i := start; i < len(nums); i++ {
			res := _nSumTarget(nums, n-1, i+1, target-nums[i])
			for _, pair := range res {
				pair = append(pair, nums[i])
				result = append(result, pair)
			}
			// ⚠️ 注意跳过相同元素
			for i < len(nums)-1 && nums[i] == nums[i+1] {
				i++
			}
		}
	}
	return result
}

func main() {
	fmt.Println(nSumTarget([]int{1, 0, -1, 0, -2, 2}, 0))
	fmt.Println(nSumTarget([]int{2, 2, 2, 2}, 8))

	//fmt.Println(threeSumTarget([]int{-1, 0, 1, 2, -1, -4}, 0))

	//fmt.Println(twoSumTarget([]int{1, 1, 1, 2, 2, 3, 3}, 4))
}
