package main

import (
	"errors"
	"fmt"
)

// 循环数组
// 🌟技巧1：左闭右开区间[x,y)表示范围，x表示第一个有效元素索引，y表示最后一个有效元素的后一个索引
// 🌟技巧2：数组尾部元素获取下标：(end-1+size)%size
// 🌟技巧3：所有涉及到start/end的更新操作，都需要在取模的的基础上操作

type CycleArray[T any] struct {
	array []T
	start int
	end   int
	count int
	size  int
}

func NewCycleArray[T any](size int) *CycleArray[T] {
	return &CycleArray[T]{
		array: make([]T, size),
		size:  size,
	}
}

// 获取数组头部元素，时间复杂度 O(1)
func (ca *CycleArray[T]) getFirst() (T, error) {
	if ca.count == 0 {
		return *new(T), errors.New("array is empty")
	}
	return ca.array[ca.start], nil
}

// 获取数组尾部元素，时间复杂度 O(1)
func (ca *CycleArray[T]) getLast() (T, error) {
	if ca.count == 0 {
		return *new(T), errors.New("array is empty")
	}
	return ca.array[(ca.end-1+ca.size)%ca.size], nil
}

// 自动扩缩容辅助函数
func (ca *CycleArray[T]) resize(newSize int) {
	newArr := make([]T, newSize)
	for i := 0; i < ca.count; i++ {
		newArr[i] = ca.array[(ca.start+i)%ca.size]
	}
	ca.start = 0
	ca.end = ca.count
	ca.array = newArr
	ca.size = newSize
}

// 在数组头部添加元素，时间复杂度 O(1)
func (ca *CycleArray[T]) addFirst(val T) {
	//是否已满
	if ca.count == ca.size {
		ca.resize(ca.size * 2)
	}
	ca.start = (ca.start - 1 + ca.size) % ca.size
	ca.array[ca.start] = val
	ca.count++
}

// 删除数组头部元素，时间复杂度 O(1)
func (ca *CycleArray[T]) removeFirst() error {
	if ca.count == 0 {
		return errors.New("array is empty")
	}
	ca.array[ca.start] = *new(T)
	ca.start = (ca.start + 1) % ca.size
	ca.count--

	// 如果数组元素数量减少到原大小的四分之一，则减小数组大小为一半
	if ca.count > 0 && ca.count == ca.size/4 {
		ca.resize(ca.size / 2)
	}
	return nil
}

// 在数组尾部添加元素，时间复杂度 O(1)
func (ca *CycleArray[T]) addLast(val T) {
	//是否已满
	if ca.count == ca.size {
		ca.resize(ca.size * 2)
	}

	ca.array[ca.end] = val
	ca.end = (ca.end + 1) % ca.size
	ca.count++
}

// 删除数组尾部元素，时间复杂度 O(1)
func (ca *CycleArray[T]) removeLast() error {
	if ca.count == 0 {
		return errors.New("array is empty")
	}

	ca.end = (ca.end - 1 + ca.size) % ca.size
	ca.array[ca.end] = *new(T)
	ca.count--
	return nil
}

func (ca *CycleArray[T]) display() {
	fmt.Println("size=", ca.size, "count=", ca.count, "start=", ca.start, "end=", ca.end, ca.array)
}

func main() {
	ca := NewCycleArray[int](5)

	ca.addLast(1)
	ca.addLast(2)
	ca.addLast(3)
	ca.display()

	ca.addFirst(4)
	ca.addFirst(5)
	ca.addFirst(6)
	ca.display()

	_ = ca.removeFirst()
	_ = ca.removeFirst()
	ca.display()

	_ = ca.removeLast()
	_ = ca.removeLast()
	ca.display()

	fmt.Println(ca.getFirst())
	fmt.Println(ca.getLast())
}
