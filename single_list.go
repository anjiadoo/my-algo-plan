package main

import (
	"fmt"
)

// 单链表实现：
// 🌟技巧1：头插法技巧 - 头插法直接返回新节点作为新的head，避免单独处理空链表
// 🌟技巧2：前驱节点操作技巧 - 插入/删除时先找前一个节点（index-1位置），通过p.next操作目标节点

// 0、func NewMyListNode(array []int) *ListNode
// 1、func insertHeadNode(head *ListNode, val int) *ListNode
// 2、func insertTailNode(head *ListNode, val int)
// 3、func insertIndexNode(head *ListNode, index, val int)
// 4、func removeIndexNode(head *ListNode, index int) *ListNode
// 5、func printListNode(head *ListNode)

type ListNode struct {
	Val  int
	Next *ListNode
}

func NewMyListNode(array []int) *ListNode {
	if array == nil || len(array) == 0 {
		return nil
	}
	head := &ListNode{Val: array[0]}
	curr := head
	for i := 1; i < len(array); i++ {
		curr.Next = &ListNode{Val: array[i]}
		curr = curr.Next
	}
	return head
}

func insertHeadNode(head *ListNode, val int) *ListNode {
	return &ListNode{Val: val, Next: head}
}

func insertTailNode(head *ListNode, val int) {
	newNode := &ListNode{Val: val}
	p := head
	for p.Next != nil {
		p = p.Next
	}
	p.Next = newNode
}

func insertIndexNode(head *ListNode, index, val int) {
	// 假设不会越界
	p := head
	for i := 0; i < index; i++ {
		p = p.Next
	}

	tail := p.Next
	p.Next = &ListNode{Val: val, Next: tail}
}

func removeIndexNode(head *ListNode, index int) *ListNode {
	// 假设不会越界
	p := head
	for i := 0; i < index; i++ {
		p = p.Next
	}
	delNode := p.Next
	if delNode != nil {
		p.Next = delNode.Next
	}
	return delNode
}

func printListNode(head *ListNode) {
	var array []int
	for p := head; p != nil; p = p.Next {
		array = append(array, p.Val)
	}
	fmt.Println("从前向后遍历单链表:", array)
}

func main1() {
	var array = []int{1, 2, 3, 4, 5}

	list := NewMyListNode(array)

	list = insertHeadNode(list, 666)
	insertTailNode(list, 999)
	insertIndexNode(list, 2, 555)
	delNode := removeIndexNode(list, 2)

	fmt.Println("del node:", delNode)

	printListNode(list)
}
