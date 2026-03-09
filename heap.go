package main

import "fmt"

// MinHeap 小顶堆，基于完全二叉树结构实现
type MinHeap struct {
	array []int
}

func (h *MinHeap) Init(nums []int) {
	// 从最后一个非叶子节点开始向上构建最小堆（索引 n/2-1）
	// ❓为啥呢
	// 因为：叶子节点没有子节点，所以它们天然满足最大堆的性质。
	h.array = nums
	for i := len(h.array)/2 - 1; i >= 0; i-- {
		h.siftDown(i, len(h.array))
	}
}

func (h *MinHeap) Push(x int) {
	// 末尾加入新元素，然后上浮到正确位置
	h.array = append(h.array, x)
	h.siftUp(len(h.array) - 1)
}

func (h *MinHeap) Pop() (int, bool) {
	if len(h.array) == 0 {
		return -1, false
	}
	minVal := h.array[0]
	last := len(h.array) - 1

	// 将堆顶与末尾元素交换，缩减堆大小
	h.array[0], h.array[last] = h.array[last], h.array[0]
	h.array = h.array[:last]

	// 下沉调整
	if len(h.array) > 0 {
		h.siftDown(0, len(h.array))
	}
	return minVal, true
}

// siftDown 下沉操作，将索引 i 处的元素下沉到合适位置
func (h *MinHeap) siftDown(i, n int) {
	minIndex := i    // 假设当前节点是最小的
	left := 2*i + 1  // 左子节点索引
	right := 2*i + 2 // 右子节点索引

	if left < n && h.array[left] < h.array[minIndex] {
		minIndex = left
	}

	if right < n && h.array[right] < h.array[minIndex] {
		minIndex = right
	}

	// 如果最小值不是当前节点，交换并递归调整
	// 交换可能破坏子树的堆性质，所以需要递归下沉
	if minIndex != i {
		h.array[i], h.array[minIndex] = h.array[minIndex], h.array[i]
		h.siftDown(minIndex, n)
	}
}

// siftUp 上浮操作，将索引 i 处的元素上浮到合适位置
func (h *MinHeap) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2 // 父节点索引
		if h.array[i] >= h.array[parent] {
			// 堆性质满足，停止上浮
			break
		}
		h.array[i], h.array[parent] = h.array[parent], h.array[i]
		i = parent
	}
}

func main() {
	heap := MinHeap{}
	heap.Init([]int{1, 8, 5, 81})

	fmt.Println(heap.Pop())
	fmt.Println(heap.Pop())
	heap.Push(6)
	heap.Push(9)
	fmt.Println(heap.Pop())
	fmt.Println(heap.Pop())
	fmt.Println(heap.Pop())
	fmt.Println(heap.Pop())
	fmt.Println(heap.Pop())
}
