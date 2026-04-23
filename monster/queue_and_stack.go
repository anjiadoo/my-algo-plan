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

// 重排链表 https://leetcode.cn/problems/reorder-list/description/
func reorderList(head *ListNode) {

}

func main() {
	fmt.Println(simplifyPath("/home/"))
	fmt.Println(simplifyPath("/home//foo/"))
	fmt.Println(simplifyPath("/home/user/Documents/../Pictures"))
	fmt.Println(simplifyPath("/../"))
}
