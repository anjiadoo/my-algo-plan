package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

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

////////////////////////////////////////////////////////////////////

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type NTreeNode struct {
	Val      int
	Children []*NTreeNode
}

func (r *TreeNode) Display() {
	var fn func(r *TreeNode, preorder, inorder, postorder *[]int)
	fn = func(root *TreeNode, preorder, inorder, postorder *[]int) {
		if root == nil {
			return
		}
		*preorder = append(*preorder, root.Val)
		fn(root.Left, preorder, inorder, postorder)
		*inorder = append(*inorder, root.Val)
		fn(root.Right, preorder, inorder, postorder)
		*postorder = append(*postorder, root.Val)
	}
	var preorder, inorder, postorder []int
	fn(r, &preorder, &inorder, &postorder)
	fmt.Printf("前序：%v 中序：%v 后序：%v\n", preorder, inorder, postorder)
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

////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////

// 最大二叉树 https://leetcode.cn/problems/maximum-binary-tree/
func constructMaximumBinaryTree(nums []int) *TreeNode {
	var build func(nums []int, lo, hi int) *TreeNode

	build = func(nums []int, lo, hi int) *TreeNode {
		if lo == hi {
			return nil
		}
		maxVal := math.MinInt
		index := -1
		for i := lo; i < hi; i++ {
			if nums[i] > maxVal {
				index = i
				maxVal = nums[i]
			}
		}
		newNode := &TreeNode{Val: maxVal}
		newNode.Left = build(nums, lo, index)
		newNode.Right = build(nums, index+1, hi)
		return newNode
	}

	return build(nums, 0, len(nums))
}

// 从前序与中序遍历序列构造二叉树 https://leetcode.cn/problems/construct-binary-tree-from-preorder-and-inorder-traversal/
func buildTree(preorder []int, inorder []int) *TreeNode {
	var build func(preorder []int, preStart, preEnd int, inorder []int, inStart, inEnd int) *TreeNode

	build = func(preorder []int, preStart, preEnd int, inorder []int, inStart, inEnd int) *TreeNode {
		if preStart > preEnd {
			return nil
		}

		rootVal := preorder[preStart]
		index := 0
		for i := inStart; i <= inEnd; i++ {
			if inorder[i] == rootVal {
				index = i
				break
			}
		}
		// 注意区间范围:[left,index)，index已经被占用了
		leftSize := index - inStart

		root := &TreeNode{Val: rootVal}
		root.Left = build(preorder, preStart+1, preStart+leftSize, inorder, inStart, index-1)
		root.Right = build(preorder, preStart+leftSize+1, preEnd, inorder, index+1, inEnd)
		return root
	}

	return build(preorder, 0, len(preorder)-1, inorder, 0, len(inorder)-1)
}

// 从中序与后序遍历序列构造二叉树 https://leetcode.cn/problems/construct-binary-tree-from-inorder-and-postorder-traversal/
func buildTree2(inorder []int, postorder []int) *TreeNode {
	var build func(inorder []int, inStart, inEnd int, postorder []int, postStart, postEnd int) *TreeNode

	build = func(inorder []int, inStart, inEnd int, postorder []int, postStart, postEnd int) *TreeNode {
		if postStart > postEnd {
			return nil
		}

		// 前序遍历 [rootVal，左子树，右子树]
		// 中序遍历 [左子树，rootVal，右子树]
		// 后序遍历 [左子树，右子树，rootVal]

		rootVal := postorder[postEnd]
		index := 0
		for i := inStart; i <= inEnd; i++ {
			if inorder[i] == rootVal {
				index = i
				break
			}
		}
		// 注意区间范围:[left,index)，index已经被占用了
		leftSize := index - inStart

		root := &TreeNode{Val: rootVal}
		root.Left = build(inorder, inStart, index-1, postorder, postStart, postStart+leftSize-1)
		root.Right = build(inorder, index+1, inEnd, postorder, postStart+leftSize, postEnd-1)
		return root
	}

	return build(inorder, 0, len(inorder)-1, postorder, 0, len(postorder)-1)
}

// 根据前序和后序遍历构造二叉树 https://leetcode.cn/problems/construct-binary-tree-from-preorder-and-postorder-traversal/description/
func constructFromPrePost(preorder []int, postorder []int) *TreeNode {
	var build func(preorder []int, preStart, preEnd int, postorder []int, postStart, postEnd int) *TreeNode

	build = func(preorder []int, preStart, preEnd int, postorder []int, postStart, postEnd int) *TreeNode {
		if preStart > preEnd {
			return nil
		}

		if preStart == preEnd {
			return &TreeNode{
				Val:   preorder[preStart],
				Left:  nil,
				Right: nil,
			}
		}

		// 前序遍历 [rootVal，左子树，右子树]
		// 中序遍历 [左子树，rootVal，右子树]
		// 后序遍历 [左子树，右子树，rootVal]

		rootVal := preorder[preStart]
		leftRootVal := preorder[preStart+1]

		index := 0
		for i := postStart; i <= postEnd; i++ {
			if postorder[i] == leftRootVal {
				index = i
				break
			}
		}
		// 注意区间范围:[left,index]，index未被占用了
		leftSize := index - postStart + 1

		root := &TreeNode{Val: rootVal}
		root.Left = build(preorder, preStart+1, preStart+leftSize, postorder, postStart, index)
		root.Right = build(preorder, preStart+leftSize+1, preEnd, postorder, index+1, postEnd-1)
		return root
	}

	return build(preorder, 0, len(preorder)-1, postorder, 0, len(postorder)-1)
}

// 寻找重复的子树 https://leetcode.cn/problems/find-duplicate-subtrees/
func findDuplicateSubtrees(root *TreeNode) []*TreeNode {
	mapping := make(map[string]int)
	var res []*TreeNode

	var traverse func(root *TreeNode) string
	traverse = func(root *TreeNode) string {
		if root == nil {
			return "#"
		}

		left := traverse(root.Left)
		right := traverse(root.Right)

		str := left + "," + right + "," + strconv.Itoa(root.Val)
		if cnt, ok := mapping[str]; ok && cnt == 1 {
			res = append(res, root)
		}
		mapping[str]++
		return str
	}
	traverse(root)
	return res
}

// 二叉树的序列化与反序列化 https://leetcode.cn/problems/serialize-and-deserialize-binary-tree/

type Codec struct{}

func NewCodec() Codec { return Codec{} }

func (c *Codec) serialize(root *TreeNode) string {
	var _serialize func(root *TreeNode)

	var res []string
	_serialize = func(root *TreeNode) {
		if root == nil {
			res = append(res, "#")
			return
		}
		res = append(res, strconv.Itoa(root.Val))
		_serialize(root.Left)
		_serialize(root.Right)
	}
	_serialize(root)
	return strings.Join(res, ",")
}

func (c *Codec) deserialize(data string) *TreeNode {
	var _deserialize func(nodes *[]string) *TreeNode

	_deserialize = func(nodes *[]string) *TreeNode {
		if len(*nodes) == 0 {
			return nil
		}
		rootValStr := (*nodes)[0]
		*nodes = (*nodes)[1:]
		if rootValStr == "#" {
			return nil
		}
		rootVal, _ := strconv.Atoi(rootValStr)
		root := &TreeNode{Val: rootVal}
		root.Left = _deserialize(nodes)
		root.Right = _deserialize(nodes)
		return root
	}

	nodes := strings.Split(data, ",")
	return _deserialize(&nodes)
}

// 二叉搜索树中第 K 小的元素 https://leetcode.cn/problems/kth-smallest-element-in-a-bst/description/
func kthSmallest1(root *TreeNode, k int) int {
	var traverse func(root *TreeNode, )

	res := 0  // 记录结果
	rank := 0 // 记录当前元素的排名

	traverse = func(root *TreeNode, ) {
		if root == nil {
			return
		}
		traverse(root.Left)
		rank++
		if rank == k {
			res = root.Val
			return
		}
		traverse(root.Right)
	}
	traverse(root)
	return res
}

// 把二叉搜索树转换为累加树 https://leetcode.cn/problems/convert-bst-to-greater-tree/
func convertBST(root *TreeNode) *TreeNode {
	var traverse func(root *TreeNode)
	sum := 0 // 记录累加和
	traverse = func(root *TreeNode) {
		if root == nil {
			return
		}
		traverse(root.Right)
		sum = sum + root.Val
		root.Val = sum
		traverse(root.Left)
	}
	traverse(root) // BST先遍历右子树，再遍历左子树就是降序了
	return root
}

// 验证二叉搜索树 https://leetcode.cn/problems/validate-binary-search-tree/
func isValidBST(root *TreeNode) bool {
	var _isValidBST func(root, min, max *TreeNode) bool

	_isValidBST = func(root, min, max *TreeNode) bool {
		if root == nil {
			return true
		}
		// 若 root.Val 不符合 max 和 min 的限制，说明不是合法 BST
		if min != nil && min.Val >= root.Val {
			return false
		}
		if max != nil && max.Val <= root.Val {
			return false
		}
		// 根据定义，限定左子树的最大值是 root.Val，右子树的最小值是 root.Val
		return _isValidBST(root.Left, min, root) && _isValidBST(root.Right, root, max)
	}
	return _isValidBST(root, nil, nil)
}

// 二叉搜索树中的搜索 https://leetcode.cn/problems/search-in-a-binary-search-tree/
func searchBST(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Val == val {
		return root
	}
	if root.Val > val {
		return searchBST(root.Left, val)
	} else {
		return searchBST(root.Right, val)
	}
}

// 二叉搜索树中的插入操作 https://leetcode.cn/problems/insert-into-a-binary-search-tree/description/
func insertIntoBST(root *TreeNode, val int) *TreeNode {
	// 定义：在以root为根的BST中插入val节点，返回插入后的根节点
	if root == nil {
		// 找到空位置插入新节点
		return &TreeNode{Val: val}
	}
	if root.Val > val {
		root.Left = insertIntoBST(root.Left, val)
	} else {
		root.Right = insertIntoBST(root.Right, val)
	}
	return root
}

// 删除二叉搜索树中的节点 https://leetcode.cn/problems/delete-node-in-a-bst/description/
func deleteNode(root *TreeNode, key int) *TreeNode {
	if root == nil {
		return nil
	}
	// 左子树为空返回右子树，右子树为空返回左子树
	// 左右子树都存在，把左子树接到右子树最小位置
	if root.Val == key {
		if root.Left == nil {
			return root.Right
		}
		if root.Right == nil {
			return root.Left
		}
		if root.Right != nil {
			p := root.Right
			for p.Left != nil {
				p = p.Left
			}
			p.Left = root.Left
			return root.Right
		}
	}
	if root.Val > key {
		root.Left = deleteNode(root.Left, key)
	} else {
		root.Right = deleteNode(root.Right, key)
	}
	return root
}

// 不同的二叉搜索树 https://leetcode.cn/problems/unique-binary-search-trees/description/
func numTrees(n int) int {
	// 定义：返回[lo, hi]范围内构造的不同BST的数量
	var count func(lo, hi int) int

	memo := make([][]int, n+1)
	for i := range memo {
		memo[i] = make([]int, n+1)
	}

	count = func(lo, hi int) int {
		if lo > hi {
			return 1
		}
		if memo[lo][hi] != 0 {
			return memo[lo][hi]
		}
		res := 0
		for i := lo; i <= hi; i++ {
			left := count(lo, i-1)
			right := count(i+1, hi)
			res += left * right
		}
		memo[lo][hi] = res
		return memo[lo][hi]
	}
	return count(1, n)
}

// 不同的二叉搜索树II https://leetcode.cn/problems/unique-binary-search-trees-ii/description/
func generateTrees(n int) []*TreeNode {
	// 定义：返回[lo, hi]范围内构造的不同BST
	var build func(lo, hi int) []*TreeNode

	build = func(lo, hi int) []*TreeNode {
		var res []*TreeNode
		if lo > hi {
			// 这里需要装一个null元素，
			// 才能让下面的两个内层for循环都能进入，正确地创建出叶子节点
			res = append(res, nil)
			return res
		}
		for i := lo; i <= hi; i++ {
			leftTree := build(lo, i-1)
			rightTree := build(i+1, hi)

			for j := 0; j < len(leftTree); j++ {
				for k := 0; k < len(rightTree); k++ {
					root := &TreeNode{Val: i}
					root.Left = leftTree[j]
					root.Right = rightTree[k]
					res = append(res, root)
				}
			}
		}
		return res
	}
	return build(1, n)
}

// 二叉搜索子树的最大键值和 https://leetcode.cn/problems/maximum-sum-bst-in-binary-tree/
func maxSumBST(root *TreeNode) int {
	var findMaxMinSum func(root *TreeNode) []int
	maxSum := 0

	// 计算以root为根的二叉树是否是BST
	// 以及它的的最大值、最小值、节点和

	findMaxMinSum = func(root *TreeNode) []int {
		if root == nil {
			return []int{1, math.MinInt, math.MaxInt, 0}
		}

		left := findMaxMinSum(root.Left)
		right := findMaxMinSum(root.Right)

		res := make([]int, 4)

		if left[0] == 1 && right[0] == 1 &&
			root.Val > left[1] && root.Val < right[2] {
			res[0] = 1
			res[1] = max(root.Val, right[1])
			res[2] = min(root.Val, left[2])
			res[3] = left[3] + root.Val + right[3]
			maxSum = max(maxSum, res[3])
		}
		return res
	}
	findMaxMinSum(root)
	return maxSum
}

// 二叉树的所有路径 https://leetcode.cn/problems/binary-tree-paths/description/
func binaryTreePaths(root *TreeNode) []string {
	var traverse func(root *TreeNode, path *[]int, res *[][]int)

	traverse = func(root *TreeNode, path *[]int, res *[][]int) {
		if root == nil {
			return
		}

		if root.Left == nil && root.Right == nil {
			temp := make([]int, len(*path))
			copy(temp, *path)
			temp = append(temp, root.Val)
			*res = append(*res, temp)
			return
		}
		*path = append(*path, root.Val)
		traverse(root.Left, path, res)
		traverse(root.Right, path, res)
		*path = (*path)[:len(*path)-1]
	}

	var path []int
	var res [][]int

	traverse(root, &path, &res)
	var result []string
	for _, itr := range res {
		sPath := ""
		for _, val := range itr {
			sPath += strconv.Itoa(val) + "->"
		}
		result = append(result, strings.TrimRight(sPath, "->"))
	}
	return result
}

// 求根节点到叶节点数字之和 https://leetcode.cn/problems/sum-root-to-leaf-numbers/
func sumNumbers(root *TreeNode) int {
	var traverse func(root *TreeNode, path *string)
	var res int

	traverse = func(root *TreeNode, path *string) {
		if root == nil {
			return
		}
		if root.Left == nil && root.Right == nil {
			temp := *path + strconv.Itoa(root.Val)
			num, _ := strconv.Atoi(temp)
			res += num
		}

		// 题目说明了val范围0～9
		*path += strconv.Itoa(root.Val)
		traverse(root.Left, path)
		traverse(root.Right, path)
		*path = (*path)[:len(*path)-1]
	}

	var path string
	traverse(root, &path)
	return res
}

// 二叉树的右视图 https://leetcode.cn/problems/binary-tree-right-side-view/description/
func rightSideView(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	queue := []*TreeNode{root}
	var res []int

	for len(queue) > 0 {
		sz := len(queue)
		res = append(res, queue[sz-1].Val)
		for i := 0; i < sz; i++ {
			curr := queue[0]
			queue = queue[1:]
			if curr.Left != nil {
				queue = append(queue, curr.Left)
			}
			if curr.Right != nil {
				queue = append(queue, curr.Right)
			}
		}
	}
	return res
}

// 二叉树中的伪回文路径 https://leetcode.cn/problems/pseudo-palindromic-paths-in-a-binary-tree/description/
func pseudoPalindromicPaths(root *TreeNode) int {
	var traverse func(root *TreeNode, path *map[int]int)
	var res int

	isPalindrome := func(path *map[int]int) bool {
		cnt := 0
		for _, val := range *path {
			if val%2 == 1 {
				cnt++
			}
			if cnt > 1 {
				return false
			}
		}
		return true
	}

	traverse = func(root *TreeNode, path *map[int]int) {
		if root == nil {
			return
		}
		if root.Left == nil && root.Right == nil {
			(*path)[root.Val]++
			if isPalindrome(path) {
				res++
			}
			(*path)[root.Val]--
			return
		}
		(*path)[root.Val]++
		traverse(root.Left, path)
		traverse(root.Right, path)
		(*path)[root.Val]--
		if _, ok := (*path)[root.Val]; !ok {
			delete(*path, root.Val)
		}
	}

	path := map[int]int{}
	traverse(root, &path)
	return res
}

// 左叶子之和 https://leetcode.cn/problems/sum-of-left-leaves/
func sumOfLeftLeaves(root *TreeNode) int {
	var traverse func(root *TreeNode)
	var sum int

	traverse = func(root *TreeNode) {
		if root == nil {
			return
		}
		if root.Left != nil &&
			root.Left.Left == nil &&
			root.Left.Right == nil {
			sum += root.Left.Val
		}
		traverse(root.Left)
		traverse(root.Right)
	}
	traverse(root)
	return sum
}

// 在二叉树中增加一行 https://leetcode.cn/problems/add-one-row-to-tree/description/
func addOneRow(root *TreeNode, val int, depth int) *TreeNode {
	if depth == 1 {
		newRoot := &TreeNode{Val: val}
		newRoot.Left = root
		return newRoot
	}

	q := []*TreeNode{root}
	step := 1
	for len(q) > 0 {
		if step == depth-1 {
			for _, node := range q {
				left := &TreeNode{Val: val, Left: node.Left}
				right := &TreeNode{Val: val, Right: node.Right}

				node.Left = left
				node.Right = right
			}
			break
		}
		sz := len(q)
		for i := 0; i < sz; i++ {
			cur := q[0]
			q = q[1:]
			if cur.Left != nil {
				q = append(q, cur.Left)
			}
			if cur.Right != nil {
				q = append(q, cur.Right)
			}
		}
		step++
	}
	return root
}

// 翻转二叉树以匹配先序遍历 https://leetcode.cn/problems/flip-binary-tree-to-match-preorder-traversal/
func flipMatchVoyage(root *TreeNode, voyage []int) []int {
	var traverse func(root *TreeNode, i *int, canMatch *bool, res *[]int)

	traverse = func(root *TreeNode, i *int, canMatch *bool, res *[]int) {
		if root == nil || !*canMatch {
			return
		}
		if root.Val != voyage[*i] {
			*canMatch = false
			return
		}
		*i++
		if root.Left != nil && root.Left.Val != voyage[*i] {
			root.Left, root.Right = root.Right, root.Left
			*res = append(*res, root.Val)
		}
		traverse(root.Left, i, canMatch, res)
		traverse(root.Right, i, canMatch, res)
	}

	var i = 0
	var res []int
	var canMatch = true

	traverse(root, &i, &canMatch, &res)
	if canMatch {
		return res
	}
	return []int{-1}
}

// 二叉树的垂序遍历 https://leetcode.cn/problems/vertical-order-traversal-of-a-binary-tree/
func verticalTraversal(root *TreeNode) [][]int {
	var traverse func(root *TreeNode, row, col int)

	type Triple struct {
		row, col, val int
	}

	var nodes []Triple

	traverse = func(root *TreeNode, row, col int) {
		if root == nil {
			return
		}
		nodes = append(nodes, Triple{row: row, col: col, val: root.Val})
		traverse(root.Left, row+1, col-1)
		traverse(root.Right, row+1, col+1)
	}

	traverse(root, 0, 0)

	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].col != nodes[j].col {
			return nodes[i].col < nodes[j].col
		}
		if nodes[i].row != nodes[j].row {
			return nodes[i].row < nodes[j].row
		}
		return nodes[i].val < nodes[j].val
	})

	var res [][]int
	preCol := math.MinInt
	for _, cur := range nodes {
		if cur.col != preCol {
			// 开始记录新的一列
			res = append(res, []int{})
			preCol = cur.col
		}
		res[len(res)-1] = append(res[len(res)-1], cur.val)
	}

	return res
}

