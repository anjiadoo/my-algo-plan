package main

import (
	"fmt"
)

// 单链表实现：
// 1、createListNode
// 2、insertHeadNode
// 3、insertTailNode
// 4、insertIndexNode
// 5、removeIndexNode
// 6、printListNode

type ListNode struct {
	val  int
	next *ListNode
}

func createListNode(array []int) *ListNode {
	if array == nil || len(array) == 0 {
		return nil
	}
	head := &ListNode{val: array[0]}
	curr := head
	for i := 1; i < len(array); i++ {
		curr.next = &ListNode{val: array[i]}
		curr = curr.next
	}
	return head
}

func insertHeadNode(head *ListNode, val int) *ListNode {
	return &ListNode{val: val, next: head}
}

func insertTailNode(head *ListNode, val int) {
	newNode := &ListNode{val: val}
	p := head
	for p.next != nil {
		p = p.next
	}
	p.next = newNode
}

func insertIndexNode(head *ListNode, index, val int) {
	// 假设不会越界
	p := head
	for i := 0; i < index; i++ {
		p = p.next
	}

	tail := p.next
	p.next = &ListNode{val: val, next: tail}
}

func removeIndexNode(head *ListNode, index int) *ListNode {
	// 假设不会越界
	p := head
	for i := 0; i < index; i++ {
		p = p.next
	}
	delNode := p.next
	if delNode != nil {
		p.next = delNode.next
	}
	return delNode
}

func printListNode(head *ListNode) {
	var array []int
	for p := head; p != nil; p = p.next {
		array = append(array, p.val)
	}
	fmt.Println("从前向后遍历单链表:", array)
}

func main() {
	var array = []int{1, 2, 3, 4, 5}

	list := createListNode(array)

	list = insertHeadNode(list, 666)
	insertTailNode(list, 999)
	insertIndexNode(list, 2, 555)
	delNode := removeIndexNode(list, 2)

	fmt.Println("del node:", delNode)

	printListNode(list)
}
