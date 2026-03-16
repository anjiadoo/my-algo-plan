package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// 跳表（Skip List）实现：
// 🌟核心原理：在有序链表基础上增加多层索引，每向上一层索引节点减半、间隔翻倍，形成类似二分查找的效果
// 🌟时间复杂度：查找、插入、删除均为 O(logN)，N 为元素个数
// 🌟空间复杂度：O(N)，多层索引的额外空间总和约为 N（等比级数 N/2 + N/4 + ... ≈ N）
// 🌟随机化层数：通过抛硬币（概率 1/2）决定新节点的层数，期望层数为 2，确保索引层高度期望为 O(logN)
// 🌟与平衡BST对比：跳表实现更简单，通过随机化而非严格旋转来维持平衡，Redis的有序集合(ZSet)底层就用了跳表

// ⚠️易错点1：查找路径记录 - 插入/删除时需要记录每层的前驱节点（update数组），才能正确修改各层指针
// ⚠️易错点2：层数维护 - 插入新节点可能增加最大层数，删除节点可能降低最大层数，需同步更新
// ⚠️易错点3：哨兵头节点 - 头节点不存储实际数据，其forward数组长度等于最大层数，简化边界处理
// ⚠️易错点4：随机层数上限 - 需要设置maxLevel防止极端情况下层数过高，一般取 log2(N) 即可

// 何时运用跳表：
// ❓1、是否需要有序集合的快速增删查？跳表提供 O(logN) 的有序操作，比哈希表多了有序性
// ❓2、是否需要范围查询？跳表的底层是有序链表，天然支持高效的范围遍历
// ❓3、是否想要比平衡BST更简单的实现？跳表通过随机化替代复杂的旋转/着色操作

// API列表：
// 0、func NewSkipList() *SkipList                      // 创建跳表
// 1、func (sl *SkipList) Search(key int) (int, bool)   // 查找key对应的value
// 2、func (sl *SkipList) Insert(key int, val int)      // 插入/更新键值对
// 3、func (sl *SkipList) Delete(key int) bool          // 删除key对应的节点
// 4、func (sl *SkipList) Size() int                    // 返回元素个数
// 5、func (sl *SkipList) Display()                     // 可视化打印跳表

const (
	maxLevel    = 16  // 最大层数，支持 2^16 = 65536 个元素的理想跳表
	probability = 0.5 // 晋升概率，每层节点数期望为下一层的一半
)

// skipNode 跳表节点
type skipNode struct {
	key     int         // 键，用于排序
	val     int         // 值
	forward []*skipNode // forward[i] 表示第 i 层的下一个节点
}

// SkipList 跳表结构
type SkipList struct {
	head  *skipNode // 哨兵头节点，不存储实际数据
	level int       // 当前跳表的最大层数（从1开始计数）
	size  int       // 元素个数
}

// NewSkipList 创建一个空的跳表
func NewSkipList() *SkipList {
	// 头节点拥有 maxLevel 层，初始所有层的 forward 都是 nil
	head := newSkipNode(0, 0, maxLevel)
	return &SkipList{
		head:  head,
		level: 1,
	}
}

// newSkipNode 创建指定层数的跳表节点
func newSkipNode(key, val, level int) *skipNode {
	return &skipNode{
		key:     key,
		val:     val,
		forward: make([]*skipNode, level),
	}
}

// randomLevel 随机生成新节点的层数
// 以 probability 的概率向上晋升，期望层数为 1/(1-probability) = 2
func randomLevel() int {
	lvl := 1
	for rand.Float64() < probability && lvl < maxLevel {
		lvl++
	}
	return lvl
}

// Search 查找key对应的value
// 从最高层开始逐层向下搜索，每层尽量向右移动
// 时间复杂度：O(logN)
func (sl *SkipList) Search(key int) (int, bool) {
	curr := sl.head

	// 从最高层往下搜索
	for i := sl.level - 1; i >= 0; i-- {
		// 在当前层尽量向右移动，直到下一个节点的key >= 目标key
		for curr.forward[i] != nil && curr.forward[i].key < key {
			curr = curr.forward[i]
		}
	}

	// 此时 curr 是第0层中 key 小于目标的最后一个节点
	// curr.forward[0] 就是可能等于目标key的节点
	curr = curr.forward[0]
	if curr != nil && curr.key == key {
		return curr.val, true
	}
	return 0, false
}

