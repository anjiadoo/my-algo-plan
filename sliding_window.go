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
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *
 *     0. slidingWindow(s string)                         // 伪码框架
 *     1. minWindow(s, t string) string                   // 最小覆盖子串，可变窗口·求最短
 *     2. checkInclusion(s1, s2 string) bool              // 字符串的排列，固定窗口
 *     3. findAnagrams(s, p string) []int                 // 字母异位词，固定窗口
 *     4. lengthOfLongestSubstring(s string) int          // 最长无重复子串，可变窗口·求最长
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

// 最小覆盖子串
func minWindow(s string, t string) string {
	// 给定两个字符串s和t，长度分别是m和n，返回s中的「最短窗口」子串
	// 使得该子串包含t中的每一个字符（包括重复字符）。如果没有这样的子串，返回空字符串 ""。

	window, needs := make(map[byte]int), make(map[byte]int)
	for i := 0; i < len(t); i++ {
		needs[t[i]]++
	}

	left, right := 0, 0
	valid := 0                      // valid记录的是「满足条件的字符种类数」，不是字符总数
	start, length := 0, math.MaxInt // ⚠️length初始化为MaxInt，用于求最小值

	for right < len(s) {
		// 增大窗口
		ch := s[right]
		right++

		// 进行窗口内数据的一系列更新
		if _, ok := needs[ch]; ok {
			window[ch]++
			if needs[ch] == window[ch] { // ⚠️只在window[ch]恰好等于needs[ch]时valid++，多出来的不算
				valid++
			}
		}

		//fmt.Println("window[", left, ",", right, "]")

		for valid == len(needs) { // 当所有字符种类都满足时，尝试缩小窗口
			// ⚠️求最短：在缩小窗口之前更新答案
			if right-left < length {
				start = left
				length = right - left // ⚠️窗口长度是right-left，不是right-left+1（左闭右开）
			}

			// 缩小窗口
			d := s[left]
			left++

			// 进行窗口内数据的一系列更新
			// ⚠️缩小窗口时：先判断window[d]==needs[d]再valid--，最后才window[d]--（与扩大窗口的顺序相反）
			if _, ok := needs[d]; ok {
				if window[d] == needs[d] {
					valid--
				}
				window[d]--
			}
		}
	}
	if length == math.MaxInt { // 没找到满足条件的子串
		return ""
	}
	return s[start : start+length]
}

// 字符串的排列
func checkInclusion(s1 string, s2 string) bool {
	// 给你两个字符串s1和s2 ，写一个函数来判断s2是否包含s1的「排列」。
	// 如果是，返回true；否则，返回false 。换句话说，s1的排列之一是s2的「子串」。

	window, needs := make(map[byte]int), make(map[byte]int)
	for i := 0; i < len(s1); i++ {
		needs[s1[i]]++
	}

	left, right := 0, 0
	valid := 0

	for right < len(s2) {
		ch := s2[right]
		right++

		if _, ok := needs[ch]; ok {
			window[ch]++
			if window[ch] == needs[ch] {
				valid++
			}
		}

		//fmt.Printf("s1.size=%d left=%d right=%d\n", len(s1), left, right)

		// ⚠️固定窗口：当窗口大小等于s1长度时，先检查答案再缩窗口
		if right-left == len(s1) {
			if valid == len(needs) {
				return true
			}

			d := s2[left]
			left++ // ⚠️缩小窗口也是++，不是--

			// ⚠️同样注意顺序：先判断再减
			if _, ok := needs[d]; ok {
				if window[d] == needs[d] {
					valid--
				}
				window[d]--
			}
		}
	}

	return false
}

// 找到字符串中所有字母异位词
func findAnagrams(s string, p string) []int {
	// 给定两个字符串s和p，找到s中所有p的「异位词」的子串
	// 返回这些子串的起始索引。不考虑答案输出的顺序。

	window, needs := make(map[byte]int), make(map[byte]int)
	for i := 0; i < len(p); i++ {
		needs[p[i]]++
	}

	left, right := 0, 0
	valid := 0
	var res []int

	for right < len(s) {
		ch := s[right]
		right++

		if _, ok := needs[ch]; ok {
			window[ch]++
			if window[ch] == needs[ch] {
				valid++
			}
		}

		// ⚠️固定窗口：与checkInclusion逻辑相同，只是收集所有起始索引而非返回true/false
		if right-left == len(p) {
			if valid == len(needs) {
				res = append(res, left)
			}

			d := s[left]
			left++ // ⚠️缩小窗口也是++

			// ⚠️同样注意顺序：先判断再减
			if _, ok := needs[d]; ok {
				if window[d] == needs[d] {
					valid--
				}
				window[d]--
			}
		}
	}

	return res
}

// 无重复字符的最长子串
func lengthOfLongestSubstring(s string) int {
	// 给定一个字符串s，请你找出其中不含有重复字符的「最长子串」的长度。

	window := make(map[byte]int)

	left, right := 0, 0
	res := 0

	for right < len(s) {
		ch := s[right]
		right++
		window[ch]++

		// 出现重复字符，缩小窗口直到没有重复
		for window[ch] > 1 {
			d := s[left]
			left++
			window[d]--
		}

		// ⚠️求最长：在缩完窗口之后更新答案（与最小覆盖子串相反）
		if right-left > res {
			res = right - left
		}
	}
	return res
}

func main() {
	//fmt.Println(minWindow("ADOBECODEBANC", "ABC"))
	//fmt.Println(minWindow("a", "a"))
	//fmt.Println(minWindow("a", "aa"))
	//fmt.Println(minWindow("aaaaaaaaaaaabbbbbcdd", "abcdd"))

	//fmt.Println(checkInclusion("ab", "eidbaooo"))
	//fmt.Println(checkInclusion("ab", "eidboaoo"))

	//fmt.Println(findAnagrams("cbaebabacd", "abc"))
	//fmt.Println(findAnagrams("abab", "ab"))

	fmt.Println(lengthOfLongestSubstring("abcabcbb"))
	fmt.Println(lengthOfLongestSubstring("bbbbb"))
	fmt.Println(lengthOfLongestSubstring("pwwkew"))
}
