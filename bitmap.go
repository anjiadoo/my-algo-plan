/*
 * ============================================================================
 *                        📘 位图算法 · 核心记忆框架
 * ============================================================================
 * 【一句话理解位图】
 *
 *   位图：用一个比特位的 0/1 标记某个数值「是否存在」，
 *         一个 uint64 能同时标记 64 个数值，内存效率是 bool 切片的 1/8。
 * ────────────────────────────────────────────────────────────────────────────
 * 【位图核心：三步定位 + 三种位运算】
 *
 *   定位公式：
 *     bucketIndex = num / 64    → 找到第几个 uint64
 *     bitIndex    = num % 64    → 找到该 uint64 的第几位（bit）
 *
 *   三种位运算（全部 O(1)）：
 *     置位（标记存在）：buckets[bucketIndex] |=  1 << bitIndex
 *     清零（标记不在）：buckets[bucketIndex] &^= 1 << bitIndex   ← Go 专属 &^ 操作符
 *     检查（是否存在）：buckets[bucketIndex] &  (1 << bitIndex) != 0
 * ────────────────────────────────────────────────────────────────────────────
 * 【固定大小 vs 动态扩容】
 *
 *   固定大小（NewBitmap + set/clear/has）：
 *     构造时分配 maxNum/64+1 个 uint64，size = maxNum+1（闭区间元素个数公式）。
 *     访问时严格检查 num < 0 || num >= b.size，越界直接报错。
 *
 *   动态扩容（dynamicSet）：
 *     写入前计算 requiredBucket = num/64+1，若超出则重新分配并 copy，
 *     扩容后 size = requiredBucket * 64（按 bucket 边界对齐，不是 num+1）。
 * ────────────────────────────────────────────────────────────────────────────
 * 【核心技巧详解】
 *
 *   🌟技巧1：位索引定位技巧
 *       bucketIndex = num/64，bitIndex = uint(num%64)
 *       除法 → 第几个桶，取模 → 桶内第几位，两步精确定位任意比特。
 *
 *   🌟技巧2：&^ 按位清零技巧
 *       Go 特有的 &^（AND NOT）操作符，等价于 C 语言的 & ~。
 *       &^= 1 << bitIndex 可将目标位清 0，其余位不变。
 *
 *   🌟技巧3：内存压缩技巧
 *       bool 切片每个元素占 1 字节（8 位），位图每个元素仅占 1 位，
 *       内存比 []bool 节省 8 倍，比 []int 节省 64 倍。
 *
 *   🌟技巧4：栅栏问题技巧
 *       区间 [0, maxNum] 共 maxNum+1 个元素，所以 size = maxNum+1。
 *       类比：n 段栅栏需要 n+1 根柱子，「元素个数 = 上界 - 下界 + 1」。
 *
 *   🌟技巧5：bucket数量 = 索引 + 1 技巧
 *       数组下标从 0 开始，因此 bucket 数量 = 最大 bucketIndex + 1，
 *       即 bucketCount = num/64 + 1，这是索引到数量的通用转换公式。
 *
 *   🌟技巧6：动态扩容按 bucket 对齐技巧
 *       扩容后 size = requiredBucket * 64，而非 num+1，
 *       按桶边界对齐避免下次扩容时出现桶内「空洞」导致边界检查错误。
 * ────────────────────────────────────────────────────────────────────────────
 * 【易错点】
 *
 *   ⚠️ 易错点1：bitIndex 必须转成 uint 类型再做移位
 *       Go 不允许 int 类型直接做移位量，须写 uint(num % 64)，否则编译报错。
 *
 *   ⚠️ 易错点2：清零用 &^=，不能用 &= ~（Go 中没有按位取反 ~ 操作符）
 *       Go 的按位取反须写 ^x，清零正确写法是 &^= 1 << bitIndex。
 *
 *   ⚠️ 易错点3：dynamicSet 扩容后 size 要用桶对齐值，不是 num+1
 *       size = requiredBucket * 64，而非 num+1，否则桶末尾的合法位会被边界检查拦截。
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次手写位图后对照检查
 *
 *     ✅ 位索引：bitIndex 有没有转成 uint(num % 64)？
 *     ✅ 置位：使用 |= 1 << bitIndex，不是 = 1 << bitIndex（会清掉其他位）？
 *     ✅ 清零：使用 &^= 1 << bitIndex，不是 &= ~（Go 无 ~ 操作符）？
 *     ✅ 检查：使用 & (1 << bitIndex) != 0，不是 == 1（结果是 2^bitIndex 不是 1）？
 *     ✅ 构造：bucketCount = maxNum/64 + 1，size = maxNum+1（栅栏问题）？
 *     ✅ 动态扩容：size = requiredBucket * 64（按桶对齐），不是 num+1？
 * ────────────────────────────────────────────────────────────────────────────
 * 【API 速查】
 *
 *     0. NewBitmap(maxNum int) (*Bitmap, error)        // 创建固定大小位图，O(n)
 *     1. (b *Bitmap) set(num int) error                // 置位，O(1)
 *     2. (b *Bitmap) clear(num int) error              // 清零，O(1)
 *     3. (b *Bitmap) has(num int) (bool, error)        // 检查存在性，O(1)
 *     4. (b *Bitmap) count() int                       // 统计置位数量，O(n/64)
 *     5. (b *Bitmap) dynamicSet(num int) error         // 动态置位（自动扩容），均摊O(1)
 * ============================================================================
 */

