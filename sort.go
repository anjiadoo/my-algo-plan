/*
 * ============================================================================
 *                        📘 排序算法全集 · 核心记忆框架
 * ============================================================================
 * 【一句话理解各排序】
 *
 *   选择排序：每轮从未排序区「选最小」放到已排序区末尾，像选牌。
 *   冒泡排序：每轮相邻两两比较，小的「上浮」，大的「下沉」，像冒泡。
 *   插入排序：每次取一张新牌，在手牌（已排序区）中找位置「插入」。
 *   快速排序：选 pivot 分区，左边 ≤ pivot，右边 > pivot，递归左右子区间。
 *   归并排序：先递归拆成单元素，再两两「合并」有序子数组，后序遍历框架。
 *   堆  排序：建最大堆后，反复将堆顶（最大值）换到末尾，堆范围逐步缩小。
 *   计数排序：统计每个值出现次数，前缀和推算位置，非比较排序。
 *   桶  排序：将元素按范围分桶，桶内用插入排序，再合并各桶。
 *   希尔排序：插入排序的改进版，先 h 有序，逐步缩小 h 到 1。
 *   基数排序：计数排序的扩展版，按位（个十百…）依次稳定排序。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【算法横向对比】
 *
 *   算法      时间复杂度(均)  时间复杂度(最差)  空间复杂度  原地   稳定
 *   ──────────────────────────────────────────────────────────────────────
 *   选择排序    O(n²)          O(n²)             O(1)       ✅     ❌
 *   冒泡排序    O(n²)          O(n²)             O(1)       ✅     ✅
 *   插入排序    O(n²)          O(n²)             O(1)       ✅     ✅
 *   快速排序    O(n logn)      O(n²) 已排序时      O(logn)   ✅     ❌
 *   归并排序    O(n logn)      O(n logn)          O(n)      ❌     ✅
 *   堆  排序    O(n logn)      O(n logn)          O(1)      ✅     ❌
 *   计数排序    O(n+k)         O(n+k)             O(n+k)    ❌     ✅  k=值域范围
 *   桶  排序    O(n+k)         O(n²) 全在一桶     O(n+k)    ❌     ✅  k=桶数
 *   希尔排序    O(n logn)      O(n²)              O(1)      ✅     ❌
 *   基数排序    O(d*(n+k))     O(d*(n+k))         O(n+k)    ❌     ✅  d=位数
 *
 *   口诀：稳定的有「冒插归计桶基」，原地的有「选冒插快堆希」。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【选择 & 冒泡 & 插入排序的核心区别】
 *
 *   三者都是 O(n²)，区别在于「内层循环扫描的区域」不同：
 *     选择排序：内层循环在「未排序区」，找最小值（直接选出来）
 *     冒泡排序：内层循环在「未排序区」，相邻交换（让最小浮上来）
 *     插入排序：内层循环在「已排序区」，找插入点（找到就 break）
 *
 *   ⚠️ 易错点1：插入排序内层一旦 nums[j] >= nums[j-1] 就要 break
 *      因为已排序区本来就有序，不满足条件说明已找到正确位置，继续比较是浪费。
 *
 *   ⚠️ 易错点2：冒泡排序可加 swap 标志提前终止，但插入排序不能提前终止外层
 *      插入排序 break 的是内层（找到插入点），外层仍需遍历所有未排序元素。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【快速排序的核心：partition 分区】
 *
 *   partition(nums, lo, hi) 做三件事：
 *     ① 选 pivot = nums[lo]
 *     ② 双指针从两端向中间扫描：i 跳过 ≤pivot 的，j 跳过 >pivot 的，相遇前交换
 *     ③ 循环结束后，将 pivot（nums[lo]）与 j 交换，返回 j
 *
 *   ⚠️ 易错点3：i 用 <=，j 用 >，必须不对称
 *      若两边都用 <，当元素等于 pivot 时 i/j 都不动 → 死循环。
 *      i 跳过「≤pivot」，j 跳过「>pivot」，保证相遇后能停下。
 *
 *   ⚠️ 易错点4：循环结束只能将 pivot 与 j 交换，不能与 i 交换
 *      循环结束时：j 左边（含 j）所有元素 ≤ pivot，j 右边所有元素 > pivot。
 *      j 就是 pivot 的最终位置，换 i 是错的（i 此时已越过 j，指向 >pivot 的元素）。
 *
 *   ⚠️ 易错点5：快排递归终止条件是 lo >= hi，不能写 lo == hi
 *      当 partition 返回 j=lo 时，sort(nums, lo, j-1) 传入的是 (lo, lo-1)，
 *      lo > hi 这种空区间必须被 >= 捕获，== 会漏掉导致越界。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【归并排序的核心：merge 合并】
 *
 *   merge(nums, lo, mid, hi, temp) 做三件事：
 *     ① 双指针 i=lo, j=mid+1，比较两边较小值依次写入临时数组 temp
 *     ② 将 i/j 中剩余的一侧直接追加到 temp
 *     ③ 将 temp[0..k-1] 写回 nums[lo..hi]
 *
 *   ⚠️ 易错点6：mid 计算用 lo+(hi-lo)/2，不用 (lo+hi)/2
 *      lo+hi 在数值很大时会整型溢出（Go 的 int 是 64 位通常不触发，但好习惯）。
 *
 *   ⚠️ 易错点7：归并排序递归终止条件是 lo == hi，不能写 lo >= hi
 *      快排是前序（先分区再递归），空区间是合法输入需要拦截；
 *      归并是后序（先递归再合并），调用方保证 lo ≤ hi，lo==hi 即单元素直接返回。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【堆排序的核心：heapify 下沉调整】
 *
 *   heapify(nums, n, i) 将以 i 为根的子树调整为最大堆（堆范围是 nums[0..n-1]）：
 *     ① 找出 i、left=2i+1、right=2i+2 中的最大值索引 maxIndex
 *     ② 若 maxIndex != i，交换 nums[i] 与 nums[maxIndex]，递归 heapify(maxIndex)
 *
 *   ⚠️ 易错点8：建堆从 n/2-1 开始向上，不是从 0 开始
 *      索引 n/2 及以后的节点都是叶子节点（无子节点），天然满足堆性质，无需 heapify。
 *      从 n/2-1 开始能减少约一半的无效调用。
 *
 *   ⚠️ 易错点9：排序阶段调用 heapify(nums, i, 0)，n 传 i 不传 len(nums)
 *      每次将堆顶换到末尾后，末尾元素已排好序，不能再被 heapify 破坏。
 *      传 len(nums) 会把已归位的元素再次纳入堆，结果错误。
 *
 *   ⚠️ 易错点10：heapify 交换后必须递归向下调整
 *      交换使 maxIndex 位置的值变小，可能破坏了以 maxIndex 为根的子树的堆性质，
 *      必须递归修复，不能只调整一层。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【计数排序的核心：前缀和定位】
 *
 *   三步：统计次数 → count 数组前缀和 → 逆序遍历原数组放置元素
 *
 *   前缀和之后 count[i] 含义变了：
 *     原来：count[i] = 值 i 出现的次数
 *     之后：count[i] = 值 i 在排序后数组中的「结束位置」（右边界，1-based）
 *     即 sorted[count[i]-1] 是值 i 最后一次出现的位置，放完后 count[i]--。
 *
 *   ⚠️ 易错点11：放置元素时必须从后往前遍历原数组（保证稳定性）
 *      相同值的元素在原数组中有先后顺序，从后往前遍历能保证靠后的元素放到靠后位置，
 *      维持相对顺序不变（稳定排序的定义）。
 *
 *   ⚠️ 易错点12：需要用 offset = -minNum 做下标平移，避免负数索引
 *      元素值直接做数组下标，遇到负数会越界。offset 将值域 [minNum, maxNum]
 *      平移到 [0, maxNum-minNum]，count 数组大小为 maxNum-minNum+1。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【桶排序的核心：映射函数 + 桶大小】
 *
 *   三步：分桶（除法向下取整）→ 桶内排序（插入排序）→ 合并各桶
 *
 *   桶索引计算：index = (nums[i] + offset) / bucketSize
 *   桶大小公式：bucketSize = (maxNum - minNum) / bucketCount + 1
 *
 *   ⚠️ 易错点13：bucketSize 必须 +1，防止最大值的桶索引越界
 *      若不加 1，最大值 maxNum 计算出的 index = (maxNum-minNum)/bucketSize = bucketCount，
 *      恰好等于桶数组长度，导致越界。+1 后最大值 index < bucketCount，安全。
 *
 *   ⚠️ 易错点14：桶内排序应用插入排序，不能递归调用 bucketSort 自身
 *      桶内元素少，插入排序在小数据量下常数因子小，性能好。
 *      递归调用自身会因再次计算 sqrt 分桶导致逻辑混乱甚至死递归。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次手写排序后对照检查
 *
 *     ✅ 插入排序：内层找到正确位置后有没有 break？
 *     ✅ 冒泡排序（优化版）：有没有在一轮无交换时提前 break 外层？
 *     ✅ 快速排序：i 用 <=、j 用 > 的不对称比较？
 *     ✅ 快速排序：分区结束后是与 j 交换 pivot，不是 i？
 *     ✅ 快速排序：递归终止条件是 lo >= hi，不是 lo == hi？
 *     ✅ 归并排序：mid 用 lo+(hi-lo)/2 防溢出？
 *     ✅ 归并排序：递归终止条件是 lo == hi，不是 lo >= hi？
 *     ✅ 堆排序：建堆从 n/2-1 开始，不是 0？
 *     ✅ 堆排序：排序阶段 heapify 的 n 传 i，不是 len(nums)？
 *     ✅ 计数排序：放置元素时从后往前遍历原数组（保稳定性）？
 *     ✅ 计数排序：count 前缀和后是「结束位置」，放置后要 count--？
 *     ✅ 计数/桶排序：有没有用 offset=-minNum 处理负数？
 *     ✅ 桶排序：bucketSize 有没有 +1 防最大值越界？
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *
 *     0. selectSort(nums []int) []int        // O(n²)，不稳定，原地
 *     1. bubbleSort(nums []int) []int         // O(n²)，稳定，原地
 *     2. bubbleSortV1(nums []int) []int       // O(n²)，稳定，原地，提前终止优化
 *     3. insertSort(nums []int) []int         // O(n²)，稳定，原地
 *     4. quickSort(nums []int) []int          // O(n logn)，不稳定，原地
 *     5. partition(nums []int, lo, hi int) int
 *     6. mergeSort(nums []int) []int          // O(n logn)，稳定，非原地 O(n)
 *     7. merge(nums []int, lo, mid, hi int, temp []int)
 *     8. heapSort(nums []int) []int           // O(n logn)，不稳定，原地
 *     9. heapify(nums []int, n, i int)
 *    10. countSort(nums []int) []int          // O(n+k)，稳定，非原地
 *    11. bucketSort(nums []int) []int         // O(n+k)，稳定，非原地
 * ============================================================================
 */

