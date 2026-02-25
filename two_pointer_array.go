package main

import "fmt"

// 双指针数组算法实现：
// 🌟技巧1：快慢指针原地修改技巧 - 原地修改数组时，慢指针指向有效区域末尾，快指针遍历查找新元素，遇到新元素时slow++并赋值（删除有序数组重复项）
// 🌟技巧2：左右指针相向而行技巧 - 二分查找、两数之和、数组反转等场景，left从最左开始，right从最右开始，逐步向中间靠拢（二分查找、两数之和、反转数组）
// 🌟技巧3：中心向外扩散技巧 - 回文问题常用技巧，从中心点（或中心两点）向两端扩散，比较字符是否相等（最长回文子串）

// 0、func removeDuplicates(nums []int) int
// 1、func binarySearch(nums []int, target int) int
// 2、func twosum(numbers []int, target int) []int
// 3、func reversestring(s []byte)
// 4、func longestPalindrome(s string) string

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

// 二分查找，一左一右两个指针相向而行
func binarySearch(nums []int, target int) int {
	left := 0
	right := len(nums) - 1

	for left <= right {
		mid := (left + right) / 2
		if nums[mid] == target {
			return mid
		}
		if nums[mid] > target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return -1
}

// 两数之和 II - 输入有序数组
func twoSum(numbers []int, target int) []int {
	// 一左一右两个指针相向而行
	left := 0
	right := len(numbers) - 1

	for left < right {
		sum := numbers[left] + numbers[right]
		if sum == target {
			// 题目要求的索引是从 1 开始的
			return []int{left + 1, right + 1}
		}
		if sum > target {
			right--
		} else {
			left++
		}
	}
	return []int{-1, -1}
}

// 反转数组
func reverseString(s []byte) {
	left := 0
	right := len(s) - 1

	for left < right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
}

// 最长回文子串
func longestPalindrome(s string) string {
	// 从中心向两端扩散的双指针技巧
	maxSubStr := ""

	palindrome := func(s string, l, r int) string {
		// 防止索引越界
		for l >= 0 && r < len(s) && s[l] == s[r] {
			l--
			r++
		}
		return s[l+1 : r]
	}

	// 一般涉及到“最长”的，都需要穷举所有可能
	for i := 0; i < len(s); i++ {
		s1 := palindrome(s, i, i)
		if len(s1) > len(maxSubStr) {
			maxSubStr = s1
		}

		s2 := palindrome(s, i, i+1)
		if len(s2) > len(maxSubStr) {
			maxSubStr = s2
		}
	}
	return maxSubStr
}

func main() {
	//array := []int{1, 2, 3, 4, 5, 6}
	//fmt.Println(twoSum(array, 12))

	//ch := []byte{'a', 'b', 'c', 'd', 'e'}
	//reverseString(ch)
	//fmt.Println(string(ch))

	str := "aba"
	fmt.Println(longestPalindrome(str))
}
