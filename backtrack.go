/*
 * ============================================================================
 *                      📘 回溯算法全集 · 核心记忆框架
 * ============================================================================
 * 【一句话理解回溯】
 *
 *   回溯 = N 叉树的遍历 + 在前序位置「做选择」 + 在后序位置「撤销选择」。
 *   所有结果都是从根到某个节点的路径，穷举决策树即穷举所有答案。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【标准框架模板（背熟！）】
 *
 *   func backtrack(路径 track, 选择列表, ...) {
 *       if 满足结束条件 {
 *           result = append(result, deepCopy(track))  // ← 必须深拷贝！
 *           return
 *       }
 *       for _, 选择 := range 选择列表 {
 *           // 剪枝（可选）
 *           做选择                    // ← 前序位置：加入 track，标记 used
 *           backtrack(路径, 选择列表)
 *           撤销选择                  // ← 后序位置：移出 track，还原 used
 *       }
 *   }
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【三要素定位法：解任何回溯题的第一步】
 *
 *   ① 路径 (track)     ：已经做出的选择集合，即当前节点到根的路径
 *   ② 选择列表         ：当前节点可以做的所有选择（随层数变化）
 *   ③ 结束条件         ：到达决策树的「叶子节点」时，收集 track 进入结果集
 *
 *   只要把这三项写清楚，框架代码几乎可以无脑套用。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【6 类问题分类矩阵（元素特征 × 复选规则）】
 *
 *                    │ 元素无重复            │ 元素有重复
 *   ─────────────────┼─────────────────────┼──────────────────────────
 *   排列-顺序有关      │ permute             │ permuteUnique
 *   不可复选          │ used标记，从0遍历     │ 先排序 + !used[i-1] 树枝去重
 *   ─────────────────┼─────────────────────┼──────────────────────────
 *   组合/子集         │ combine / subsets   │ subsetsWithDup / combinationSum2
 *   不可复选          │ start参数，从start    │ 先排序 + i>start 同层去重
 *   ─────────────────┼─────────────────────┼──────────────────────────
 *   组合/子集         │ combinationSum      │ 先排序 + i>start 去重 + 传 i
 *   可复选            │ 传 i（不传 i+1）   │（同左侧，加去重即可）
 *   ─────────────────┴─────────────────────┴──────────────────────────
 *   口诀：排列用 used+从0，组合用 start+从start；可复选传 i，不可复选传 i+1。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【决策树可视化（帮助理解三类问题的本质差异）】
 *
 *   排列树（nums=[1,2,3]）：每层从 0 开始，用 used 防重，叶子收集结果
 *              []
 *         /    |    \
 *       [1]   [2]   [3]        ← 每层都可以选所有未用元素
 *      /  \   / \   / \
 *   [1,2][1,3]...            ← 叶子（长度==n）时收集
 *
 *   子集树（nums=[1,2,3]）：每层从 start 开始，每个节点都收集结果
 *     [] ←收集
 *     ├─[1] ←收集
 *     │  ├─[1,2] ←收集
 *     │  │  └─[1,2,3] ←收集
 *     │  └─[1,3] ←收集
 *     ├─[2] ←收集
 *     │  └─[2,3] ←收集
 *     └─[3] ←收集
 *
 *   组合树（n=4,k=2）：每层从 start 开始，只在 len(track)==k 时收集
 *     (start=1): [1,2] [1,3] [1,4] [2,3] [2,4] [3,4]
 *
 *   核心差异：子集在「每个节点」收集，排列/组合只在「叶子/满足条件节点」收集。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【去重剪枝原理深析：同层去重 vs 树枝去重】
 *
 *   场景：nums=[1,1,2]（已排序），求所有不重复子集/排列。
 *
 *   ★ 同层去重（组合/子集，i>start && nums[i]==nums[i-1]）
 *
 *     子集树同层展开：
 *       start=0: [1a, 1b, 2]  → 选 1a 展开一棵树；选 1b 与 1a 完全相同 → 跳过！
 *       条件 i>start 保证「只在同一层跳过」，不影响不同层中的同值元素。
 *       i>0 是错的：会把不同层的合法分支也剪掉。必须用 i>start。
 *
 *   ★ 树枝去重（排列，nums[i]==nums[i-1] && !used[i-1]）
 *
 *     排列树中，相同元素有多种出现顺序：
 *       1a先用、1b后用  → [1a,1b,2]
 *       1b先用、1a后用  → [1b,1a,2]  ← 与上面等价，需剪掉
 *     !used[i-1] 表示「i-1 不在当前路径中」，说明此刻是想先选 i 跳过 i-1，
 *     强制规定：相同元素必须按下标从小到大选取（先选 i-1 再选 i），
 *     保证同值元素的相对顺序固定，消灭重复排列。
 *
 *   ⚠️ 两种去重绝对不能混用：
 *      组合/子集无 used 数组（无法用 !used[i-1]），
 *      排列无 start 参数（无法用 i>start），
 *      用错会导致「过度剪枝（缺答案）」或「剪枝不足（有重复）」。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【结果收集时机对比】
 *
 *   题型          收集时机              收集位置
 *   ──────────────────────────────────────────────────────────────────────
 *   全排列        len(track)==len(nums) 函数顶部的 if 判断内
 *   组合          len(track)==k         函数顶部的 if 判断内
 *   组合总和      trackSum==target      函数顶部的 if 判断内
 *   子集          每次进入函数          函数顶部，if 之前（无条件收集）
 *
 *   子集的「无条件收集」对应决策树的「每个节点都是答案」，
 *   排列/组合的「条件收集」对应「只有叶子节点才是答案」。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【经典扩展题型（超出当前代码，结合知识补充）】
 *
 *   N 皇后：        二维棋盘回溯，额外维护 col/diag1/diag2 三个冲突集合剪枝
 *   解数独：        对每个空格枚举1-9，维护行/列/宫格三维冲突集合，找到第一解即返回
 *   括号生成：      维护 open/close 计数器，open<n 可加 '('，close<open 可加 ')'
 *   单词搜索(DFS)：  二维网格回溯，visited 矩阵防止重复访问同一格，四方向扩展
 *   分割回文串：    回溯枚举切割点，结合 DP 或双指针预处理回文判断加速
 *   复原 IP 地址：  回溯枚举切割位置（每段1-3位），维护段数计数器作为结束条件
 *
 *   时间复杂度参考：
 *     全排列  O(n! · n)       ← n! 个叶子，每个叶子复制路径 O(n)
 *     组合    O(C(n,k) · k)   ← C(n,k) 个叶子，复制路径 O(k)
 *     子集    O(2^n · n)      ← 2^n 个子集，复制路径最长 O(n)
 *     N皇后   O(n! · n)       ← 实际因剪枝远小于此
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次手写回溯后对照检查
 *
 *     ✅ 收集结果时用 copy(temp, track) 深拷贝，而非直接 append(res, track)？
 *     ✅ 做选择（append + used=true）与撤销选择（截断 + used=false）完全对称？
 *     ✅ 有重复元素的题，调用前先 sort.Ints(nums)，去重条件才能正确生效？
 *     ✅ Go 中 track/res 用指针传递（*[]int / *[][]int），防止 append 换底层数组后失效？
 *     ✅ 组合/子集去重用 i>start（同层），排列去重用 !used[i-1]（树枝），没混用？
 *     ✅ 可复选（combinationSum）传 i，不可复选（combinationSum2）传 i+1，没搞反？
 *     ✅ 维护了 trackSum 的题，回溯时 trackSum-=nums[i] 有没有漏写？
 *     ✅ 子集题在函数顶部「无条件收集」，排列/组合在「满足条件时收集」，没搞反？
 *     ✅ N 皇后 / 数独等二维题，四方向扩展时有没有检查边界（防数组越界）？
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *
 *     1. permute(nums []int) [][]int                           // 全排列，元素无重不可复选
 *     2. permuteUnique(nums []int) [][]int                     // 全排列II，元素可重不可复选
 *     3. combine(n, k int) [][]int                             // 组合，元素无重不可复选
 *     4. subsetsWithDup(nums []int) [][]int                    // 子集II，元素可重不可复选
 *     5. combinationSum(candidates []int, target int) [][]int  // 组合总和，元素无重可复选
 *     6. combinationSum2(candidates []int, target int) [][]int // 组合总和II，元素可重不可复选
 * ============================================================================
 */

