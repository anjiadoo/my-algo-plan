/*
 * ============================================================================
 *                      📘 跳表（Redis SortedSet 简化版）· 核心记忆框架
 * ============================================================================
 * 【一句话理解 SortedSet】
 *
 *   SortedSet = 跳表（按 score 有序）+ 哈希表（member → score 的 O(1) 映射）
 *   member string 是唯一标识，score float64 是排序依据，两者分离。
 *
 *   Redis 命令对照：
 *     ZADD   key s m  →  ZAdd(member, score)
 *     ZREM   key m    →  ZRem(member)
 *     ZSCORE key m    →  ZScore(member) (float64, bool)
 *     ZRANK  key m    →  ZRank(member) (int, bool)       // 0-based 排名
 *     ZRANGE key s e  →  ZRangeByScore(min, max) []ZNode // 按 score 范围
 *     ZCARD  key      →  ZCard() int
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【SortedSet 的双索引结构】
 *
 *   Redis SortedSet 实际上维护两个索引：
 *     ① 跳表（skiplist）：按 score 排序，支持范围查询和排名
 *     ② 哈希表（dict）  ：member → score 映射，支持 O(1) 查分
 *
 *   本实现也用同样的双索引：
 *     z.dict  map[string]float64   // O(1) 查分 / 判断存在
 *     z.head  *skipNode            // 跳表头，按 score 有序
 *
 *   ⚠️ 易错点1：ZAdd 时 member 已存在需要先删后插（不能只更新 score）
 *      score 变了意味着节点在跳表中的位置变了，必须从旧位置摘除再插入新位置。
 *      流程：① dict 查到旧 score → ② 用(旧score, member)定位并删除跳表节点
 *            → ③ 插入新节点(新score, member) → ④ 更新 dict
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【SortedSet 跳表的排序规则（双键排序）】
 *
 *   跳表节点排序：先按 score 升序，score 相同时按 member 字典序升序
 *   这保证了每个(score, member)组合在跳表中唯一且位置确定。
 *
 *   比较函数（核心）：
 *     func less(aScore float64, aMember string, bScore float64, bMember string) bool {
 *         if aScore != bScore { return aScore < bScore }
 *         return aMember < bMember
 *     }
 *
 *   ⚠️ 易错点2：所有跳表操作（插入/删除/查找前驱）都必须用双键比较
 *      只用 score 比较会导致 score 相同的多个节点无法精确定位，
 *      删除时可能删错节点（删掉同 score 的另一个 member）。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【ZRank 排名的实现原理】
 *
 *   跳表天然支持排名：从 head 出发，沿底层（forward[0]）线性数到目标节点。
 *   更高效的实现（Redis 真实做法）是在每个 forward 指针上额外记录 span（跨度），
 *   但本简化版直接在底层链表上计数，时间复杂度 O(N)。
 *
 *   ⚠️ 易错点3：ZRank 必须先用 dict 确认 member 存在，再去跳表计数
 *      若 member 不存在，沿底层遍历完也找不到，需要返回 (0, false)。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【ZRangeByScore 的正确姿势】
 *
 *   流程：
 *     ① 用多层索引跳到第一个 score >= min 的前驱（和普通跳表范围查询一样）
 *     ② 落底（forward[0]），沿底层遍历直到 score > max
 *
 *   ⚠️ 易错点4：定位前驱时，内层循环条件用 score 单键比较即可（不需要双键）
 *      因为我们只需要找到"score < min"的最后一个节点，
 *      只要 forward[i].score < min 就继续向右，不涉及 member 的比较。
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【ZRem 删除的关键步骤】
 *
 *   ① 从 dict 查出该 member 的 score（得知跳表中的精确位置）
 *   ② 用双键(score, member)在跳表中定位前驱，执行链表删除
 *   ③ 从 dict 中删除该 member
 *
 *   ⚠️ 易错点5：跳表定位时内层循环必须用双键比较（less 函数）
 *      不能只比较 score，否则相同 score 的多个 member 会定位到错误节点。
 *      具体条件：forward[i].score < score || (forward[i].score == score && forward[i].member < member)
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次写完 SortedSet 操作后对照检查
 *
 *     ✅ ZAdd：member 已存在时，是否先删跳表旧节点再插新节点？
 *     ✅ 所有跳表定位：score 相同时是否用了 member 字典序作为第二排序键？
 *     ✅ ZRem：是否先查 dict 得到 score，再去跳表精确定位？
 *     ✅ ZRangeByScore：定位前驱后是否落到 forward[0] 再线性遍历？
 *     ✅ ZRank：是否先判断 member 在 dict 中存在再去跳表计数？
 *     ✅ 插入：两步链接顺序是"先接后继，再接前驱"？
 *     ✅ 插入：newLevel > sl.level 时有没有补填 update[i] = sl.head？
 *     ✅ 删除：删完后有没有收缩 sl.level（去掉空的最高层）？
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *
 *     0. NewSortedSet() *SortedSet
 *     1. ZAdd(member string, score float64)              // O(logN)，已存在则更新
 *     2. ZRem(member string) bool                        // O(logN)
 *     3. ZScore(member string) (float64, bool)           // O(1)
 *     4. ZRank(member string) (int, bool)                // O(N)，0-based 升序排名
 *     5. ZRangeByScore(min, max float64) []ZNode         // O(logN+M)，闭区间
 *     6. ZCard() int                                     // O(1)
 *     7. Display()                                       // 可视化打印各层结构
 * ============================================================================
 */

