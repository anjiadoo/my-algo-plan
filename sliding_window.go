/*
 * ============================================================================
 *                      📘 滑动窗口算法全集 · 核心记忆框架
 * ============================================================================
 * 【一句话理解滑动窗口】
 *
 *   滑动窗口 = 双指针维护的 [left, right) 左闭右开区间，
 *   right++ 向右扩大窗口纳入新元素，left++ 向右缩小窗口移出旧元素。
 *   本质：用 O(n) 时间枚举所有「满足条件的子串/子数组」，避免暴力 O(n²)。
 *
 *   判断能否用滑动窗口，回答三个问题：
 *     ❓ 什么时候扩大窗口（right++）？
 *     ❓ 什么时候缩小窗口（left++）？
 *     ❓ 什么时候更新答案？
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【三种窗口模式对比】
 *
 *   模式            代表题目              缩窗条件              更新答案时机
 *   ──────────────────────────────────────────────────────────────────────
 *   可变窗口·求最短  最小覆盖子串          valid==len(needs)      缩窗「之前」
 *   可变窗口·求最长  无重复字符最长子串     window[ch] > 1         缩窗「之后」
 *   固定窗口        字符串排列/异位词      right-left==len(s1)   缩窗「之前」判断
 *
 *   口诀：求最短在缩前更新，求最长在缩后更新，固定窗口在缩前判断。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【valid 的核心语义】
 *
 *   valid 记录的是「已完全满足频次要求的字符种类数」，不是字符总数。
 *   当 valid == len(needs) 时，窗口内包含了目标串的所有字符（含重复）。
 *
 *   valid++ 的条件：window[ch] == needs[ch]（恰好满足，多一个不算）
 *   valid-- 的条件：window[d]  == needs[d] （移出前恰好满足，移出后就不满足了）
 *
 *   ⚠️ 易错点1：valid++ 用 ==，不能用 >=
 *      window[ch] 每次只加 1，只有从 needs[ch]-1 变成 needs[ch] 的那一刻才算「刚好满足」。
 *      若用 >=，窗口内超出需求的字符也会触发 valid++，导致 valid 虚高、提前认为满足条件。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【扩窗 & 缩窗操作的对称性】
 *
 *   扩窗（加入 ch）：先 window[ch]++，再判断是否触发 valid++
 *   缩窗（移出 d） ：先判断是否触发 valid--，再 window[d]--
 *
 *   两者顺序「镜像相反」，原因：
 *     扩窗时：要用更新后的 window[ch] 来判断是否刚好达到 needs[ch] → 先加后判
 *     缩窗时：要用移出前的 window[d]  来判断是否恰好等于 needs[d]  → 先判后减
 *
 *   ⚠️ 易错点2：缩窗时不能先 window[d]-- 再判断 valid--
 *      如果先减，window[d] 已经低于 needs[d]，判断条件永远不等，valid 永远不会减，
 *      导致窗口在不满足条件时仍认为 valid==len(needs)，陷入死循环或结果错误。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【可变窗口·求最短：minWindow 最小覆盖子串】
 *
 *   框架：外层 for 扩窗，内层 for (valid==len(needs)) 缩窗
 *
 *   答案更新在缩窗「之前」，原因：
 *     进入内层 for 说明当前窗口已满足条件，此时记录最小长度，
 *     然后再缩窗继续寻找更短的满足窗口。
 *     若先缩再更新，left 已移动，记录的起点 start 就错了。
 *
 *   ⚠️ 易错点3：length 初始化为 math.MaxInt，不能为 0
 *      用 right-left < length 来取最小值，初始值必须足够大才能被第一个合法窗口覆盖。
 *      初始化为 0 会导致任何窗口都不满足 < 条件，答案永远不更新。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【可变窗口·求最长：lengthOfLongestSubstring 无重复字符最长子串】
 *
 *   框架：外层 for 扩窗，内层 for (window[ch] > 1) 缩窗
 *
 *   答案更新在缩窗「之后」，原因：
 *     缩窗的目的是消除刚加入的重复字符，缩完后窗口内必然无重复，
 *     此时 [left, right) 是以 right-1 结尾的最长无重复子串，才适合更新答案。
 *     若在缩窗前更新，窗口内还有重复字符，记录的是无效状态。
 *
 *   ⚠️ 易错点4：内层缩窗条件是 window[ch] > 1，不是 window[ch] >= 1
 *      加入 ch 后 window[ch] 变为 2 才说明有重复，需要缩窗；
 *      等于 1 说明 ch 在窗口中只出现一次，是合法状态，不需要缩。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【固定窗口：checkInclusion / findAnagrams 字符串排列 & 异位词】
 *
 *   框架：外层 for 扩窗，if (right-left == len(s1)) 时判断答案 + 缩窗一格
 *
 *   固定窗口不用内层 for 缩窗，而是 if 缩窗，每次只缩一格，原因：
 *     窗口大小固定为 len(s1)，每扩一格就必须缩一格以维持大小，
 *     不存在「缩到满足条件为止」的逻辑，所以用 if 而非 for。
 *
 *   ⚠️ 易错点5：固定窗口先判断答案，再缩窗
 *      进入 if 时窗口大小恰好等于 len(s1)，是检查答案的时机；
 *      缩窗后窗口大小变为 len(s1)-1，等下一轮 right++ 再恢复，
 *      若先缩再判断，窗口大小已经不等于 len(s1)，判断的是错误状态。
 *
 *   checkInclusion 与 findAnagrams 的唯一区别：
 *     checkInclusion：找到一个即返回 true
 *     findAnagrams  ：收集所有满足条件的 left 下标到结果数组
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【进阶技巧】—— 新增题目揭示的高级模式
 *
 *   ▶ 问题转换（minOperations）
 *     原题是"从两端取元素使和为 x"，等价转换为"找中间和为 sum-x 的最长子数组"。
 *     当题目操作不直接适配滑动窗口时，尝试取补集/反面。
 *
 *   ▶ 子数组计数（numSubarrayProductLessThanK）
 *     每次扩窗+缩窗后，以 right-1 结尾的合法子数组个数 = right - left。
 *     原因：[left..right-1], [left+1..right-1], ..., [right-1..right-1] 都合法。
 *
 *   ▶ 替换类问题的通用模式（longestOnes / characterReplacement）
 *     核心思路：窗口长度 - 窗口内「不需要替换的元素数」> k 时缩窗。
 *     longestOnes：     right-left - oneCount > k     （0 的个数超过 k）
 *     characterReplacement：right-left - maxCharCount > k（非最频字符数超过 k）
 *     注意 characterReplacement 中 maxCharCount 不需要在缩窗时精确维护（只增不减），
 *     因为只有 maxCharCount 增大时窗口才可能变大，旧的较小值不会产生更优解。
 *
 *   ▶ 桶排序 + 滑动窗口（containsNearbyAlmostDuplicate）
 *     当需要在窗口内 O(1) 判断「是否存在差值 ≤ t 的元素」时，
 *     用桶宽 t+1 把数轴切桶，只需检查当前桶 + 左右相邻桶。
 *     负数桶号计算：(x+1)/bucketSize - 1（避免 0 桶被正负共用）。
 *
 *   ▶ 枚举创造缩窗条件（longestSubstring）
 *     原题"每种字符至少出现 k 次"没有天然缩窗条件。
 *     技巧：外层枚举窗口内允许的字符种类数 count∈[1,26]，
 *     当 len(window) > count 时缩窗，将问题转化为标准滑动窗口。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次手写滑动窗口后对照检查
 *
 *     ✅ 窗口区间是 [left, right) 左闭右开，窗口长度是 right-left？
 *     ✅ left、right 初始值都是 0，缩小窗口是 left++，不是 left--？
 *     ✅ valid 记录的是「满足的字符种类数」，不是字符总数？
 *     ✅ valid++ 的条件是 window[ch] == needs[ch]（恰好等于，不是 >=）？
 *     ✅ 缩窗时顺序是「先判断 valid--，再 window[d]--」？
 *     ✅ 求最短子串（minWindow）：答案在缩窗「之前」更新？
 *     ✅ 求最长子串（lengthOfLongestSubstring）：答案在缩窗「之后」更新？
 *     ✅ 固定窗口（checkInclusion/findAnagrams）：答案在缩窗「之前」判断，缩窗用 if 不用 for？
 *     ✅ length 初始化为 math.MaxInt，最终用 length == math.MaxInt 判断无解？
 *     ✅ 子数组计数题：缩窗后 count += right-left，不是 count++？
 *     ✅ 替换题缩窗条件：right-left-xxx > k（窗口长度减去不需要替换的 > k）？
 *     ✅ 问题转换题：有没有把原问题正确转换为求最长/最短子数组？
 *     ✅ 桶排序题：负数桶号用 (x+1)/bucketSize-1 而非直接 x/bucketSize？
 *     ✅ 枚举+滑窗题：外层枚举的边界是否正确（如 1~26）？
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *
 *     0. slidingWindow(s string)                                // 伪码框架
 *     1. minWindow(s, t string) string                          // 最小覆盖子串，可变窗口·求最短
 *     2. checkInclusion(s1, s2 string) bool                     // 字符串的排列，固定窗口
 *     3. findAnagrams(s, p string) []int                        // 字母异位词，固定窗口
 *     4. lengthOfLongestSubstring(s string) int                 // 最长无重复子串，可变窗口·求最长
 *     5. minOperations(nums []int, x int) int                   // 将x减到0的最小操作数，问题转换·求最长
 *     6. numSubarrayProductLessThanK(nums []int, k int) int     // 乘积小于K的子数组，可变窗口·计数
 *     7. longestOnes(nums []int, k int) int                     // 最大连续1的个数III，替换约束·求最长
 *     8. characterReplacement(s string, k int) int              // 替换后的最长重复字符，替换约束·求最长
 *     9. containsNearbyDuplicate(nums []int, k int) bool        // 存在重复元素II，固定窗口·集合判重
 *    10. containsNearbyAlmostDuplicate(nums, indexDiff, valueDiff) // 存在重复元素III，固定窗口·桶排序
 *    11. minSubArrayLen(target int, nums []int) int              // 长度最小的子数组，可变窗口·求最短
 *    12. longestSubstring(s string, k int) int                   // 至少有K个重复字符的最长子串，枚举+可变窗口
 * ============================================================================
 */

