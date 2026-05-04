/*
 * ============================================================================
 *                   📘 数组双指针算法全集 · 核心记忆框架
 * ============================================================================
 * 【一句话理解数组双指针】
 *
 *   数组双指针 = 利用两个索引变量协作遍历，把 O(n²) 暴力降为 O(n)。
 *   核心在于：利用数组的有序性或问题的单调性，让两个指针「有策略地移动」，
 *   避免重复扫描已排除的区域。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【三种指针模式对比】
 *
 *   模式            初始位置             移动方向           代表题目
 *   ──────────────────────────────────────────────────────────────────────
 *   对撞指针        left=0, right=n-1   相向而行(→ ←)     两数之和II、反转数组、接雨水
 *   快慢指针        slow=0, fast=0      同向而行(→ →)     删除重复项、移除元素、移动零
 *   中心扩散        left=i, right=i     背向而行(← →)     最长回文子串、回文验证
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【对撞指针的核心思想】
 *
 *   前提：数组有序（或问题具有「缩小搜索空间」的单调性）。
 *
 *   策略：
 *     - 当 sum > target → right--（让和变小）
 *     - 当 sum < target → left++（让和变大）
 *     - 当 sum == target → 找到答案
 *
 *   为什么正确？每次移动都「排除」了一整行/列的搜索空间，不会错过答案。
 *
 *   ⚠️ 易错点1：对撞指针要求有序，无序数组须先排序或用哈希表
 *      无序数组不具备单调性，left++/right-- 无法确定方向，必须排序后再用。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【快慢指针的核心思想 —— 原地修改数组】
 *
 *   slow 指针维护「结果数组」的边界，fast 指针负责扫描原数组。
 *   slow 只在遇到「有效元素」时前进，天然过滤掉不需要的元素。
 *
 *   两种边界约定（初学极易混淆）：
 *
 *   ┌──────────────────────────────────────────────────────────────────┐
 *   │ 约定A [0..slow] 有效区（含slow位置）                             │
 *   │   初始：slow=0, fast=0                                          │
 *   │   操作：先 slow++，再 nums[slow] = nums[fast]                   │
 *   │   返回：slow + 1（长度）                                         │
 *   │   代表：removeDuplicates（删除有序数组重复项）                     │
 *   ├──────────────────────────────────────────────────────────────────┤
 *   │ 约定B [0..slow) 有效区（不含slow位置）                            │
 *   │   初始：slow=0, fast=0                                          │
 *   │   操作：先 nums[slow] = nums[fast]，再 slow++                   │
 *   │   返回：slow（长度，即下一个待填位置的下标）                       │
 *   │   代表：removeElement（移除指定元素）、moveZeroes（移动零）        │
 *   └──────────────────────────────────────────────────────────────────┘
 *
 *   ⚠️ 易错点2：约定A 中 slow++ 在赋值前，约定B 中 slow++ 在赋值后
 *      搞反会导致第一个元素被覆盖（约定A）或者多保留一个无效元素（约定B）。
 *
 *   ⚠️ 易错点3：removeDuplicates 判断条件是 nums[slow] != nums[fast]
 *      这里比较的是 slow 位置（已确认的最后一个有效元素），不是 fast-1 位置。
 *      用 nums[fast] != nums[fast-1] 虽然在有序数组中结果相同，但语义不同，
 *      扩展到「最多保留K个」时就不适用了。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【中心扩散的核心思想】
 *
 *   从某个中心点同时向左向右扩展，直到不满足条件为止。
 *   用于回文问题：如果 s[l] == s[r]，则向两端继续扩展。
 *
 *   关键：必须枚举两种中心：
 *     - 奇数长度回文：中心是单个字符 palindrome(s, i, i)
 *     - 偶数长度回文：中心是两字符间隙 palindrome(s, i, i+1)
 *
 *   ⚠️ 易错点4：只枚举一种中心会漏掉一半的回文
 *      "cbbd" 的最长回文是 "bb"（偶数长度），只枚举 (i,i) 会得到 "b"。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【荷兰国旗问题（三指针分区）】
 *
 *   三个指针 p0、p、p2 把数组分成四个区域：
 *     [0, p0)       → 全是 0
 *     [p0, p)       → 全是 1
 *     [p, p2]       → 未处理区
 *     (p2, n-1]     → 全是 2
 *
 *   移动规则：
 *     nums[p]==0 → 交换到左区(p0)，p0++
 *     nums[p]==2 → 交换到右区(p2)，p2--
 *     nums[p]==1 → p++（保持在中间区）
 *
 *   ⚠️ 易错点5：nums[p]==0 交换后 p 可以前进，nums[p]==2 交换后 p 不能前进
 *      从 p0 换来的值一定是 1（因为 p0 <= p，之前已被 p 扫过），所以可以 p++。
 *      从 p2 换来的值未知（p2 在 p 后面，还没被扫过），必须再次检查，不能 p++。
 *      实现中用 if p < p0 { p = p0 } 来保证 p 不落后于 p0。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【逆向双指针（从尾部开始）】
 *
 *   适用场景：结果数组从大到小填充，或原地合并时避免覆盖。
 *
 *   代表题目：
 *     - 合并两个有序数组：从 nums1 尾部开始填充，避免覆盖 nums1 前面的有效元素
 *     - 有序数组的平方：  绝对值最大的在两端，从结果数组尾部开始填充
 *
 *   ⚠️ 易错点6：合并有序数组不能从前向后填充
 *      nums1 前半部分存放有效数据，从前向后会覆盖还未处理的元素。
 *      从后向前填充，被覆盖的位置（尾部的0）不影响数据完整性。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【田忌赛马（贪心+双指针）】
 *
 *   核心贪心策略：排序后，尽可能用「刚好比对手大」的去赢；赢不了就用最差的去送。
 *
 *   实现：nums1 排序，nums2 带下标排序（因为要恢复原始位置）。
 *   从 nums2 最大的开始匹配：
 *     - nums1[right] > nums2[i] → 用最强的赢，right--
 *     - 否则 → 用最弱的送，left++
 *
 *   ⚠️ 易错点7：必须记住 nums2 的原始下标
 *      排序后 nums2 的顺序变了，但结果数组要按 nums2 的原始位置填入。
 *      所以 nums2 必须带下标排序（用 pair 结构或索引数组）。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次手写数组双指针后对照检查
 *
 *     ✅ 对撞指针：数组是否有序？循环条件是 left < right 还是 left <= right？
 *     ✅ 快慢指针：slow 维护的是 [0..slow] 还是 [0..slow)？返回 slow 还是 slow+1？
 *     ✅ 快慢指针：是先移动 slow 再赋值，还是先赋值再移动 slow？
 *     ✅ 中心扩散：是否同时枚举了奇数和偶数长度的情况？
 *     ✅ 三指针分区：与左区交换后 p 要前进，与右区交换后 p 不要前进？
 *     ✅ 逆向填充：是否从结果数组的最后一个位置开始？循环条件是否涵盖剩余元素？
 *     ✅ 边界条件：空数组、单元素数组、全相同元素是否正确处理？
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *
 *     对撞指针：
 *       1. reverseString(s []byte)                           // 反转数组
 *       2. twoSum(numbers []int, target int) []int           // 两数之和II（有序）
 *       3. isPalindrome1(s string) bool                      // 验证回文串
 *       4. sortedSquares(nums []int) []int                   // 有序数组的平方（逆向对撞）
 *       5. merge1(nums1 []int, m int, nums2 []int, n int)   // 合并两个有序数组（逆向）
 *
 *     快慢指针：
 *       6. removeDuplicates(nums []int) int                  // 删除有序数组重复项
 *       7. removeDuplicates2(nums []int) int                 // 删除有序数组重复项II（最多保留2个）
 *       8. removeElement(nums []int, val int) int            // 移除指定元素
 *       9. moveZeroes(nums []int)                            // 移动零到末尾
 *
 *     中心扩散：
 *      10. longestPalindrome(s string) string                // 最长回文子串
 *
 *     三指针分区：
 *      11. sortColors(nums []int) []int                      // 颜色分类（荷兰国旗）
 *
 *     贪心+双指针：
 *      12. advantageCount(nums1, nums2 []int) []int          // 优势洗牌（田忌赛马）
 *
 *     其他：
 *      13. longestCommonPrefix(strs []string) string         // 最长公共前缀（纵向扫描）
 * ============================================================================
 */

