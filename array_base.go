/*
 * ============================================================================
 *                   📘 数组基础算法全集 · 核心记忆框架
 * ============================================================================
 * 【一句话理解本文件】
 *
 *   本文件涵盖三大数组经典技巧：N数之和、前缀和/积、差分数组。
 *   核心思想：通过「预处理」将暴力 O(n²) 或 O(n³) 降为 O(n) 或 O(n·logn)。
 *
 * ════════════════════════════════════════════════════════════════════════════
 *                        第一部分：N 数之和
 * ════════════════════════════════════════════════════════════════════════════
 *
 * 【N数之和的通用框架】
 *
 *   所有 nSum 问题都可以递归分解：
 *     nSum → 枚举第一个数 + (n-1)Sum
 *     最终递归到 2Sum，用对撞指针 O(n) 解决。
 *
 *   总体时间复杂度：O(n^(N-1))，其中排序 O(n·logn) 被覆盖。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【去重的核心手法】
 *
 *   前提：数组必须先排序（相同元素相邻）。
 *
 *   去重手法：在确定某个数后，用 for 循环跳过与该数相同的后续元素。
 *
 *   ┌────────────────────────────────────────────────────────────────────────┐
 *   │  两处去重位置（以 twoSum 为例）：                                       │
 *   │                                                                        │
 *   │  ① 找到答案后，lo 和 hi 都要跳过重复：                                  │
 *   │     for lo < hi && nums[lo] == left { lo++ }                           │
 *   │     for lo < hi && nums[hi] == right { hi-- }                          │
 *   │                                                                        │
 *   │  ② 未找到答案时，移动指针也要跳过重复：                                  │
 *   │     sum < target → for lo < hi && nums[lo] == left { lo++ }           │
 *   │     sum > target → for lo < hi && nums[hi] == right { hi-- }          │
 *   │                                                                        │
 *   │  ③ nSum 枚举第一个数时，也要跳过重复：                                   │
 *   │     for i < len(nums)-1 && nums[i] == nums[i+1] { i++ }              │
 *   └────────────────────────────────────────────────────────────────────────┘
 *
 *   ⚠️ 易错点1：去重时必须先保存当前值（left, right），再用 for 跳过
 *      如果直接写 for nums[lo] == nums[lo+1] { lo++ }，当 lo 和 hi 相邻时
 *      可能越界，且逻辑更难正确处理边界。先存值再比较更安全。
 *
 *   ⚠️ 易错点2：nSum 枚举第一个数的去重位置在循环体末尾
 *      必须先处理当前 nums[i]，再跳过后续相同的 nums[i+1]。
 *      如果把去重放在循环体开头，会跳过第一个合法的 nums[i]。
 *
 *   ⚠️ 易错点3：去重循环中必须保持 lo < hi 的条件
 *      跳过重复元素时如果不检查 lo < hi，lo 可能越过 hi 导致越界或重复计算。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【nSumTarget 递归框架要点】
 *
 *   递归函数签名：_nSumTarget(nums []int, n int, start int, target int)
 *
 *   - n：当前需要几个数相加
 *   - start：搜索起点（避免重复使用前面的元素）
 *   - target：每层递归时 target 减去当前选定的数
 *
 *   base case（n==2）：标准对撞指针 twoSum
 *   递归步骤：枚举 nums[i]，递归 _nSumTarget(nums, n-1, i+1, target-nums[i])
 *
 *   ⚠️ 易错点4：递归调用的 start 参数是 i+1，不是 start+1
 *      每次选定 nums[i] 后，下一层只能从 i+1 开始搜索，避免重复使用同一元素。
 *
 * ════════════════════════════════════════════════════════════════════════════
 *                        第二部分：前缀和
 * ════════════════════════════════════════════════════════════════════════════
 *
 * 【一句话理解前缀和】
 *
 *   前缀和 = 用 O(n) 预处理，换取 O(1) 的区间和查询。
 *   preSum[i] = nums[0] + nums[1] + ... + nums[i-1]
 *   区间和 sum(l, r) = preSum[r+1] - preSum[l]
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【一维前缀和】
 *
 *   构建：preSum[0] = 0, preSum[i] = preSum[i-1] + nums[i-1]
 *   查询：sum(left, right) = preSum[right+1] - preSum[left]
 *
 *   为什么 preSum 长度是 n+1？
 *     preSum[0] = 0 代表「空前缀」，这样 sum(0, r) = preSum[r+1] - preSum[0]
 *     不需要对 left==0 做特殊处理。
 *
 *   ⚠️ 易错点5：preSum 下标和 nums 下标差 1
 *      preSum[i] 对应 nums[0..i-1] 的和，不是 nums[0..i] 的和。
 *      查询 nums[left..right] 的和 = preSum[right+1] - preSum[left]。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【二维前缀和】
 *
 *   构建：preSum[i][j] = preSum[i-1][j] + preSum[i][j-1]
 *                       - preSum[i-1][j-1] + matrix[i-1][j-1]
 *
 *   查询（左上角(x1,y1)到右下角(x2,y2)）：
 *     sum = preSum[x2+1][y2+1] - preSum[x1][y2+1]
 *         - preSum[x2+1][y1] + preSum[x1][y1]
 *
 *   记忆口诀：「大减两边加左上角」（容斥原理）
 *
 *   ⚠️ 易错点6：二维前缀和的 +1 偏移容易搞混
 *      构建时 matrix[i-1][j-1] 对应 preSum[i][j]（因为 preSum 多了一行一列空行列）。
 *      查询时坐标要 +1 转换到 preSum 的坐标系。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【前缀积（prefix product）】
 *
 *   与前缀和类似，preProduct[i] = nums[0] * nums[1] * ... * nums[i-1]
 *   区间积 product(l, r) = preProduct[r+1] / preProduct[l]
 *
 *   特殊处理「含零」的情况（ProductOfNumbers）：
 *     遇到 0 时重置 preProduct 为 [1]，因为含 0 的任何区间积都是 0。
 *     查询时若 k > len(preProduct)-1，说明最近 k 个数中有 0，直接返回 0。
 *
 *   「除自身以外的乘积」（productExceptSelf）：
 *     不能用除法（可能有 0），所以用 prefix[i] * suffix[i] 的思路：
 *     prefix[i] = nums[0] * ... * nums[i]（从左到右的前缀积）
 *     suffix[i] = nums[i] * ... * nums[n-1]（从右到左的后缀积）
 *     answer[i] = prefix[i-1] * suffix[i+1]
 *
 *   ⚠️ 易错点7：ProductOfNumbers 遇 0 必须重置整个 preProduct
 *      不能只存一个 0 进去，否则后续除法会出错（除以0）。
 *      重置后长度变短，查询时通过长度判断是否覆盖到了 0。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【前缀和 + 哈希表 —— 子数组问题的终极模式】
 *
 *   核心转化：子数组 [i, j] 的和 = preSum[j+1] - preSum[i]
 *   如果要找「和为 k 的子数组」，等价于找满足 preSum[j] - preSum[i] = k 的配对，
 *   即找 preSum[j] - k 是否在之前出现过 → 用哈希表记录前缀和出现的次数/位置。
 *
 *   ┌─────────────────────────────────────────────────────────────────────┐
 *   │ 变体1：和为K的子数组个数（subarraySum）                              │
 *   │   哈希表记录：preSum 值 → 出现次数                                   │
 *   │   查找：mapSumToCounter[preSum[i] - k] 有多少个                     │
 *   │                                                                     │
 *   │ 变体2：和为0的最长子数组（findMaxLength）                             │
 *   │   将 0 视为 -1，问题变为「和为 0 的最长子数组」                       │
 *   │   哈希表记录：preSum 值 → 第一次出现的下标（越早越好→越长）            │
 *   │   查找：mapSumToIdx[preSum[i]] 是否存在，存在则 i - idx 是长度        │
 *   │                                                                     │
 *   │ 变体3：和为K倍数的子数组（checkSubarraySum / subarraysDivByK）        │
 *   │   preSum[j] - preSum[i] 是 k 的倍数                                 │
 *   │   ⟺ preSum[j] % k == preSum[i] % k                                │
 *   │   哈希表记录：preSum 的余数 → 第一次出现的下标 / 出现次数             │
 *   │                                                                     │
 *   │ 变体4：子数组和>0的最长（longestWPI）                                 │
 *   │   preSum 每次只变化 ±1 → 具有连续性                                  │
 *   │   若 preSum[i] > 0 → [0, i-1] 整段都满足                            │
 *   │   否则找 mapSumToIdx[preSum[i]-1] → 最早的 preSum 恰好比当前小 1     │
 *   └─────────────────────────────────────────────────────────────────────┘
 *
 *   ⚠️ 易错点8：哈希表中只记录「第一次出现」的下标（求最长时）
 *      因为要最长子数组，所以同一前缀和第一次出现的位置越靠前越好。
 *      后续再出现相同前缀和时不更新哈希表。
 *
 *   ⚠️ 易错点9：Go 的 % 运算对负数结果为负，必须手动修正
 *      数学上 -1 mod 5 = 4，但 Go 中 -1 % 5 = -1。
 *      修正：remainder := preSum[i] % k; if remainder < 0 { remainder += k }
 *
 *   ⚠️ 易错点10：checkSubarraySum 要求子数组长度 >= 2
 *      找到相同余数后还要检查 i - idx >= 2，不能只有一个元素。
 *
 *   ⚠️ 易错点11：longestWPI 中查找 preSum[i]-1 而非 preSum[i]-k
 *      因为 preSum 每次只变 ±1（每小时要么 +1 要么 -1），
 *      preSum[i]-1 就是比当前值「刚好小 1」的历史前缀和。
 *      由于连续性，它一定是最早出现的比 preSum[i] 小的值（无需遍历更小的值）。
 *
 * ════════════════════════════════════════════════════════════════════════════
 *                        第三部分：差分数组
 * ════════════════════════════════════════════════════════════════════════════
 *
 * 【一句话理解差分数组】
 *
 *   差分数组 = 前缀和的逆运算。
 *   用于「频繁对区间 [i, j] 整体加减同一个值」的场景，
 *   将每次 O(n) 的区间更新降为 O(1)。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【差分数组的构建与还原】
 *
 *   构建：diff[0] = nums[0], diff[i] = nums[i] - nums[i-1]  (i >= 1)
 *   还原：对 diff 做前缀和即还原原数组 → res[i] = res[i-1] + diff[i]
 *
 *   性质：diff 数组的前缀和 = 原数组。原数组的差分 = diff 数组。互为逆运算。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【区间加减操作】
 *
 *   对原数组 nums[i..j] 全部加上 val：
 *     diff[i] += val     （从 i 开始影响后续所有元素）
 *     diff[j+1] -= val   （从 j+1 开始取消影响）
 *
 *   批量操作完成后，对 diff 做一次前缀和即得到最终结果。
 *
 *   ⚠️ 易错点12：diff[j+1] -= val 必须检查 j+1 < len(diff)
 *      如果 j 是数组最后一个下标，j+1 越界，此时不需要 -= val
 *      （因为后面没有元素需要取消影响了）。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【差分数组的经典应用】
 *
 *   ┌─────────────────────────────────────────────────────────────────────┐
 *   │ 航班预订统计（corpFlightBookings）：                                  │
 *   │   直接对 diff[start-1] += seats, diff[end] -= seats（注意1-indexed） │
 *   │                                                                     │
 *   │ 拼车（carPooling）：                                                 │
 *   │   乘客在第 i 站上车、第 j 站下车 → 影响区间是 [i, j-1]              │
 *   │   （第 j 站下车意味着第 j 站不占座位）                               │
 *   │   diff[i] += val, diff[j] -= val（不是 j+1！）                      │
 *   │   还原后检查每一站是否超过 capacity                                   │
 *   └─────────────────────────────────────────────────────────────────────┘
 *
 *   ⚠️ 易错点13：拼车的区间是 [上车站, 下车站-1]，不是 [上车站, 下车站]
 *      「第 j 站下车」= 到第 j 站时人已经不在车上了，所以 diff[j] -= val。
 *      如果写成 diff[j+1] -= val，会多算一站的乘客。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【前缀和 vs 差分数组 的选择】
 *
 *   ┌──────────────────────────────────────────────────────────────────────┐
 *   │         前缀和                          差分数组                      │
 *   │  ───────────────────────────     ─────────────────────────────       │
 *   │  场景：频繁「查询」区间和          场景：频繁「修改」区间值            │
 *   │  预处理 O(n)，查询 O(1)           修改 O(1)，还原 O(n)               │
 *   │  原数组不变，多次查询              多次修改，最后一次性还原             │
 *   └──────────────────────────────────────────────────────────────────────┘
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次手写前缀和/差分后对照检查
 *
 *     ✅ 前缀和数组长度是 n+1，preSum[0] = 0？
 *     ✅ 查询区间和时下标偏移是否正确？sum(l,r) = preSum[r+1] - preSum[l]？
 *     ✅ 二维前缀和构建用容斥：+上 +左 -左上 +当前？
 *     ✅ 二维前缀和查询用容斥：+右下 -上 -左 +左上？
 *     ✅ 前缀和+哈希表：是否先查哈希表再更新（避免用当前值匹配自己）？
 *     ✅ 取模时是否处理了负余数？
 *     ✅ 差分操作：diff[j+1] 是否越界检查？
 *     ✅ 拼车场景：下车站不占座，影响区间是 [i, j-1]？
 *     ✅ N数之和：排序了吗？去重的 for 循环中有 lo < hi 保护吗？
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *
 *     N数之和：
 *       1. twoSumTarget(nums []int, target int) [][]int       // 两数之和（去重）
 *       2. threeSumTarget(nums []int, target int) [][]int     // 三数之和（去重）
 *       3. nSumTarget(nums []int, target int) [][]int         // N数之和（通用递归）
 *
 *     一维前缀和：
 *       4. NumArray / SumRange(left, right int) int           // 区间和查询
 *       5. pivotIndex(nums []int) int                         // 寻找中心下标
 *
 *     二维前缀和：
 *       6. NumMatrix / SumRegion(x1,y1,x2,y2 int) int        // 二维区间和
 *       7. matrixBlockSum(mat [][]int, k int) [][]int         // 矩阵区域和
 *
 *     前缀积：
 *       8. productExceptSelf(nums []int64) []int64            // 除自身以外乘积
 *       9. ProductOfNumbers / GetProduct(k int) int           // 最后K个数乘积
 *
 *     前缀和 + 哈希表：
 *      10. findMaxLength(nums []int) int                      // 连续数组（0/1最长等量）
 *      11. checkSubarraySum(nums []int, k int) bool           // 子数组和为K倍数
 *      12. subarraySum(nums []int, k int) int                 // 和为K的子数组个数
 *      13. longestWPI(hours []int) int                        // 表现良好最长时间段
 *      14. subarraysDivByK(nums []int, k int) int             // 和可被K整除的子数组
 *
 *     差分数组：
 *      15. DiffArray / Add(i, j, val int)                     // 差分数组（通用）
 *      16. corpFlightBookings(bookings [][]int, n int) []int  // 航班预订统计
 *      17. carPooling(trips [][]int, capacity int) bool       // 拼车
 * ============================================================================
 */

