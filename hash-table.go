package main

import (
	"container/list"
	"fmt"
)

// 哈希表实现（链地址法解决冲突）：
// 🌟技巧1：hash函数用取模法 key % capacity，简单高效
// 🌟技巧2：使用container/list标准双向链表作为桶，避免手写链表
// 0、func NewMyHashTable(capacity int) *MyHashTable
// 1、func (m *MyHashTable) Get(key int) (int, bool)
// 2、func (m *MyHashTable) Put(key, val int)
// 3、func (m *MyHashTable) Remove(key int)
// 4、func (m *MyHashTable) Display()

type KVNode struct {
	key   int
	value int
}

type MyHashTable struct {
	table []*list.List
}

func NewMyHashTable(capacity int) *MyHashTable {
	return &MyHashTable{
		table: make([]*list.List, capacity),
	}
}

func (m *MyHashTable) hash(key int) int {
	return key % len(m.table)
}

func (m *MyHashTable) Get(key int) (int, bool) {
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

func (m *MyHashTable) Put(key, val int) {
	hashCode := m.hash(key)
	if m.table[hashCode] == nil {
		m.table[hashCode] = list.New()
		m.table[hashCode].PushFront(&KVNode{key, val})
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
	m.table[hashCode].PushFront(&KVNode{key: key, value: val})
}

func (m *MyHashTable) Remove(key int) {
	hashCode := m.hash(key)
	if m.table[hashCode] == nil {
		return
	}
	for e := m.table[hashCode].Front(); e != nil; e = e.Next() {
		node := e.Value.(*KVNode)
		if node.key == key {
			m.table[hashCode].Remove(e)
			break
		}
	}
}
func (m *MyHashTable) Display() {
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
	hashTable := NewMyHashTable(10)

	hashTable.Put(1, 10)
	hashTable.Put(5, 50)
	hashTable.Put(9, 90)
	hashTable.Display()

	hashTable.Put(11, 100)
	hashTable.Put(55, 500)
	hashTable.Put(99, 900)
	hashTable.Display()

	fmt.Println(hashTable.Get(10))
	fmt.Println(hashTable.Get(5))

	hashTable.Remove(100)
	hashTable.Remove(1)
	hashTable.Display()
}