package main

import (
	"errors"
	"fmt"
)

// Bitmap 位图结构，用uint64切片存储比特位
type Bitmap struct {
	buckets []uint64 // 每个uint64存储64个比特位
	size    int      // 位图可存储的最大数值+1（可选，用于边界检查）
}

func NewBitmap(maxNum int) (*Bitmap, error) {
	if maxNum < 0 {
		return nil, errors.New("maxNum必须是非负数")
	}

	bucketCount := maxNum/64 + 1 // 计算需要的bucket数量(数量=索引+1)

	return &Bitmap{
		buckets: make([]uint64, bucketCount),
		// 用户指定的精确上界
		// 区间[0,maxNum]元素个数=maxNum-0+1，闭区间元素个数公式，栅栏问题
		size: maxNum + 1,
	}, nil
}

// 置位：标记num存在
func (b *Bitmap) set(num int) error {
	if num < 0 || num >= b.size {
		return errors.New("num超出位图范围")
	}
	bucketIndex := num / 64                 // 计算所在的uint64索引
	bitIndex := uint(num % 64)              // 计算所在的比特位索引
	b.buckets[bucketIndex] |= 1 << bitIndex // 位运算置1
	return nil
}

func (b *Bitmap) dynamicSet(num int) error {
	if num < 0 {
		return errors.New("num必须是非负数")
	}

	requiredBucket := num/64 + 1 // 计算需要的bucket数量(数量=索引+1)

	if requiredBucket > len(b.buckets) { // 是否需要扩容
		newBits := make([]uint64, requiredBucket)
		copy(newBits, b.buckets)
		b.buckets = newBits
		b.size = requiredBucket * 64 // 扩容后按bucket对齐的实际容量
	}
	return b.set(num)
}

// 清零：标记num不存在
func (b *Bitmap) clear(num int) error {
	if num < 0 || num >= b.size {
		return errors.New("num超出位图范围")
	}
	bucketIndex := num / 64
	bitIndex := uint(num % 64)

	// 位运算置0（&^是Go的按位清零）
	b.buckets[bucketIndex] &^= 1 << bitIndex
	return nil
}

// 检查：判断num是否存在
func (b *Bitmap) has(num int) (bool, error) {
	if num < 0 || num >= b.size {
		return false, errors.New("num超出位图范围")
	}
	bucketIndex := num / 64
	bitIndex := uint(num % 64)
	// 位运算检查该位是否为1
	return (b.buckets[bucketIndex] & (1 << bitIndex)) != 0, nil
}

// 统计：计算已置位的数量（存在的数值个数）
func (b *Bitmap) count() int {
	count := 0
	// 遍历每个uint64，统计其中1的个数
	for _, bucket := range b.buckets {
		count += popcount(bucket)
	}
	return count
}

// 统计一个uint64中1的个数（内置函数优化）
func popcount(x uint64) int {
	if x == 0 {
		return 0
	}
	// 递归简化版，生产环境可用runtime.buckets.OnesCount64
	return int(x&1) + popcount(x>>1)
}

// 测试示例
func main() {
	// 创建可存储0~100的位图
	bitmap, err := NewBitmap(100)
	if err != nil {
		fmt.Println("创建位图失败：", err)
		return
	}

	// 置位操作
	_ = bitmap.set(10)
	_ = bitmap.set(20)
	_ = bitmap.set(65) // 跨uint64的数值（65=64+1，在第二个uint64的第1位）

	// 检查存在性
	has10, _ := bitmap.has(10)
	has20, _ := bitmap.has(20)
	has30, _ := bitmap.has(30)
	has65, _ := bitmap.has(65)
	fmt.Println("10是否存在：", has10) // true
	fmt.Println("20是否存在：", has20) // true
	fmt.Println("30是否存在：", has30) // false
	fmt.Println("65是否存在：", has65) // true

	// 清零操作
	_ = bitmap.clear(20)
	has20AfterClear, _ := bitmap.has(20)
	fmt.Println("20清零后是否存在：", has20AfterClear) // false

	// 统计数量
	fmt.Println("已置位的数量：", bitmap.count()) // 2（10和65）

	// 内存占用对比：存储0~100的3个数值
	// 传统切片（[]int）：3*8=24字节
	// 位图：2个uint64=16字节（可存储0~127），实际仅用3个位
	fmt.Println("位图内存占用（字节）：", len(bitmap.buckets)*8) // 16
}