package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

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

// 删除有序数组中的重复项 https://leetcode.cn/problems/remove-duplicates-from-sorted-array/description/
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	slow, fast := 0, 0
	for fast < len(nums) {
		if nums[slow] != nums[fast] {
			slow++
			nums[slow] = nums[fast] // 维护nums[0..slow]无重复
		}
		fast++
	}
	return slow + 1
}

// 删除有序数组中的重复项II https://leetcode.cn/problems/remove-duplicates-from-sorted-array-ii/description/
func removeDuplicates2(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	slow, fast, count := 0, 0, 0
	for fast < len(nums) {
		// 1, 1, 1, 2, 2, 2, 3, 4, 5, 5, 6
		// s
		// f
		// 维护nums[0..slow]无重复
		// slow++在fast++前面，所以需要判断 slow < fast
		if nums[slow] != nums[fast] {
			slow++
			nums[slow] = nums[fast]
		} else if slow < fast && count < 2 {
			slow++
			nums[slow] = nums[fast]
		}
		fast++
		count++
		if fast < len(nums) && nums[fast] != nums[fast-1] {
			count = 0
		}
	}
	return slow + 1
}

// 移除指定元素 https://leetcode.cn/problems/remove-element/description/
func removeElement(nums []int, val int) int {
	slow := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != val {
			nums[slow] = nums[fast] // 维护nums[0..slow)无val
			slow++
		}
	}
	return slow
}

