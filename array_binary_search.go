/*
 * ============================================================================
 *                    📘 数组二分搜索 · 核心记忆框架
 * ============================================================================
 * 【六大二分模式】
 *
 *   ① 基础二分：在有序数组中精确查找 target
 *   ② 边界二分：在有序数组中找 target 的左/右边界
 *   ③ 抽象二分：在单调函数 f(x) 上求满足约束的 x 最值（二分答案）
 *   ④ 矩阵搜索：2D 展平为 1D / 右上角淘汰法
 *   ⑤ 峰值搜索：利用局部单调性（比较 mid 与 mid+1）
 *   ⑥ 旋转数组：先定位有序半区，再判断 target 归属
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式一：基础二分 —— 精确查找】
 *
 *   核心框架（闭区间 [left, right]）：
 *     left, right := 0, len(nums)-1
 *     for left <= right {
 *         mid := left + (right-left)/2
 *         if nums[mid] == target { return mid }
 *         else if nums[mid] > target { right = mid - 1 }
 *         else { left = mid + 1 }
 *     }
 *     return -1
 *
 *   ⚠️ 易错点1：mid 的计算必须用 left + (right-left)/2 防止溢出
 *      不能写 (left+right)/2，当 left+right 超过 int 最大值时会溢出。
 *
 *   ⚠️ 易错点2：循环条件是 left <= right（不是 <）
 *      因为搜索区间是闭区间 [left, right]，当 left == right 时区间内还有一个元素，
 *      必须检查。若用 left < right，会漏掉最后一个元素。
 *
 *   ⚠️ 易错点3：搜索区间的收缩必须排除 mid 本身
 *      已确认 nums[mid] != target，所以下一步搜 [left, mid-1] 或 [mid+1, right]。
 *      若写 right = mid 或 left = mid 会导致死循环（left==right==mid 时不收缩）。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式二：左右边界二分 —— 找第一个/最后一个 target】
 *
 *   三种写法的唯一差异在 nums[mid] == target 时的处理：
 *     · 精确查找：return mid             （找到就返回）
 *     · 左边界：  right = mid - 1        （继续向左收缩，逼近第一个）
 *     · 右边界：  left = mid + 1         （继续向右收缩，逼近最后一个）
 *
 *   记忆口诀：
 *     左边界 → "找到了还要往左挤" → right = mid-1 → 结束时 left 停在左边界
 *     右边界 → "找到了还要往右挤" → left = mid+1  → 结束时 right 停在右边界
 *
 *   ⚠️ 易错点4：循环结束后必须做两层检查
 *      第一层：越界检查（left >= len 或 right < 0）
 *      第二层：值检查（nums[left] != target 或 nums[right] != target）
 *      因为即使循环正常结束，left/right 指向的值也可能不是 target（target 不存在时）。
 *
 *   ⚠️ 易错点5：左边界返回 left，右边界返回 right
 *      循环结束时 left = right + 1：
 *        · 左边界：right 一直被"向左压"，循环结束时 left 是第一个 ≥ target 的位置
 *        · 右边界：left 一直被"向右推"，循环结束时 right 是最后一个 ≤ target 的位置
 *
 *   ⚠️ 易错点6：leftBound 去掉值检查后的语义 = "大于等于 target 的最小索引"（插入位置）
 *      这就是 searchInsert 和 numMatchingSubseq 中用到的形态——
 *      不关心 target 是否存在，只关心"第一个 ≥ target 的位置在哪"。
 *      同理 rightBound 去掉值检查 = "小于等于 target 的最大索引"。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式三：抽象二分（二分答案）—— 单调函数求最值】
 *
 *   适用条件：
 *     1. 能抽象出自变量 x、单调函数 f(x)、目标值 target
 *     2. f(x) 关于 x 单调（递增或递减）
 *     3. 求满足 f(x) ≤ target（或 ≥ target）的 x 的最小值/最大值
 *
 *   建模三步走：
 *     Step1：x 是什么？f(x) 是什么？target 是什么？
 *     Step2：x 的范围 [left, right] 是什么？
 *     Step3：求左边界（最小 x）还是右边界（最大 x）？
 *
 *   "连续分段"类问题的统一模型（shipWithinDays / splitArray）：
 *     · x = 每段容量上限
 *     · f(x) = 按容量 x 分段，最少需要几段（f 关于 x 单调递减）
 *     · target = 段数限制
 *     · 求满足 f(x) ≤ target 的最小 x → 左边界
 *     · x 范围：[max(arr), sum(arr)]
 *
 *   ⚠️ 易错点7：x 的下界是 max(arr) 而不是 0 或 1
 *      如果 x < max(arr)，则数组中最大的那个元素单独一个都放不下，f(x) 无意义。
 *      例外：minEatingSpeed 的 x 是速度，下界为 1（每堆可以分多次吃完）。
 *
 *   ⚠️ 易错点8：f(x) 单调递减时的收缩方向
 *      f(x) 递减意味着 x 越大，f(x) 越小：
 *        f(mid) > target → x 太小，需增大 → left = mid + 1
 *        f(mid) < target → x 太大，可减小 → right = mid - 1
 *        f(mid) == target → 还能更小吗？→ right = mid - 1（求左边界/最小 x）
 *      三者合并：f(mid) <= target 时都 right = mid - 1，最终返回 left。
 *
 *   ⚠️ 易错点9：f(x) 中分段计数的 i-- 技巧
 *      当 sum + weights[i] > x 时，weights[i] 放不进当前段，需要：
 *        1. 开启新段（days++）
 *        2. i-- 把当前元素"退回去"，下一轮重新处理
 *        3. sum 归零
 *      如果不 i--，该元素会被跳过，导致段数计算偏小。
 *
 *   ⚠️ 易错点10：f(x) 循环结束后需要 +1 计入最后一段
 *      最后一段不会触发 sum > x（因为循环正常结束），所以 days/num 漏掉了它，
 *      必须在循环外补 days++ 或 num++。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式四：矩阵搜索】
 *
 *   场景A —— 严格有序矩阵（searchMatrix）：
 *     每行递增，且下一行首元素 > 上一行末元素 → 整体可视为一维有序数组
 *     技巧：一维索引 mid 映射到二维：row = mid/n, col = mid%n（n 是列数）
 *
 *   场景B —— 行列各自有序矩阵（searchMatrixII）：
 *     每行递增，每列递增，但行间不严格 → 不能展平
 *     技巧：从右上角出发，利用"向左变小，向下变大"的淘汰法
 *       matrix[i][j] > target → j--（太大，排除该列）
 *       matrix[i][j] < target → i++（太小，排除该行）
 *
 *   ⚠️ 易错点11：展平时除以的是列数 n，不是行数 m
 *      matrix[mid/n][mid%n] 中 n 是列数。若误用 m，行列映射错乱。
 *
 *   ⚠️ 易错点12：searchMatrixII 必须从右上角（或左下角）出发
 *      从左上角出发：向右和向下都是增大，无法判断走哪个方向。
 *      从右上角出发：向左减小，向下增大，每步都能排除一行或一列。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式五：峰值搜索 —— 局部单调性二分】
 *
 *   核心洞察：不需要全局有序，只需局部"方向"信息：
 *     nums[mid] < nums[mid+1] → 右侧有更大值 → 峰在右边 → left = mid + 1
 *     nums[mid] > nums[mid+1] → mid 可能是峰或峰在左边 → right = mid - 1
 *
 *   ⚠️ 易错点13：必须先判断 mid+1 是否越界
 *      if mid+1 == len(nums) { return mid }
 *      当 mid 是最后一个元素时，nums[mid+1] 越界。
 *      此时 mid 就是峰值（题目保证 nums[-1] = nums[n] = -∞）。
 *
 *   ⚠️ 易错点14：峰值搜索中 right = mid-1（不是 right = mid）
 *      当 nums[mid] > nums[mid+1] 时，mid 本身可能是峰，但因为 left <= right
 *      的循环条件，left 最终会追上来"踩到"峰值位置，不会漏掉。
 *      若用 right = mid，在 left == right == mid 时不收缩，死循环。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式六：旋转数组搜索】
 *
 *   核心思路：旋转后数组被"断崖"分成两段有序部分。
 *   每次二分先判断 mid 落在哪一段，再确定 target 在有序的那半区间内还是外。
 *
 *   判断有序半区：
 *     nums[mid] >= nums[left] → 左半 [left, mid] 有序
 *     否则                    → 右半 [mid, right] 有序
 *
 *   在有序半区内定位 target：
 *     左半有序时：target ∈ [nums[left], nums[mid]] → 搜左半，否则搜右半
 *     右半有序时：target ∈ [nums[mid], nums[right]] → 搜右半，否则搜左半
 *
 *   ⚠️ 易错点15：判断有序半区时用 >= 而非 >
 *      nums[mid] >= nums[left]（含等号）。
 *      当 left == mid 时（只有两个元素），左半就是 [left, mid] 这一个元素，仍然有序。
 *      若用 >，当 left==mid 时会错误地认为"左半无序"，进入错误分支。
 *
 *   ⚠️ 易错点16：含重复元素时必须先跳过重复（searchII）
 *      重复元素会导致 nums[left] == nums[mid] == nums[right]，此时无法判断
 *      mid 落在断崖的哪一侧。解决方法：循环开头先收缩 left/right 跳过边界重复。
 *      代价：最坏时间复杂度退化为 O(n)。
 *
 *   ⚠️ 易错点17：findMin 的循环条件是 left < right（不是 <=）且 right = mid（不是 mid-1）
 *      findMin 不是在"找 target"，而是在"缩小候选区间"：
 *        · left < right 保证区间非空，退出时 left == right 指向最小值
 *        · nums[mid] > nums[right] → 最小值在 (mid, right] → left = mid + 1
 *        · nums[mid] <= nums[right] → 最小值在 [left, mid] → right = mid
 *      如果用 right = mid-1，会跳过 mid 本身是最小值的情况。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式七：二分 + 扩展 —— 找 K 个最接近元素】
 *
 *   步骤：
 *     1. 用 leftBound 找到 x 的插入位置 p（第一个 ≥ x 的位置）
 *     2. 以 p 为中心，向两侧扩展，每次选距离 x 更近的那一侧
 *     3. 扩展 k 次后，开区间 (left, right) 内即为结果
 *
 *   ⚠️ 易错点18：初始化用开区间 (p-1, p) 而非闭区间
 *      p 本身可能越界（p == len(arr) 时 arr[p] 不存在）。
 *      用开区间 left=p-1, right=p，区间内元素数 = right-left-1 = 0，天然处理边界。
 *
 *   ⚠️ 易错点19：距离相等时优先扩展左侧（题目要求返回较小值）
 *      代码中 x-arr[left] > arr[right]-x 用严格大于，相等时走 else 即 left--。
 *      这保证结果中相同距离时取较小的数。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式八：二分加速子序列匹配】
 *
 *   numMatchingSubseq：判断 word 是否为 s 的子序列
 *   朴素做法 O(|s|·|words|) 超时，用"字符索引表 + 二分"优化为 O(|word|·log|s|)。
 *
 *   步骤：
 *     1. 预处理 s：为每个字符建立其所有出现位置的升序列表 mapIndexs
 *     2. 对 word 的每个字符 word[i]，在 mapIndexs[word[i]] 中二分搜索 ≥ j 的最小位置
 *        （j 是上一次匹配位置 +1，保证子序列的顺序性）
 *     3. 若找不到（pos == len(indexList)），说明 word 不是子序列
 *
 *   ⚠️ 易错点20：这里 leftBound 不做值检查，直接返回 left（插入位置语义）
 *      我们要找的是"≥ j 的最小索引"，不关心是否恰好等于 j。
 *      返回 left 后需判断 left == len(indexList)（所有位置都 < j，匹配失败）。
 *
 *   ⚠️ 易错点21：匹配成功后 j 要更新为 indexList[pos]+1（不是 pos+1）
 *      pos 是索引列表中的下标，indexList[pos] 才是 s 中的实际位置。
 *      下一轮要从 s 中该位置的下一个字符开始，所以 j = indexList[pos] + 1。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次写完二分题后对照检查
 *
 *     ✅ mid 用 left+(right-left)/2 防溢出了吗？
 *     ✅ 循环条件：是 left <= right（精确/边界/旋转搜索）还是 left < right（findMin）？
 *     ✅ 边界收缩：是 mid-1/mid+1（排除 mid）还是 mid（保留 mid）？两者必须配对循环条件。
 *     ✅ 左/右边界：找到 target 时是压 right 还是压 left？返回的是 left 还是 right？
 *     ✅ 循环结束后：做了越界检查 + 值检查吗？（除非需要插入位置语义）
 *     ✅ 抽象二分 x 范围：下界用 max(arr)，上界用 sum(arr)？
 *     ✅ f(x) 分段计数：有 i-- 退回当前元素吗？循环后有 +1 补最后一段吗？
 *     ✅ 2D 矩阵：展平用 mid/n 和 mid%n（n 是列数）？
 *     ✅ 峰值搜索：先判断 mid+1 越界了吗？
 *     ✅ 旋转数组：判断有序半区用 >= nums[left]（含等号）？
 *     ✅ findMin：循环条件 left < right + 收缩用 right = mid（不是 mid-1）？
 *     ✅ 最接近元素：用开区间 (left, right) 初始化了吗？
 * ============================================================================
 */