package main

import (
	"fmt"
	"sort"
)

// 回溯算法实现：
// 🌟技巧1：决策树三要素技巧 - 解题前先明确：①路径track（已做选择）②选择列表（当前层可选项）③结束条件（收集结果的时机）；三要素清晰后框架直接套用，无需另想
// 🌟技巧2：排列 vs 组合/子集 结构区别 - 排列用 used 标记索引 + 每层从0遍历（顺序有关）；组合/子集用 start 参数 + 每层从start遍历（顺序无关）；两者核心区别是「选择列表的起点」
// 🌟技巧3：可复选 vs 不可复选 - 可复选（combinationSum）递归传 i（允许再选自身）；不可复选传 i+1（跳过自身）；一个参数之差决定是否可重复选取
// 🌟技巧4：同层去重技巧（i>start） - 含重复元素的组合/子集，先排序后用 i>start && nums[i]==nums[i-1] 跳过同层相同分支；i>0 是错的，会跨层剪枝导致缺答案
// 🌟技巧5：树枝去重技巧（!used[i-1]） - 含重复元素的排列，先排序后用 nums[i]==nums[i-1] && !used[i-1] 剪枝；强制相同元素按下标从小到大使用，固定相对顺序消灭重复排列
// 🌟技巧6：子集无条件收集技巧 - 子集在函数顶部「无条件收集」（每个节点都是答案）；排列/组合在满足长度或总和条件时才收集（只有叶子节点是答案）；收集时机搞反会缺答案
// 🌟技巧7：深拷贝 + 对称撤销技巧 - 收集用 copy(temp,track) 深拷贝防引用污染；做选择与撤销必须完全对称（append↔截断，used=true↔false，sum+=↔sum-=）；Go中用指针传递*[]int防append换底层数组

