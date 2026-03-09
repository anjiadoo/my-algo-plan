package main

import (
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func NewMyListNode(nums []int) *ListNode {
	dummy := &ListNode{}
	p := dummy
	for i := 0; i < len(nums); i++ {
		p.Next = &ListNode{Val: nums[i]}
		p = p.Next
	}
	return dummy.Next
}

func (h *ListNode) Display() {
	p := h
	str := "打印单链表: "
	for p != nil {
		str += fmt.Sprintf("%d->", p.Val)
		p = p.Next
	}
	fmt.Println(str + "nil")
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := &ListNode{}
	p := dummy
	p1 := list1
	p2 := list2

	for p1 != nil && p2 != nil {
		if p1.Val > p2.Val {
			p.Next = p2
			p2 = p2.Next
		} else {
			p.Next = p1
			p1 = p1.Next
		}
		p = p.Next
	}
	if p1 != nil {
		p.Next = p1
	}
	if p2 != nil {
		p.Next = p2
	}
	return dummy.Next
}

func partition(head *ListNode, x int) *ListNode {
	dummy1 := &ListNode{}
	p1 := dummy1

	dummy2 := &ListNode{Next: head}
	p2 := dummy2

	for p2.Next != nil {
		if p2.Next.Val < x {
			p1.Next = p2.Next
			p1 = p1.Next

			//删除 p2
			p2.Next = p2.Next.Next
		} else {
			//移动 p2
			p2 = p2.Next
		}
	}

	p1.Next = dummy2.Next
	return dummy1.Next
}

func mergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	return merge(lists, 0, len(lists)-1)
}

func merge(lists []*ListNode, lo, hi int) *ListNode {
	if lo == hi {
		return lists[lo]
	}
	mid := lo + (hi-lo)/2
	l1 := merge(lists, lo, mid)
	l2 := merge(lists, mid+1, hi)
	return mergeTwoLists(l1, l2)
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}
	p1 := dummy
	for i := 0; i < n+1; i++ {
		if p1 != nil {
			p1 = p1.Next
		}
	}

	p2 := dummy
	for p1 != nil {
		p1 = p1.Next
		p2 = p2.Next
	}

	p2.Next = p2.Next.Next
	return dummy.Next
}

func middleNode(head *ListNode) *ListNode {
	fast := head
	slow := head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}
	return slow
}

func detectCycle(head *ListNode) *ListNode {
	// 1、快慢指针同时出发，相遇时停止
	fast, slow := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if fast == slow {
			break
		}
	}

	if fast == nil || fast.Next == nil {
		return nil
	}

	// 2、快慢指针相同步数出发，相遇时停止
	slow = head
	for slow != fast {
		slow = slow.Next
		fast = fast.Next
	}
	return slow
}

func getIntersectionNode(headA, headB *ListNode) *ListNode {
	// 1,2,3,4,5,6,7,8,9,10|11,5,6,7,8,9,10
	// 11,5,6,7,8,9,10|1,2,3,4,5,6,7,8,9,10
	p1 := headA
	p2 := headB

	for p1 != p2 {
		if p1 != nil {
			p1 = p1.Next
		} else {
			p1 = headB
		}
		if p2 != nil {
			p2 = p2.Next
		} else {
			p2 = headA
		}
	}
	return p1
}

func hasCycle(head *ListNode) bool {
	fast, slow := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if fast == slow {
			return true
		}
	}
	return false
}

func trainingPlan(head *ListNode, cnt int) *ListNode {
	p1 := head
	for i := 0; i < cnt; i++ {
		if p1 != nil {
			p1 = p1.Next
		}
	}
	p2 := head
	for p1 != nil {
		p1 = p1.Next
		p2 = p2.Next
	}
	return p2
}

func deleteDuplicates(head *ListNode) *ListNode {
	dummy := &ListNode{Next: head}
	p := dummy
	for p.Next != nil && p.Next.Next != nil {
		if p.Next.Val == p.Next.Next.Val {
			rmVal := p.Next.Val
			for p.Next != nil && p.Next.Val == rmVal {
				p.Next = p.Next.Next
			}
		} else {
			p = p.Next
		}
	}
	return dummy.Next
}

// 思路：把矩阵的每一行看作一个已排序的链表，然后做 K 路归并。每一行只需要关心自己的"下一个元素"（即 j+1），不需要关心其他行。
func kthSmallest(matrix [][]int, k int) int {
	// 实现最小堆
	heap := &minHeap{}

	// 初始化：把每行的第一个元素（第0列）都推入堆
	for i := 0; i < len(matrix); i++ {
		heap.push([]int{matrix[i][0], i, 0})
	}

	// for k-- 直到 k 为 0 即为答案
	n := len(matrix)
	res := -1
	for k > 0 {
		cur := heap.pop()
		res = cur[0]
		i, j := cur[1], cur[2]

		if j+1 < n {
			heap.push([]int{matrix[i][j+1], i, j + 1})
		}
		k--
	}
	return res
}

type minHeap struct {
	// array[i]代表一个节点
	// array[i][0]=val
	// array[i][1]=行index
	// array[i][2]=列index
	array [][]int
}

func (h *minHeap) push(x []int) {
	h.array = append(h.array, x)
	h.siftUp(len(h.array) - 1)
}

func (h *minHeap) pop() []int {
	if len(h.array) == 0 {
		return []int{-1, -1, -1}
	}
	minVal := h.array[0]
	last := len(h.array) - 1

	h.array[0], h.array[last] = h.array[last], h.array[0]
	h.array = h.array[:last]
	if len(h.array) > 0 {
		h.siftDown(len(h.array), 0)
	}
	return minVal
}

// 下沉
func (h *minHeap) siftDown(n, i int) {
	minIndex := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && h.array[left][0] < h.array[minIndex][0] {
		minIndex = left
	}
	if right < n && h.array[right][0] < h.array[minIndex][0] {
		minIndex = right
	}
	if minIndex != i {
		h.array[minIndex], h.array[i] = h.array[i], h.array[minIndex]
		h.siftDown(n, minIndex)
	}
}

// 上浮
func (h *minHeap) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.array[parent][0] <= h.array[i][0] {
			break
		}
		h.array[parent], h.array[i] = h.array[i], h.array[parent]
		i = parent
	}
}

func kSmallestPairs(nums1 []int, nums2 []int, k int) [][]int {
	//nums1 = [1,7,11], nums2 = [2,4,6], k = 3
	//[1, 2] -> [1, 4] -> [1, 6]
	//[7, 2] -> [7, 4] -> [7, 6]
	//[11, 2] -> [11, 4] -> [11, 6]
	var res [][]int

	return res
}

func main() {
	fmt.Println(kSmallestPairs([]int{1, 7, 11}, []int{2, 4, 6}, 3))

	//l1 := NewMyListNode([]int{1, 2, 3, 4, 5})
	//l0 := middleNode(l1)
	//l0.Display()

	//martix := [][]int{{1, 5, 9}, {10, 11, 13}, {12, 13, 15}}
	//fmt.Println(kthSmallest(martix, 8))
	//fmt.Println(kthSmallest([][]int{{-5}}, 1))
}