package main

import (
	"fmt"
	"math"
)

// 二分搜索-基础框架
func binarySearch(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		}
	}
	return -1
}

// 二分搜索-寻找左侧边界
func leftBound(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			right = mid - 1
		} else if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		}
	}

	// 如果越界target肯定不存在，返回-1
	if left >= len(nums) || left < 0 {
		return -1
	}

	// 结束时，left=right+1
	if nums[left] != target {
		return -1
	}

	// 去掉上面判断，target不存在时返回的是「大于target的最小索引」
	return left
}

// 二分搜索-寻找右侧边界
func rightBound(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		}
	}

	// 如果越界target肯定不存在，返回-1
	if right < 0 || right >= len(nums) {
		return -1
	}

	// 结束时，left=right+1 => right=left-1
	if nums[right] != target {
		return -1
	}

	// 去掉上面判断，target不存在时返回的是「小于target的最大索引」
	return right
}

// 在排序数组中查找元素的第一个和最后一个位置
// https://leetcode.cn/problems/find-first-and-last-position-of-element-in-sorted-array/
func searchRange(nums []int, target int) []int {
	leftIdx := leftBound(nums, target)
	if leftIdx == -1 {
		return []int{-1, -1}
	}
	rightIdx := rightBound(nums, target)
	return []int{leftIdx, rightIdx}
}