package main

import (
	"fmt"
	"sort"
)

// 两数之和，返回和为target的两个数，注意不能返回重复数对儿
func twoSumTarget(nums []int, target int) [][]int {
	// 解题核心是：先排序，再++和--时使用for循环跳过相同元素
	sort.Ints(nums)
	var lo, hi = 0, len(nums) - 1
	var res [][]int

	for lo < hi {
		left, right := nums[lo], nums[hi]
		sum := nums[lo] + nums[hi]

		if sum < target {
			for lo < hi && nums[lo] == left {
				lo++
			}
		} else if sum > target {
			for lo < hi && nums[hi] == right {
				hi--
			}
		} else {
			res = append(res, []int{left, right})
			for lo < hi && nums[lo] == left {
				lo++
			}
			for lo < hi && nums[hi] == right {
				hi--
			}
		}
	}
	return res
}

// 3数之和，返回三元组nums[i]+nums[j]+nums[k]=target且i!=j!=k的所有元素对儿
func threeSumTarget(nums []int, target int) [][]int {
	sort.Ints(nums)
	var result [][]int

	twoSumTarget := func(nums []int, start, target int) [][]int {
		sort.Ints(nums)
		var lo, hi = start, len(nums) - 1
		var res [][]int
		for lo < hi {
			left, right := nums[lo], nums[hi]
			sum := nums[lo] + nums[hi]

			if sum < target {
				for lo < hi && nums[lo] == left {
					lo++
				}
			} else if sum > target {
				for lo < hi && nums[hi] == right {
					hi--
				}
			} else {
				res = append(res, []int{left, right})
				for lo < hi && nums[lo] == left {
					lo++
				}
				for lo < hi && nums[hi] == right {
					hi--
				}
			}
		}
		return res
	}

	// 穷举 threeSum 的第一个数
	for i := 0; i < len(nums); i++ {
		tuples := twoSumTarget(nums, i+1, target-nums[i])
		for _, tuple := range tuples {
			tuple = append(tuple, nums[i])
			result = append(result, tuple)
		}
		// 跳过第一个数字重复的情况，否则会出现重复结果
		for i < len(nums)-1 && nums[i] == nums[i+1] {
			i++
		}
	}

	return result
}