package main

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	maxLevel    = 16  // 最大层数，支持 2^16 = 65536 个元素的理想跳表
	probability = 0.5 // 晋升概率，每层节点数期望为下一层的一半
)

// skipNode 跳表节点，对应 Redis SortedSet 中的一个元素
type skipNode struct {
	member  string      // 唯一成员名，对应 Redis 的 member
	score   float64     // 排序分值，对应 Redis 的 score
	forward []*skipNode // forward类似单链表中的next指针，只不过跳表有多个不同next指针，用数组来表示而已，下标表示第几层
}

// SortedSet Redis ZSet 简化版：跳表 + 哈希表双索引
type SortedSet struct {
	head  *skipNode          // 哨兵头节点，不存储实际数据
	dict  map[string]float64 // member → score，O(1) 查分
	level int                // 当前跳表最大层数（从1开始）
	size  int                // 元素个数
}

// NewSortedSet 创建一个空的 SortedSet
func NewSortedSet() *SortedSet {
	head := &skipNode{forward: make([]*skipNode, maxLevel)}
	return &SortedSet{
		head:  head,
		dict:  make(map[string]float64),
		level: 1,
	}
}

// randomLevel 随机生成新节点的层数
func randomLevel() int {
	lvl := 1
	for rand.Float64() < probability && lvl < maxLevel {
		lvl++
	}
	return lvl
}

// less 跳表的双键比较：先按 score 升序，score 相同时按 member 字典序升序
// 所有需要在跳表中精确定位节点的操作都必须用这个比较，而不是只比较 score
func less(aScore float64, aMember string, bScore float64, bMember string) bool {
	if aScore != bScore {
		return aScore < bScore
	}
	return aMember < bMember
}

// insertNode 在跳表中插入节点（内部方法，不操作 dict）
func (z *SortedSet) insertNode(member string, score float64) {
	update := make([]*skipNode, maxLevel)
	curr := z.head

	// 从最高层往下，找每层的前驱：前驱满足 (score,member) < 新节点
	for i := z.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil && less(curr.forward[i].score, curr.forward[i].member, score, member) {
			curr = curr.forward[i]
		}
		update[i] = curr
	}

	newLvl := randomLevel()
	if newLvl > z.level {
		for i := z.level; i < newLvl; i++ {
			update[i] = z.head // 新层的前驱只有 head
		}
		z.level = newLvl
	}

	node := &skipNode{
		member:  member,
		score:   score,
		forward: make([]*skipNode, newLvl),
	}
	for i := 0; i < newLvl; i++ {
		node.forward[i] = update[i].forward[i] // ① 先接后继
		update[i].forward[i] = node            // ② 再接前驱
	}
}

// deleteNode 从跳表中删除节点（内部方法，不操作 dict）
func (z *SortedSet) deleteNode(member string, score float64) bool {
	update := make([]*skipNode, maxLevel)
	curr := z.head

	for i := z.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil && less(curr.forward[i].score, curr.forward[i].member, score, member) {
			curr = curr.forward[i]
		}
		update[i] = curr
	}

	// 精确匹配：score 和 member 都要相同
	target := curr.forward[0]
	if target == nil || target.score != score || target.member != member {
		return false
	}

	for i := 0; i < z.level; i++ {
		if update[i].forward[i] != target {
			break // 该层及更高层都不含此节点
		}
		update[i].forward[i] = target.forward[i]
	}

	// 收缩空层
	for z.level > 1 && z.head.forward[z.level-1] == nil {
		z.level--
	}
	return true
}

// ZAdd 添加或更新成员分值，对应 Redis ZADD
// 若 member 已存在，先删除旧跳表节点，再插入新位置（因为 score 变了，位置也变了）
// 时间复杂度：O(logN)
func (z *SortedSet) ZAdd(member string, score float64) {
	if oldScore, exists := z.dict[member]; exists {
		if oldScore == score {
			return // score 没变，无需任何操作
		}
		z.deleteNode(member, oldScore) // 从跳表旧位置摘除
		z.size--
	}
	z.insertNode(member, score)
	z.dict[member] = score
	z.size++
}

// ZRem 删除成员，对应 Redis ZREM
// 时间复杂度：O(logN)
func (z *SortedSet) ZRem(member string) bool {
	score, exists := z.dict[member]
	if !exists {
		return false
	}
	z.deleteNode(member, score)
	delete(z.dict, member)
	z.size--
	return true
}

