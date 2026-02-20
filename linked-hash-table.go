package main

import (
	"container/list"
	"fmt"
)

// 哈希表实现（链地址法解决冲突）：
// 🌟技巧1：hash函数用取模法 key % capacity，简单高效
// 🌟技巧2：使用container/list标准双向链表作为桶，避免手写链表
// 🌟技巧3：新增head/tail节点，把普通节点链接在一起形成一个双链表，就可以顺序访问
// 0、func NewMyLinkedHashTable(capacity int) *MyLinkedHashTable
// 1、func (m *MyLinkedHashTable) Get(key int) (int, bool)
// 2、func (m *MyLinkedHashTable) Put(key, val int)
// 3、func (m *MyLinkedHashTable) Remove(key int)
// 4、func (m *MyLinkedHashTable) Display()
// 5、func (m *MyLinkedHashTable) Keys() []int
// 6、func (m *MyLinkedHashTable) Values() []int

type KVNode struct {
	key   int
	value int
	prev  *KVNode // 新增-链式哈希
	next  *KVNode // 新增-链式哈希
}

type MyLinkedHashTable struct {
	head  *KVNode // 新增-链式哈希
	tail  *KVNode // 新增-链式哈希
	table []*list.List
}

func NewMyLinkedHashTable(capacity int) *MyLinkedHashTable {
	head, tail := &KVNode{}, &KVNode{}
	head.next, tail.prev = tail, head
	return &MyLinkedHashTable{
		head:  head,
		tail:  tail,
		table: make([]*list.List, capacity),
	}
}

func (m *MyLinkedHashTable) hash(key int) int {
	return key % len(m.table)
}

func (m *MyLinkedHashTable) Get(key int) (int, bool) {
	hashCode := m.hash(key)
	if m.table[hashCode] == nil {
		return -1, false
	}
	for e := m.table[hashCode].Front(); e != nil; e = e.Next() {
		node := e.Value.(*KVNode)
		if node.key == key {
			return node.value, true
		}
	}
	return -1, false
}

func (m *MyLinkedHashTable) Put(key, val int) {
	hashCode := m.hash(key)
	if m.table[hashCode] == nil {
		m.table[hashCode] = list.New()

		node := &KVNode{key: key, value: val}
		node.prev = m.tail.prev
		node.next = m.tail

		m.tail.prev.next = node
		m.tail.prev = node

		m.table[hashCode].PushFront(node)
		return
	}
	for e := m.table[hashCode].Front(); e != nil; e = e.Next() {
		node := e.Value.(*KVNode)
		if node.key == key {
			node.value = val
			return
		}
	}
	// 链表中没有目标 key，添加新节点
	node := &KVNode{key: key, value: val}
	node.prev = m.tail.prev
	node.next = m.tail

	m.tail.prev.next = node
	m.tail.prev = node

	m.table[hashCode].PushFront(node)
}

func (m *MyLinkedHashTable) Remove(key int) {
	hashCode := m.hash(key)
	if m.table[hashCode] == nil {
		return
	}
	for e := m.table[hashCode].Front(); e != nil; e = e.Next() {
		node := e.Value.(*KVNode)
		if node.key == key {
			prev := node.prev
			next := node.next

			prev.next = next
			next.prev = prev

			m.table[hashCode].Remove(e)
			break
		}
	}
}

func (m *MyLinkedHashTable) Keys() []int {
	var keys []int
	p := m.head.next
	for p != nil {
		keys = append(keys, p.key)
		p = p.next
	}
	return keys
}

func (m *MyLinkedHashTable) Values() []int {
	var values []int
	p := m.head.next
	for p != nil {
		values = append(values, p.value)
		p = p.next
	}
	return values
}

func (m *MyLinkedHashTable) Display() {
	for i := 0; i < len(m.table); i++ {
		if m.table[i] == nil {
			continue
		}

		var values []int
		var keys []int

		for e := m.table[i].Front(); e != nil; e = e.Next() {
			node := e.Value.(*KVNode)
			keys = append(keys, node.key)
			values = append(values, node.value)
		}
		if len(values) > 0 {
			fmt.Printf("hash=%d keys=%v values=%+v\n", i, keys, values)
		}
	}
	fmt.Println("=>")
}

func main() {
	hashTable := NewMyLinkedHashTable(10)

	hashTable.Put(1, 10)
	hashTable.Put(5, 50)
	hashTable.Put(9, 90)
	hashTable.Display()
	fmt.Println(hashTable.Keys(), hashTable.Values())

	hashTable.Put(11, 100)
	hashTable.Put(55, 500)
	hashTable.Put(99, 900)
	hashTable.Display()
	fmt.Println(hashTable.Keys(), hashTable.Values())

	fmt.Println(hashTable.Get(10))
	fmt.Println(hashTable.Get(5))

	hashTable.Remove(100)
	hashTable.Remove(1)
	hashTable.Display()
	fmt.Println(hashTable.Keys(), hashTable.Values())
}