// n数之和，返回n元组nums[i]+nums[j]...nums[n]=target且i!=j...n的所有元素对儿
func nSumTarget(nums []int, target int) [][]int {
	sort.Ints(nums)
	return _nSumTarget(nums, 4, 0, target)
}

func _nSumTarget(nums []int, n int, start int, target int) [][]int {
	// nums 必须是升序数组
	var result [][]int
	if n == 2 {
		lo, hi := start, len(nums)-1
		for lo < hi {
			left, right := nums[lo], nums[hi]
			sum := nums[lo] + nums[hi]
			if sum < target {
				for lo < hi && nums[lo] == left {
					lo++
				}
			} else if sum > target {
				for lo < hi && nums[hi] == right {
					hi--
				}
			} else {
				result = append(result, []int{left, right})
				for lo < hi && nums[lo] == left {
					lo++
				}
				for lo < hi && nums[hi] == right {
					hi--
				}
			}
		}
	} else {
		for i := start; i < len(nums); i++ {
			res := _nSumTarget(nums, n-1, i+1, target-nums[i])
			for _, pair := range res {
				pair = append(pair, nums[i])
				result = append(result, pair)
			}
			// ⚠️ 注意跳过相同元素
			for i < len(nums)-1 && nums[i] == nums[i+1] {
				i++
			}
		}
	}
	return result
}

