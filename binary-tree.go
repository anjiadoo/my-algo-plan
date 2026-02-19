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