// 祖父节点值为偶数的节点和 https://leetcode.cn/problems/sum-of-nodes-with-even-valued-grandparent/description/
func sumEvenGrandparent(root *TreeNode) int {
	var traverse func(root *TreeNode)
	var sum int

	traverse = func(root *TreeNode) {
		if root == nil {
			return
		}
		if root.Val%2 == 0 {
			// 累加左子树孙子节点的值
			if root.Left != nil {
				if root.Left.Left != nil {
					sum += root.Left.Left.Val
				}
				if root.Left.Right != nil {
					sum += root.Left.Right.Val
				}
			}
			// 累加右子树孙子节点的值
			if root.Right != nil {
				if root.Right.Left != nil {
					sum += root.Right.Left.Val
				}
				if root.Right.Right != nil {
					sum += root.Right.Right.Val
				}
			}
		}
		traverse(root.Left)
		traverse(root.Right)
	}
	traverse(root)
	return sum
}

// 路径总和III https://leetcode.cn/problems/path-sum-iii/
func pathSum(root *TreeNode, targetSum int) int {
	var traverse func(root *TreeNode, pathSum int, preSumCount map[int]int) int

	traverse = func(root *TreeNode, pathSum int, preSumCount map[int]int) int {
		if root == nil {
			return 0
		}

		pathSum += root.Val

		// 从二叉树的根节点开始，路径和为 pathSum - targetSum 的路径条数
		// 就是路径和为 targetSum 的路径条数
		res := preSumCount[pathSum-targetSum]
		preSumCount[pathSum]++

		res += traverse(root.Left, pathSum, preSumCount)
		res += traverse(root.Right, pathSum, preSumCount)

		preSumCount[pathSum]--
		pathSum -= root.Val
		return res
	}

	// 定义：从二叉树的根节点开始，路径和为 pathSum 的路径有 preSumCount[pathSum] 个
	preSumCount := make(map[int]int)
	preSumCount[0] = 1

	return traverse(root, 0, preSumCount)
}

