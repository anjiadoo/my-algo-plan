/*
 * ============================================================================
 *                   📘 链表双指针算法全集 · 核心记忆框架
 * ============================================================================
 * 【一句话理解链表双指针】
 *
 *   链表双指针 = 利用两个节点指针的「速度差」或「起点差」来推导位置关系，
 *   在无法随机访问的链表结构中实现 O(n) 时间的定位、检测与变换。
 *
 *   判断能否用链表双指针，回答两个问题：
 *     ❓ 两个指针的移动速度分别是多少？（快慢指针）
 *     ❓ 两个指针的起始位置差了多少？（间隔指针）
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【四种链表双指针模式】
 *
 *   模式           速度/间距          代表题目
 *   ──────────────────────────────────────────────────────────────────────
 *   快慢指针       slow=1步,fast=2步  环检测、找中点、找环入口
 *   间隔指针       先走n步再同步      倒数第K个节点、删除倒数第N个
 *   拉链合并       p1,p2各走一步      合并有序链表、K路归并
 *   虚拟头指针     dummy节点辅助      分隔链表、删除节点、反转子链
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【快慢指针 —— 环检测与找中点】
 *
 *   核心原理：slow 每次走 1 步，fast 每次走 2 步。
 *
 *   ┌─────────────────────────────────────────────────────────────────────┐
 *   │ 应用1：判断链表是否有环（hasCycle）                                  │
 *   │   若有环 → fast 必追上 slow（相遇）                                 │
 *   │   若无环 → fast 先到 nil                                            │
 *   │                                                                     │
 *   │ 应用2：找环入口（detectCycle）                                       │
 *   │   第一步：快慢指针相遇                                              │
 *   │   第二步：slow 回到 head，两指针同速前进，再次相遇即为环入口             │
 *   │   数学证明：设头到环入口为 a，入口到相遇点为 b，环长为 c                 │
 *   │            slow 走了 a+b，fast 走了 a+b+nc                          │
 *   │            2(a+b) = a+b+nc → a = nc-b = (n-1)c + (c-b)            │
 *   │            所以从 head 和相遇点同步走 a 步就会在入口相遇                │
 *   │                                                                     │
 *   │ 应用3：找链表中点（middleNode）                                      │
 *   │   fast 到末尾时，slow 恰好在中间                                     │
 *   │   奇数长度：slow 在正中间                                            │
 *   │   偶数长度：slow 在中间偏右（第二个中间节点）                           │
 *   └─────────────────────────────────────────────────────────────────────┘
 *
 *   ⚠️ 易错点1：环检测的循环条件是 fast != nil && fast.Next != nil
 *      只判断 fast != nil 不够，fast.Next 为 nil 时 fast.Next.Next 会空指针。
 *      必须同时保证 fast 和 fast.Next 都非空。
 *
 *   ⚠️ 易错点2：找环入口第二步中，两指针必须同速（都走1步）
 *      如果 fast 仍走2步，数学关系不成立，不会在入口相遇。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【间隔指针 —— 倒数第K个节点】
 *
 *   核心思想：p1 先走 k 步，然后 p1、p2 同步前进。
 *   当 p1 到达末尾（nil）时，p2 恰好在倒数第 k 个节点。
 *
 *   原理：p1 和 p2 始终保持 k 步距离，p1 走了 n 步到达 nil，
 *         p2 走了 n-k 步，即正数第 n-k+1 个节点 = 倒数第 k 个。
 *
 *   变体「删除倒数第N个」：
 *     让 p1 先走 n+1 步（从 dummy 出发），这样 p2 最终停在「待删节点的前驱」，
 *     执行 p2.Next = p2.Next.Next 即可删除。
 *
 *   ⚠️ 易错点3：删除倒数第N个要从 dummy 节点出发
 *      从 head 出发的话，如果要删除的是 head 本身（倒数第n个 = 第1个），
 *      p2 无法找到前驱节点。用 dummy → head 可以统一处理。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【拉链合并 —— 有序链表归并】
 *
 *   核心思想：p1、p2 分别指向两个链表，每次取较小的接到结果链表尾部。
 *
 *   两路归并（mergeTwoLists）：
 *     while p1 && p2 → 取小的接到 p.Next，循环结束后接上剩余部分。
 *
 *   K路归并（mergeKLists）两种实现：
 *     ① 最小堆：将 K 个头节点入堆，每次弹出最小的，把它的 Next 入堆
 *     ② 分治：递归地两两合并，T(K) = 2T(K/2) + O(nK) → O(nK·logK)
 *
 *   ⚠️ 易错点4：堆中弹出节点后，必须将其 Next（而非节点本身）入堆
 *      每个链表的节点是按顺序入堆的。弹出一个节点后，该链表的「下一个候选」
 *      是 node.Next，如果 Next 为 nil 则该链表已耗尽，不再入堆。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【虚拟头节点（dummy node）的使用时机】
 *
 *   ✅ 何时必须用 dummy：
 *     1. 头节点可能被删除（删除操作类题目）
 *     2. 需要在头部之前插入节点（分隔链表、合并链表）
 *     3. 结果链表从空开始构建
 *
 *   好处：统一「对头节点的操作」和「对中间节点的操作」，避免 if head == nil 特判。
 *   最终返回 dummy.Next 即为真正的头节点。
 *
 *   ⚠️ 易错点5：分隔链表时必须断开原始 Next 指针
 *      将节点分配到 p1 或 p2 后，若不断开原始 Next，会在两个子链中形成环。
 *      正确做法：先保存 next := p.Next，再 p.Next = nil，再 p = next。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【链表反转 —— 核心操作】
 *
 *   递归反转（reverse）：
 *     base case：head==nil || head.Next==nil → 返回 head（最后一个节点=新头）
 *     递归后：head.Next.Next = head（让下一个节点指回自己）
 *              head.Next = nil（断尾，防止环）
 *
 *   迭代反转（reverse1）：
 *     三指针 pre=nil, cur=head, nxt=head.Next
 *     每步：cur.Next = pre → pre=cur → cur=nxt → nxt=nxt.Next
 *     结束后 pre 是新头节点
 *
 *   反转前N个（reverseN）：
 *     与全部反转的区别：到第N个节点时记录 tail = head.Next（第N+1个节点），
 *     递归回溯时 head.Next = tail（接尾）而非 nil（断尾）。
 *
 *   反转区间 [m, n]（reverseBetween）：
 *     找到第 m-1 个节点作为前驱，对其后面的子链调用 reverseN(pre.Next, n-m+1)。
 *
 *   K个一组反转（reverseKGroup）：
 *     先数够 K 个节点 [a, b)，对这 K 个调用 reverseN(a, k)，
 *     然后递归处理 b 开头的剩余链表，将 a.Next（反转后的尾）指向递归结果。
 *
 *   ⚠️ 易错点6：递归反转必须 head.Next = nil（断尾）
 *      不断尾的话，原来的 head 和 head.Next 会互相指向，形成环。
 *
 *   ⚠️ 易错点7：reverseN 中 tail 变量的作用域
 *      tail 记录的是「第 N+1 个节点」，在所有递归层之间共享。
 *      如果用局部变量，递归返回后 tail 值会丢失。通常用全局变量或闭包。
 *
 *   ⚠️ 易错点8：reverseKGroup 不足 K 个时不反转
 *      如果剩余节点不足 K 个，直接返回 head，保持原序。
 *      先用循环数 K 个，中途遇 nil 则返回，这是 base case。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【相交链表（getIntersectionNode）】
 *
 *   核心思想：两个指针分别走 A+B 和 B+A 的路径，最终对齐。
 *
 *   设链表 A 长度为 a+c，链表 B 长度为 b+c（c 为公共部分长度）。
 *   p1 走完 A 后转到 B 头部，走 a+c+b 步到达交点。
 *   p2 走完 B 后转到 A 头部，走 b+c+a 步到达交点。
 *   a+c+b == b+c+a，所以两者同时到达交点（或同时到达 nil，即无交点）。
 *
 *   ⚠️ 易错点9：循环中 p1 到 nil 时应指向 headB（不是 headA）
 *      容易写反。p1 走的是「A链 → B链」的路径，所以 A 走完后接 B。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【Floyd 判圈算法在数组中的应用（findDuplicate）】
 *
 *   将数组 nums 视为隐式链表：索引 i → nums[i] 构成一条链。
 *   如果存在重复数字，就存在「多个索引指向同一个值」→ 即链表中存在环。
 *   环的入口就是重复的数字。
 *
 *   做法与链表找环入口完全相同：
 *     第一步：slow = nums[slow], fast = nums[nums[fast]]，直到相遇
 *     第二步：slow 回到起点 0，两者同速，再次相遇即为答案
 *
 *   ⚠️ 易错点10：数组版的「走一步」是 slow = nums[slow]，不是 slow++
 *      这里不是线性遍历，而是顺着「值→索引」的映射走，模拟链表的 Next 指针。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【回文链表判断（isPalindrome）】
 *
 *   方法一（递归 —— 本文实现）：
 *     利用递归的「后序遍历」特性，递归到底后从尾部回溯，
 *     同时用一个外部指针 left 从头部前进，实现「双向比较」。
 *
 *   方法二（经典做法）：
 *     1. 快慢指针找中点
 *     2. 反转后半部分
 *     3. 从头和中点同时比较
 *
 *   ⚠️ 易错点11：递归法中 left 必须是外部变量（闭包捕获）
 *      如果 left 作为参数传递，递归回溯时 left 的值不会随之更新，
 *      因为 Go 中指针变量传递的是指针的副本，需要用闭包或二级指针。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次手写链表双指针后对照检查
 *
 *     ✅ 环检测：循环条件是 fast != nil && fast.Next != nil？
 *     ✅ 找环入口：第二步中两个指针都是走 1 步？
 *     ✅ 间隔指针：先走的步数是 k 还是 k+1（取决于是否用 dummy）？
 *     ✅ 虚拟头节点：是否需要 dummy？最终返回的是 dummy.Next？
 *     ✅ 链表分隔：分配节点后是否断开了原始 Next？
 *     ✅ 归并：弹出节点后入堆的是 node.Next 而非 node？
 *     ✅ 反转：递归反转是否执行了 head.Next = nil（全部反转）或 head.Next = tail（前N个）？
 *     ✅ K组反转：不足 K 个是否直接返回不反转？
 *     ✅ 相交链表：p1 走完 A 后接的是 headB，p2 走完 B 后接的是 headA？
 *     ✅ 两数相加：循环条件是否包含 flag > 0（处理最后的进位）？
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *
 *     合并与分隔：
 *       1. mergeTwoLists(list1, list2 *ListNode) *ListNode       // 合并两个有序链表
 *       2. mergeKLists(lists []*ListNode) *ListNode              // K路归并（分治法）
 *       3. mergeKLists_(lists []*ListNode) *ListNode             // K路归并（最小堆）
 *       4. partitionList(head *ListNode, x int) *ListNode        // 分隔链表
 *
 *     快慢指针：
 *       5. hasCycle(head *ListNode) bool                         // 判断链表是否有环
 *       6. detectCycle(head *ListNode) *ListNode                 // 找环入口
 *       7. middleNode(head *ListNode) *ListNode                  // 链表中点
 *       8. getIntersectionNode(headA, headB *ListNode) *ListNode // 相交链表
 *
 *     间隔指针：
 *       9. removeNthFromEnd(head *ListNode, n int) *ListNode     // 删除倒数第N个
 *      10. trainingPlan(head *ListNode, cnt int) *ListNode       // 返回倒数第K个
 *
 *     链表反转：
 *      11. reverse(head *ListNode) *ListNode                     // 递归反转全部
 *      12. reverse1(head *ListNode) *ListNode                    // 迭代反转全部
 *      13. reverseN(head *ListNode, n int) *ListNode             // 反转前N个
 *      14. reverseBetween(head *ListNode, left, right int)       // 反转区间[left,right]
 *      15. reverseKGroup(head *ListNode, k int) *ListNode        // K个一组反转
 *      16. swapPairs(head *ListNode) *ListNode                   // 两两交换
 *
 *     其他经典：
 *      17. addTwoNumbers(l1, l2 *ListNode) *ListNode             // 两数相加
 *      18. isPalindrome(head *ListNode) bool                     // 回文链表
 *      19. deleteDuplicates(head *ListNode) *ListNode            // 删除重复元素II
 *      20. findDuplicate(nums []int) int                         // 寻找重复数（Floyd判圈）
 *      21. reorderList(head *ListNode)                           // 重排链表
 *      22. kthSmallest(matrix [][]int, k int) int                // 有序矩阵第K小（K路归并）
 *      23. kSmallestPairs(nums1, nums2 []int, k int) [][]int     // 最小K对数字
 * ============================================================================
 */