package main

import (
	"fmt"
	"math"
)

// 排序算法实现：
// 🌟技巧1：提前终止技巧 - 冒泡排序可用 swap 标志提前终止，若一轮无交换说明已有序，内层循环在未排序区（冒泡排序）
// 🌟技巧2：及时break技巧 - 插入排序内层找到正确位置后要 break，不需要继续往前比较，内层循环在已排序区（插入排序）
// 🌟技巧3：不对称比较技巧 - 分区时 i 用 <=、j 用 >，不对称才能保证 i/j 相遇时停下，否则死循环（快速排序）
// 🌟技巧4：pivot交换技巧 - 分区结束后 pivot 只能与 j 交换而非 i，因为 j 及其左边都 <= pivot（快速排序）
// 🌟技巧5：递归终止条件技巧 - 快排递归终止条件是 lo >= hi，归并是 lo == hi，快排用 == 会漏掉空区间（快排 vs 归并）
// 🌟技巧6：防溢出技巧 - 计算中间位置时使用 lo + (hi-lo)/2 防止整型溢出（归并排序）
// 🌟技巧7：非叶子节点起始技巧 - 构建堆时从最后一个非叶子节点（n/2-1）开始，叶子节点天然满足堆性质（堆排序）
// 🌟技巧8：堆范围缩小技巧 - 堆排序交换堆顶到末尾后，heapify 的 n 要缩小（传 i 而非 len），否则已排好的元素会被破坏（堆排序）
// 🌟技巧9：递归下沉技巧 - heapify 交换后必须递归向下调整，因为交换可能破坏子树的堆性质（堆排序）
// 🌟技巧10：前缀和定位技巧 - count 数组做前缀和累加后，存的是"结束位置"而非"起始位置"，放置元素后要 count--（计数排序）
// 🌟技巧11：逆序遍历保稳定技巧 - 从后往前遍历原数组放置元素，才能保证排序的稳定性（计数排序）
// 🌟技巧12：偏移量平移技巧 - 使用 offset=-minNum 将数据范围平移到 [0, max-min]，避免负数索引（计数排序、桶排序）
// 🌟技巧13：桶大小+1技巧 - 桶大小用 (max-min)/bucketCount+1，必须 +1 防止最大值越界（桶排序）
// 🌟技巧14：h有序数组技巧 - 希尔排序是对「插入排序」的简单改进，关键概念：h有序数组（希尔排序）
// 🌟技巧15：按位排序技巧 - 基数排序是由「计数排序」扩展而来，原理按(位)排序（基数排序）