// NumArray 一维区域和检索-数组不可变（一维前缀和） https://leetcode.cn/problems/range-sum-query-immutable/description/
type NumArray struct {
	PreSum []int
}

func NewNumArray(nums []int) NumArray {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] += preSum[i-1] + nums[i-1]
	}
	return NumArray{
		PreSum: preSum,
	}
}

func (n *NumArray) SumRange(left int, right int) int {
	return n.PreSum[right+1] - n.PreSum[left]
}

// NumMatrix 二维区域和检索-矩阵不可变（二维前缀和） https://leetcode.cn/problems/range-sum-query-2d-immutable/
type NumMatrix struct {
	PreSum [][]int
}

func NewNumMatrix(matrix [][]int) NumMatrix {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return NumMatrix{}
	}

	preSum := make([][]int, len(matrix)+1)
	for i := range preSum {
		preSum[i] = make([]int, len(matrix[0])+1)
	}

	for i := 1; i < len(preSum); i++ {
		for j := 1; j < len(preSum[0]); j++ {
			// [ 3, 0, 1, 4, 2 ]
			// [ 5, 6, 3, 2, 1 ]
			// [ 1, 2, 0, 1, 5 ]
			// [ 4, 1, 0, 1, 7 ]
			// [ 1, 0, 3, 0, 5 ]
			preSum[i][j] = preSum[i-1][j] + preSum[i][j-1] - preSum[i-1][j-1] + matrix[i-1][j-1]
		}
	}

	return NumMatrix{PreSum: preSum}
}