// Insert 插入或更新键值对
// 如果key已存在，更新其value；否则创建新节点插入
// 时间复杂度：O(logN)
func (sl *SkipList) Insert(key int, val int) {
	// update[i] 记录第 i 层中，新节点的前驱节点
	update := make([]*skipNode, maxLevel)
	curr := sl.head

	// 从最高层往下，找到每层的前驱节点
	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil && curr.forward[i].key < key {
			curr = curr.forward[i]
		}
		update[i] = curr
	}

	// 检查key是否已存在
	curr = curr.forward[0]
	if curr != nil && curr.key == key {
		// key已存在，直接更新value
		curr.val = val
		return
	}

	// key不存在，创建新节点
	newLevel := randomLevel()

	// 如果新节点的层数超过当前最大层数，需要更新update数组
	// 超出部分的前驱节点都是head
	if newLevel > sl.level {
		for i := sl.level; i < newLevel; i++ {
			update[i] = sl.head
		}
		sl.level = newLevel
	}

	// 创建新节点并在每层中插入
	node := newSkipNode(key, val, newLevel)
	for i := 0; i < newLevel; i++ {
		// 经典链表插入：新节点指向前驱的下一个，前驱指向新节点
		node.forward[i] = update[i].forward[i]
		update[i].forward[i] = node
	}

	sl.size++
}

// Delete 删除key对应的节点
// 返回是否成功删除（key不存在返回false）
// 时间复杂度：O(logN)
func (sl *SkipList) Delete(key int) bool {
	// update[i] 记录第 i 层中，待删除节点的前驱节点
	update := make([]*skipNode, maxLevel)
	curr := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil && curr.forward[i].key < key {
			curr = curr.forward[i]
		}
		update[i] = curr
	}

	// 检查key是否存在
	curr = curr.forward[0]
	if curr == nil || curr.key != key {
		return false
	}

	// 从每层中移除该节点
	for i := 0; i < sl.level; i++ {
		if update[i].forward[i] != curr {
			// 当前层没有该节点（该节点的层数小于当前层），后续更高层也不会有
			break
		}
		// 经典链表删除：前驱直接指向被删节点的下一个
		update[i].forward[i] = curr.forward[i]
	}

	// 如果删除节点后最高层变空了，降低层数
	for sl.level > 1 && sl.head.forward[sl.level-1] == nil {
		sl.level--
	}

	sl.size--
	return true
}

// Size 返回跳表中的元素个数
func (sl *SkipList) Size() int {
	return sl.size
}

// Display 可视化打印跳表的每一层
// 从最高层到最低层依次打印，直观展示跳表的索引结构
func (sl *SkipList) Display() {
	fmt.Printf("跳表（元素数: %d, 层数: %d）:\n", sl.size, sl.level)

	for i := sl.level - 1; i >= 0; i-- {
		fmt.Printf("Level %2d: head", i)
		curr := sl.head.forward[i]
		for curr != nil {
			fmt.Printf(" -> [%d:%d]", curr.key, curr.val)
			curr = curr.forward[i]
		}
		fmt.Println(" -> nil")
	}

	// 打印底层有序链表的key序列
	fmt.Print("有序序列: ")
	curr := sl.head.forward[0]
	keys := make([]string, 0, sl.size)
	for curr != nil {
		keys = append(keys, fmt.Sprintf("%d", curr.key))
		curr = curr.forward[0]
	}
	fmt.Println(strings.Join(keys, " -> "))
	fmt.Println()
}

func main() {
	sl := NewSkipList()
	num := 100

	// 插入测试
	fmt.Println("===== 插入元素 =====")
	for i := 1; i <= num; i++ {
		sl.Insert(i, i*10)
	}
	sl.Display()

	// 查找测试
	fmt.Println("===== 查找元素 =====")
	searchKeys := []int{7, 19, 100}
	for _, key := range searchKeys {
		if val, ok := sl.Search(key); ok {
			fmt.Printf("查找 key=%d -> val=%d ✓\n", key, val)
		} else {
			fmt.Printf("查找 key=%d -> 不存在 ✗\n", key)
		}
	}
	fmt.Println()

	// 更新测试
	fmt.Println("===== 更新元素 =====")
	sl.Insert(7, 777)
	val, _ := sl.Search(7)
	fmt.Printf("更新 key=7, newVal=777 -> 查找结果: %d\n", val)
	fmt.Println()

	// 删除测试
	fmt.Println("===== 删除元素 =====")
	deleteKeys := []int{6, 19, 100}
	for _, key := range deleteKeys {
		if sl.Delete(key) {
			fmt.Printf("删除 key=%d -> 成功 ✓\n", key)
		} else {
			fmt.Printf("删除 key=%d -> 不存在 ✗\n", key)
		}
	}
	fmt.Printf("删除后元素个数: %d\n", sl.Size())
	sl.Display()
}