func main() {

	root := NewTreeNode("6,7,8,2,7,1,3,9,null,1,4,null,null,null,5")
	fmt.Println(sumEvenGrandparent(root))

	//root := NewTreeNode("4,2,6,3,1,5,null")
	//printTree(addOneRow(root, 1, 2))

	//root := NewTreeNode("3,9,20,null,null,15,7")
	//fmt.Println(sumOfLeftLeaves(root))

	//root := NewTreeNode("2,3,1,3,1,null,1")
	//fmt.Println(pseudoPalindromicPaths(root))

	//root := NewTreeNode("1,2,3,null,5,null,4")
	//fmt.Println(rightSideView(root))

	//fmt.Println(numTrees(3))
	//fmt.Println(numTrees(1))

	//root := buildTree([]int{3, 9, 20, 15, 7}, []int{9, 3, 15, 20, 7})
	//root.Display()

	//root := buildTree2([]int{9, 3, 15, 20, 7}, []int{9, 15, 7, 20, 3})
	//root.Display()

	//root := constructFromPrePost([]int{1, 2, 4, 5, 3, 6, 7}, []int{4, 5, 2, 6, 7, 3, 1})
	//root.Display()

	//root := buildTree([]int{3, 9, 20, 15, 7}, []int{9, 3, 15, 20, 7})
	//cc := NewCodec()
	//str := cc.serialize(root)
	//fmt.Println(str)
	//newRoot := cc.deserialize(str)
	//newRoot.Display()

	//root := buildTree([]int{3, 9, 20, 15, 7}, []int{9, 3, 15, 20, 7})
	//fmt.Println(kthSmallest1(root, 1))
}

