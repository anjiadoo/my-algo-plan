package main

import (
	"errors"
	"fmt"
)

// 循环数组
// 🌟技巧1：左闭右开区间[x,y)表示范围，x表示第一个有效元素索引，y表示最后一个有效元素的后一个索引
// 🌟技巧2：数组尾部元素获取下标：(end-1+size)%size
// 🌟技巧3：所有涉及到start/end的更新操作，都需要在取模的的基础上操作
// 1、NewCycleArray[T any](size int) *MyCycleArray[T]
// 2、func (m *MyCycleArray[T]) getFirst() (T, error)
// 3、func (m *MyCycleArray[T]) getLast() (T, error)
// 4、func (m *MyCycleArray[T]) resize(newSize int)
// 5、func (m *MyCycleArray[T]) addFirst(val T)
// 6、func (m *MyCycleArray[T]) addLast(val T)
// 7、func (m *MyCycleArray[T]) removeFirst() error
// 8、func (m *MyCycleArray[T]) removeLast() error
// 9、func (m *MyCycleArray[T]) display()

type MyCycleArray[T any] struct {
	array []T
	start int
	end   int
	count int
	size  int
}

func NewMyCycleArray[T any](size int) *MyCycleArray[T] {
	return &MyCycleArray[T]{
		array: make([]T, size),
		size:  size,
	}
}

func (m *MyCycleArray[T]) getFirst() (T, error) {
	if m.len() == 0 {
		return *new(T), errors.New("array is empty")
	}
	return m.array[m.start], nil
}

func (m *MyCycleArray[T]) getLast() (T, error) {
	if m.len() == 0 {
		return *new(T), errors.New("array is empty")
	}
	return m.array[(m.end-1+m.size)%m.size], nil
}

func (m *MyCycleArray[T]) resize(newSize int) {
	newArray := make([]T, newSize)
	for i := 0; i < m.count; i++ {
		newArray[i] = m.array[(m.start+i+m.size)%m.size]
	}
	m.start = 0
	m.end = m.count
	m.size = newSize
	m.array = newArray
}

func (m *MyCycleArray[T]) addFirst(val T) {
	if m.len() == m.size {
		m.resize(m.size * 2)
	}
	m.start = (m.start - 1 + m.size) % m.size
	m.array[m.start] = val
	m.count++
}

func (m *MyCycleArray[T]) addLast(val T) {
	if m.len() == m.size {
		m.resize(m.size * 2)
	}
	m.array[m.end] = val
	m.end = (m.end + 1 + m.size) % m.size
	m.count++
}

func (m *MyCycleArray[T]) removeFirst() error {
	if m.len() == 0 {
		return errors.New("array is empty")
	}
	m.array[m.start] = *new(T)
	m.start = (m.start + 1 + m.size) % m.size
	m.count--
	return nil
}

func (m *MyCycleArray[T]) removeLast() error {
	if m.len() == 0 {
		return errors.New("array is empty")
	}
	m.end = (m.end - 1 + m.size) % m.size
	m.array[m.end] = *new(T)
	m.count--
	return nil
}

func (m *MyCycleArray[T]) len() int {
	return m.count
}

func (m *MyCycleArray[T]) display() {
	fmt.Printf("start=%d end=%d count=%d size=%d array=%v\n", m.start, m.end, m.count, m.size, m.array)
}

func main() {
	//start=0 end=3 count=3 size=5 array=[1 2 3 0 0]
	//start=9 end=5 count=6 size=10 array=[5 4 1 2 3 0 0 0 0 6]
	//start=1 end=5 count=4 size=10 array=[0 4 1 2 3 0 0 0 0 0]
	//start=1 end=3 count=2 size=10 array=[0 4 1 0 0 0 0 0 0 0]
	//4 <nil>
	//1 <nil>

	m := NewMyCycleArray[int](5)

	m.addLast(1)
	m.addLast(2)
	m.addLast(3)
	m.display()

	m.addFirst(4)
	m.addFirst(5)
	m.addFirst(6)
	m.display()

	_ = m.removeFirst()
	_ = m.removeFirst()
	m.display()

	_ = m.removeLast()
	_ = m.removeLast()
	m.display()

	fmt.Println(m.getFirst())
	fmt.Println(m.getLast())
}