// 1、选择排序-顺着从头遍历到尾-原地排序-不稳定排序
func selectSort(nums []int) []int {
	// 算法的核心思想分两步，如下
	//	a. 每轮从「未排序区」中找到最小元素，将其放到「已排序区」的末尾（直接交换到 sortedIndex 位置）。
	//	b. 重复上述过程，已排序区不断扩大，直到所有元素有序。

	for sortedIndex := 0; sortedIndex < len(nums); sortedIndex++ {
		for j := sortedIndex + 1; j < len(nums); j++ {
			if nums[sortedIndex] > nums[j] {
				nums[sortedIndex], nums[j] = nums[j], nums[sortedIndex]
			}
		}
	}
	return nums
}

// 2、冒泡排序-逆着从尾遍历到头-原地排序-稳定排序
func bubbleSort(nums []int) []int {
	// 算法的核心思想分两步，如下
	//	a. 每轮从数组末尾开始，相邻元素两两比较，将较小的元素像气泡一样「浮」到「未排序区」的最前面。
	//	b. 重复上述过程，已排序区从头部不断扩大，直到所有元素有序。

	for sortedIndex := 0; sortedIndex < len(nums); sortedIndex++ {
		for j := len(nums) - 1; j > sortedIndex; j-- {
			if nums[j-1] > nums[j] {
				nums[j-1], nums[j] = nums[j], nums[j-1]
			}
		}
	}
	return nums
}

