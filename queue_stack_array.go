package main

import "fmt"

// 数组队列和栈实现（使用泛型）：
// 🌟技巧1：切片截取出队技巧 - 队列出队用切片重新分配 q.queue[1:]，简单直观
// 🌟技巧2：切片动态入出栈技巧 - 栈用append入栈，出栈取最后一个元素后缩减切片，利用slice的动态特性

// 队列：
// 0、func NewMyArrayQueue[T any]() *MyArrayQueue[T]
// 1、func (q *MyArrayQueue[T]) Push(val T)
// 2、func (q *MyArrayQueue[T]) Pop() T
// 3、func (q *MyArrayQueue[T]) Peek() T
// 4、func (q *MyArrayQueue[T]) Size() int
// 5、func (q *MyArrayQueue[T]) Display()

// 栈：
// 6、func NewMyArrayStack[T any]() *MyArrayStack[T]
// 7、func (s *MyArrayStack[T]) Push(val T)
// 8、func (s *MyArrayStack[T]) Pop() T
// 9、func (s *MyArrayStack[T]) Peek() T
// 10、func (s *MyArrayStack[T]) Size() int
// 11、func (s *MyArrayStack[T]) Display()

type MyArrayQueue[T any] struct {
	queue []T
}

func NewMyArrayQueue[T any]() *MyArrayQueue[T] {
	return &MyArrayQueue[T]{
		queue: make([]T, 0),
	}
}

func (q *MyArrayQueue[T]) Push(val T) {
	q.queue = append(q.queue, val)
}

func (q *MyArrayQueue[T]) Pop() T {
	if q.Size() == 0 {
		return *new(T)
	}
	val := q.queue[0]
	q.queue = q.queue[1:]
	return val
}

func (q *MyArrayQueue[T]) Peek() T {
	if q.Size() == 0 {
		return *new(T)
	}
	return q.queue[0]
}

func (q *MyArrayQueue[T]) Size() int {
	return len(q.queue)
}

func (q *MyArrayQueue[T]) Display() {
	fmt.Printf("size=%d array=%v\n", q.Size(), q.queue)
}

type MyArrayStack[T any] struct {
	stack []T
}

func NewMyArrayStack[T any]() *MyArrayStack[T] {
	return &MyArrayStack[T]{
		stack: make([]T, 0),
	}
}

func (s *MyArrayStack[T]) Push(val T) {
	s.stack = append(s.stack, val)
}

func (s *MyArrayStack[T]) Pop() T {
	if s.Size() == 0 {
		return *new(T)
	}
	e := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return e
}

func (s *MyArrayStack[T]) Peek() T {
	if s.Size() == 0 {
		return *new(T)
	}
	e := s.stack[len(s.stack)-1]
	return e
}

func (s *MyArrayStack[T]) Size() int {
	return len(s.stack)
}

func (s *MyArrayStack[T]) Display() {
	fmt.Printf("size=%d array=%v\n", s.Size(), s.stack)
}

func main() {
	fmt.Println("--------Queue---------")
	queue := NewMyArrayQueue[int]()
	queue.Push(1)
	queue.Push(2)
	queue.Push(3)
	queue.Push(40)
	queue.Push(50)
	queue.Display()

	fmt.Println(queue.Peek())
	fmt.Println(queue.Size())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())

	queue.Display()

	fmt.Println("--------Stack---------")
	stack := NewMyArrayStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Push(40)
	stack.Push(50)

	stack.Display()

	fmt.Println(stack.Peek())
	fmt.Println(stack.Size())
	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop())

	stack.Display()
}