/*================================================================*/

func NewTreeNode(nodes string) *TreeNode {
	parts := strings.Split(nodes, ",")
	if len(parts) == 0 || parts[0] == "null" {
		return nil
	}

	rootVal, _ := strconv.Atoi(parts[0])
	parts = parts[1:]
	root := &TreeNode{Val: rootVal}
	queue := []*TreeNode{root}

	for len(queue) > 0 && len(parts) > 0 {
		cur := queue[0]
		queue = queue[1:]

		// 构建左子节点
		left := parts[0]
		parts = parts[1:]
		if left != "null" {
			val, _ := strconv.Atoi(left)
			cur.Left = &TreeNode{Val: val}
			queue = append(queue, cur.Left)
		}

		// 构建右子节点
		right := parts[0]
		parts = parts[1:]
		if right != "null" {
			val, _ := strconv.Atoi(right)
			cur.Right = &TreeNode{Val: val}
			queue = append(queue, cur.Right)
		}
	}

	// 打印二叉树结构
	fmt.Println("--------------------")
	printTree(root)
	fmt.Println("--------------------")
	return root
}

func printTree(root *TreeNode) {
	if root == nil {
		fmt.Println("<空树>")
		return
	}
	shiftLines := func(lines []string, shift int) []string {
		prefix := strings.Repeat(" ", shift)
		result := make([]string, len(lines))
		for i, line := range lines {
			result[i] = prefix + line
		}
		return result
	}

	// maxLineWidth 返回行列表中最长行的宽度
	maxLineWidth := func(lines []string) int {
		maxW := 0
		for _, line := range lines {
			if len(line) > maxW {
				maxW = len(line)
			}
		}
		return maxW
	}

	findRootLen := func(line string, start int) int {
		length := 0
		for i := start; i < len(line); i++ {
			if line[i] == ' ' {
				break
			}
			length++
		}
		return length
	}

	findRootStart := func(line string) int {
		for i, ch := range line {
			if ch != ' ' {
				return i
			}
		}
		return 0
	}

	var buildTreeLines func(node *TreeNode) []string
	buildTreeLines = func(node *TreeNode) []string {
		if node == nil {
			return []string{}
		}

		rootStr := strconv.Itoa(node.Val)
		rootLen := len(rootStr)

		// 叶子节点
		if node.Left == nil && node.Right == nil {
			return []string{rootStr}
		}

		// 只有左子树
		if node.Right == nil {
			leftLines := buildTreeLines(node.Left)
			leftRootStart := findRootStart(leftLines[0])
			leftRootMid := leftRootStart + findRootLen(leftLines[0], leftRootStart)/2

			// 根节点放在左子树根的右上方
			rootStart := leftRootMid + 1
			// 构建连接线
			var lines []string
			// 第一行：根节点
			line1 := strings.Repeat(" ", rootStart) + rootStr
			lines = append(lines, line1)
			// 第二行：/ 连接符
			slashPos := rootStart - 1
			line2 := strings.Repeat(" ", slashPos) + "/"
			lines = append(lines, line2)
			// 后续行：左子树内容
			for _, l := range leftLines {
				lines = append(lines, l)
			}
			return lines
		}

		// 只有右子树
		if node.Left == nil {
			rightLines := buildTreeLines(node.Right)
			rightRootStart := findRootStart(rightLines[0])
			rightRootMid := rightRootStart + findRootLen(rightLines[0], rightRootStart)/2

			// 根节点放在右子树根的左上方
			rootEnd := rightRootMid // 根节点结束位置对齐到右子树根的中间左边
			rootStart := rootEnd - rootLen
			if rootStart < 0 {
				// 需要右移整个右子树
				shift := -rootStart
				rootStart = 0
				rightLines = shiftLines(rightLines, shift)
				rightRootMid += shift
			}

			var lines []string
			// 第一行：根节点
			line1 := strings.Repeat(" ", rootStart) + rootStr
			lines = append(lines, line1)
			// 第二行：\ 连接符
			backslashPos := rootStart + rootLen
			line2 := strings.Repeat(" ", backslashPos) + "\\"
			lines = append(lines, line2)
			// 后续行：右子树内容
			for _, l := range rightLines {
				lines = append(lines, l)
			}
			return lines
		}

		// 左右子树都存在
		leftLines := buildTreeLines(node.Left)
		rightLines := buildTreeLines(node.Right)

		leftRootStart := findRootStart(leftLines[0])
		leftRootLen := findRootLen(leftLines[0], leftRootStart)
		leftRootMid := leftRootStart + leftRootLen/2

		rightRootStart := findRootStart(rightLines[0])
		rightRootLen := findRootLen(rightLines[0], rightRootStart)
		rightRootMid := rightRootStart + rightRootLen/2

		// 根节点位置：在左子树 / 和右子树 \ 之间
		// / 的位置在 leftRootMid + 1 的上方
		// \ 的位置在 rightRootMid - 1 + gap 的上方
		slashPos := leftRootMid + 1
		// 根节点起始位置
		rootStart := slashPos + 1
		rootEndPos := rootStart + rootLen
		// \ 的位置
		backslashPos := rootEndPos
		// 右子树需要的偏移量：使得右子树的根中点对齐到 backslashPos + 1
		rightShift := backslashPos + 1 - rightRootMid
		// 确保右子树不与左子树重叠
		leftWidth := maxLineWidth(leftLines)
		gap := 3 // 左右子树之间最小间距
		if rightShift < leftWidth+gap {
			rightShift = leftWidth + gap
			// 重新计算根节点和连接符位置
			backslashPos = rightShift + rightRootMid - 1
			rootStart = (slashPos + backslashPos) / 2
			rootEndPos = rootStart + rootLen
			if backslashPos < rootEndPos {
				backslashPos = rootEndPos
			}
		}

		var lines []string
		// 第一行：根节点
		line1 := strings.Repeat(" ", rootStart) + rootStr
		lines = append(lines, line1)
		// 第二行：/ 和 \ 连接符
		line2 := strings.Repeat(" ", slashPos) + "/" + strings.Repeat(" ", backslashPos-slashPos-1) + "\\"
		lines = append(lines, line2)
		// 后续行：合并左右子树
		maxLines := len(leftLines)
		if len(rightLines) > maxLines {
			maxLines = len(rightLines)
		}
		for i := 0; i < maxLines; i++ {
			leftPart := ""
			if i < len(leftLines) {
				leftPart = leftLines[i]
			}
			rightPart := ""
			if i < len(rightLines) {
				rightPart = rightLines[i]
			}
			// 合并行：左子树内容 + 间距 + 右子树内容
			if rightPart != "" {
				// 右子树内容需要偏移
				neededWidth := rightShift
				if len(leftPart) < neededWidth {
					merged := leftPart + strings.Repeat(" ", neededWidth-len(leftPart)) + rightPart
					lines = append(lines, merged)
				} else {
					merged := leftPart + "   " + rightPart
					lines = append(lines, merged)
				}
			} else {
				lines = append(lines, leftPart)
			}
		}
		return lines
	}

	lines := buildTreeLines(root)
	for _, line := range lines {
		fmt.Println(line)
	}
}

/*================================================================*/