// 全排列 - 元素无重不可复选
func permute(nums []int) [][]int {
	// 给定一个不含重复数字的数组nums，返回其「所有可能的全排列」

	// 1、路径：已经做出的选择
	// 2、选择列表：还可以选择的列表
	// 3、结束条件：到达决策底层，无法再做选择的条件

	track := make([]int, 0, len(nums))
	res := make([][]int, 0)
	used := make(map[int]bool)

	backtrack(nums, used, &track, &res)
	return res
}

func backtrack(nums []int, used map[int]bool, track *[]int, res *[][]int) {
	if len(*track) == len(nums) {
		temp := make([]int, len(*track))
		copy(temp, *track)
		*res = append(*res, temp)
		return
	}

	// 排列问题下标从0开始做选择
	for i := 0; i < len(nums); i++ {
		if used[i] {
			continue
		}
		*track = append(*track, nums[i])
		used[i] = true

		backtrack(nums, used, track, res)

		used[i] = false
		*track = (*track)[:len(*track)-1]
	}
}

// 全排列II - 元素可重不可复选
func permuteUnique(nums []int) [][]int {
	// 给定一个可包含重复数字的序列nums，按任意顺序返回所有不重复的全排列。

	used := make(map[int]bool, len(nums))
	var res [][]int
	var track []int

	sort.Ints(nums)

	backtrack4(nums, used, &track, &res)
	return res
}

func backtrack4(nums []int, used map[int]bool, track *[]int, res *[][]int) {
	if len(*track) == len(nums) {
		temp := make([]int, len(*track))
		copy(temp, *track)
		*res = append(*res, temp)
		return
	}

	for i := 0; i < len(nums); i++ {
		if used[i] {
			continue
		}
		// 新添加的剪枝逻辑，固定相同的元素在排列中的相对位置
		if i > 0 && nums[i] == nums[i-1] && !used[i-1] {
			continue
		}

		*track = append(*track, nums[i])
		used[i] = true
		backtrack4(nums, used, track, res)
		used[i] = false
		*track = (*track)[:len(*track)-1]
	}
}

// 组合 - 元素无重不可复选
func combine(n, k int) [][]int {
	// 给定两个整数n和k，返回范围[1, n]中所有可能的k个数的组合

	var track []int
	var res [][]int

	backtrack1(n, k, 1, &track, &res)
	return res
}

