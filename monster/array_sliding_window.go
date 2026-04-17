package main

import (
	"fmt"
	"math"
)

// 最小覆盖子串 https://leetcode.cn/problems/minimum-window-substring/description/
func minWindow(s string, t string) string {
	window, need := make(map[byte]int), make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}

	left, right := 0, 0
	start, length := 0, math.MaxInt
	valid := 0

	for right < len(s) {
		ch := s[right]
		right++

		if _, ok := need[ch]; ok {
			window[ch]++
			if window[ch] == need[ch] {
				valid++
			}
		}

		for valid == len(need) {
			if right-left < length {
				start = left
				length = right - left
			}

			d := s[left]
			left++

			if _, ok := need[d]; ok {
				if window[d] == need[d] {
					valid--
				}
				window[d]--
			}
		}
	}
	if length == math.MaxInt {
		return ""
	}
	return s[start : start+length]
}

// 字符串的排列 https://leetcode.cn/problems/permutation-in-string/description/
func checkInclusion(s1 string, s2 string) bool {
	window, need := make(map[byte]int), make(map[byte]int)
	for i := range s1 {
		need[s1[i]]++
	}

	left, right := 0, 0
	valid := 0

	for right < len(s2) {
		ch := s2[right]
		right++

		if _, ok := need[ch]; ok {
			window[ch]++
			if window[ch] == need[ch] {
				valid++
			}
		}

		if right-left == len(s1) {
			if valid == len(need) {
				return true
			}

			d := s2[left]
			left++

			if _, ok := need[d]; ok {
				if window[d] == need[d] {
					valid--
				}
				window[d]--
			}
		}
	}
	return false
}

// 找到字符串中所有字母异位词 https://leetcode.cn/problems/find-all-anagrams-in-a-string/description/
func findAnagrams(s string, p string) []int {
	window, need := make(map[byte]int), make(map[byte]int)
	for i := range p {
		need[p[i]]++
	}

	left, right := 0, 0
	valid := 0
	var res []int

	for right < len(s) {
		ch := s[right]
		right++
		if _, ok := need[ch]; ok {
			window[ch]++
			if window[ch] == need[ch] {
				valid++
			}
		}

		if right-left == len(p) {
			if valid == len(need) {
				res = append(res, left)
			}

			d := s[left]
			left++
			if _, ok := need[d]; ok {
				if window[d] == need[d] {
					valid--
				}
				window[d]--
			}
		}
	}
	return res
}

// 无重复字符的最长子串 https://leetcode.cn/problems/longest-substring-without-repeating-characters/description/
func lengthOfLongestSubstring(s string) int {
	window := make(map[byte]int)
	left, right := 0, 0
	res := 0

	for right < len(s) {
		ch := s[right]
		right++
		window[ch]++

		for window[ch] > 1 {
			d := s[left]
			left++
			window[d]--
		}
		res = max(res, right-left)
	}
	return res
}

// 将x减到0的最小操作数 https://leetcode.cn/problems/minimum-operations-to-reduce-x-to-zero/description/
func minOperations(nums []int, x int) int {
	// 等价于寻找nums中和为sum(nums)-x的最长子数组
	sum := 0
	for _, num := range nums {
		sum += num
	}

	left, right, windowSum := 0, 0, 0
	target := sum - x
	res := -1

	for right < len(nums) {
		windowSum += nums[right]
		right++
		for left < right && windowSum > target {
			windowSum -= nums[left]
			left++
		}
		if windowSum == target {
			res = max(res, right-left)
		}
	}
	if res == -1 {
		return -1
	}
	return len(nums) - res
}

// 乘积小于K的子数组 https://leetcode.cn/problems/subarray-product-less-than-k/
func numSubarrayProductLessThanK(nums []int, k int) int {
	left, right := 0, 0
	windowProduct := 1
	count := 0

	for right < len(nums) {
		windowProduct *= nums[right]
		right++

		for left < right && windowProduct >= k {
			windowProduct /= nums[left]
			left++
		}
		// 窗口中的子数组个数怎么计算?
		// 比方说left=1,right=4划定了[1, 2, 3]这个窗口(right是开区间)
		// 但不止[left..right)是合法的子数组，[left+1..right),[left+2..right)等都是合法子数组
		// 需要把[3], [2,3], [1,2,3]这right-left个子数组都加上，至于[1][1,2][2]等之前已经算过了
		count += right - left
	}
	return count
}

