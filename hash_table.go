/*
 * ============================================================================
 *                📘 哈希表（LinkedHashMap 简化版）· 核心记忆框架
 * ============================================================================
 * 【一句话理解 LinkedHashMap】
 *
 *   LinkedHashMap = 哈希表（O(1) 查改）+ 双向链表（保持插入顺序）
 *   桶负责 O(1) 定位，顺序链负责按序遍历，KVNode 同时挂在两个结构上。
 * ────────────────────────────────────────────────────────────────────────────
 * 【核心数据结构】
 *   table（哈希桶数组）                   head/tail（插入顺序链）
 *   ┌─────┐                             head ↔ node1 ↔ node2 ↔ node3 ↔ tail
 *   │  0  │→ nil
 *   │  1  │→ [node1] → [node3] → nil    ← 冲突链（链地址法）
 *   │  2  │→ nil
 *   │  3  │→ [node2] → nil
 *   └─────┘
 *   每个 KVNode 同时挂在两个地方：
 *     ① table[hash] 的 container/list 冲突链 —— 负责 O(1) 查找
 *     ② head ↔ tail 双向链表 —— 负责按插入顺序遍历
 * ────────────────────────────────────────────────────────────────────────────
 * 【哨兵节点 + 尾插法】
 *
 *   head 和 tail 是两个空哨兵节点，初始状态：head.next = tail，tail.prev = head
 *
 *   插入新节点（统一插在 tail 前）：
 *     node.prev = tail.prev   →   node.next = tail
 *     tail.prev.next = node   →   tail.prev = node
 *
 *   无论链表是否为空，这 4 行代码永远正确，无需判空——这就是哨兵的价值。
 * ────────────────────────────────────────────────────────────────────────────
 * 【Put 操作的两条路径】
 *
 *   路径①：桶为空 → 懒初始化桶 → 创建节点 → 挂入顺序链 → 挂入桶链
 *   路径②：桶非空 → 遍历冲突链
 *               ├── key 存在 → 只更新 value，不动顺序链（保持原始插入位置）
 *               └── key 不存在 → 创建节点 → 挂入顺序链 → 挂入桶链
 *
 *   ⚠️ 易错点1：key 已存在时只改 value，不重新挂入顺序链
 *      若重新插入，节点位置会变成"最新"，破坏 LinkedHashMap 的插入顺序语义，
 *      且原节点仍挂在顺序链上无法被访问，造成内存泄漏。
 * ────────────────────────────────────────────────────────────────────────────
 * 【Remove 操作的双重摘除】
 *
 *   必须同时从两个结构中删除节点：
 *     ① 从顺序链摘除：prev.next = next；next.prev = prev
 *     ② 从桶链摘除：table[hash].Remove(e)
 *
 *   ⚠️ 易错点2：顺序链摘除时，哨兵保证 prev/next 永不为 nil
 *      即使删除第一个或最后一个真实节点，prev 是 head，next 是 tail，
 *      操作依然安全，无需特殊处理边界。
 * ────────────────────────────────────────────────────────────────────────────
 * 【Keys / Values 遍历走顺序链】
 *
 *   遍历时走 head → ... → tail 双向链表，不走哈希桶，原因：
 *     ① 保证按插入顺序输出
 *     ② 时间复杂度 O(n)，无需跳过大量空桶
 *
 *   ⚠️ 易错点3：循环终止条件必须是 p != tail，不能是 p != nil
 *      tail 是哨兵节点（tail.key = 0, tail.next = nil），若用 p != nil，
 *      循环会把 tail 的零值混入结果，每次输出都多出一个错误的 0。
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次手写 LinkedHashMap 后对照检查
 *
 *     ✅ 哨兵初始化：head.next = tail，tail.prev = head？
 *     ✅ 尾插入：新节点插在 tail 前，4 行代码赋值顺序正确？
 *     ✅ Put 更新：key 存在时只改 value，不重新挂入顺序链？
 *     ✅ Remove：同时从顺序链和桶链双重摘除？
 *     ✅ Keys/Values：循环终止条件是 p != tail，而非 p != nil？
 *     ✅ 桶懒初始化：Get/Put/Remove 均有 nil 判断再访问桶链？
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *     0. NewMyLinkedHashTable(capacity int) *MyLinkedHashTable
 *     1. Get(key int) (int, bool)     // O(1) 均摊
 *     2. Put(key, val int)            // O(1) 均摊
 *     3. Remove(key int)              // O(1) 均摊
 *     4. Display()                    // O(n)，打印哈希桶结构
 *     5. Keys() []int                 // O(n)，按插入顺序返回 key
 *     6. Values() []int               // O(n)，按插入顺序返回 value
 * ============================================================================
 */

package main

import (
	"container/list"
	"fmt"
)

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
	for p != m.tail {
		keys = append(keys, p.key)
		p = p.next
	}
	return keys
}

func (m *MyLinkedHashTable) Values() []int {
	var values []int
	p := m.head.next
	for p != m.tail {
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