package main

import (
	"fmt"
	"math"
)

// 滑动窗口算法伪码框架
func slidingWindow(s, t string) {
	// 用合适的数据结构记录窗口中的数据，根据具体场景变通
	window, needs := make(map[byte]int), make(map[byte]int)
	for i := 0; i < len(t); i++ {
		needs[t[i]]++
	}

	left, right := 0, 0 // ⚠️窗口区间是[left, right)，初始时窗口为空
	valid := 0          // valid记录的是「满足条件的字符种类数」，不是字符总数

	for right < len(s) {

		ch := s[right]
		right++ // 增大窗口，right向右移动

		// 窗口内数据更新
		// ...

		if _, ok := needs[ch]; ok {
			window[ch]++
			if needs[ch] == window[ch] { // ⚠️只在window[ch]恰好等于needs[ch]时valid++，多出来的不算
				valid++
			}
		}

		// *** debug 输出的位置 ***
		//fmt.Println("window: [", left, ", ", right, ")")

		// 判断左侧窗口是否要收缩
		for left < right /* && valid == len(needs) */ {

			// ⚠️通常在这里判断目标值，注意更新答案的位置取决于求最长还是最短

			d := s[left]
			left++ // ⚠️缩小窗口也是++，不是--！left也是向右移动

			// 窗口内数据更新
			// ...

			// ⚠️缩小窗口时：先判断window[d]==needs[d]再valid--，最后才window[d]--（与扩大窗口的顺序相反）
			if _, ok := needs[d]; ok {
				if window[d] == needs[d] {
					valid--
				}
				window[d]--
			}
		}
	}
}

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
