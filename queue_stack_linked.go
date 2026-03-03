package main

import "fmt"

// 链式队列和栈实现（基于MyLinkedList双链表）：
// 队列：
// 0、func NewMyLinkedQueue() *MyLinkedQueue
// 1、func (q *MyLinkedQueue) Push(val int)
// 2、func (q *MyLinkedQueue) Pop() int
// 3、func (q *MyLinkedQueue) Peek() int
// 4、func (q *MyLinkedQueue) Size() int
// 5、func (q *MyLinkedQueue) Display()
// 栈：
// 6、func NewMyLinkedStack() *MyLinkedStack
// 7、func (s *MyLinkedStack) Push(val int)
// 8、func (s *MyLinkedStack) Pop() int
// 9、func (s *MyLinkedStack) Peek() int
// 10、func (s *MyLinkedStack) Size() int
// 11、func (s *MyLinkedStack) Display()

// MyLinkedQueue 链式队列
type MyLinkedQueue struct {
	list *MyLinkedList
}

func NewMyLinkedQueue() *MyLinkedQueue {
	return &MyLinkedQueue{
		list: NewMyLinkedList([]int{}),
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
		list: NewMyLinkedList([]int{}),
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