// 2.1、冒泡排序-逆着从尾遍历到头-原地排序-稳定排序（优化：提前终止）
func bubbleSortV1(nums []int) []int {
	// 算法的核心思想与 bubbleSort 相同，在其基础上加了一个优化：
	//	若某轮遍历中没有发生任何交换，说明数组已经有序，可以提前终止，避免多余的循环。

	for sortedIndex := 0; sortedIndex < len(nums); sortedIndex++ {
		swap := false
		for j := len(nums) - 1; j > sortedIndex; j-- {
			if nums[j-1] > nums[j] {
				swap = true
				nums[j-1], nums[j] = nums[j], nums[j-1]
			}
		}
		if !swap { //冒泡排序-可提前终止，不过时间复杂度仍然是O(n^2)
			break
		}
	}
	return nums
}

// 3、插入排序-逆着从当前位置到头遍历-原地排序-稳定排序
func insertSort(nums []int) []int {
	// 算法的核心思想分两步，如下
	//	a. 取「未排序区」的第一个元素，在「已排序区」中从后往前找到它的正确插入位置（逐个向前比较并交换）。
	//	b. 插入后已排序区扩大一位，重复上述过程，直到所有元素有序。

	for sortedIndex := 0; sortedIndex < len(nums)-1; sortedIndex++ {
		if nums[sortedIndex] > nums[sortedIndex+1] {
			nums[sortedIndex], nums[sortedIndex+1] = nums[sortedIndex+1], nums[sortedIndex]
			for j := sortedIndex; j > 0; j-- {
				if nums[j] < nums[j-1] {
					nums[j], nums[j-1] = nums[j-1], nums[j]
				} else {
					break
				}
			}
		}
	}
	return nums
}