/* 🌟什么时候可以运用二分搜索算法技巧❓
 *
 *   首先要能抽象出一个自变量 x，一个关于 x 的函数 f(x)，以及一个目标值 target。
 *   同时 x, f(x), target 还要满足以下条件：
 *   1、f(x) 必须是在 x 上的单调函数（单调增单调减都可以）。
 *   2、题目是计算满足约束条件 f(x) == target 时的 x 的最值。
 *
 * 🌟具体操作步骤如下👇：
 *
 *   1、确定 x, f(x), target 分别是什么，并写出函数 f 的代码，即：
 *    - x 是什么？（我能控制的变量，注意：通常x本身就是一个'最值'）
 *    - f(x) 是什么？（它是关于x单调函数，x变化时(穷举所有可能)，结果如何变化）
 *    - target 是什么？求左边界还是右边界？（优化方向，也就是穷举的过程）
 *   2、找到 x 的取值范围作为二分搜索的搜索区间，初始化 left 和 right 变量。
 *   3、根据题目的要求，确定应该使用「搜索左侧」还是「搜索右侧」的二分搜索算法，写出解法代码。
 *
 * 🌟如何训练这种抽象能力？
 *   🔑 第一步：忘掉故事，提取"操作"
 *   不要被"船"、"包裹"、"天数"这些词干扰。问自己：
 *    - 我在操作什么？ → 一个有序序列
 *    - 我能做什么操作？ → 按顺序切分成连续段
 *    - 约束是什么？ → 每段有个上限 / 段数有个上限
 *    - 优化什么？ → 在满足一个约束的前提下，最小化另一个
 *
 *   🔑 第二步：识别"对偶关系"
 *   这类题目的核心洞察是——两个变量之间存在单调的对偶关系：
 *    - 每段上限 x ↑  ⟹  需要的段数 f(x) ↓
 *    - 每段上限 x ↓  ⟹  需要的段数 f(x) ↑
 *   这就是你代码注释里写的 "f(x) 必须是在 x 上的单调函数"。
 *   一旦识别出这种单调对偶关系，就可以用二分搜索。
 *
 *   🔑 第三步：统一建模模板
 *    "连续分段"类二分问题的统一模型：
 *     - 输入：数组 arr，限制值 target
 *     - x：每段的容量上限
 *     - f(x)：按容量 x 分段，最少需要几段
 *     - 单调性：x 越大 → f(x) 越小
 *     - 求：满足 f(x) ≤ target 的最小 x（左边界）
 *    x 的范围：[max(arr), sum(arr)]
 */