// 移动零 https://leetcode.cn/problems/move-zeroes/
func moveZeroes(nums []int) {
	slow := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != 0 {
			nums[slow] = nums[fast] // 维护nums[0..slow)无val
			slow++
		}
	}
	for ; slow < len(nums); slow++ {
		nums[slow] = 0
	}
}

// 两数之和II - 输入有序数组 https://leetcode.cn/problems/two-sum-ii-input-array-is-sorted/
func twoSum(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1
	for left < right {
		sum := numbers[left] + numbers[right]
		if sum == target {
			return []int{left + 1, right + 1}
		} else if sum > target {
			right--
		} else {
			left++
		}
	}
	return []int{-1, -1}
}

// 最长回文子串 https://leetcode.cn/problems/longest-palindromic-substring/
func longestPalindrome(s string) string {
	// 从中心向两端扩散的双指针技巧
	maxLongestStr := ""

	palindrome := func(s string, l, r int) string {
		// 防止索引越界
		for l >= 0 && r < len(s) && s[l] == s[r] {
			l--
			r++
		}
		return s[l+1 : r]
	}

	for i := 0; i < len(s); i++ {
		str1 := palindrome(s, i, i)
		if len(str1) > len(maxLongestStr) {
			maxLongestStr = str1
		}

		str2 := palindrome(s, i, i+1)
		if len(str2) > len(maxLongestStr) {
			maxLongestStr = str2
		}
	}
	return maxLongestStr
}

// 验证回文串 https://leetcode.cn/problems/valid-palindrome/
func isPalindrome1(s string) bool {
	sb := strings.Builder{}
	for i := 0; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' {
			sb.WriteByte(s[i])
		} else if s[i] >= 'A' && s[i] <= 'Z' {
			sb.WriteByte(s[i] + 32)
		} else if s[i] >= '0' && s[i] <= '9' {
			sb.WriteByte(s[i])
		}
	}

	s = sb.String()

	left, right := 0, len(s)-1
	for left < right {
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}
	return true
}

// 颜色分类 https://leetcode.cn/problems/sort-colors/
func sortColors(nums []int) []int {
	// 双指针p0/p2=>[0,p0)存放0，(p2,len(nums)-1]存放2
	p0, p2 := 0, len(nums)-1
	p := 0

	for p <= p2 {
		if nums[p] == 0 {
			nums[p], nums[p0] = nums[p0], nums[p]
			p0++
		} else if nums[p] == 2 {
			nums[p], nums[p2] = nums[p2], nums[p]
			p2--
		} else {
			p++
		}
		if p < p0 {
			p = p0
		}
	}
	return nums
}

// 合并两个有序数组 https://leetcode.cn/problems/merge-sorted-array/description/
func merge1(nums1 []int, m int, nums2 []int, n int) {
	p1 := m - 1
	p2 := n - 1
	p := len(nums1) - 1
	for p1 >= 0 && p2 >= 0 {
		if nums1[p1] >= nums2[p2] {
			nums1[p] = nums1[p1]
			p1--
			p--
		} else {
			nums1[p] = nums2[p2]
			p2--
			p--
		}
	}
	for p2 >= 0 {
		nums1[p] = nums2[p2]
		p2--
		p--
	}
}

// 有序数组的平方 https://leetcode.cn/problems/squares-of-a-sorted-array/description/
func sortedSquares(nums []int) []int {
	abs := func(a int) int {
		if a >= 0 {
			return a
		}
		return -a
	}

	left, right := 0, len(nums)-1
	res := make([]int, len(nums))
	p := len(res) - 1

	for left <= right {
		if abs(nums[left]) > abs(nums[right]) {
			res[p] = nums[left] * nums[left]
			p--
			left++
		} else {
			res[p] = nums[right] * nums[right]
			p--
			right--
		}
	}
	return res
}

// 最长公共前缀 https://leetcode.cn/problems/longest-common-prefix/
func longestCommonPrefix(strs []string) string {
	//输入：strs = ["flower","flow","flight"]
	//输出："fl"

	baseStr := strs[0]

	for chIdx := 0; chIdx < len(baseStr); chIdx++ {
		// 挨个比较strs中的字符串
		for i := 1; i < len(strs); i++ {
			thisStr, prevStr := strs[i], strs[i-1]
			if chIdx >= len(thisStr) || chIdx >= len(prevStr) || thisStr[chIdx] != prevStr[chIdx] {
				return baseStr[0:chIdx]
			}
		}
	}
	return baseStr
}