func backtrack1(n, k, start int, track *[]int, res *[][]int) {
	if len(*track) == k {
		temp := make([]int, len(*track))
		copy(temp, *track)
		*res = append(*res, temp)
		return
	}

	//组合/子集问题下标从start开始做选择
	for i := start; i <= n; i++ {
		*track = append(*track, i)
		backtrack1(n, k, i+1, track, res)
		*track = (*track)[:len(*track)-1]
	}
}

// 子集II - 元素可重不可复选
func subsetsWithDup(nums []int) [][]int {
	// 给你一个整数数组nums ，其中可能包含重复元素，请你返回该数组所有可能的子集（幂集），解集「不能」包含重复的子集。

	var track []int
	var res [][]int

	sort.Ints(nums)

	backtrack2(nums, 0, &track, &res)
	return res
}

func backtrack2(nums []int, start int, track *[]int, res *[][]int) {
	temp := make([]int, len(*track))
	copy(temp, *track)
	*res = append(*res, temp)

	for i := start; i < len(nums); i++ {
		// 剪枝逻辑，值相同的相邻树枝，只遍历第一条
		if i > start && nums[i] == nums[i-1] {
			continue
		}

		*track = append(*track, nums[i])
		backtrack2(nums, i+1, track, res)
		*track = (*track)[:len(*track)-1]
	}
}

// 组合总和 - 元素无重可复选
func combinationSum(candidates []int, target int) [][]int {
	// 给你一个无重复元素的整数数组「candidates」和一个目标整数「target」，找出candidates中可以使数字和为目标数target的所有不同组合。
	// candidates中的同一个数字可以无限制重复被选取。如果至少一个数字的被选数量不同，则两种组合是不同的。

	var res [][]int
	var track []int
	targetSum := 0

	backtrack5(candidates, 0, targetSum, target, &track, &res)
	return res
}

func backtrack5(nums []int, start, targetSum, target int, track *[]int, res *[][]int) {
	if targetSum == target {
		temp := make([]int, len(*track))
		copy(temp, *track)
		*res = append(*res, temp)
		return
	}
	if targetSum > target {
		return
	}

	for i := start; i < len(nums); i++ {
		*track = append(*track, nums[i])
		targetSum += nums[i]

		// ‼️重点在在索引下标「i」，可复选不+1
		backtrack5(nums, i, targetSum, target, track, res)

		targetSum -= nums[i]
		*track = (*track)[:len(*track)-1]
	}
}

// 组合总和II - 元素可重不可复选
func combinationSum2(candidates []int, target int) [][]int {
	// 给定一个候选人编号的集合candidates和一个目标数target，找出candidates中所有可以使数字和为target的组合
	// candidates中的每个数字在每个组合中只能使用「一次」

	var track []int
	var res [][]int
	trackSum := 0

	sort.Ints(candidates)

	backtrack3(candidates, 0, target, trackSum, &track, &res)
	return res
}

func backtrack3(nums []int, start, target, trackSum int, track *[]int, res *[][]int) {
	if trackSum == target {
		temp := make([]int, len(*track))
		copy(temp, *track)
		*res = append(*res, temp)
		return
	}
	if trackSum > target { // 提前剪枝
		return
	}

	for i := start; i < len(nums); i++ {
		if i > start && nums[i] == nums[i-1] {
			continue
		}

		*track = append(*track, nums[i])
		trackSum += nums[i]

		backtrack3(nums, i+1, target, trackSum, track, res)

		trackSum -= nums[i]
		*track = (*track)[:len(*track)-1]
	}
}

func main() {
	//fmt.Println("全排列:", permute([]int{1, 2, 3}))
	//fmt.Println("全排列II:", permuteUnique([]int{1, 1, 2}))
	//fmt.Println("全排列II:", permuteUnique([]int{1, 2, 3}))

	//fmt.Println("组合:", combine(4, 2))
	//fmt.Println("字集2:", subsetsWithDup([]int{1, 1, 2, 2}))

	fmt.Println("组合总和I:", combinationSum([]int{2, 3, 5}, 8))
	//fmt.Println("组合总和II:", combinationSum2([]int{10, 1, 2, 7, 6, 1, 5}, 8))
}
