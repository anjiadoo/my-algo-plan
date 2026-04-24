package main

import (
	"fmt"
	"strings"
)

// 简化路径 https://leetcode.cn/problems/simplify-path/description/
func simplifyPath(path string) string {
	// 一定要明确栈中存储的是啥
	// 这里存的是文件夹组成路径
	var stack []string

	for _, str := range strings.Split(path, "/") {
		switch str {
		case "", ".":
			continue
		case "..":
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		default:
			stack = append(stack, str)
		}
	}
	if len(stack) == 0 {
		return "/"
	}
	return "/" + strings.Join(stack, "/")
}

// 有效的括号 https://leetcode.cn/problems/valid-parentheses/description/
func isValid(s string) bool {
	var stack []byte
	for _, ch := range s {
		switch ch {
		case '(', '{', '[':
			stack = append(stack, byte(ch))
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false
			}
			stack = stack[:len(stack)-1]
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				return false
			}
			stack = stack[:len(stack)-1]
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	if len(stack) > 0 {
		return false
	}
	return true
}

// 逆波兰表达式求值 https://leetcode.cn/problems/evaluate-reverse-polish-notation/
func evalRPN(tokens []string) int {
	return -1
}

func main() {

	fmt.Println(isValid("()"))
	//fmt.Println(isValid("()[]{}"))
	//fmt.Println(isValid("([])"))
	//fmt.Println(isValid("([)]"))

	//fmt.Println(simplifyPath("/home/"))
	//fmt.Println(simplifyPath("/home//foo/"))
	//fmt.Println(simplifyPath("/home/user/Documents/../Pictures"))
	//fmt.Println(simplifyPath("/../"))
}
