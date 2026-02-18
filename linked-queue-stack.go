package main

import "fmt"

// MyLinkedQueue 链式队列
type MyLinkedQueue struct {
	list *MyLinkedList
}

func NewMyLinkedQueue() *MyLinkedQueue {
	return &MyLinkedQueue{
		list: NewMyLinkedList(),
	}
}

func (q *MyLinkedQueue) Push(val int) {
	q.list.AddAtTail(val)
}

func (q *MyLinkedQueue) Pop() int {
	return q.list.DeleteAtIndex(0)
}

func (q *MyLinkedQueue) Peek() int {
	return q.list.Get(0)
}

func (q *MyLinkedQueue) Size() int {
	return q.list.Size()
}

func (q *MyLinkedQueue) Display() {
	q.list.Display()
}

// MyLinkedStack 链式栈
type MyLinkedStack struct {
	list *MyLinkedList
}

func NewMyLinkedStack() *MyLinkedStack {
	return &MyLinkedStack{
		list: NewMyLinkedList(),
	}
}

func (s *MyLinkedStack) Push(val int) {
	s.list.AddAtHead(val)
}

func (s *MyLinkedStack) Pop() int {
	return s.list.DeleteAtIndex(0)
}

func (s *MyLinkedStack) Peek() int {
	return s.list.Get(0)
}

func (s *MyLinkedStack) Size() int {
	return s.list.Size()
}

func (s *MyLinkedStack) Display() {
	s.list.Display()
}

func main() {
	fmt.Println("--------Queue---------")
	queue := NewMyLinkedQueue()
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
	stack := NewMyLinkedStack()
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
