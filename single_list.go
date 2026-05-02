package main

func insertHeadNode(head *ListNode, val int) *ListNode {
	return &ListNode{Val: val, Next: head}
}

func insertTailNode(head *ListNode, val int) *ListNode {
	newNode := &ListNode{Val: val}
	if head == nil {
		return newNode
	}
	p := head
	for p.Next != nil {
		p = p.Next
	}
	p.Next = newNode
	return head
}

func insertIndexNode(head *ListNode, index, val int) *ListNode {
	// 边界判断
	if index < 0 {
		return head
	}
	if index == 0 {
		return &ListNode{Val: val, Next: head}
	}

	// 检查索引是否越界
	p := head
	for i := 0; i < index-1; i++ {
		if p == nil {
			break
		}
		p = p.Next
	}

	if p == nil {
		return head // 索引越界
	}

	tail := p.Next
	p.Next = &ListNode{Val: val, Next: tail}
	return head
}

func removeIndexNode(head *ListNode, index int) *ListNode {
	// 边界判断
	if head == nil || index < 0 {
		return head
	}
	if index == 0 {
		return head.Next
	}

	// 检查索引是否越界
	p := head
	for i := 0; i < index-1; i++ {
		if p == nil || p.Next == nil {
			break
		}
		p = p.Next
	}

	if p == nil || p.Next == nil {
		return head // 索引越界
	}

	delNode := p.Next
	p.Next = delNode.Next
	return head
}

// 递归反转
func _reverse(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	newHead := _reverse(head.Next)
	// head->next1->*
	// newHead => next3->next2->next1
	// 把head放到next1的下一个节点：newHead => next3->next2->next1->head->*
	// 需要把head的next指针置为空
	head.Next.Next = head
	head.Next = nil
	return newHead
}

// 迭代反转
func reverseList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	// 由于单链表的结构，至少要用三个指针才能完成迭代反转
	// pre->cur->nxt
	// cur->pre
	// pre = cur
	// cur = nxt
	// nxt = nxt.Next
	var pre, cur, nxt *ListNode
	pre, cur, nxt = nil, head, head.Next
	for cur != nil {
		// 逐个结点反转(指针指向pre节点)
		cur.Next = pre
		// 更新指针位置
		pre = cur
		cur = nxt
		if nxt != nil {
			nxt = nxt.Next
		}
	}
	return pre
}

var _tail *ListNode

// 反转链表前n个节点 - 递归
func _reverseN(head *ListNode, n int) *ListNode {
	if n == 1 {
		_tail = head.Next
		return head
	}
	newHead := _reverseN(head.Next, n-1)
	// head->next1->next2->next3
	// newHead => next3->next2->next1
	// 把head放到next1的下一个节点：newHead => next3->next2->next1->head->*
	// 需要把head的next指针置为tail：newHead => next3->next2->next1->head->tail
	head.Next.Next = head
	head.Next = _tail
	return newHead
}

// 反转链表前n个节点 - 迭代
func _reverseN1(head *ListNode, n int) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	// head->node1->node2->node3->node4->node5->nil
	// 由于单链表的结构，至少要用三个指针才能完成迭代反转
	// pre->cur->nxt
	// cur->pre
	// pre = cur
	// cur = nxt
	// nxt = nxt.Next
	var pre, cur, nxt *ListNode
	pre, cur, nxt = nil, head, head.Next
	for n > 0 {
		// 逐个结点反转(指针指向pre节点)
		cur.Next = pre
		// 更新指针位置
		pre = cur
		cur = nxt
		if nxt != nil {
			nxt = nxt.Next
		}
		n--
	}
	// 此时的cur是第n+1个节点，head是反转后的尾结点
	head.Next = cur
	// 此时的 pre 是反转后的头结点
	return pre
}

func main1() {
	var nums = []int{1, 2, 3, 4, 5}

	list := NewMyListNode(nums)
	list.Display()

	list = insertHeadNode(list, 666)
	list = insertTailNode(list, 999)
	list = insertIndexNode(list, 2, 555)
	list = removeIndexNode(list, 2)
	list.Display()

	list = _reverse(list)
	list.Display()
}