func (n *NumMatrix) SumRegion(x1, y1 int, x2, y2 int) int {
	return n.PreSum[x2+1][y2+1] - n.PreSum[x1][y2+1] - n.PreSum[x2+1][y1] + n.PreSum[x1][y1]
}

// 矩阵区域和 https://leetcode.cn/problems/matrix-block-sum/description/
func matrixBlockSum(mat [][]int, k int) [][]int {
	preSum := make([][]int, len(mat)+1)
	for i := range preSum {
		preSum[i] = make([]int, len(mat[0])+1)
	}

	for i := 1; i < len(preSum); i++ {
		for j := 1; j < len(preSum[0]); j++ {
			preSum[i][j] = preSum[i-1][j] + preSum[i][j-1] - preSum[i-1][j-1] + mat[i-1][j-1]
		}
	}

	answer := make([][]int, len(mat))
	for i := range answer {
		answer[i] = make([]int, len(mat[0]))
	}

	for i := 0; i < len(answer); i++ {
		for j := 0; j < len(answer[0]); j++ {
			// i - k <= r <= i + k, j - k <= c <= j + k 且 (r, c) 在矩阵内。
			x1 := max(i-k, 0)
			y1 := max(j-k, 0)
			x2 := min(i+k, len(answer)-1)
			y2 := min(j+k, len(answer[0])-1)
			answer[i][j] = preSum[x2+1][y2+1] - preSum[x1][y2+1] - preSum[x2+1][y1] + preSum[x1][y1]
		}
	}
	return answer
}

