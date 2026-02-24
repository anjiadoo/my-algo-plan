package main

import (
	"container/heap"
)

// 🌟技巧1：当需要创造一条新链表的时候，可以使用「虚拟头结点」简化边界情况的处理

// 合并两个有序链表
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	//「虚拟头结点」技巧
	dummy := &ListNode{}
	p := dummy

	p1 := l1
	p2 := l2

	//合并
	for p1 != nil && p2 != nil {
		if p1.Val <= p2.Val {
			p.Next = p1
			p1 = p1.Next
		} else {
			p.Next = p2
			p2 = p2.Next
		}
		p = p.Next
	}

	//收尾
	if p1 != nil {
		p.Next = p1
	}
	if p2 != nil {
		p.Next = p2
	}
	return dummy.Next
}

// 单链表的分解
func partitionList(head *ListNode, x int) *ListNode {
	// 存放小于 x 的链表的虚拟头结点
	dummy1 := &ListNode{}
	p1 := dummy1

	// 存放大于等于 x 的链表的虚拟头结点
	dummy2 := &ListNode{} // 大于等于x
	p2 := dummy2

	p := head
	for p != nil {
		if p.Val < x {
			p1.Next = p
			p1 = p1.Next
		} else {
			p2.Next = p
			p2 = p2.Next
		}

		// 不能直接让 p 指针前进，
		// p = p.Next
		// 断开原链表中的每个节点的 next 指针
		next := p.Next
		p.Next = nil
		p = next
	}

	p1.Next = dummy2.Next
	return dummy1.Next
}

type PriorityQueue []*ListNode

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].Val < pq[j].Val }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*ListNode)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

// 合并 K 个升序链表
func mergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}

	// 1、实现一个优先级队列
	// 2、把K个升序链表头节点入队
	// 3、每次取最小节点放入结果列表

	pq := &PriorityQueue{}
	heap.Init(pq)
	for _, list := range lists {
		if list != nil {
			heap.Push(pq, list)
		}
	}

	dummy := &ListNode{}
	p := dummy

	for pq.Len() > 0 {
		x := heap.Pop(pq)
		node := x.(*ListNode)

		next := node.Next
		node.Next = nil

		p.Next = node
		p = p.Next

		if next != nil {
			heap.Push(pq, next)
		}
	}

	return dummy.Next
}

func main() {
	//array1 := []int{-9, 3}
	//array2 := []int{5, 7}
	//l1 := NewMyListNode(array1)
	//l2 := NewMyListNode(array2)
	//list := mergeTwoLists(l1, l2)

	//array := []int{1, 4, 3, 2, 5, 2}
	//head := NewMyListNode(array)
	//list := partitionList(head, 3)

	array1 := []int{1, 4, 5}
	array2 := []int{1, 3, 4}
	array3 := []int{2, 6}
	l1 := NewMyListNode(array1)
	l2 := NewMyListNode(array2)
	l3 := NewMyListNode(array3)
	list := mergeKLists([]*ListNode{l1, l2, l3})

	printListNode(list)
}
