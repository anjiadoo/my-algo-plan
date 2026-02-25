package main

import (
	"fmt"
	"math"
)

// 滑动窗口算法实现：
// 🌟技巧1：窗口框架技巧 - 滑动窗口本质是双指针的[left, right)左闭右开区间，right++扩大窗口，left++缩小窗口（伪码框架）
// 🌟技巧2：needs/valid计数技巧 - 用needs记录目标字符频次，valid记录已满足的字符种类数，valid==len(needs)表示窗口满足条件（最小覆盖子串）
// 🌟技巧3：先扩后缩技巧 - 扩大窗口时先right++再处理数据，缩小窗口时先处理数据再left++，注意扩缩操作的对称性（最小覆盖子串）
// 🌟技巧4：固定窗口技巧 - 当需要查找固定长度的子串时，用right-left==len(target)判断窗口大小，满足时先检查答案再缩窗口（字符串排列/异位词）
// 🌟技巧5：最长子串技巧 - 求最长子串时，在缩完窗口之后更新答案；求最短子串时，在缩窗口之前更新答案（无重复字符最长子串 vs 最小覆盖子串）
//
// ⚠️易错点1：valid的更新时机 - 扩大窗口时window[ch]==needs[ch]才valid++，缩小窗口时必须先判断window[d]==needs[d]再valid--，最后才window[d]--，顺序不能反
// ⚠️易错点2：缩小窗口也是left++ - 初学者容易写成left--，缩小窗口方向和扩大窗口一致，都是向右移动
// ⚠️易错点3：窗口区间是左闭右开[left, right) - right指向的是下一个待加入的元素，所以窗口长度是right-left而非right-left+1
// ⚠️易错点4：更新答案的位置 - 固定窗口在right-left==len时更新，可变窗口(最小覆盖)在满足条件的for循环内更新，最长子串在缩完窗口后更新

// 何时运用滑动窗口算法，可通过下面三个问题来判断：
// ❓1、什么时候应该扩大窗口？
// ❓2、什么时候应该缩小窗口？
// ❓3、什么时候应该更新答案？
// 只要能回答这三个问题，就说明可以使用滑动窗口技巧解题。

// 0、func slidingWindow(s string)                          // 伪码框架
// 1、func minWindow(s string, t string) string             // 最小覆盖子串
// 2、func checkInclusion(s1 string, s2 string) bool        // 字符串的排列
// 3、func findAnagrams(s string, p string) []int           // 找到字符串中所有字母异位词
// 4、func lengthOfLongestSubstring(s string) int           // 无重复字符的最长子串

// 滑动窗口算法伪码框架
func slidingWindow(s string) {
	// 用合适的数据结构记录窗口中的数据，根据具体场景变通
	var window = map[byte]int{}

	left, right := 0, 0 // ⚠️窗口区间是[left, right)，初始时窗口为空
	for right < len(s) {

		c := s[right]
		right++     // 增大窗口，right向右移动
		window[c]++ // c 是将移入窗口的字符

		// 窗口内数据更新

		// *** debug 输出的位置 ***
		//fmt.Println("window: [", left, ", ", right, ")")

		// 判断左侧窗口是否要收缩
		for left < right /* && window needs shrink*/ {

			// ⚠️通常在这里判断目标值，注意更新答案的位置取决于求最长还是最短

			d := s[left]
			left++      // ⚠️缩小窗口也是++，不是--！left也是向右移动
			window[d]-- // d 是将移出窗口的字符

			// 窗口内数据更新
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