// 寻找数组的中心下标 https://leetcode.cn/problems/find-pivot-index/
func pivotIndex(nums []int) int {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] += preSum[i-1] + nums[i-1]
	}

	for i := 1; i < len(preSum); i++ {
		// 判断：左侧之和 与 右侧之和
		if preSum[i-1] == preSum[len(nums)]-preSum[i] {
			return i - 1
		}
	}
	return -1
}

// 除了自身以外数组的乘积 https://leetcode.cn/problems/product-of-array-except-self/
func productExceptSelf(nums []int64) []int64 {
	n := len(nums)

	// prefix[i] 表示num[0...i]的乘积
	prefix := make([]int64, n)
	prefix[0] = nums[0]
	for i := 1; i < n; i++ {
		prefix[i] = prefix[i-1] * nums[i]
	}

	// prefix[i] 表示num[i...n-1]的乘积
	suffix := make([]int64, n)
	suffix[n-1] = nums[n-1]
	for i := n - 2; i >= 0; i-- {
		suffix[i] = suffix[i+1] * nums[i]
	}

	answer := make([]int64, n)
	answer[0] = suffix[1]
	answer[n-1] = prefix[n-2]

	for i := 1; i <= n-2; i++ {
		answer[i] = prefix[i-1] * suffix[i+1]
	}
	return answer
}

// ProductOfNumbers 最后K个数的乘积 https://leetcode.cn/problems/product-of-the-last-k-numbers/
type ProductOfNumbers struct {
	preProduct []int
}

func NewProductOfNumbers() ProductOfNumbers {
	return ProductOfNumbers{preProduct: []int{1}}
}

func (p *ProductOfNumbers) Add(num int) {
	if num == 0 {
		p.preProduct = []int{1}
		return
	}
	n := len(p.preProduct)
	p.preProduct = append(p.preProduct, p.preProduct[n-1]*num)
}

func (p *ProductOfNumbers) GetProduct(k int) int {
	n := len(p.preProduct)
	// 不足k个元素，是因为最后k个元素存在0，非0情况下:k+1=n
	if k > n-1 {
		return 0
	}
	return p.preProduct[n-1] / p.preProduct[n-1-k]
}

// 连续数组 https://leetcode.cn/problems/contiguous-array/
func findMaxLength(nums []int) int {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		if nums[i-1] == 0 {
			preSum[i] = preSum[i-1] - 1
		} else {
			preSum[i] = preSum[i-1] + 1
		}
	}

	res := 0

	// 子数组和为零，说明preSum[i]==preSum[j]
	mapSumToIdx := make(map[int]int)
	for i := 0; i < len(preSum); i++ {
		if idx, ok := mapSumToIdx[preSum[i]]; !ok {
			mapSumToIdx[preSum[i]] = i
		} else {
			res = max(res, i-idx)
		}
	}
	return res
}

// 连续的子数组和 https://leetcode.cn/problems/continuous-subarray-sum/
func checkSubarraySum(nums []int, k int) bool {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	// 子数组元素总和为k的倍数,说明(preSum[j]-preSum[i])%k==0
	// 推导：preSum[j]%k == preSum[i]%k

	mapSumToIdx := make(map[int]int)
	for i := 0; i < len(preSum); i++ {
		if idx, ok := mapSumToIdx[preSum[i]%k]; ok {
			if i-idx >= 2 {
				return true
			}
		} else {
			mapSumToIdx[preSum[i]%k] = i
		}
	}
	return false
}

// 和为K的子数组 https://leetcode.cn/problems/subarray-sum-equals-k/
func subarraySum(nums []int, k int) int {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	res := 0

	// 该数组中和为k的子数组的个数: preSum[i]-preSum[j] = k
	// preSum[j] = preSum[i]-k
	// preSum[i]-k 是否存在，存在多少次就加多少

	mapSumToCounter := make(map[int]int)
	for i := 0; i < len(preSum); i++ {
		if val, ok := mapSumToCounter[preSum[i]-k]; ok {
			res += val
		}
		mapSumToCounter[preSum[i]]++
	}
	return res
}

