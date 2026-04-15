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

func main() {
	fmt.Println(checkInclusion("bac", "labcdefg"))
	fmt.Println(checkInclusion("ab", "eidbaooo"))
	fmt.Println(checkInclusion("ab", "eidboaoo"))

	//fmt.Println(minWindow("ADOBECODEBANC", "ABC"))
	//fmt.Println(minWindow("anjiadoo", "jia"))
	//fmt.Println(minWindow("abcdefgh", "cde"))
}