// 4、快速排序-二叉树前序遍历框架-原地排序-不稳定排序
func quickSort(nums []int) []int {
	// 算法的核心思想分三步，如下
	//	a. 从数组中选一个元素作为基准值 pivot（这里取区间第一个元素）。
	//	b. 分区(partition)：将数组重新排列，使 pivot 左侧的元素都 <= pivot，右侧的元素都 > pivot，pivot 落到其最终位置。
	//	c. 递归地对 pivot 左、右两个子区间重复上述过程，直到区间只剩一个元素为止（前序遍历框架：先分区，再递归左右）。

	var sort func(nums []int, lo, hi int)

	sort = func(nums []int, lo, hi int) {
		if lo >= hi {
			return
		}
		p := partition(nums, lo, hi)
		sort(nums, lo, p-1)
		sort(nums, p+1, hi)
	}

	sort(nums, 0, len(nums)-1)
	return nums
}

// 找一个参考值pivot，对于[lo,hi]，从两边往中间移动，找到左侧比pivot大，右侧比pivot小，然后两者交换♻️，最后再交换pivot。
func partition(nums []int, lo, hi int) int {
	pivot := nums[lo]
	i := lo + 1
	j := hi
	for i <= j {
		// 🌟技巧1：注意这里的 " <= 和 > "，原因如下:
		// i会跳过所有小于等于pivot的元素
		// j会跳过所有大于pivot的元素
		// 这样i和j最终会在某个位置相遇，保证分区能够完成，否则可能陷入死循环
		for i < hi && nums[i] <= pivot {
			i++
		}
		for j > lo && nums[j] > pivot {
			j--
		}
		if i >= j {
			break
		}
		nums[i], nums[j] = nums[j], nums[i]
	}
	// 🌟技巧2：这里只能是lo和j交换，原因如下:
	// 在循环中，所有j右边的元素都大于pivot
	// 所有j左边的元素（包括j本身）都小于等于pivot
	// 因此j是pivot的正确位置
	nums[lo], nums[j] = nums[j], nums[lo]
	return j
}

// 5、归并排序-二叉树后序遍历框架-非原地-稳定排序
func mergeSort(nums []int) []int {
	// 算法的核心思想分三步，如下
	//	a. 将数组从中间一分为二，递归地对左、右两个子数组分别排序（后序遍历框架：先递归左右，再合并）。
	//	b. 合并(merge)：借助临时数组，用双指针将两个已有序的子数组按大小顺序合并成一个有序数组。
	//	c. 将合并结果写回原数组，递归归并直到整个数组有序。

	var sort func(nums []int, lo, hi int)

	temp := make([]int, len(nums))

	sort = func(nums []int, lo, hi int) {
		// 单个元素不用排序
		if lo == hi {
			return
		}
		// 相比(lo+hi)/2来说，这种方式可防止溢出
		mid := lo + (hi-lo)/2
		sort(nums, lo, mid)
		sort(nums, mid+1, hi)
		merge(nums, lo, mid, hi, temp)
	}

	sort(nums, 0, len(nums)-1)
	return nums
}

func merge(nums []int, lo, mid, hi int, temp []int) {
	// 临时数组起始索引
	k := 0

	// 双指针分别指向左右子数组的起始位置
	i, j := lo, mid+1

	// 比较并合并两个有序子数组
	for i <= mid && j <= hi {
		if nums[i] <= nums[j] {
			temp[k] = nums[i]
			i++
		} else {
			temp[k] = nums[j]
			j++
		}
		k++
	}

	// 将剩余元素复制到临时数组
	for i <= mid {
		temp[k] = nums[i]
		i++
		k++
	}

	for j <= hi {
		temp[k] = nums[j]
		j++
		k++
	}

	// 将临时数组复制回原数组
	for idx := 0; idx < k; idx++ {
		nums[lo+idx] = temp[idx]
	}
}

