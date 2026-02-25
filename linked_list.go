package main

import "fmt"

// 双链表实现：
// 🌟技巧1：虚拟头尾节点技巧 - 设置虚拟头节点和尾节点
// 🌟技巧2：新节点先链接技巧 - 新增节点可先设置prev、next，没有副作用，已有节点修改prev、next时需注意先后顺序

// 0、func NewMyLinkedList() *MyLinkedList
// 1、func (m *MyLinkedList) Get(index int) int
// 2、func (m *MyLinkedList) AddAtHead(val int)
// 3、func (m *MyLinkedList) AddAtTail(val int)
// 4、func (m *MyLinkedList) AddAtIndex(index int, val int)
// 5、func (m *MyLinkedList) DeleteAtIndex(index int)
// 6、func (m *MyLinkedList) Display()

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

func NewMyLinkedList() *MyLinkedList {
	head := &Node{}
	tail := &Node{}
	head.next, tail.prev = tail, head
	return &MyLinkedList{
		head: head,
		tail: tail,
		size: 0,
	}
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
	newNode := &Node{val: val}
	newNode.prev = m.head
	newNode.next = m.head.next

	m.head.next.prev = newNode
	m.head.next = newNode

	m.size++
}

func (m *MyLinkedList) AddAtTail(val int) {
	newNode := &Node{val: val}
	newNode.prev = m.tail.prev
	newNode.next = m.tail

	m.tail.prev.next = newNode
	m.tail.prev = newNode

	m.size++
}

func (m *MyLinkedList) AddAtIndex(index int, val int) {
	if !(index >= 0 && index <= m.size) {
		return
	}

	if index == 0 {
		m.AddAtHead(val)
		return
	}
	if index == m.size {
		m.AddAtTail(val)
		return
	}

	p := m.head.next
	for i := 0; i < index-1; i++ {
		p = p.next
	}

	newNode := &Node{val: val}
	newNode.prev = p
	newNode.next = p.next

	p.next.prev = newNode
	p.next = newNode

	m.size++
}

func (m *MyLinkedList) DeleteAtIndex(index int) int {
	if m.size == 0 || index < 0 || index >= m.size {
		return -1
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

	del.prev = nil
	del.next = nil
	return del.val
}
func (m *MyLinkedList) Size() int {
	return m.size
}

func (m *MyLinkedList) Display() {
	var arr1 []int
	p := m.head.next
	for p != nil && p != m.tail {
		arr1 = append(arr1, p.val)
		p = p.next
	}

	var arr2 []int
	q := m.tail.prev
	for q != nil && q != m.head {
		arr2 = append(arr2, q.val)
		q = q.prev
	}
	fmt.Printf("size=%d 前序遍历：%v 后序遍历：%v\n", m.size, arr1, arr2)
}

func _main() {
	list := NewMyLinkedList()
	list.Display()

	list.AddAtHead(100)
	list.AddAtHead(200)
	list.AddAtHead(300)
	list.Display()

	list.AddAtIndex(0, 10)
	list.Display()

	list.AddAtIndex(0, 20)
	list.Display()

	list.AddAtIndex(1, 30)
	list.Display()

	val := list.Get(0)
	fmt.Println("期望是20，实际是：", val)

	list.AddAtTail(400)
	list.AddAtTail(500)
	list.AddAtTail(600)
	list.Display()

	list.DeleteAtIndex(0)
	list.Display()

	list.DeleteAtIndex(5)
	list.Display()

	list.DeleteAtIndex(0)
	list.Display()

}
