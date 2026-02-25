package main

import (
	"container/heap"
)

// 双指针链表算法实现：
// 🌟技巧1：虚拟头结点技巧 - 当需要创造一条新链表的时候，可以使用「虚拟头结点」简化边界情况的处理
// 🌟技巧2：双指针合并技巧 - 合并链表时使用双指针分别遍历，较小的先接上，最后处理剩余部分（合并两个有序链表）
// 🌟技巧3：优先级队列技巧 - 合并K个链表时，使用优先级队列（最小堆）每次取最小节点，空间O(k)时间O(nlogk)（合并K个有序链表）
// 🌟技巧4：链表分解技巧 - 链表分解时创建两个虚拟头结点，分别接小于x和大于等于x的节点，最后链接（单链表的分解）
// 🌟技巧5：快慢指针删除技巧 - 删除倒数第N个节点时，快指针先走N+1步，然后快慢指针同步前进，慢指针指向要删除节点的前一个（删除链表的倒数第N个结点）
// 🌟技巧6：快慢指针中间节点技巧 - 快慢指针找中间节点，快指针走两步慢指针走一步，快指针到末尾时慢指针在中间（链表的中间结点）
// 🌟技巧7：环形检测技巧 - 环形链表检测，快慢指针相遇说明有环，快指针到末尾说明无环（环形链表I）
// 🌟技巧8：环起点定位技巧 - 找环起点时，快慢指针相遇后，慢指针从头开始，快慢指针同步前进，再次相遇点就是环起点（环形链表II）
// 🌟技巧9：相交链表检测技巧 - 相交链表检测，两个指针分别遍历两个链表，到达末尾时切换到另一个链表头，相遇点就是交点（相交链表）

// 0、func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode
// 1、func mergeKLists(lists []*ListNode) *ListNode
// 2、func mergeKLists_(lists []*ListNode) *ListNode
// 3、func partitionList(head *ListNode, x int) *ListNode
// 4、func removeNthFromEnd(head *ListNode, n int) *ListNode
// 5、func middleNode(head *ListNode) *ListNode
// 6、func hasCycle(head *ListNode) bool
// 7、func detectCycle(head *ListNode) *ListNode
// 8、func getIntersectionNode(headA, headB *ListNode) *ListNode

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

// 合并 K 个升序链表
func mergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	l0 := lists[0]
	for i := 1; i < len(lists); i++ {
		l0 = mergeTwoLists(l0, lists[i])
	}
	return l0
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
func mergeKLists_(lists []*ListNode) *ListNode {
	// 1、实现一个优先级队列(最小堆)
	// 2、把K个升序链表头节点入队
	// 3、每次取最小节点放入结果链表

	if len(lists) == 0 {
		return nil
	}

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

// 单链表的分解
func partitionList(head *ListNode, x int) *ListNode {
	// 遍历链表，把小的放p1，大的放p2，最后链接p1,p2

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

// 删除链表的倒数第 N 个结点
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	// 双指针技巧：两个指针p1，p2，p1先走 n+1 步

	// 虚拟头节可点避免删除第一个节点的特色情况
	dummy := &ListNode{Next: head}
	p1 := dummy

	for i := 0; i < n+1; i++ {
		if p1 != nil {
			p1 = p1.Next
		} else {
			if i < n+1 {
				return head
			}
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

// 链表的中间结点
func middleNode(head *ListNode) *ListNode {
	// 快慢指针技巧：慢指针走一步，快指针就走两步

	slow, fast := head, head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}

// 环形链表 I 是否有环
func hasCycle(head *ListNode) bool {
	slow := head
	fast := head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			return true
		}
	}
	return false
}

// 环形链表 II 寻找环起点
func detectCycle(head *ListNode) *ListNode {
	slow, fast := head, head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if fast == slow {
			break
		}
	}

	// fast 遇到空指针说明没有环
	if fast == nil || fast.Next == nil {
		return nil
	}

	// 重新指向头结点
	slow = head

	// 快慢指针同步前进，相交点就是环起点
	for fast != slow {
		fast = fast.Next
		slow = slow.Next
	}
	return slow
}

// 相交链表，把两个链表拼接在一起遍历
func getIntersectionNode(headA, headB *ListNode) *ListNode {
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

func main() {
	//array1 := []int{-9, 3}
	//array2 := []int{5, 7}
	//l1 := NewMyListNode(array1)
	//l2 := NewMyListNode(array2)
	//list := mergeTwoLists(l1, l2)

	//array := []int{1, 4, 3, 2, 5, 2}
	//head := NewMyListNode(array)
	//list := partitionList(head, 3)

	//array1 := []int{1, 4, 5}
	//array2 := []int{1, 3, 4}
	//array3 := []int{2, 6}
	//l1 := NewMyListNode(array1)
	//l2 := NewMyListNode(array2)
	//l3 := NewMyListNode(array3)
	//list := mergeKLists([]*ListNode{l1, l2, l3})

	//array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//head := NewMyListNode(array)
	//list := removeNthFromEnd(head, 10)

	//array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//head := NewMyListNode(array)
	//list := middleNode(head)

	array1 := []int{1, 2, 3, 4, 9, 10}
	array2 := []int{5, 6, 7, 8, 9, 10}
	l1 := NewMyListNode(array1)
	l2 := NewMyListNode(array2)
	list := getIntersectionNode(l1, l2)

	printListNode(list)
}