// 最大连续1的个数III https://leetcode.cn/problems/max-consecutive-ones-iii/description/
func longestOnes(nums []int, k int) int {
	left, right := 0, 0
	windowOneCount := 0
	res := math.MinInt

	for right < len(nums) {
		if nums[right] == 1 {
			windowOneCount++
		}
		right++

		// 窗口中替换的0的数量大于k时缩小窗口
		for right-left-windowOneCount > k {
			if nums[left] == 1 {
				windowOneCount--
			}
			left++
		}

		// 求最大一般都在这个位置
		res = max(res, right-left)
	}
	return res
}

// 替换后的最长重复字符 https://leetcode.cn/problems/longest-repeating-character-replacement/description/
func characterReplacement(s string, k int) int {
	windowCharCount := make(map[byte]int)
	left, right := 0, 0
	res := -1

	// 记录窗口中的字符的最大重复次数
	// 因为最划算的替换方法是：把其他字符替换成出现次数最多的那个字符
	maxCharCount := 0

	for right < len(s) {
		ch := s[right]
		right++
		windowCharCount[ch]++

		maxCharCount = max(maxCharCount, windowCharCount[ch])

		if right-left-maxCharCount > k {
			d := s[left]
			left++
			windowCharCount[d]--
		}

		// 求最大一般都在这个位置
		res = max(res, right-left)
	}

	return res
}

// 存在重复元素II https://leetcode.cn/problems/contains-duplicate-ii/
func containsNearbyDuplicate(nums []int, k int) bool {
	window := make(map[int]bool)
	left, right := 0, 0

	for right < len(nums) {
		num := nums[right]
		if window[num] {
			return true
		}

		window[num] = true
		right++

		if right-left > k {
			d := nums[left]
			left++
			delete(window, d)
		}
	}
	return false
}

// 存在重复元素III https://leetcode.cn/problems/contains-duplicate-iii/description/
func containsNearbyAlmostDuplicate(nums []int, indexDiff int, valueDiff int) bool {
	bucketSize := valueDiff + 1
	window := make(map[int]int)
	left, right := 0, 0

	getBucketIdx := func(x int) int {
		if x >= 0 {
			return x / bucketSize
		}
		return (x+1)/bucketSize - 1
	}

	// 解题关键：如何在窗口 [left,right] 中快速判断是否有元素之差小于 t 的两个元素呢❓
	// 这需要利用：桶排序的「地板元素」&「天花板元素」的特性，具体如下：
	// 把数轴切成宽度为 valueDiff+1 的桶，将"是否存在差值 ≤ valueDiff 的元素"这个问题，
	// 转化为"只需检查当前桶和左右相邻桶"这个 O(1) 操作。配合滑动窗口维护下标约束，整体时间复杂度 O(n)。

	for right < len(nums) {
		bucketIdx := getBucketIdx(nums[right])

		// 情况1: 当前桶存在元素之差小于等于indexDiff的元素
		if _, ok := window[bucketIdx]; ok {
			return true
		}
		// 情况2: 检查左边桶，可能存在，须判断
		if num, ok := window[bucketIdx-1]; ok && math.Abs(float64(nums[right]-num)) <= float64(valueDiff) {
			return true
		}
		// 情况3: 检查右边桶，可能存在，须判断
		if num, ok := window[bucketIdx+1]; ok && math.Abs(float64(nums[right]-num)) <= float64(valueDiff) {
			return true
		}

		window[bucketIdx] = nums[right]
		right++

		// 缩小窗口
		if right-left > indexDiff {
			delete(window, getBucketIdx(nums[left]))
			left++
		}
	}
	return false
}

// 长度最小的子数组 https://leetcode.cn/problems/minimum-size-subarray-sum/description/
func minSubArrayLen(target int, nums []int) int {
	left, right := 0, 0
	windowSum := 0
	res := math.MaxInt

	for right < len(nums) {
		windowSum += nums[right]
		right++

		for windowSum >= target && left < right {
			res = min(res, right-left)
			windowSum -= nums[left]
			left++
		}
	}
	if res == math.MaxInt {
		return 0
	}
	return res
}

