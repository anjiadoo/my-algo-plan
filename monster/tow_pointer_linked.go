/*
 * ============================================================================
 *                    📘 链表指针操作 · 核心记忆框架
 * ============================================================================
 * 【三种基本操作】—— 所有链表题都是这三种的组合
 *
 *     1. 断链：node.Next = nil          // 切断连接
 *     2. 接链：node.Next = target       // 建立连接
 *     3. 移动：node = node.Next         // 移动指针
 *
 *     🔑 关键原则：先保存，再断链，最后接链
 *        一旦覆盖了 node.Next，原来的下一个节点就丢了！
 * ────────────────────────────────────────────────────────────────────────────
 * 【五大模式与口诀】
 *
 *   ① 虚拟头结点（Dummy Node）
 *      口诀：凡是头结点可能变的，都加 dummy
 *      模板：dummy := &ListNode{Next: head}
 *            ... 操作 ...
 *            return dummy.Next
 *      适用：合并、分隔、删除
 *
 *   ② 快慢指针
 *      口诀：找中点用 2倍速，找倒数第k用 k步差
 *        · 找中点 ：fast 每次走2步，slow走1步，fast到底 slow就在中间
 *        · 找倒数k：fast 先走k步，然后一起走，fast到底 slow就在倒数第k
 *        · 找环入口：快慢相遇后，slow回起点，同速再走，再次相遇就是环入口
 *
 *   ③ 反转链表
 *      递归口诀：后面已经反转好了，我只管把自己接上去
 *        head.Next.Next = head   // 让下一个节点指回我
 *        head.Next = nil         // 我指向空（断尾）
 *
 *      迭代口诀：三指针，逐个翻，pre-cur-nxt 往右搬
 *        步骤图：
 *            pre    cur    nxt
 *             ↓      ↓      ↓
 *            nil ← [ 1 ] → [ 2 ] → [ 3 ] → nil
 *                   ① cur.Next = pre    // 反指
 *                   ② pre = cur         // 右移
 *                   ③ cur = nxt         // 右移
 *                   ④ nxt = nxt.Next    // 右移
 *
 *   ④ 穿针引线（拼接 / 拆分）
 *      口诀：拆成多条链，最后缝起来
 *      典型：partition —— 拆成 "小于x" 和 "大于等于x" 两条链，最后拼接
 *
 *   ⑤ 递归思维
 *      口诀：相信子问题已解决，我只处理当前节点
 *      典型：reverseBetween2 —— 相信 head.Next 后面已经反转好了
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次写完链表题后对照检查
 *
 *     ✅ 是否处理了空链表？           → head == nil
 *     ✅ 是否处理了单节点？           → head.Next == nil
 *     ✅ 是否有指针丢失？             → 覆盖 Next 前没保存
 *     ✅ 尾结点的 Next 是否为 nil？   → 反转后忘记断尾
 *     ✅ 返回值是 dummy.Next 还是 head？ → 头结点被移动时应返回 dummy.Next
 * ────────────────────────────────────────────────────────────────────────────
 * 【万能画图法】—— 与其记代码，不如记画图的习惯
 *
 *     操作前:  pre → cur → nxt → ...
 *     操作后:  pre → nxt → ...         (删除 cur)
 *              pre ← cur    nxt → ...  (反转 cur)
 *
 *     推荐符号：→ 表示 Next 指针  ✕ 表示断开  ①②③ 表示操作顺序
 * ────────────────────────────────────────────────────────────────────────────
 *   通用防错清单：
 *      □ 操作 p.Next 前，先确认 p != nil
 *      □ 需要删除/跳过节点时，用前驱节点操作：pre.Next = pre.Next.Next
 *      □ 断开节点时记得把 node.Next = nil，防止出现环
 *      □ 构建新链表时，用 dummy 开头，返回 dummy.Next
 *      □ 反转后原来的 head 变成了 tail，别忘了处理 tail.Next 的指向
 * ============================================================================
 */

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

// 合并两个有序链表 https://leetcode.cn/problems/merge-two-sorted-lists/description/
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

// 分隔链表 https://leetcode.cn/problems/partition-list/description/
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

// 合并 K 个升序链表 https://leetcode.cn/problems/merge-k-sorted-lists/description/
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

// 删除链表的倒数第 N 个结点 https://leetcode.cn/problems/remove-nth-node-from-end-of-list/description/
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

// 链表的中间结点 https://leetcode.cn/problems/middle-of-the-linked-list/description/
func middleNode(head *ListNode) *ListNode {
	fast := head
	slow := head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}
	return slow
}

// 环形链表 II https://leetcode.cn/problems/linked-list-cycle-ii/description/
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

// 相交链表 https://leetcode.cn/problems/intersection-of-two-linked-lists/description/
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

// 环形链表 https://leetcode.cn/problems/linked-list-cycle/description/
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

// 返回倒数第cnt个节点 https://leetcode.cn/problems/lian-biao-zhong-dao-shu-di-kge-jie-dian-lcof/
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

// 删除排序链表中的重复元素II https://leetcode.cn/problems/remove-duplicates-from-sorted-list-ii/description/
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