// 6、堆排序-完全二叉树结构-原地排序-不稳定排序
func heapSort(nums []int) []int {
	// 算法的核心思想分三步，如下
	//	a. 建堆：从最后一个非叶子节点（n/2-1）开始，向上逐个调用 heapify，将数组构建为最大堆（堆顶为最大值）。
	//	b. 排序：将堆顶（最大值）与末尾元素交换，使最大值归位；然后对缩小了范围的堆重新 heapify，恢复堆性质。
	//	c. 重复步骤 b，堆的有效范围每次缩小 1，直到所有元素有序。

	// 从最后一个非叶子节点开始向上构建（索引 n/2-1）
	// ❓为啥呢
	// 因为：叶子节点没有子节点，所以它们天然满足最大堆的性质。
	for i := len(nums)/2 - 1; i >= 0; i-- {
		heapify(nums, len(nums), i)
	}

	for i := len(nums) - 1; i > 0; i-- {
		nums[0], nums[i] = nums[i], nums[0]
		heapify(nums, i, 0)
	}

	return nums
}

func heapify(nums []int, n, i int) {
	maxIndex := i    // 假设当前节点是最大的
	left := 2*i + 1  // 左子节点索引
	right := 2*i + 2 // 右子节点索引

	// 如果左子节点存在且大于根节点
	if left < n && nums[left] > nums[maxIndex] {
		maxIndex = left
	}

	// 如果右子节点存在且大于当前最大值
	if right < n && nums[right] > nums[maxIndex] {
		maxIndex = right
	}

	// 如果最大值不是根节点，交换并递归调整
	// ❓为啥要递归调整
	// 因为，交换破坏了原本符合堆性质的子树，所以需要递归调整
	if maxIndex != i {
		nums[i], nums[maxIndex] = nums[maxIndex], nums[i]
		heapify(nums, n, maxIndex)
	}
}

// 7、计数排序-非原地-稳定排序-非比较排序
func countSort(nums []int) []int {
	// 算法的核心思想分三步，如下
	//	a. 统计：用 count 数组统计每个元素出现的次数（以元素值映射为下标，需用 offset 平移处理负数）。
	//	b. 累加：对 count 数组做前缀和，count[i] 变为「值 i 在排序后数组中的结束位置（右边界）」。
	//	c. 放置：从后往前遍历原数组，按 count 中的位置将元素放入结果数组，每放一个 count-- 保证稳定性。
	// 时间和空间复杂度都是：O(n+max−min)，其中 n 是待排序数组长度，max−min 是待排序数组的元素范围

	// 找到最大和最小元素
	// 计算索引偏移量和 count 数组大小
	minNum, maxNum := nums[0], nums[0]
	for _, num := range nums {
		if num < minNum {
			minNum = num
		}
		if num > maxNum {
			maxNum = num
		}
	}

	// 根据最大值和最小值，将元素映射到从 0 开始的索引值
	offset := -minNum
	count := make([]int, maxNum-minNum+1)

	// 统计每个元素出现的次数
	for _, num := range nums {
		index := num + offset
		count[index]++
	}

	// 累加 count 数组，得到的是 nums[i] 在排序后的数组中的结束位置
	// count[0]=3 => sorted[0]=x，sorted[1]=x，sorted[2]=x
	for i := 1; i < len(count); i++ {
		count[i] += count[i-1]
	}

	// 根据每个元素排序后的索引位置，完成排序
	// 这里注意，我们从后往前遍历 nums，是为了保证排序的稳定性
	sorted := make([]int, len(nums))
	for i := len(nums) - 1; i >= 0; i-- {
		index := nums[i] + offset
		sorted[count[index]-1] = nums[i]
		count[index]--
	}

	copy(nums, sorted)
	return nums
}

