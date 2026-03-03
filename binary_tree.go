package main

// 二叉树算法实现：
// 🌟技巧1：两种思维模式 - 解二叉树问题有两种视角：「遍历」思路（用外部变量收集结果，类似回溯）和「分解子问题」思路（函数返回值代表子问题结果，类似分治）；优先考虑后者，代码更简洁
// 🌟技巧2：前序 vs 后序位置 - 前序位置的代码在进入节点时执行（能获取父节点传下来的信息），后序位置在离开节点时执行（能获取左右子树返回的信息）；需要子树信息时必须用后序
// 🌟技巧3：层序遍历框架 - 层序遍历用队列实现，外层for控制层数，内层for遍历当前层所有节点（sz=len(q)固定当前层大小），depth在内层结束后递增，与BFS框架完全一致
// 🌟技巧4：后序位置的妙用 - 需要同时知道左右子树信息才能计算当前节点结果时（如树的直径、路径和），必须在后序位置用闭包变量或返回值收集，不能在前序位置提前计算

// ⚠️易错点1：base case 为 nil 而非叶子节点 - 递归的终止条件通常是 root==nil 返回0或nil，而非判断叶子节点；判断叶子节点（left==nil && right==nil）会导致少处理一层
// ⚠️易错点2：分解子问题时不要用外部变量 - 「分解子问题」思路下函数语义靠返回值传递，混用外部变量会导致语义混乱；若必须用外部变量收集答案，应统一切换为「遍历」思路
// ⚠️易错点3：层序遍历的 sz 必须在内层循环前固定 - sz:=len(q) 必须在内层 for 之前赋值，不能写成 i<len(q)，否则每次入队新节点都会扩大边界，导致跨层处理
// ⚠️易错点4：直径不等于最大深度 - 树的直径是「经过某节点的左右子树深度之和」，不是树的最大深度；需要在后序位置用 leftMax+rightMax 更新全局最大值，而非直接返回深度

// 何时运用二叉树算法：
// ❓1、能否将问题分解为「左子树结果 + 右子树结果 → 当前节点结果」？优先用后序位置的分解子问题思路
// ❓2、是否需要父节点向子节点传递信息？（如路径前缀和）用前序位置 + 函数参数传递
// ❓3、是否需要按层处理或求最短路径？用层序遍历（BFS）框架而非递归

// 0、func traverse(root *TreeNode)                           // 二叉树递归遍历框架
// 1、func traverseNary(root *NTreeNode)                      // N叉树递归遍历框架
// 2、func levelOrderTraverse(root *TreeNode)                 // 二叉树层序遍历框架
// 3、func levelOrderTraverseNary(root *NTreeNode)            // N叉树层序遍历框架
// 4、func maxDepth(root *TreeNode) int                       // 二叉树的最大深度
// 5、func preorderTraversal(root *TreeNode) []int            // 二叉树前序遍历
// 6、func diameterOfBinaryTree(root *TreeNode) int           // 二叉树的直径

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

	return 1 + max(leftDepth, rightDepth)
}

// 二叉树的前序遍历 - 分解子问题思路
func preorderTraversal(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}
	res := []int{root.Val}
	res = append(res, preorderTraversal(root.Left)...)
	res = append(res, preorderTraversal(root.Right)...)
	return res
}

// 二叉树的直径 - 分解子问题思路
func diameterOfBinaryTree(root *TreeNode) int {
	// 二叉树的最大直径 = 某个节点的（左子树最大直径 + 右子树最大直径）
	// 某个节点的最大直径 = 左子树最大直径+1 or 右子树最大直径+1

	var depthFun func(root *TreeNode) int
	finalMax := 0

	depthFun = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		leftMax := depthFun(root.Left)
		rightMax := depthFun(root.Right)
		finalMax = max(finalMax, leftMax+rightMax)
		return 1 + max(leftMax, rightMax)
	}

	depthFun(root)
	return finalMax
}