// 有序矩阵中第K小的元素 https://leetcode.cn/problems/kth-smallest-element-in-a-sorted-matrix/
func kthSmallest(matrix [][]int, k int) int {
	// 思路：把矩阵的每一行看作一个已排序的链表，然后做 K 路归并。每一行只需要关心自己的"下一个元素"（即 j+1），不需要关心其他行。
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

// 查找和最小的K对数字 https://leetcode.cn/problems/find-k-pairs-with-smallest-sums/description/
func kSmallestPairs(nums1 []int, nums2 []int, k int) [][]int {
	heap := &MinHeap{}

	for i := 0; i < len(nums1); i++ {
		heap.Push([]int{nums1[i] + nums2[0], i, 0})
	}

	var res [][]int

	for k > 0 {
		x := heap.Pop()
		i := x[1]
		j := x[2]
		res = append(res, []int{nums1[i], nums2[j]})
		if j+1 < len(nums2) {
			heap.Push([]int{nums1[i] + nums2[j+1], i, j + 1})
		}
		k--
	}

	return res
}

type MinHeap struct {
	// array[i][0] minSum
	// array[i][1] i
	// array[i][2] j
	array [][]int
}

func (h *MinHeap) Push(x []int) {
	h.array = append(h.array, x)
	h.siftUp(len(h.array) - 1)
}

func (h *MinHeap) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.array[parent][0] <= h.array[i][0] {
			break
		}
		h.array[parent], h.array[i] = h.array[i], h.array[parent]
		i = parent
	}
}

func (h *MinHeap) Pop() []int {
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

func (h *MinHeap) siftDown(n, i int) {
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
		h.array[i], h.array[minIndex] = h.array[minIndex], h.array[i]
		h.siftDown(n, minIndex)
	}
}

// 两数相加 https://leetcode.cn/problems/add-two-numbers/
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	p := dummy
	p1 := l1
	p2 := l2
	flag := 0
	// 两条链表走完且没有进位时才能结束循环
	for p1 != nil || p2 != nil || flag > 0 {
		// 先加上次的进位
		val := flag
		if p1 != nil {
			val += p1.Val
			p1 = p1.Next
		}
		if p2 != nil {
			val += p2.Val
			p2 = p2.Next
		}
		flag = val / 10
		val = val % 10
		p.Next = &ListNode{Val: val}
		p = p.Next
	}
	return dummy.Next
}

// 两数相加II https://leetcode.cn/problems/add-two-numbers-ii/
func addTwoNumbers2(l1 *ListNode, l2 *ListNode) *ListNode {
	l1 = reverse(l1)
	l2 = reverse(l2)
	return reverse(addTwoNumbers(l1, l2))
}

// 两数相加II https://leetcode.cn/problems/add-two-numbers-ii/
func addTwoNumbers21(l1 *ListNode, l2 *ListNode) *ListNode {
	var traverse func(head *ListNode) []int

	traverse = func(head *ListNode) []int {
		if head == nil {
			return []int{}
		}
		res := traverse(head.Next)
		res = append(res, head.Val)
		return res
	}

	res1 := traverse(l1)
	res2 := traverse(l2)

	var res []int
	i1 := 0
	i2 := 0
	flag := 0

	for i1 < len(res1) || i2 < len(res2) || flag > 0 {
		val := flag
		if i1 < len(res1) {
			val += res1[i1]
			i1++
		}
		if i2 < len(res2) {
			val += res2[i2]
			i2++
		}
		flag = val / 10
		val = val % 10
		res = append(res, val)
	}
	dummy := &ListNode{}
	p := dummy
	for i := len(res) - 1; i >= 0; i-- {
		p.Next = &ListNode{Val: res[i]}
		p = p.Next
	}
	return dummy.Next
}

// 寻找重复数 https://leetcode.cn/problems/find-the-duplicate-number/description/
func findDuplicate(nums []int) int {
	// 以下是Floyd判圈算法，把数组当作隐式链表
	slow, fast := 0, 0
	for {
		slow = nums[slow]
		fast = nums[nums[fast]]
		if slow == fast {
			break
		}
	}

	// 重新指向头结点（索引 0）
	slow = 0

	// 快慢指针同步前进，相交点就是环入口
	for {
		slow = nums[slow]
		fast = nums[fast]
		if fast == slow {
			break
		}
	}
	return slow
}

// 反转链表-递归
func reverse(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	newHead := reverse(head.Next)
	// head->next1->next2->next3
	// newHead => next3->next2->next1
	// 把head放到next1的下一个节点：newHead => next3->next2->next1->head->*
	// 需要把head的next指针置为空：newHead => next3->next2->next1->head->nil
	head.Next.Next = head
	head.Next = nil
	return newHead
}