// 8、桶排序-非原地-稳定排序
func bucketSort(nums []int) []int {
	// 算法的核心思想分三步，如下
	//	a. 将待排序数组中的元素使用映射函数(除法下取整)分配到若干个「桶」中。
	//	b. 对每个桶中的元素进行排序(应该使用简单的排序算法，如插入排序)。
	//	c. 最后将这些排好序的桶进行合并，得到排序结果。

	// 找到最大和最小
	// 计算桶的大小和offset
	maxNum, minNum := math.MinInt, math.MaxInt
	for i := 0; i < len(nums); i++ {
		if nums[i] > maxNum {
			maxNum = nums[i]
		}
		if nums[i] < minNum {
			minNum = nums[i]
		}
	}

	// 一个经验值
	bucketCount := int(math.Sqrt(float64(len(nums))))
	fmt.Println(bucketCount)

	// offset具体作用：
	//	1、范围平移：将原始数据范围 [minNum, maxNum] 平移到 [0, maxNum-minNum]
	//	2、索引计算：确保桶索引从0开始，避免负数索引或索引越界
	offset := -minNum

	// 计算理论上每个桶需要装的元素个数
	bucketSize := (maxNum-minNum)/bucketCount + 1

	// 初始化桶
	bucket := make([][]int, bucketCount)

	// 把元素分配到桶里
	// 用除法向下取整的方式计算桶的索引
	for i := 0; i < len(nums); i++ {
		index := (nums[i] + offset) / bucketSize
		bucket[index] = append(bucket[index], nums[i])
	}

	// 对桶中每个元素进行排序（使用插入排序，而不是递归调用自身）
	for i := 0; i < bucketCount; i++ {
		if len(bucket[i]) > 0 {
			insertSort(bucket[i])
		}
	}

	// 合并每个桶中的元素
	index := 0
	for i := 0; i < bucketCount; i++ {
		for j := 0; j < len(bucket[i]); j++ {
			nums[index] = bucket[i][j]
			index++
		}
	}
	return nums
}

// 9、希尔排序，是对「插入排序」的简单改进，关键概念：h有序数组

//一个数组是 h 有序的，是指这个数组中任意间隔为 h（或者说间隔元素的个数为 h-1）的元素都是有序的。
//这个概念用文字不好描述清楚，直接看个例子吧。比方说 h=3 时，一个 3 有序数组是这样的：
// nums:
// [1, 2, 4, 3, 5, 7, 8, 6, 10, 9, 12, 11]
//  ^--------^--------^---------^
// 	   ^--------^--------^---------^
//        ^--------^--------^----------^
//  1--------3--------8---------9
//     2--------5--------6---------12
//        4--------7--------10---------11
//可以看到，[1,3,8,9]、[2,5,6,12]、[4,7,10,11] 这三个数组都是有序的，且元素 1, 3，元素 3, 8，元素 8, 9，元素 2, 5 等等，
//每一对儿元素的间隔是 3（或者说间隔元素的个数是 2），当一个数组完成排序的时候，其实就是 1 有序数组。

// 10、基数排序，是由「计数排序」扩展而来，原理按(位)排序

// 基数排序的主要思路是对待排序元素的每一位(个,十,百...)依次进行计数排序，由于计数排序是稳定的，所以对每一位完成排序后，所有元素就完成了排序。
// 比方说输入的数组都是三位数 nums = [329, 457, 839, 439, 720, 355, 350]
// a. 先按照个位数排序
// b. 然后按照十位数排序
// c. 然后按照百位数排序
// 最终就完成了整个数组的排序。
// 这里面的关键在于，对每一位的排序都必须是稳定排序，否则最终结果就不对了。

func main() {
	nums := []int{5, 4, 10, 3, 8, 1, 2, 6, 7, 9}
	fmt.Println("原始数组:", nums)

	//fmt.Println("选择排序:", selectSort(nums))
	//fmt.Println("冒泡排序:", bubbleSort(nums))
	//fmt.Println("冒泡排序:", bubbleSortV1(nums))
	//fmt.Println("插入排序:", insertSort(nums))
	//fmt.Println("快速排序:", quickSort(nums))
	fmt.Println("归并排序:", mergeSort(nums))
	//fmt.Println("堆排序:", heapSort(nums))
	//fmt.Println("计数排序:", countSort(nums))
	//fmt.Println("桶排序:", bucketSort(nums))
}