// 二分搜索题目的思维框架：在 f(x)==target 的约束下求 x 的最值
func binarySearchSolution(nums []int, target int) int {
	// 函数 f 是关于自变量 x 的单调函数
	f := func(x int) int {
		// 这是内层for循环，x 已经固定了
		for i := 0; i < len(nums); {
			// ...
		}
		return -1
	}

	left := -1  // 自变量 x 的最小值是多少❓
	right := -1 // 自变量 x 的最大值是多少？❓

	// 外侧for循环，目的是穷举 x 所有可能
	for left <= right {

		mid := left + (right-left)/2
		fv := f(mid)

		if fv == target {
			// 要求解的最值在左边界❓
			left = mid + 1
			// 还是在右边界❓
			right = mid - 1

		} else if fv > target {
			// 怎么让 f(x) 小一点❓

		} else if fv < target {
			// 怎么让 f(x) 大一点❓

		}
	}
	return left
}

// 爱吃香蕉的珂珂 https://leetcode.cn/problems/koko-eating-bananas/
func minEatingSpeed(piles []int, h int) int {
	// x: 吃香蕉的速度 x
	// f(x): 吃完全部香蕉需要的时间 f(x)
	// target: 给定的时间限制 h
	f := func(x int) (hours int) {
		for _, pile := range piles {
			hours += pile / x
			if pile%x > 0 {
				hours++
			}
		}
		return hours
	}

	// 吃香蕉速度最小:1根/h
	left := 1
	// 吃香蕉速度最大:max(piles)根/h
	right := int(math.Pow10(9) + 1)

	// 穷举自变量 x 的所有可能
	for left <= right {
		mid := left + (right-left)/2
		fv := f(mid)
		if fv == h {
			right = mid - 1
		} else if fv > h {
			left = mid + 1
		} else if fv < h {
			right = mid - 1
		}
	}
	return left
}