// 至少有K个重复字符的最长子串 https://leetcode.cn/problems/longest-substring-with-at-least-k-repeating-characters/description/
func longestSubstring(s string, k int) int {
	// 原题没有缩小缩小窗口的时机，那么就自己创造缩窗的时机，题目改写成：
	// 在s中寻找仅含有count种字符，且每种字符出现次数都大于k的最长子串，count取值范围[1~26]，因为题目说了只含小写字符
	var res int
	for i := 1; i <= 26; i++ {
		res = max(res, _longestSubstring(s, k, i))
	}
	return res
}

func _longestSubstring(s string, k, count int) int {
	window := make(map[byte]int)
	left, right := 0, 0
	valid := 0
	res := math.MinInt

	for right < len(s) {
		ch := s[right]
		right++

		window[ch]++
		if window[ch] == k {
			valid++
		}

		for len(window) > count {
			d := s[left]
			left++

			if window[d] == k {
				valid--
			}
			window[d]--
			if window[d] == 0 {
				delete(window, d)
			}
		}
		if valid == len(window) {
			res = max(res, right-left)
		}
	}
	if res == math.MinInt {
		return 0
	}
	return res
}

func main() {

	fmt.Println(longestSubstring("aaabb", 3))
	fmt.Println(longestSubstring("ababbc", 2))

	//fmt.Println(minSubArrayLen(11, []int{1, 2, 3, 4, 5}))
	//fmt.Println(minSubArrayLen(7, []int{2, 3, 1, 2, 4, 3}))
	//fmt.Println(minSubArrayLen(4, []int{1, 4, 4}))
	//fmt.Println(minSubArrayLen(11, []int{1, 1, 1, 1, 1, 1, 1, 1}))

	//fmt.Println(containsNearbyAlmostDuplicate([]int{1, 2, 3, 1}, 3, 0))
	//fmt.Println(containsNearbyAlmostDuplicate([]int{1, 5, 9, 1, 5, 9}, 2, 3))

	//fmt.Println(containsNearbyAlmostDuplicate([]int{1, 2, 3, 1}, 3, 0))
	//fmt.Println(containsNearbyAlmostDuplicate([]int{1, 5, 9, 1, 5, 9}, 2, 3))

	//fmt.Println(containsNearbyDuplicate([]int{1, 2, 3, 1}, 3))
	//fmt.Println(containsNearbyDuplicate([]int{1, 0, 1, 1}, 1))
	//fmt.Println(containsNearbyDuplicate([]int{1, 2, 3, 1, 2, 3}, 2))

	//fmt.Println(characterReplacement("ABAB", 2))
	//fmt.Println(characterReplacement("AABABBA", 1))

	//fmt.Println(longestOnes([]int{0, 0, 1, 1, 1, 0, 0}, 0))
	//fmt.Println(longestOnes([]int{1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0}, 2))
	//fmt.Println(longestOnes([]int{0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1}, 3))

	//fmt.Println(numSubarrayProductLessThanK([]int{10, 5, 2, 6}, 100))
	//fmt.Println(numSubarrayProductLessThanK([]int{1, 2, 3}, 0))

	//fmt.Println(minOperations([]int{1, 1}, 3))
	//fmt.Println(minOperations([]int{1, 1, 4, 2, 3}, 5))
	//fmt.Println(minOperations([]int{5, 6, 7, 8, 9}, 4))
	//fmt.Println(minOperations([]int{3, 2, 20, 1, 1, 3}, 10))

	//fmt.Println(lengthOfLongestSubstring(""))
	//fmt.Println(lengthOfLongestSubstring("bbbbb"))
	//fmt.Println(lengthOfLongestSubstring("pwwkew"))

	//fmt.Println(findAnagrams("cbaebabacd", "abc"))
	//fmt.Println(findAnagrams("abab", "ab"))
	//fmt.Println(findAnagrams("abaacbabc", "abc"))

	//fmt.Println(checkInclusion("bac", "labcdefg"))
	//fmt.Println(checkInclusion("ab", "eidbaooo"))
	//fmt.Println(checkInclusion("ab", "eidboaoo"))

	//fmt.Println(minWindow("ADOBECODEBANC", "ABC"))
	//fmt.Println(minWindow("anjiadoo", "jia"))
	//fmt.Println(minWindow("abcdefgh", "cde"))
}