// 优势洗牌 https://leetcode.cn/problems/advantage-shuffle/
func advantageCount(nums1 []int, nums2 []int) []int {
	type pair struct {
		index int
		value int
	}

	array2 := make([]pair, len(nums2))
	for i := 0; i < len(nums1); i++ {
		array2[i] = pair{index: i, value: nums2[i]}
	}

	// 按战力排序,用最快的比,比得过就比,比不过就用最差的糊弄
	sort.Ints(nums1)
	sort.Slice(array2, func(i, j int) bool {
		return array2[i].value < array2[j].value
	})

	left, right := 0, len(nums1)-1
	res := make([]int, len(nums1))

	for i := len(nums1) - 1; i >= 0; i-- {
		if nums1[right] > array2[i].value {
			res[array2[i].index] = nums1[right]
			right--
		} else {
			res[array2[i].index] = nums1[left]
			left++
		}
	}
	return res
}

// 接雨水 https://leetcode.cn/problems/trapping-rain-water/
func trap(height []int) int {
	// 类似单调栈的思路，找到中间最高的墙
	maxIndex := -1
	maxNum := math.MinInt
	for i := 0; i < len(height); i++ {
		if height[i] > maxNum {
			maxIndex = i
			maxNum = height[i]
		}
	}

	res := 0

	// 保留[0,maxIndex]区间递增的元素
	var incr []int
	for i := 0; i <= maxIndex; i++ {
		if len(incr) > 0 && incr[len(incr)-1] > height[i] {
			res += incr[len(incr)-1] - height[i]
		} else {
			incr = append(incr, height[i])
		}
	}

	// 保留[maxIndex,end]区间递减的元素
	var decr []int
	for i := len(height) - 1; i > maxIndex; i-- {
		if len(decr) > 0 && decr[len(decr)-1] >= height[i] {
			res += decr[len(decr)-1] - height[i]
		} else {
			decr = append(decr, height[i])
		}
	}
	return res
}

func trap2(height []int) int {
	// 双指针解法，左右两边各维护一堵高墙
	left, right := 0, len(height)-1
	lMax, rMax := 0, 0
	res := 0

	for left < right {
		lMax = max(lMax, height[left])
		rMax = max(rMax, height[right])
		if lMax < rMax {
			res += lMax - height[left]
			left++
		} else {
			res += rMax - height[right]
			right--
		}
	}
	return res
}

// 盛最多水的容器 https://leetcode.cn/problems/container-with-most-water/description/
func maxArea(height []int) int {
	left, right := 0, len(height)-1
	res := 0
	for left < right {
		res = max(res, (right-left)*min(height[left], height[right]))
		if height[left] > height[right] {
			right--
		} else {
			left++
		}
	}
	return res
}

func main() {

	fmt.Println(maxArea([]int{1, 8, 6, 2, 5, 4, 8, 3, 7}))
	fmt.Println(maxArea([]int{1, 1}))

	//fmt.Println(trap2([]int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}))
	//fmt.Println(trap2([]int{4, 2, 0, 3, 2, 5}))

	//fmt.Println(advantageCount([]int{2, 7, 11, 15}, []int{1, 10, 4, 11}))
	//fmt.Println(advantageCount([]int{12, 24, 8, 32}, []int{13, 25, 32, 11}))

	//fmt.Println(longestCommonPrefix([]string{"flower", "flow", "flight"}))
	//fmt.Println(longestCommonPrefix([]string{"anjiadoo", "anji", "anjido"}))

	//fmt.Println(sortedSquares([]int{-4, -1, 0, 3, 10}))

	//nums1 := []int{1, 2, 3, 0, 0, 0}
	//nums2 := []int{4, 5, 6}
	//merge1(nums1, 3, nums2, 3)
	//fmt.Println(nums1)

	//fmt.Println(sortColors([]int{2, 0, 2, 1, 1, 0}))
	//fmt.Println(sortColors([]int{2, 0, 1}))
	//fmt.Println(sortColors([]int{0, 2, 1, 2, 1, 0, 2, 0, 1}))

	//fmt.Println(isPalindrome1("A man, a plan, a canal: Panama"))
	//fmt.Println(isPalindrome1("race a car"))
	//fmt.Println(isPalindrome1("0P"))

	//fmt.Println(longestPalindrome("babad"))
	//fmt.Println(longestPalindrome("cbbd"))

	//fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
	//fmt.Println(twoSum([]int{2, 3, 4}, 6))
	//fmt.Println(twoSum([]int{-1, 0}, -1))

	//nums := []int{1, 0, 2, 0, 3, 0, 4, 5}
	//moveZeroes(nums)
	//fmt.Println(nums)

	//fmt.Println(removeElement([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 5))

	//nums = []int{1, 1, 1, 2, 2, 2, 3, 4, 5, 5, 6}
	//fmt.Println(removeDuplicates2(nums))
	//fmt.Println(nums)
}
