package main

import "fmt"

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

func (s *MyArrayStack[T]) Push(val T) {}

func (s *MyArrayStack[T]) Pop() T {
	return *new(T)
}

func (s *MyArrayStack[T]) Peek() T {
	return *new(T)
}

func (s *MyArrayStack[T]) Size() T {
	return *new(T)
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
