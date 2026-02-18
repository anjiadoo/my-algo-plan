package main

import (
	"errors"
	"fmt"
)

// 位图（Bitmap）实现，用uint64切片高效存储比特位：
// 🌟技巧1：一个uint64存储64个比特位，bucketIndex = num/64，bitIndex = num%64
// 🌟技巧2：置位用 |=，清零用 &^=（Go按位清零），检查用 & != 0，都是O(1)操作
// 🌟技巧3：相比bool切片，内存占用仅为1/8，适合海量数据的布尔标记场景
// 0、func NewBitmap(maxNum int) (*Bitmap, error)
// 1、func (b *Bitmap) set(num int) error
// 2、func (b *Bitmap) clear(num int) error
// 3、func (b *Bitmap) has(num int) (bool, error)
// 4、func (b *Bitmap) count() int

// Bitmap 位图结构，用uint64切片存储比特位
type Bitmap struct {
	bits []uint64 // 每个uint64存储64个比特位
	size int      // 位图可存储的最大数值+1（可选，用于边界检查）
}

// NewBitmap 创建位图实例，指定可存储的最大数值
func NewBitmap(maxNum int) (*Bitmap, error) {
	if maxNum < 0 {
		return nil, errors.New("maxNum必须是非负数")
	}
	// 计算需要的uint64数量：(maxNum + 64) / 64
	bucketCount := (maxNum + 64) / 64
	return &Bitmap{
		bits: make([]uint64, bucketCount),
		size: maxNum + 1,
	}, nil
}

// set 置位：标记num存在
func (b *Bitmap) set(num int) error {
	if num < 0 || num >= b.size {
		return errors.New("num超出位图范围")
	}
	bucketIndex := num / 64              // 计算所在的uint64索引
	bitIndex := uint(num % 64)           // 计算所在的比特位索引
	b.bits[bucketIndex] |= 1 << bitIndex // 位运算置1
	return nil
}

// clear 清零：标记num不存在
func (b *Bitmap) clear(num int) error {
	if num < 0 || num >= b.size {
		return errors.New("num超出位图范围")
	}
	bucketIndex := num / 64
	bitIndex := uint(num % 64)
	b.bits[bucketIndex] &^= 1 << bitIndex // 位运算置0（&^是Go的按位清零）
	return nil
}

// has 检查：判断num是否存在
func (b *Bitmap) has(num int) (bool, error) {
	if num < 0 || num >= b.size {
		return false, errors.New("num超出位图范围")
	}
	bucketIndex := num / 64
	bitIndex := uint(num % 64)
	// 位运算检查该位是否为1
	return (b.bits[bucketIndex] & (1 << bitIndex)) != 0, nil
}

// count 统计：计算已置位的数量（存在的数值个数）
func (b *Bitmap) count() int {
	count := 0
	// 遍历每个uint64，统计其中1的个数
	for _, bucket := range b.bits {
		count += popcount(bucket)
	}
	return count
}

// popcount 统计一个uint64中1的个数（内置函数优化）
// 也可以用手写位运算，这里用Go内置的快速实现
func popcount(x uint64) int {
	if x == 0 {
		return 0
	}
	return int(x&1) + popcount(x>>1) // 递归简化版，生产环境可用runtime.bits.OnesCount64
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
	fmt.Println("位图内存占用（字节）：", len(bitmap.bits)*8) // 16
}
