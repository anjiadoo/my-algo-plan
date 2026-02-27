package main

import "fmt"

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

func main() {
	nums := []int{1, 2, 3}
	fmt.Println(permute(nums))
}