// 在D天内送达包裹的能力 https://leetcode.cn/problems/capacity-to-ship-packages-within-d-days/
// 本质：给定一个数组，将其按顺序分成若干连续段，每段的"总量"不超过 x。求使得段数 ≤ target 时，x 的最小值
func shipWithinDays(weights []int, days int) int {
	// x: 船的载能力 x
	// f(x): 按运载能力x运送，需要f(x)天运完
	// target: 给定的时间限制 days

	f := func(x int) (days int) {
		sum := 0
		for i := 0; i < len(weights); i++ {
			sum += weights[i]
			if sum > x {
				days++
				i-- //已经放进来的weights[i]要重新放回去
				sum = 0
			}
		}
		days++
		return days
	}

	// 船的最小载能力 max(weights)
	left := -1
	// 船的最大载能力 sum(weights)
	right := -1

	for i := range weights {
		left = max(left, weights[i])
		right += weights[i]
	}

	// 穷举自变量 x 的所有可能
	for left <= right {
		mid := left + (right-left)/2
		fv := f(mid)
		if fv == days {
			right = mid - 1
		} else if fv > days {
			left = mid + 1
		} else if fv < days {
			right = mid - 1
		}
	}
	return left
}

// 分割数组的最大值 https://leetcode.cn/problems/split-array-largest-sum/
// 本质：给定一个数组，将其按顺序分成若干连续段，每段的"总量"不超过 x。求使得段数 ≤ target 时，x 的最小值
func splitArray(nums []int, k int) int {
	// x: 子数组和的上限 x
	// f(x): 按子数组和sum为x来拆分，需要f(x)个连续子数组
	// target: 给定的子数组个数限制 k

	f := func(x int) (num int) {
		sum := 0
		for i := 0; i < len(nums); i++ {
			sum += nums[i]
			if sum > x {
				num++
				i-- //已经放进来的nums[i]要重新放回去
				sum = 0
			}
		}
		num++
		return num
	}

	// 子数组和最小 max(weights)
	left := -1
	// 子数组和最大 sum(weights)
	right := -1

	for i := range nums {
		left = max(left, nums[i])
		right += nums[i]
	}

	// 穷举自变量 x 的所有可能
	for left <= right {
		mid := left + (right-left)/2
		fv := f(mid)
		if fv == k {
			right = mid - 1
		} else if fv > k {
			left = mid + 1
		} else if fv < k {
			right = mid - 1
		}
	}
	return left
}

