package main

import (
	"fmt"
)

// 单链表实现：
// 🌟技巧1：头插法技巧 - 头插法直接返回新节点作为新的head，避免单独处理空链表，时间复杂度O(1)
// 🌟技巧2：前驱节点操作技巧 - 插入/删除时先找前一个节点（index-1位置），通过p.Next操作目标节点，避免丢失链表连接
// 🌟技巧3：虚拟头节点(dummy)技巧 - NewMyListNode中使用dummy节点简化链表构建，避免对第一个节点做特殊判断，构建完成后返回dummy.Next
// 🌟技巧4：递归反转技巧 - reverse递归到链表末尾后，利用head.Next.Next=head反转指针方向，再head.Next=nil断开旧连接避免成环
// 🌟技巧5：尾插法遍历技巧 - insertTailNode遍历到最后一个节点时判断条件是p.Next!=nil（停在最后一个节点），而不是p!=nil（会越过最后一个节点）

// ⚠️易错点1：空链表处理 - insertTailNode/reverse等函数必须先判断head==nil，否则对nil指针操作会panic
// ⚠️易错点2：index==0特殊处理 - insertIndexNode/removeIndexNode在index==0时操作的是头节点，需要单独处理并返回新的head，不能走通用的前驱节点逻辑
// ⚠️易错点3：前驱节点循环次数 - 找第index位置的前驱节点循环index-1次（而非index次），差一会导致操作位置偏移一个节点
// ⚠️易错点4：递归反转断链 - reverse中head.Next.Next=head后必须head.Next=nil断开原来的指向，否则会形成环导致无限循环
// ⚠️易错点5：索引越界检查 - 遍历过程中需要检查p是否为nil（insertIndexNode）或p.Next是否为nil（removeIndexNode），否则越界访问会panic
//
// 何时运用单链表：
// ❓1、是否需要频繁在头部插入/删除？链表头部操作O(1)，优于数组的O(n)
// ❓2、是否不需要随机访问？链表只能顺序遍历，随机访问需O(n)，不适合频繁按索引查找的场景
// ❓3、是否需要动态大小？链表无需预分配空间，适合元素数量不确定的场景

// 0、func NewMyListNode(nums []int) *ListNode
// 1、func insertHeadNode(head *ListNode, val int) *ListNode
// 2、func insertTailNode(head *ListNode, val int) *ListNode
// 3、func insertIndexNode(head *ListNode, index, val int) *ListNode
// 4、func removeIndexNode(head *ListNode, index int) *ListNode
// 5、func reverse(head *ListNode) *ListNode
// 6、func (h *ListNode) display()

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

func reverse(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	newHead := reverse(head.Next)
	head.Next.Next = head
	head.Next = nil
	return newHead
}

func (h *ListNode) display() {
	p := h
	str := "打印单链表: "
	for p != nil {
		str += fmt.Sprintf("%d->", p.Val)
		p = p.Next
	}
	fmt.Println(str + "nil")
}

func main() {
	var nums = []int{1, 2, 3, 4, 5}

	list := NewMyListNode(nums)
	list.display()

	list = insertHeadNode(list, 666)
	list = insertTailNode(list, 999)
	list = insertIndexNode(list, 2, 555)
	list = removeIndexNode(list, 2)
	list.display()

	list = reverse(list)
	list.display()
}
