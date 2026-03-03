package main

import "fmt"

// 双链表实现：
// 🌟技巧1：虚拟头尾节点技巧 - 设置虚拟头节点和尾节点，所有操作都不需要特判边界（头部/尾部插入和删除与中间节点操作完全一致），大幅简化代码
// 🌟技巧2：新节点先链接技巧 - 新增节点可先设置自身的prev、next（无副作用），再修改已有节点的指针；已有节点改指针时需注意先后顺序，避免断链后找不到原邻居
// 🌟技巧3：双向遍历技巧 - 双链表可从头或从尾遍历，利用这一特性可根据index大小选择从较近的一端出发，将平均查找复杂度从O(n)优化到O(n/2)
// 🌟技巧4：size字段维护 - 维护size字段避免每次遍历计算长度，所有增删操作都必须同步更新size，越界检查直接用size判断

// ⚠️易错点1：删除节点必须断开指针 - 删除节点后需将del.Prev和del.Next置nil，防止被删节点仍持有对链表节点的引用，导致内存泄漏或逻辑错误
// ⚠️易错点2：修改已有节点的指针顺序 - 插入新节点时，必须先修改与新节点相连的「远端」节点指针（如p.Next.Prev），再修改「近端」节点指针（如p.Next）；否则近端指针覆盖后找不到远端节点
// ⚠️易错点3：index边界语义 - 要明确index是0-based还是1-based，AddAtIndex插入到index节点「之后」还是「之前」；本实现中index=0表示插入到第一个真实节点之后
// ⚠️易错点4：虚拟节点不计入size - head和tail是虚拟节点，size只统计真实节点数量，遍历时起点是head.Next而非head，终止条件是p!=tail而非p==nil

// 何时运用双链表：
// ❓1、是否需要O(1)时间在已知节点的前后插入/删除？双链表比单链表多了Prev指针，无需从头找前驱节点
// ❓2、是否需要同时支持从头和从尾的高效访问？如实现LRU缓存（配合哈希表达到O(1)的get和put）
// ❓3、是否需要频繁移动节点（如将某节点移到头部/尾部）？双链表可O(1)完成，单链表需O(n)找前驱

// 0、func NewMyLinkedList() *MyLinkedList
// 1、func (m *MyLinkedList) Get(index int) int
// 2、func (m *MyLinkedList) AddAtHead(val int)
// 3、func (m *MyLinkedList) AddAtTail(val int)
// 4、func (m *MyLinkedList) AddAtIndex(index int, val int)
// 5、func (m *MyLinkedList) DeleteAtIndex(index int)
// 6、func (m *MyLinkedList) Size() int
// 7、func (m *MyLinkedList) Display()

type Node struct {
	Val  int
	Prev *Node
	Next *Node
}

type MyLinkedList struct {
	head *Node
	tail *Node
	size int
}

func NewMyLinkedList(nums []int) *MyLinkedList {
	head, tail := &Node{}, &Node{}
	head.Next = tail
	tail.Prev = head

	l := &MyLinkedList{
		head: head,
		tail: tail,
	}

	for i := 0; i < len(nums); i++ {
		l.AddAtTail(nums[i])
	}
	return l
}

func (m *MyLinkedList) Get(index int) int {
	if !(index >= 0 && index < m.size) {
		return -1
	}
	p := m.head.Next
	for i := 0; i < index; i++ {
		p = p.Next
	}
	return p.Val
}

func (m *MyLinkedList) AddAtHead(val int) {
	newNode := &Node{Val: val}
	newNode.Prev = m.head
	newNode.Next = m.head.Next

	m.head.Next.Prev = newNode
	m.head.Next = newNode
	m.size++
}

func (m *MyLinkedList) AddAtTail(val int) {
	newNode := &Node{Val: val}
	newNode.Prev = m.tail.Prev
	newNode.Next = m.tail

	m.tail.Prev.Next = newNode
	m.tail.Prev = newNode
	m.size++
}

func (m *MyLinkedList) AddAtIndex(index int, val int) {
	if !(index >= 0 && index < m.size) {
		return
	}
	if index == 0 {
		m.AddAtHead(val)
		return
	}
	if index == m.size-1 {
		m.AddAtTail(val)
		return
	}

	// 注意起点和边界
	p := m.head.Next
	for i := 0; i < index; i++ {
		p = p.Next
	}

	newNode := &Node{Val: val}
	newNode.Prev = p
	newNode.Next = p.Next

	p.Next.Prev = newNode //必须先修改「远端」节点指针
	p.Next = newNode      //再修改「近端」节点指针
	m.size++
}

func (m *MyLinkedList) DeleteAtIndex(index int) int {
	if !(index >= 0 && index < m.size) {
		return -1
	}

	// 注意起点和边界
	del := m.head.Next
	for i := 0; i < index; i++ {
		del = del.Next
	}

	prev := del.Prev
	next := del.Next

	prev.Next = next
	next.Prev = prev

	del.Prev = nil
	del.Next = nil
	m.size--
	return del.Val
}

func (m *MyLinkedList) Size() int {
	return m.size
}

func (m *MyLinkedList) Display() {
	str1 := "head<->"
	p := m.head.Next
	for p != m.tail {
		str1 += fmt.Sprintf("%d<->", p.Val)
		p = p.Next
	}
	str1 += "tail"

	str2 := "tail<->"
	p = m.tail.Prev
	for p != m.head {
		str2 += fmt.Sprintf("%d<->", p.Val)
		p = p.Prev
	}
	str2 += "head"

	fmt.Printf("打印双链表(%d): 顺序:%s 逆序:%s\n", m.size, str1, str2)
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	list := NewMyLinkedList(nums)
	list.Display()

	fmt.Println(list.Get(0))
	fmt.Println(list.Get(4))

	//list.AddAtHead(666)
	//list.AddAtTail(999)
	//list.Display()

	//list.DeleteAtIndex(0)
	//list.DeleteAtIndex(3)
	//list.Display()

	list.AddAtIndex(2, 222)
	list.AddAtIndex(2, 22)
	list.Display()
}
