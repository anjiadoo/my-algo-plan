package main

// 分治算法实现：
// 🌟技巧1：分治三步框架 - 分治的核心是「分解」→「递归求解」→「合并」：将大问题拆成规模更小的同类子问题，递归解决后再合并子问题的结果
// 🌟技巧2：二分拆分技巧 - 对数组/链表类分治，通常用 mid = start + (end-start)/2 将区间一分为二，分别递归处理左半部分和右半部分，避免整数溢出
// 🌟技巧3：分治 vs 递归 - 分治是一种特殊的递归：普通递归是「自顶向下」解决子问题，分治强调子问题之间相互独立且结构相同，合并步骤是分治的关键所在
// 🌟技巧4：时间复杂度分析 - 分治复杂度用主定理分析：T(n) = aT(n/b) + f(n)；合并K个链表每层合并操作为O(n)，树高为O(logK)，总复杂度O(n*logK)

// ⚠️易错点1：base case 必须完备 - start==end 时直接返回单个元素，start>end 时返回 nil；缺少任一 base case 都会导致无限递归或数组越界
// ⚠️易错点2：mid 计算方式 - 必须用 start+(end-start)/2 而非 (start+end)/2，后者在 start 和 end 都很大时会整数溢出
// ⚠️易错点3：左右子问题的边界划分 - 左半部分是 [start, mid]，右半部分是 [mid+1, end]，写成 [start, mid-1] 或 [mid, end] 都会导致边界错误或死递归
// ⚠️易错点4：合并逻辑要处理剩余部分 - 合并两个有序链表时，内层循环结束后必须将非空的剩余链表直接接上，不能只靠循环处理所有节点

// 何时运用分治算法：
// ❓1、问题能否被拆解为规模更小的「同类子问题」？即子问题结构与原问题相同
// ❓2、子问题之间是否相互独立？若子问题有重叠（如斐波那契），应优先考虑动态规划而非分治
// ❓3、子问题的解合并后能否得到原问题的解？合并步骤的正确性是分治算法的关键

// 0、func mergeKLists1(lists []*ListNode) *ListNode  // 合并K个有序链表 - 分治算法

// 合并k个有序链表-分治算法
func mergeKLists1(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}

	var mergeList func(lists []*ListNode, start, end int) *ListNode
	var merge func(l1 *ListNode, l2 *ListNode) *ListNode

	// 合并两个链表
	merge = func(l1 *ListNode, l2 *ListNode) *ListNode {
		dummy := &ListNode{}
		p := dummy
		p1, p2 := l1, l2

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
		if p1 != nil {
			p.Next = p1
		}
		if p2 != nil {
			p.Next = p2
		}
		return dummy.Next
	}

	// 递归合并K个链表
	mergeList = func(lists []*ListNode, start, end int) *ListNode {
		if start == end {
			return lists[start]
		}

		mid := start + (end-start)/2

		left := mergeList(lists, start, mid)
		right := mergeList(lists, mid+1, end)

		return merge(left, right)
	}

	return mergeList(lists, 0, len(lists)-1)
}

func main() {
	l1 := NewMyListNode([]int{1, 3, 5, 7})
	l2 := NewMyListNode([]int{2, 4, 6, 8})
	l3 := NewMyListNode([]int{9, 10, 11, 12})
	printListNode(mergeKLists1([]*ListNode{l1, l2, l3}))
}
