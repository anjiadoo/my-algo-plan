package main

import "fmt"

// 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	slow, fast := 0, 0
	for fast < len(nums) {
		if nums[slow] != nums[fast] {
			slow++
			// 维护nums[0..slow]无重复
			nums[slow] = nums[fast]
		}
		fast++
	}
	return slow + 1
}

func main() {
	array := []int{1, 1, 3, 3, 4, 5, 5, 5, 6}
	k := removeDuplicates(array)
	fmt.Println(k, array)
}
