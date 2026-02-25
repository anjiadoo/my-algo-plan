package main

import (
	"fmt"
	"math"
)

// 何时运用滑动窗口算法，可通过下面三个问题来判断：
// ❓1、什么时候应该扩大窗口？
// ❓2、什么时候应该缩小窗口？
// ❓3、什么时候应该更新答案？
// 只要能回答这三个问题，就说明可以使用滑动窗口技巧解题。

// 滑动窗口算法伪码框架
func slidingWindow(s string) {
	// 用合适的数据结构记录窗口中的数据，根据具体场景变通
	// 比如说，我想记录窗口中元素出现的次数，就用 map
	// 如果我想记录窗口中的元素和，就可以只用一个 int
	var window = map[rune]int{}

	left, right := 0, 0
	for right < len(s) {

		c := rune(s[right])
		window[c]++ // c 是将移入窗口的字符
		right++     // 增大窗口

		// TODO:进行窗口内数据的一系列更新

		// *** debug 输出的位置 ***
		//fmt.Println("window: [", left, ", ", right, ")")

		// 判断左侧窗口是否要收缩
		for left < right /* && window needs shrink*/ {

			d := rune(s[left])
			window[d]-- // d 是将移出窗口的字符
			left++      // 缩小窗口

			// TODO:进行窗口内数据的一系列更新
		}
	}
}

// 最小覆盖子串
func minWindow(s string, t string) string {
	window, needs := make(map[byte]int), make(map[byte]int)
	for i := 0; i < len(t); i++ {
		needs[t[i]]++
	}

	left, right := 0, 0
	valid := 0
	start, length := 0, math.MaxInt

	for right < len(s) {
		// 增大窗口
		ch := s[right]
		right++

		// 进行窗口内数据的一系列更新
		if _, ok := needs[ch]; ok {
			window[ch]++
			if needs[ch] == window[ch] {
				valid++
			}
		}

		//fmt.Println("window[", left, ",", right, "]")

		for valid == len(needs) {
			// 在这里更新最小覆盖子串
			if right-left < length {
				start = left
				length = right - left
			}

			// 缩小窗口
			d := s[left]
			left++

			// 进行窗口内数据的一系列更新
			if _, ok := needs[d]; ok {
				if window[d] == needs[d] {
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

// 字符串的排列
func checkInclusion(s1 string, s2 string) bool {
	// 给你两个字符串 s1 和 s2 ，写一个函数来判断 s2 是否包含 s1 的 排列。如果是，返回 true ；否则，返回 false 。
	// 换句话说，s1 的排列之一是 s2 的 子串 。
	return false
}

// 找到字符串中所有字母异位词
func findAnagrams(s string, p string) []int {
	// 给定两个字符串 s 和 p，找到 s 中所有 p 的 异位词 的子串，返回这些子串的起始索引。不考虑答案输出的顺序。
	return nil
}

// 无重复字符的最长子串
func lengthOfLongestSubstring(s string) int {
	// 给定一个字符串 s ，请你找出其中不含有重复字符的 最长 子串 的长度。
	return -1
}

func main() {
	fmt.Println(minWindow("ADOBECODEBANC", "ABC"))
	fmt.Println(minWindow("a", "a"))
	fmt.Println(minWindow("a", "aa"))
	fmt.Println(minWindow("aaaaaaaaaaaabbbbbcdd", "abcdd"))
}