// 表现良好的最长时间段 https://leetcode.cn/problems/longest-well-performing-interval/
func longestWPI(hours []int) int {
	preSum := make([]int, len(hours)+1)
	for i := 1; i < len(preSum); i++ {
		if hours[i-1] > 8 {
			preSum[i] = preSum[i-1] + 1
		} else {
			preSum[i] = preSum[i-1] - 1
		}
	}

	// preSum重要特性：每次变化只能是±1（因为每个小时要么+1要么-1），这意味着：
	// 1、前缀和是连续变化的：从 0 开始，每次只能增加或减少 1
	// 2、数值的出现顺序：较小的前缀和数值（如 -4、-5）通常出现在较大的前缀和数值（如 -3、-2）之后。

	res := 0
	mapSumToIdx := make(map[int]int)

	for i := 0; i < len(preSum); i++ {
		// 求最长子数组，那么记录第一次出现就好
		if _, exists := mapSumToIdx[preSum[i]]; !exists {
			mapSumToIdx[preSum[i]] = i
		}

		if preSum[i] > 0 {
			// preSum[i] 为正，说明 hours[0..i-1] 都是「表现良好的时间段」
			res = max(res, i)
		} else {
			// 求子数组和 >0 ，即：preSum[i] - preSum[j] > 0 表示区间 [j, i-1] 的和 > 0
			// preSum[j] < preSum[i] => preSum[j] <= preSum[i]-1
			// preSum[j] 最大值就是：preSum[i]-1，因为preSum的单调性(递减)，所以负数越大越靠前
			// 所以查找 mapSumToIdx[preSum[i]-1] 就是在查找最小的 j
			if j, found := mapSumToIdx[preSum[i]-1]; found {
				res = max(res, i-j)
			}
		}
	}
	return res
}

// 和可被K整除的子数组 https://leetcode.cn/problems/subarray-sums-divisible-by-k/
func subarraysDivByK(nums []int, k int) int {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	// (preSum[i] - preSum[j])%k == 0
	// preSum[i]%k - preSum[j]%k == 0
	// preSum[i]%k == preSum[j]%k

	res := 0
	mapCounter := make(map[int]int)

	for i := 0; i < len(preSum); i++ {

		// 数学上模运算的结果总是在 [0, m-1] 范围内，例如：
		// -1 mod 5 = 4（因为 -1 = (-1)×5 + 4）
		// -3 mod 5 = 2（因为 -3 = (-1)×5 + 2）
		// 但大多数编程语言的 % 运算符不同与此不同，Go中：
		// -1 % 5 = -1（而不是数学上的4）=> -1+5=4
		// -3 % 5 = -3（而不是数学上的2）=> -3+5=2

		// 子数组和可能为负，把负余数调整为正余数
		remainder := preSum[i] % k
		if remainder < 0 {
			remainder += k
		}

		if num, ok := mapCounter[remainder]; !ok {
			mapCounter[remainder] = 1
		} else {
			res += num
			mapCounter[remainder] = num + 1
		}
	}
	return res
}

// DiffArray 差分数组
type DiffArray struct {
	diff []int
}

func NewDiffArray(nums []int) *DiffArray {
	diff := make([]int, len(nums))
	diff[0] = nums[0]
	for i := 1; i < len(nums); i++ {
		diff[i] = nums[i] - nums[i-1]
	}
	return &DiffArray{diff: diff}
}

func (d *DiffArray) Add(i, j, val int) {
	d.diff[i] += val
	if j+1 < len(d.diff) {
		d.diff[j+1] -= val
	}
}
func (d *DiffArray) Result() []int {
	res := make([]int, len(d.diff))
	res[0] = d.diff[0]

	for i := 1; i < len(d.diff); i++ {
		res[i] = res[i-1] + d.diff[i]
	}
	return res
}

// 航班预订统计 https://leetcode.cn/problems/corporate-flight-bookings/
func corpFlightBookings(bookings [][]int, n int) []int {
	diff := make([]int, n)
	for _, booking := range bookings {
		start := booking[0]
		end := booking[1]
		seats := booking[2]

		diff[start-1] += seats
		if end < n {
			diff[end] -= seats
		}
	}
	res := make([]int, n)
	res[0] = diff[0]
	for i := 1; i < n; i++ {
		res[i] = res[i-1] + diff[i]
	}
	return res
}

