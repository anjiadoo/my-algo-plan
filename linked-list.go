package main

import "fmt"

// 双链表实现：
// 🌟技巧1：设置虚拟头节点和尾节点
// 🌟技巧2: 新增节点可先设置prev、next，没有副作用，已有节点修改prev、next时需注意先后顺序

type Node struct {
	val  int
	prev *Node
	next *Node
}

type MyLinkedList struct {
	head *Node
	tail *Node
	size int
}

func Constructor() MyLinkedList {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head
	return MyLinkedList{head: head, tail: tail, size: 0}
}

func (m *MyLinkedList) Get(index int) int {
	if !(index >= 0 && index < m.size) {
		return -1
	}
	p := m.head.next
	for i := 0; i < index; i++ {
		p = p.next
	}
	return p.val
}

func (m *MyLinkedList) AddAtHead(val int) {
	next := m.head.next

	newNode := &Node{val: val}
	newNode.prev = m.head
	newNode.next = next

	m.head.next = newNode
	next.prev = newNode

	m.size++
}

func (m *MyLinkedList) AddAtTail(val int) {
	prev := m.tail.prev

	newNode := &Node{val: val}
	newNode.prev = prev
	newNode.next = m.tail

	prev.next = newNode
	m.tail.prev = newNode

	m.size++
}

func (m *MyLinkedList) AddAtIndex(index int, val int) {
	if !(index >= 0 && index <= m.size) {
		return
	}

	// 在头节点添加
	if index == 0 {
		m.AddAtHead(val)
		return
	}
	// 在尾节点添加
	if index == m.size {
		m.AddAtTail(val)
		return
	}

	// index=0(即头节点)，下面的for循环不生效(不特殊处理头节点‼️)
	p := m.head.next
	for i := 0; i < index-1; i++ {
		p = p.next
	}

	newNode := &Node{val: val}
	newNode.prev = p
	newNode.next = p.next

	// index=size(即尾节点)，下面代码会panic(不特殊处理尾节点‼️)
	p.next = newNode
	newNode.next.prev = newNode

	m.size++
}

func (m *MyLinkedList) DeleteAtIndex(index int) {
	if !(index >= 0 && index < m.size) {
		return
	}

	del := m.head.next
	for i := 0; i < index; i++ {
		del = del.next
	}

	prev := del.prev
	next := del.next

	prev.next = next
	next.prev = prev

	m.size--
}

func (m *MyLinkedList) Display() {
	fmt.Printf("size=%d, ", m.size)
	p := m.head.next
	for p != m.tail {
		fmt.Printf("%v <-> ", p.val)
		p = p.next
	}
	fmt.Println("null")
}

func main() {
	//["MyLinkedList","addAtIndex","addAtIndex","addAtIndex","get"]
	//[[],[0,10],[0,20],[1,30],[0]]
	list := Constructor()

	list.AddAtIndex(0, 10)
	list.Display()

	list.AddAtIndex(0, 20)
	list.Display()

	list.AddAtIndex(1, 30)
	list.Display()

	val := list.Get(0)
	fmt.Println("期望是20，实际是：", val)

}
