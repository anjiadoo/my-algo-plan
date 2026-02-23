package main

import (
	"fmt"
	"math"
)

// 排序算法实现：
// 🌟技巧1：冒泡排序可用 swap 标志提前终止，若一轮无交换说明已有序，内层循环在未排序区（冒泡排序）
// 🌟技巧2：插入排序内层找到正确位置后要 break，不需要继续往前比较，内层循环在已排序区（插入排序）
// 🌟技巧3：分区时 i 用 <=、j 用 >，不对称才能保证 i/j 相遇时停下，否则死循环（快速排序）
// 🌟技巧4：分区结束后 pivot 只能与 j 交换而非 i，因为 j 及其左边都 <= pivot（快速排序）
// 🌟技巧5：快排递归终止条件是 lo >= hi，归并是 lo == hi，快排用 == 会漏掉空区间（快排 vs 归并）
// 🌟技巧6：计算中间位置时使用 lo + (hi-lo)/2 防止整型溢出（归并排序）
// 🌟技巧7：构建堆时从最后一个非叶子节点（n/2-1）开始，叶子节点天然满足堆性质（堆排序）
// 🌟技巧8：堆排序交换堆顶到末尾后，heapify 的 n 要缩小（传 i 而非 len），否则已排好的元素会被破坏（堆排序）
// 🌟技巧9：heapify 交换后必须递归向下调整，因为交换可能破坏子树的堆性质（堆排序）
// 🌟技巧10：count 数组做前缀和累加后，存的是"结束位置"而非"起始位置"，放置元素后要 count--（计数排序）
// 🌟技巧11：从后往前遍历原数组放置元素，才能保证排序的稳定性（计数排序）
// 🌟技巧12：使用 offset=-minNum 将数据范围平移到 [0, max-min]，避免负数索引（计数排序、桶排序）
// 🌟技巧13：桶大小用 (max-min)/bucketCount+1，必须 +1 防止最大值越界（桶排序）
// 🌟技巧14：希尔排序是对「插入排序」的简单改进，关键概念：h有序数组（希尔排序）
// 🌟技巧15：基数排序是由「计数排序」扩展而来，原理按(位)排序（基数排序）

// 0、func selectSort(nums []int) []int
// 1、func bubbleSort(nums []int) []int
// 2、func insertSort(nums []int) []int
// 3、func quickSort(nums []int) []int
// 4、func partition(nums []int, lo, hi int) int
// 5、func mergeSort(nums []int) []int
// 6、func merge(nums []int, lo, mid, hi int)
// 7、func heapSort(nums []int) []int
// 8、func heapify(nums []int, n, i int)
// 9、func countSort(nums []int) []int
// 10、func bucketSort(nums []int) []int

// 1、选择排序-顺着从头遍历到尾-原地排序-不稳定排序
func selectSort(nums []int) []int {
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
	var sort func(nums []int, lo, hi int)

	sort = func(nums []int, lo, hi int) {
		// 单个元素不用排序
		if lo == hi {
			return
		}
		// 相比(lo+hi)/2来说，这种方式可防止溢出
		mid := lo + (hi-lo)/2
		sort(nums, lo, mid)
		sort(nums, mid+1, hi)
		merge(nums, lo, mid, hi)
	}

	sort(nums, 0, len(nums)-1)
	return nums
}

func merge(nums []int, lo, mid, hi int) {
	//todo
}

// 6、堆排序-完全二叉树结构-原地排序-不稳定排序
func heapSort(nums []int) []int {
	// 1.构建最大堆（或最小堆）
	// 2.将堆顶元素（最大值）与最后一个元素交换
	// 3.调整堆，使其重新满足堆的性质
	// 4.重复步骤2-3，直到所有元素有序

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
	// 原理：统计每种元素出现的次数，进而推算出每个元素在排序后数组中的索引位置，最终完成排序。
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
		count[num+offset]++
	}

	// 累加 count 数组，得到的是 nums[i] 在排序后的数组中的结束位置
	for i := 1; i < len(count); i++ {
		count[i] += count[i-1]
	}

	// 根据每个元素排序后的索引位置，完成排序
	// 这里注意，我们从后往前遍历 nums，是为了保证排序的稳定性
	sorted := make([]int, len(nums))
	for i := len(nums) - 1; i >= 0; i-- {
		sorted[count[nums[i]+offset]-1] = nums[i]
		count[nums[i]+offset]--
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
	//fmt.Println("堆排序:", heapSort(nums))
	//fmt.Println("计数排序:", countSort(nums))
	fmt.Println("桶排序:", bucketSort(nums))

}