// 搜索二维矩阵 https://leetcode.cn/problems/search-a-2d-matrix/description/
// 每行中的整数从左到右按非严格递增顺序排列。
// 每行的第一个整数大于前一行的最后一个整数。
func searchMatrix(matrix [][]int, target int) bool {
	m, n := len(matrix), len(matrix[0])
	left, right := 0, m*n-1

	for left <= right {
		mid := left + (right-left)/2
		find := matrix[mid/n][mid%n]

		if find == target {
			return true
		} else if find > target {
			right = mid - 1
		} else if find < target {
			left = mid + 1
		}
	}
	return false
}

// 搜索二维矩阵II https://leetcode.cn/problems/search-a-2d-matrix-ii/
// 每行的元素从左到右升序排列。
// 每列的元素从上到下升序排列。
func searchMatrixII(matrix [][]int, target int) bool {
	// 从右上角开始，规定只能向左或向下移动
	m, n := len(matrix), len(matrix[0])
	i, j := 0, n-1

	for i < m && j >= 0 {
		//fmt.Printf("i=%d j=%d next=%d\n", i, j, matrix[i][j])

		if matrix[i][j] == target {
			return true
		} else if matrix[i][j] > target {
			j--
		} else if matrix[i][j] < target {
			i++
		}
	}
	return false
}

// 匹配子序列的单词数 https://leetcode.cn/problems/number-of-matching-subsequences/
func numMatchingSubseq(s string, words []string) int {
	mapIndexs := make(map[byte][]int)
	for i := range s {
		if mapIndexs[s[i]] == nil {
			mapIndexs[s[i]] = []int{}
		}
		mapIndexs[s[i]] = append(mapIndexs[s[i]], i)
	}

	leftBound := func(nums []int, target int) int {
		left, right := 0, len(nums)-1
		for left <= right {
			mid := left + (right-left)/2
			if nums[mid] == target {
				right = mid - 1
			} else if nums[mid] > target {
				right = mid - 1
			} else if nums[mid] < target {
				left = mid + 1
			}
		}
		return left
	}

	res := 0
	for _, word := range words {
		i := 0 // word的字符索引
		j := 0 // s的字符索引
		for i < len(word) {
			ch := word[i]

			indexList := mapIndexs[ch]
			if indexList == nil {
				break
			}

			// 二分搜索大于等于j的最小索引（s已经走到j-1的位置了）
			// 即在 s[j..] 中搜索等于 word[i] 的最小索引
			pos := leftBound(indexList, j)
			if pos == len(indexList) {
				break
			}

			// 找到 word[i] == s[j]，继续往后匹配
			j = indexList[pos]

			i++
			j++
		}
		if i == len(word) {
			res++
		}
	}
	return res
}