// 拼车 https://leetcode.cn/problems/car-pooling/description/
func carPooling(trips [][]int, capacity int) bool {
	maxSite := 1001

	diff := make([]int, maxSite)
	for _, trip := range trips {
		val := trip[0]
		i := trip[1] // 第i站上车
		j := trip[2] // 第j站下车

		// 更新区间应为：diff[i,j-1]
		diff[i] += val
		if j < len(diff)-1 {
			diff[j] -= val
		}
	}

	res := make([]int, maxSite)
	res[0] = diff[0]
	// 起点就超载了
	if res[0] > capacity {
		return false
	}

	for i := 1; i < maxSite; i++ {
		res[i] = res[i-1] + diff[i]
		if res[i] > capacity {
			return false
		}
	}
	return true
}

func main() {

	fmt.Println(carPooling([][]int{{2, 1, 5}, {3, 3, 7}}, 4))
	fmt.Println(carPooling([][]int{{2, 1, 5}, {3, 3, 7}}, 5))
	fmt.Println(carPooling([][]int{{9, 0, 1}, {3, 3, 7}}, 4))

	//fmt.Println(corpFlightBookings([][]int{{1, 2, 10}, {2, 3, 20}, {2, 5, 25}}, 5))
	//fmt.Println(corpFlightBookings([][]int{{1, 2, 10}, {2, 2, 15}}, 2))

	//da := NewDiffArray([]int{0, 0, 0, 0, 0, 0, 0, 0})
	//da.Add(1, 3, 1)
	//da.Add(4, 6, -1)
	//da.Add(2, 5, 10)
	//fmt.Println(da.Result())

	//fmt.Println(subarraysDivByK([]int{-1, 2, 9}, 2))
	//fmt.Println(subarraysDivByK([]int{4, 5, 0, -2, -3, 1}, 5))

	//fmt.Println(longestWPI([]int{8, 10, 6, 16, 5}))
	//fmt.Println(longestWPI([]int{6, 6, 6}))
	//fmt.Println(longestWPI([]int{9, 9, 6, 0, 6, 6, 9}))

	//fmt.Println(subarraySum([]int{1, -1, 0}, 0))
	//fmt.Println(subarraySum([]int{1, 2, 3}, 3))
	//fmt.Println(subarraySum([]int{1, 9, 2, 8, 3, 7, 4, 6, 5}, 10))

	//fmt.Println(checkSubarraySum([]int{5, 0, 0, 0}, 3))
	//fmt.Println(checkSubarraySum([]int{0}, 1))
	//fmt.Println(checkSubarraySum([]int{23, 2, 4, 6, 7}, 6))
	//fmt.Println(checkSubarraySum([]int{23, 2, 6, 4, 7}, 13))

	//fmt.Println(findMaxLength([]int{0, 1}))
	//fmt.Println(findMaxLength([]int{0, 1, 0}))
	//fmt.Println(findMaxLength([]int{0, 1, 1, 1, 1, 1, 0, 0, 0}))

	//pp := NewProductOfNumbers()
	//pp.Add(4)
	//pp.Add(3)
	//pp.Add(2)
	//pp.Add(1)
	//fmt.Println(pp.GetProduct(1))
	//fmt.Println(pp.GetProduct(2))
	//fmt.Println(pp.GetProduct(3))
	//fmt.Println(pp.GetProduct(4))

	//fmt.Println(productExceptSelf([]int64{1, 2, 3, 4}))
	//fmt.Println(productExceptSelf([]int64{-1, 1, 0, -3, 3}))

	//fmt.Println(pivotIndex([]int{1, 7, 3, 6, 5, 6}))
	//fmt.Println(pivotIndex([]int{2, 1, -1}))

	//// 输入：mat = [[1,2,3],[4,5,6],[7,8,9]], k = 1
	//// 输出：[[12,21,16],[27,45,33],[24,39,28]]
	//fmt.Println(matrixBlockSum([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, 1))
	//
	////输入：mat = [[1,2,3],[4,5,6],[7,8,9]], k = 2
	////输出：[[45,45,45],[45,45,45],[45,45,45]]
	//fmt.Println(matrixBlockSum([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, 2))
}