// 反转链表-迭代
func reverse1(head *ListNode) *ListNode {
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
		// 逐个结点反转
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

var tail *ListNode

// 反转链表前n个节点
func reverseN(head *ListNode, n int) *ListNode {
	if n == 1 {
		tail = head.Next
		return head
	}
	newHead := reverseN(head.Next, n-1)
	// head->next1->next2->next3
	// newHead => next3->next2->next1
	// 把head放到next1的下一个节点：newHead => next3->next2->next1->head->*
	// 需要把head的next指针置为tail：newHead => next3->next2->next1->head->tail
	head.Next.Next = head
	head.Next = tail
	return newHead
}

// 反转链表前n个节点
func reverseN1(head *ListNode, n int) *ListNode {
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

// 反转链表II-迭代解法 https://leetcode.cn/problems/reverse-linked-list-ii/description/
func reverseBetween(head *ListNode, left int, right int) *ListNode {
	if left == 1 {
		return reverseN(head, right)
	}

	// 寻找left的前一个节点
	pre := head
	for i := 1; i < left-1; i++ {
		pre = pre.Next
	}

	// 反转left后面的N个节点
	pre.Next = reverseN(pre.Next, right-left+1)
	return head
}

// 反转链表II-递归解法 https://leetcode.cn/problems/reverse-linked-list-ii/description/
func reverseBetween2(head *ListNode, m int, n int) *ListNode {
	if m == 1 {
		return reverseN(head, n)
	}
	tail := reverseBetween2(head.Next, m-1, n-1)
	head.Next = tail
	return head
}

// K个一组翻转链表 https://leetcode.cn/problems/reverse-nodes-in-k-group/
func reverseKGroup(head *ListNode, k int) *ListNode {
	if head == nil {
		return nil
	}
	// 1->2->3->4->5->6->
	// a     b
	// 区间[a, b)包含k个待反转元素
	a, b := head, head
	for i := 0; i < k; i++ {
		if b == nil {
			return head
		}
		b = b.Next
	}

	// 反转前k个元素
	newHead := reverseN(a, k)

	// 此时b指向下一组待反转的头结点
	// 递归反转后续链表并连接起来
	a.Next = reverseKGroup(b, k)
	return newHead
}

// 两两交换链表中的节点 https://leetcode.cn/problems/swap-nodes-in-pairs/
func swapPairs(head *ListNode) *ListNode {
	if head == nil {
		return head
	}
	// 1->2->3->4->5->6->
	// a     b
	// 区间[a, b)包含2个待反转元素
	a, b := head, head
	for i := 0; i < 2; i++ {
		if b == nil {
			return head
		}
		b = b.Next
	}

	// 反转前2个元素
	newHead := reverseN(a, 2)

	// 此时b指向下一组待反转的头结点
	// 递归反转后续链表并连接起来
	a.Next = swapPairs(b)

	return newHead
}

// 回文链表 https://leetcode.cn/problems/palindrome-linked-list/
func isPalindrome(head *ListNode) bool {
	var traverse func(head *ListNode)

	left := head
	res := true

	traverse = func(right *ListNode) {
		if right == nil {
			return
		}
		traverse(right.Next)
		if right.Val != left.Val {
			res = false
			return
		}
		left = left.Next
	}

	traverse(head)
	return res
}

func main() {
	l1 := NewMyListNode([]int{1, 2, 3, 4, 5, 6, 4, 3, 2, 1})
	fmt.Println(isPalindrome(l1))

	//l1 := NewMyListNode([]int{1, 2, 3, 4, 5, 6, 7, 8})
	//l0 := swapPairs(l1)
	//l0.Display()

	//fmt.Println(findDuplicate([]int{1, 3, 4, 2, 2}))
	//fmt.Println(findDuplicate([]int{3, 1, 3, 4, 2}))
	//fmt.Println(findDuplicate([]int{3, 3, 3, 3, 3}))
	//fmt.Println(findDuplicate([]int{7, 9, 7, 4, 2, 8, 7, 7, 1, 5}))

	//l1 := NewMyListNode([]int{7, 2, 4, 3})
	//l2 := NewMyListNode([]int{5, 6, 4})
	//l0 := addTwoNumbers21(l1, l2)
	//l0.Display()

	//fmt.Println(kSmallestPairs([]int{0, 0, 0}, []int{-3, 22, 35}, 9))

	//heap := &MinHeap{}
	//heap.Push([]int{1, 1, 1})
	//heap.Push([]int{10, 1, 1})
	//heap.Push([]int{19, 1, 1})
	//heap.Push([]int{4, 1, 1})
	//heap.Push([]int{2, 1, 1})
	//fmt.Println(heap.Pop())
	//fmt.Println(heap.Pop())
	//fmt.Println(heap.Pop())
	//fmt.Println(heap.Pop())
	//fmt.Println(heap.Pop())
	//fmt.Println(heap.Pop())

	//fmt.Println(kSmallestPairs([]int{1, 7, 11}, []int{2, 4, 6}, 3))

	//l1 := NewMyListNode([]int{1, 2, 3, 4, 5})
	//l0 := middleNode(l1)
	//l0.Display()

	//martix := [][]int{{1, 5, 9}, {10, 11, 13}, {12, 13, 15}}
	//fmt.Println(kthSmallest(martix, 8))
	//fmt.Println(kthSmallest([][]int{{-5}}, 1))
}