// 找到K个最接近的元素 https://leetcode.cn/problems/find-k-closest-elements/
func findClosestElements(arr []int, k int, x int) []int {
	leftBound := func(nums []int, x int) int {
		left, right := 0, len(nums)-1
		for left <= right {
			mid := left + (right-left)/2
			if nums[mid] == x {
				right = mid - 1
			} else if nums[mid] > x {
				right = mid - 1
			} else if nums[mid] < x {
				left = mid + 1
			}
		}
		return left
	}

	p := leftBound(arr, x)

	// 因为p本身可能越界，选择两端都开的区间(left, right)
	left, right := p-1, p

	// 扩展区间，直到区间内包含k个元素
	for right-left-1 < k {
		if left == -1 {
			right++
		} else if right == len(arr) {
			left--
		} else if x-arr[left] > arr[right]-x {
			right++
		} else {
			left--
		}
	}

	var res []int
	for i := left + 1; i < right; i++ {
		res = append(res, arr[i])
	}
	return res
}

// 搜索插入位置 https://leetcode.cn/problems/search-insert-position/description/
func searchInsert(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		}
	}
	return left
}

// 寻找峰值 https://leetcode.cn/problems/find-peak-element/
func findPeakElement(nums []int) int {
	// 题目前提：对于所有有效的i都有nums[i]!=nums[i+1]
	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)/2
		if mid+1 == len(nums) {
			return mid
		} else if nums[mid] > nums[mid+1] {
			right = mid - 1
		} else if nums[mid] < nums[mid+1] {
			left = mid + 1
		}
	}
	return left
}

// 山脉数组的峰顶索引 https://leetcode.cn/problems/peak-index-in-a-mountain-array/description/
func peakIndexInMountainArray(arr []int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := left + (right-left)/2
		if mid+1 == len(arr) {
			return mid
		} else if arr[mid] > arr[mid+1] {
			right = mid - 1
		} else if arr[mid] < arr[mid+1] {
			left = mid + 1
		}
	}
	return left
}

// 搜索旋转排序数组 https://leetcode.cn/problems/search-in-rotated-sorted-array/description/
func search(nums []int, target int) int {
	// 二分搜索只能在有序的数组上查找
	// 本题须先确定有序的区间范围

	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return mid
		}
		if nums[mid] >= nums[left] {
			// 在左侧, 此时[left..mid]有序
			if target >= nums[left] && target <= nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			//在右侧, 此时[mid..right]有序
			if target >= nums[mid] && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}
	return -1
}