package main

import (
	"container/heap"
	"fmt"
)

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
func _partition(head *ListNode, x int) *ListNode {
	// 存放「小于x」的节点
	dummy1 := &ListNode{}
	p1 := dummy1

	// 存放「大于等于x」的节点
	dummy2 := &ListNode{Next: head}
	p2 := dummy2

	for p2.Next != nil {
		if p2.Next.Val < x {
			p1.Next = p2.Next
			p1 = p1.Next
			p2.Next = p2.Next.Next // 把p2中小于x的删除
		} else {
			p2 = p2.Next // p2中的节点>=x，保留并移动p2
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
	return _merge(lists, 0, len(lists)-1)
}

func _merge(lists []*ListNode, lo, hi int) *ListNode {
	if lo == hi {
		return lists[lo]
	}
	mid := lo + (hi-lo)/2
	l1 := _merge(lists, lo, mid)
	l2 := _merge(lists, mid+1, hi)
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
	heap := &_MinHeap{}

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

type _MinHeap struct {
	// array[i][0] minSum
	// array[i][1] i
	// array[i][2] j
	array [][]int
}

func (h *_MinHeap) Push(x []int) {
	h.array = append(h.array, x)
	h.siftUp(len(h.array) - 1)
}

func (h *_MinHeap) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.array[parent][0] <= h.array[i][0] {
			break
		}
		h.array[parent], h.array[i] = h.array[i], h.array[parent]
		i = parent
	}
}

func (h *_MinHeap) Pop() []int {
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

func (h *_MinHeap) siftDown(n, i int) {
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
	// head == nil 空链表的情况
	// head.Next == nil 非空链表最后一个节点就是【新头节点】
	if head == nil || head.Next == nil {
		return head
	}

	// head->next1->next2->next3->nil
	// newHead => next3->next2->next1->nil
	// 让下一个节点(next1)指回我：newHead => next3->next2->next1->head->*
	// 我指向空（断尾）：newHead => next3->next2->next1->head->nil

	newHead := reverse(head.Next)
	head.Next.Next = head
	head.Next = nil
	return newHead
}

// 反转链表-迭代
func reverse1(head *ListNode) *ListNode {
	// head == nil 空链表的情况
	// head.Next == nil 只有一个节点的情况（不用反转了）

	if head == nil || head.Next == nil {
		return head
	}

	// head->node1->node2->node3->node4->node5->nil
	// 由于单链表的结构，至少要用三个指针才能完成迭代反转
	// pre(nil)，cur->nxt(初始状态)
	// cur->pre
	// pre = cur
	// cur = nxt
	// nxt = nxt.Next（注意判空）

	pre := (*ListNode)(nil)
	cur := head
	nxt := head.Next

	for cur != nil {
		// 当前结点反转
		cur.Next = pre

		// 整体向右平移
		pre = cur
		cur = nxt
		if nxt != nil { //（注意判空）
			nxt = nxt.Next
		}
	}
	// 此时的 pre 是反转后的头结点
	return pre
}

var tail *ListNode

// 反转链表前n个节点
func reverseN(head *ListNode, n int) *ListNode {
	// 遍历到第n个节点时，保留后面的节点顺序
	if n == 1 {
		tail = head.Next
		return head
	}

	// head->next1->next2->next3->nil
	// newHead => next2->next1->next3->nil
	// 让下一个节点(next1)指回我：newHead => next2->next1->head->*
	// 我指向尾（接尾）：newHead => next3->next2->next1->head->tail

	newHead := reverseN(head.Next, n-1)
	head.Next.Next = head
	head.Next = tail
	return newHead
}

// 反转链表前n个节点
func reverseN1(head *ListNode, n int) *ListNode {
	// head == nil 空链表的情况
	// head.Next == nil 只有一个节点的情况（不用反转了）
	if head == nil || head.Next == nil {
		return head
	}

	// head->node1->node2->node3->node4->node5->nil
	// 由于单链表的结构，至少要用三个指针才能完成迭代反转
	// pre(nil)，cur->nxt(初始状态)
	// cur->pre
	// pre = cur
	// cur = nxt
	// nxt = nxt.Next

	pre := (*ListNode)(nil)
	cur := head
	nxt := head.Next

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
	next := reverseBetween2(head.Next, m-1, n-1)
	head.Next = next
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

// 重排链表 https://leetcode.cn/problems/reorder-list/description/
func reorderList(head *ListNode) {
	var stack []*ListNode
	for p := head; p != nil; p = p.Next {
		stack = append(stack, p)
	}

	p := head
	for p != nil {
		tail := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// 结束条件
		if p == tail || p.Next == tail {
			tail.Next = nil
			break
		}

		next := p.Next
		p.Next = tail
		tail.Next = next
		p = next
	}
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

	//array1 := []int{1, 2, 3, 4, 9, 10}
	//array2 := []int{5, 6, 7, 8, 9, 10}
	//l1 := NewMyListNode(array1)
	//l2 := NewMyListNode(array2)
	//list := getIntersectionNode(l1, l2)

	array := []int{1, 2, 2, 2, 3, 3, 3, 4, 5, 6, 6, 6}
	head := NewMyListNode(array)
	list := deleteDuplicates(head)
	list.Display()
}
