/*
 * ============================================================================
 *                    📘 数组双指针 · 核心记忆框架
 * ============================================================================
 * 【三大指针模式】
 *
 *   ① 快慢指针（同向）：fast 探路，slow 维护"干净区间"
 *      适用：原地删除、去重、移零
 *
 *   ② 对撞指针（相向）：left/right 从两端向中间夹逼
 *      适用：有序数组两数之和、平方排序、回文验证
 *
 *   ③ 三路分区指针：p0/p/p2 同时维护三个区间
 *      适用：荷兰国旗问题（颜色分类）
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式一：快慢指针 —— 原地过滤/去重】
 *
 *   两种写法的核心区别，只有一行之差，但返回值语义完全不同：
 *
 *   写法A：先移动 slow，再赋值（用于"去重"）
 *     if 条件满足 {
 *         slow++
 *         nums[slow] = nums[fast]   // slow 先走一步再写
 *     }
 *     fast++
 *     返回 slow + 1                  // slow 是最后一个有效元素的下标
 *
 *   写法B：先赋值，再移动 slow（用于"移除"）
 *     if 条件满足 {
 *         nums[slow] = nums[fast]   // slow 原地写
 *         slow++
 *     }
 *     fast++
 *     返回 slow                      // slow 是有效元素的个数
 *
 *   ⚠️ 易错点1：去重(写法A) vs 移除(写法B) 的 slow++ 位置不同，返回值也不同
 *      · removeDuplicates  → 写法A，返回 slow+1
 *      · removeElement     → 写法B，返回 slow
 *      · moveZeroes        → 写法B（非零的用写法B移到前面，再补零）
 *      记忆法：去重时 slow 指向"已放好的最后一个"，所以先跳再放；
 *              移除时 slow 指向"下一个空位"，所以先放再跳。
 *
 *   ⚠️ 易错点2：去重 II（允许重复至多2次）的 count 重置时机
 *      count 统计的是当前值连续出现次数。
 *      重置条件：fast 指向的元素 != fast 前一个元素，即换了新值
 *        if fast < len(nums) && nums[fast] != nums[fast-1] { count = 0 }
 *      注意是 fast++ 之后立刻判断，而不是循环开头。
 *
 *   ⚠️ 易错点3：去重 II 中 slow < fast 的守卫条件
 *      slow 先于 fast 移动，起点相同(都是0)，开头 slow==fast。
 *      若不加 slow < fast，会把 nums[slow] 自己赋给自己，count 也会被误消耗。
 *      只有 slow 真的落后于 fast 时，才允许"因count<2而写入"。
 *
 *   ⚠️ 易错点4：moveZeroes 不能只用快慢指针
 *      快慢指针把非零元素移到前面后，slow 之后的位置需要手动补零：
 *        for ; slow < len(nums); slow++ { nums[slow] = 0 }
 *      不补零则尾部残留原来的数。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式二：对撞指针 —— 有序数组夹逼】
 *
 *   口诀：左右夹逼，偏大右移，偏小左移，相等命中
 *     sum > target → right--  （右侧偏大，收缩）
 *     sum < target → left++   （左侧偏小，扩张）
 *     sum == target → 找到
 *
 *   ⚠️ 易错点5：循环条件是 left < right（不是 <=）
 *      left == right 表示两指针指向同一元素，不能再组成"两数"，需停止。
 *
 *   ⚠️ 易错点6：twoSum 返回的是 1-indexed 下标
 *      return []int{left + 1, right + 1}  // 下标从1开始！
 *
 *   有序数组的平方（sortedSquares）：
 *     · 关键洞察：有序数组绝对值最大的一定在两端，从两端向中间对撞
 *     · 结果数组从后往前填（p 从末尾递减）：
 *         abs(left) > abs(right) → res[p] = nums[left]²; left++
 *         else                   → res[p] = nums[right]²; right--
 *     ⚠️ 易错点7：结果数组必须从尾部填，不能从头填
 *        因为每次取的是当前最大值，天然是降序，尾部填入才能保证升序结果。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式三：逆向双指针 —— 倒序合并】
 *
 *   merge1（合并有序数组）：
 *     · 关键洞察：nums1 尾部有空位，从后往前填不会覆盖未读数据
 *     · p1 指向 nums1 有效末尾(m-1)，p2 指向 nums2 末尾(n-1)，p 指向填充位置(m+n-1)
 *     · 每次取两者中较大的填入 p 位置，对应指针后退
 *
 *   ⚠️ 易错点8：循环结束后只需处理 p2 剩余，不需要处理 p1
 *      p1 剩余说明 nums1 的剩余部分已在原位，无需移动。
 *      p2 剩余说明 nums2 还有更小的元素，需要填入 nums1 前部。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式四：三路分区（荷兰国旗）】
 *
 *   sortColors：将 [0,1,2] 分到三个区间，一次遍历
 *     p0 = 0：[0, p0) 全是 0
 *     p2 = n-1：(p2, n-1] 全是 2
 *     p：当前扫描指针，[p0, p2] 是待处理区间
 *
 *   三路处理：
 *     nums[p] == 0 → swap(p, p0); p0++  （0 放到左区，p0 扩张）
 *     nums[p] == 2 → swap(p, p2); p2--  （2 放到右区，p2 收缩）
 *     nums[p] == 1 → p++                 （1 留中间，直接跳过）
 *
 *   ⚠️ 易错点9：遇到 0 时 p 不递增，遇到 2 时 p 也不递增！
 *      · 遇到 2：swap 后 p2-- ，但 p 不动。原因：从右边换来的元素未被检查过，
 *        需要在下一轮继续判断 nums[p] 是否为 0/1/2。
 *      · 遇到 0：swap 后 p0++ ，p 也不能盲目递增。
 *        代码用 if p < p0 { p = p0 } 来同步 p 跟上 p0，
 *        因为从左边换来的值一定是 1（p0 左边都是已确认的0），可安全跳过。
 *      · 只有遇到 1 时才 p++，因为 1 就待在中间不需要任何交换。
 *
 *   ⚠️ 易错点10：终止条件是 p <= p2（不是 p < len(nums)）
 *      p2 右边已全是 2，一旦 p > p2 说明未处理区间为空，可以停止。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式五：中心扩展 —— 回文问题】
 *
 *   最长回文子串（longestPalindrome）：
 *     · 从每个位置向两侧扩展，同时处理奇偶两种情况
 *     · 奇数长度：palindrome(s, i, i)     （单字符为中心）
 *     · 偶数长度：palindrome(s, i, i+1)   （两字符间隙为中心）
 *
 *   ⚠️ 易错点11：必须同时枚举奇偶两种中心，缺一会漏掉答案
 *      "cbbd" 的最长回文是 "bb"（偶数），只枚举奇数中心会错误返回 "b"。
 *
 *   ⚠️ 易错点12：扩展停止后截取的下标是 s[l+1 : r]（不是 s[l:r+1]）
 *      循环退出时 l 和 r 都越界了一格（s[l] != s[r] 或越界才退出），
 *      所以实际回文范围是 [l+1, r-1]，Go 切片写法是 s[l+1 : r]。
 *
 *   验证回文串（isPalindrome1）：
 *     · 先过滤：只保留字母和数字，大写转小写（+32）
 *     · 再对撞：left/right 向中间夹逼
 *
 *   ⚠️ 易错点13：大写转小写用 s[i] + 32（ASCII 偏移），而非调用库函数
 *      'A'=65, 'a'=97，差值恰好是 32，可直接加偏移。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式六：纵向扫描 —— 最长公共前缀】
 *
 *   longestCommonPrefix：以第一个字符串为基准，逐列比较所有字符串同一位置
 *     外层：列下标 chIdx 从 0 到 len(baseStr)-1
 *     内层：对 strs[1..n-1] 与前一个字符串比较 strs[i][chIdx] vs strs[i-1][chIdx]
 *     一旦不一致或越界，立即返回 baseStr[0:chIdx]
 *
 *   ⚠️ 易错点14：必须先判断长度越界，再判断字符相等
 *      if chIdx >= len(thisStr) || chIdx >= len(prevStr) || thisStr[chIdx] != prevStr[chIdx]
 *      短路求值：越界判断在前，防止越界访问 thisStr[chIdx]。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次写完双指针题后对照检查
 *
 *     ✅ 快慢指针：slow++ 在赋值前还是后？返回 slow 还是 slow+1？
 *     ✅ 去重II：count 重置条件写对了吗？slow < fast 的守卫加了吗？
 *     ✅ moveZeroes：快慢结束后有没有补零？
 *     ✅ 对撞指针：循环条件是 left < right（不是 <=）？
 *     ✅ 平方排序：结果数组是从尾部填的吗？
 *     ✅ 逆向合并：循环后只处理 p2 剩余（不处理 p1 剩余）？
 *     ✅ 三路分区：遇到2时 p 不动，遇到0时用 if p<p0 同步，遇到1才 p++？
 *     ✅ 中心扩展：奇偶中心都枚举了吗？切片取 s[l+1:r]（不是 s[l:r+1]）？
 *     ✅ 公共前缀：越界判断在字符比较之前？
 * ============================================================================
 */

package main

import (
	"fmt"
	"sort"
	"strings"
)

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

func main() {
	fmt.Println(advantageCount([]int{2, 7, 11, 15}, []int{1, 10, 4, 11}))
	fmt.Println(advantageCount([]int{12, 24, 8, 32}, []int{13, 25, 32, 11}))

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