// 搜索旋转排序数组II https://leetcode.cn/problems/search-in-rotated-sorted-array-ii/description/
func searchII(nums []int, target int) bool {
	// 0、跳过重复元素
	// 1、确定 mid 中点落在「断崖」左侧还是右侧。
	// 2、在第 1 步确定的结果之上，根据target和nums[left],nums[right],nums[mid]的相对大小收缩搜索区间。

	left, right := 0, len(nums)-1
	for left <= right {
		for left < len(nums)-1 && nums[left] == nums[left+1] {
			left++
		}
		for right > 0 && nums[right] == nums[right-1] {
			right--
		}

		mid := left + (right-left)/2
		if nums[mid] == target {
			return true
		}
		if nums[mid] >= nums[left] {
			// 在左侧, 此时[left..mid]有序
			if target >= nums[left] && target <= nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			//在右侧, 此时[mid..right]有序
			if target >= nums[mid] && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}
	return false
}

// 寻找旋转排序数组中的最小值 https://leetcode.cn/problems/find-minimum-in-rotated-sorted-array/description/
func findMin(nums []int) int {
	left, right := 0, len(nums)-1
	for left < right {
		mid := left + (right-left)/2
		if nums[mid] > nums[right] {
			// mid落在断崖左边,最小值在[mid+1,right]
			left = mid + 1
		} else {
			// mid落在断崖右边(或本身就是最小值),最小值在[left,mid]
			right = mid
		}
	}
	return nums[left]
}

func main() {

	fmt.Println(findMin([]int{5, 6, 7, 8, 1, 2, 3, 4}))
	fmt.Println(findMin([]int{1, 2, 3, 4}))

	//fmt.Println(searchII([]int{1, 0, 1, 1, 1}, 0))
	//fmt.Println(searchII([]int{2, 5, 6, 0, 0, 1, 2}, 0))
	//fmt.Println(searchII([]int{2, 5, 6, 0, 0, 1, 2}, 3))

	//fmt.Println(search([]int{6, 7, 8, 9, 10, 1, 2, 3, 4, 5}, 10))
	//fmt.Println(search([]int{6, 7, 8, 9, 10, 1, 2, 3, 4, 5}, 7))
	//fmt.Println(search([]int{6, 7, 8, 9, 10, 1, 2, 3, 4, 5}, 4))
	//fmt.Println(search([]int{6, 7, 8, 9, 10, 1, 2, 3, 4, 5}, 100))

	//fmt.Println(peakIndexInMountainArray([]int{1, 2, 3, 4, 5, 4, 3, 2, 1}))
	//fmt.Println(peakIndexInMountainArray([]int{10}))
	//fmt.Println(peakIndexInMountainArray([]int{-1, 3, 5, 9, 3, 2, -10}))

	//fmt.Println(findPeakElement([]int{1, 2, 3, 4, 5}))
	//fmt.Println(findPeakElement([]int{5, 4, 3, 2, 1}))
	//fmt.Println(findPeakElement([]int{1, 2, 3, 1}))
	//fmt.Println(findPeakElement([]int{1, 2, 1, 3, 5, 6, 4}))

	//fmt.Println(searchInsert([]int{1, 3, 5, 6}, 5))
	//fmt.Println(searchInsert([]int{1, 3, 5, 6}, 2))
	//fmt.Println(searchInsert([]int{1, 3, 5, 6}, 7))

	//fmt.Println(findClosestElements([]int{3, 5, 8, 10}, 2, 15))
	//fmt.Println(findClosestElements([]int{1, 2, 3, 4, 5}, 4, 3))
	//fmt.Println(findClosestElements([]int{1, 2, 3, 4, 5}, 4, -1))

	//fmt.Println(numMatchingSubseq("abcde", []string{"a", "bb", "acd", "ace"}))
	//fmt.Println(numMatchingSubseq("dsahjpjauf", []string{"ahjpjau", "ja", "ahbwzgqnuk", "tnmlanowax"}))

	//fmt.Println(searchMatrixII([][]int{
	//	{1, 4, 7, 11, 15},
	//	{2, 5, 8, 12, 19},
	//	{3, 6, 9, 16, 22},
	//	{10, 13, 14, 17, 24},
	//	{18, 21, 23, 26, 30},
	//}, 13))

	//fmt.Println(searchMatrix([][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 60))
	//fmt.Println(searchMatrix([][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 13))

	//fmt.Println(splitArray([]int{7, 2, 5, 10, 8}, 2))
	//fmt.Println(splitArray([]int{1, 2, 3, 4, 5}, 2))
	//fmt.Println(splitArray([]int{1, 4, 4}, 3))

	//fmt.Println(shipWithinDays([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5))
	//fmt.Println(shipWithinDays([]int{3, 2, 2, 4, 1, 4}, 3))
	//fmt.Println(shipWithinDays([]int{1, 2, 3, 1, 1}, 4))

	//fmt.Println(minEatingSpeed([]int{3, 6, 7, 11}, 8))
	//fmt.Println(minEatingSpeed([]int{30, 11, 23, 4, 20}, 5))
	//fmt.Println(minEatingSpeed([]int{30, 11, 23, 4, 20}, 6))

	//fmt.Println(searchRange([]int{5, 7, 7, 8, 8, 10}, 8))
	//fmt.Println(searchRange([]int{5, 7, 7, 8, 8, 10}, 6))

	//fmt.Println(rightBound([]int{-1, 1, 2, 3, 3, 3, 3, 3, 5, 6}, 3))
	//fmt.Println(leftBound([]int{-1, 1, 2, 3, 3, 3, 3, 3, 5, 6}, 3))
	//fmt.Println(binarySearch([]int{-1, 1, 2, 3, 3, 3, 3, 3, 5, 6}, 3))
	fmt.Println("------END------")
}
