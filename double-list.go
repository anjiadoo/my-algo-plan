package main

import "fmt"

// 双链表实现：
// 1、createDListNode
// 2、insertHeadNode2
// 3、insertTailNode2
// 4、insertIndexNode2
// 5、removeIndexNode2
// 6、printDListNode

type dListNode struct {
	val  int
	prev *dListNode
	next *dListNode
}

func createDListNode(array []int) *dListNode {
	if array == nil || len(array) == 0 {
		return nil
	}
	head := &dListNode{val: array[0]}
	p := head
	for i := 1; i < len(array); i++ {
		newNode := &dListNode{val: array[i]}
		p.next = newNode
		newNode.prev = p
		p = newNode
	}
	return head
}

func insertHeadNode2(head *dListNode, val int) *dListNode {
	newNode := &dListNode{val: val}
	newNode.next = head
	head.prev = newNode
	head = newNode
	return head
}

func insertTailNode2(head *dListNode, val int) {
	newNode := &dListNode{val: val}

	tail := head
	for tail.next != nil {
		tail = tail.next
	}

	tail.next = newNode
	newNode.prev = tail
}

func insertIndexNode2(head *dListNode, index, val int) {
	// 假设不会越界
	p := head
	for i := 0; i < index; i++ {
		p = p.next
	}

	// 组装新节点
	newNode := &dListNode{val: val}
	newNode.next = p.next
	newNode.prev = p

	// 插入新节点
	p.next.prev = newNode
	p.next = newNode
}

func removeIndexNode2(head *dListNode, index int) *dListNode {
	// 假设不会越界
	p := head
	for i := 0; i < index; i++ {
		p = p.next
	}

	delNode := p.next
	if delNode == nil {
		return nil
	}

	p.next = delNode.next
	if delNode.next != nil {
		delNode.next.prev = p
	}

	// 后指针都置为nil是个好习惯
	delNode.next = nil
	delNode.prev = nil

	return delNode
}

func printDListNode(head *dListNode) {
	var tail *dListNode

	var array []int
	for p := head; p != nil; p = p.next {
		array = append(array, p.val)
		tail = p
	}
	fmt.Println("从头节点向后遍历双链表:", array)

	array = []int{}
	for p := tail; p != nil; p = p.prev {
		array = append(array, p.val)
	}
	fmt.Println("从尾节点向前遍历双链表:", array)
}

func main() {
	var array = []int{1, 2, 3, 4, 5}

	dList := createDListNode(array)

	dList = insertHeadNode2(dList, 666)
	insertTailNode2(dList, 999)
	insertIndexNode2(dList, 2, 555)
	delNode := removeIndexNode2(dList, 7)

	fmt.Println("del node:", delNode)

	printDListNode(dList)
}