// ZScore 获取成员的分值，对应 Redis ZSCORE
// 时间复杂度：O(1)，直接查哈希表
func (z *SortedSet) ZScore(member string) (float64, bool) {
	score, ok := z.dict[member]
	return score, ok
}

// ZRank 获取成员的升序排名（0-based），对应 Redis ZRANK
// 时间复杂度：O(N)（简化版，完整版需在 forward 上维护 span 跨度才能 O(logN)）
func (z *SortedSet) ZRank(member string) (int, bool) {
	if _, exists := z.dict[member]; !exists {
		return 0, false
	}
	rank := 0
	curr := z.head.forward[0] // 从底层链表第一个节点开始计数
	for curr != nil {
		if curr.member == member {
			return rank, true
		}
		rank++
		curr = curr.forward[0]
	}
	return 0, false
}

// ZNode 范围查询的返回结果
type ZNode struct {
	Member string
	Score  float64
}

// ZRangeByScore 按 score 范围查找，返回 score 在 [min, max] 闭区间内的所有节点，按 score 升序
// 对应 Redis ZRANGEBYSCORE key min max
// 时间复杂度：O(logN + M)，M 为结果数
func (z *SortedSet) ZRangeByScore(min, max float64) []ZNode {
	if min > max {
		return nil
	}

	var result []ZNode
	curr := z.head

	// 用多层索引跳到第一个 score >= min 的前驱（只用 score 单键定位即可）
	for i := z.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil && curr.forward[i].score < min {
			curr = curr.forward[i]
		}
	}

	// 落底，沿底层遍历直到 score > max
	curr = curr.forward[0]
	for curr != nil && curr.score <= max {
		result = append(result, ZNode{Member: curr.member, Score: curr.score})
		curr = curr.forward[0]
	}
	return result
}

// ZCard 返回成员总数，对应 Redis ZCARD
// 时间复杂度：O(1)
func (z *SortedSet) ZCard() int {
	return z.size
}

// Display 可视化打印跳表各层结构（调试用）
func (z *SortedSet) Display() {
	fmt.Printf("SortedSet（元素数: %d, 层数: %d）:\n", z.size, z.level)
	for i := z.level - 1; i >= 0; i-- {
		fmt.Printf("Level %2d: head", i)
		curr := z.head.forward[i]
		for curr != nil {
			fmt.Printf(" -> [%s:%v]", curr.member, curr.score)
			curr = curr.forward[i]
		}
		fmt.Println(" -> nil")
	}
	fmt.Print("有序序列: ")
	curr := z.head.forward[0]
	parts := make([]string, 0, z.size)
	for curr != nil {
		parts = append(parts, fmt.Sprintf("%s(%v)", curr.member, curr.score))
		curr = curr.forward[0]
	}
	fmt.Println(strings.Join(parts, " -> "))
	fmt.Println()
}

func main() {
	z := NewSortedSet()

	// 插入测试
	fmt.Println("===== ZADD =====")
	members := map[string]float64{
		"jay":      63,
		"alice":    88,
		"bob":      72,
		"charlie":  95,
		"dave":     72, // 与 bob 同分，按 member 字典序排在 bob 后
		"eve":      60,
		"jayden":   100,
		"andy":     21,
		"suki":     89,
		"anjiadoo": 45,
	}
	for member, score := range members {
		z.ZAdd(member, score)
	}
	z.Display()

	// ZSCORE 测试
	fmt.Println("===== ZSCORE =====")
	for member := range members {
		if s, ok := z.ZScore(member); ok {
			fmt.Printf("ZSCORE %s -> %v ✓\n", member, s)
		}
	}
	fmt.Println()

	// ZRANK 测试
	fmt.Println("===== ZRANK =====")
	for member := range members {
		if r, ok := z.ZRank(member); ok {
			fmt.Printf("ZRANK %s -> %d\n", member, r)
		}
	}
	fmt.Println()

	// ZADD 更新 score 测试
	fmt.Println("===== ZADD 更新 score =====")
	z.ZAdd("eve", 99.0) // eve 从 60 升到 99，跳表位置变化
	z.Display()

	// ZRANGEBYSCORE 测试
	fmt.Println("===== ZRANGEBYSCORE [70, 95] =====")
	results := z.ZRangeByScore(70, 95)
	fmt.Printf("共 %d 个结果:\n", len(results))
	for _, r := range results {
		fmt.Printf("  %s -> %v\n", r.Member, r.Score)
	}
	fmt.Println()

	// ZREM 测试
	fmt.Println("===== ZREM =====")
	fmt.Printf("ZREM bob -> %v\n", z.ZRem("bob"))
	fmt.Printf("ZREM notexist -> %v\n", z.ZRem("notexist"))
	fmt.Printf("ZCARD -> %d\n", z.ZCard())
	z.Display()
}
