package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type NTreeNode struct {
	Val      int
	Children []*NTreeNode
}

// 二叉树递归遍历框架
func traverse(root *TreeNode) {
	if root == nil {
		return
	}
	// 前序遍历
	traverse(root.Left)
	// 中序遍历
	traverse(root.Right)
	// 后续遍历
}

// N 叉树的遍历框架
func traverseNary(root *NTreeNode) {
	if root == nil {
		return
	}
	// 前序位置
	for _, child := range root.Children {
		traverseNary(child)
	}
	// 后序位置
}

// 二叉树层序遍历框架
func levelOrderTraverse(root *TreeNode) {
	if root == nil {
		return
	}
	q := []*TreeNode{root}
	// 记录当前遍历到的层数（根节点视为第 1 层）
	depth := 1

	for len(q) > 0 {
		// 获取当前队列长度
		sz := len(q)
		for i := 0; i < sz; i++ {
			// 弹出队列头
			cur := q[0]
			q = q[1:]

			// 访问 cur 节点，同时知道它所在的层数

			// 把 cur 的左右子节点加入队列
			if cur.Left != nil {
				q = append(q, cur.Left)
			}
			if cur.Right != nil {
				q = append(q, cur.Right)
			}
		}
		depth++
	}
}

// N 叉树层序遍历框架
func levelOrderTraverseNary(root *NTreeNode) {
	if root == nil {
		return
	}
	q := []*NTreeNode{root}
	// 记录当前遍历到的层数（根节点视为第 1 层）
	depth := 1

	for len(q) > 0 {
		sz := len(q)
		for i := 0; i < sz; i++ {
			cur := q[0]
			q = q[1:]

			// 访问 cur 节点，同时知道它所在的层数

			for _, child := range cur.Children {
				q = append(q, child)
			}
		}
		depth++
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 二叉树的最大深度 - 分解子问题思路
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	leftDepth := maxDepth(root.Left)
	rightDepth := maxDepth(root.Right)

	if leftDepth > rightDepth {
		return leftDepth + 1
	} else {
		return rightDepth + 1
	}
}

// 二叉树的前序遍历 - 分解子问题思路
func preorderTraversal(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}
	res := []int{root.Val}
	leftVals := preorderTraversal(root.Left)
	rightVals := preorderTraversal(root.Right)
	res = append(res, leftVals...)
	res = append(res, rightVals...)
	return res
}

// 二叉树的直径 - 分解子问题思路
func diameterOfBinaryTree(root *TreeNode) int {
	// 二叉树的最大直径 = 某个节点的（左子树最大直径 + 右子树最大直径）
	// 某个节点的最大直径 = 左子树最大直径+1 or 右子树最大直径+1

	var depthFun func(root *TreeNode) int
	finalMaxNum := 0

	maxFun := func(x, y int) int {
		if x > y {
			return x
		}
		return y
	}

	depthFun = func(root *TreeNode) int {
		if root == nil {
			return 0
		}

		leftMax := depthFun(root.Left)
		rightMax := depthFun(root.Right)

		finalMaxNum = maxFun(finalMaxNum, leftMax+rightMax)

		return 1 + maxFun(leftMax, rightMax)
	}

	depthFun(root)
	return finalMaxNum
}
