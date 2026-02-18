package main

import (
	"fmt"
	"unsafe"
)

// 内存对齐示例
// 🌟技巧1：对齐大小 和 类型大小 是两码事
// 🌟技巧2：string，slice 的类型大小是看底层的数据结构的字段大小

func printMemSize() {
	fmt.Println("基础类型内存大小:")
	fmt.Println("bool:          ", unsafe.Sizeof(true), "（字节）")           // 1字节
	fmt.Println("int8/uint8:    ", unsafe.Sizeof(int8(0)), "（字节）")        // 1字节
	fmt.Println("int16/uint16:  ", unsafe.Sizeof(int16(0)), "（字节）")       // 2字节
	fmt.Println("int32/uint32:  ", unsafe.Sizeof(int32(0)), "（字节）")       // 4字节
	fmt.Println("int64/uint64:  ", unsafe.Sizeof(int64(0)), "（字节）")       // 8字节
	fmt.Println("int/uint:      ", unsafe.Sizeof(int(0)), "（字节）")         // 8字节 (64位系统)
	fmt.Println("float32:       ", unsafe.Sizeof(float32(0)), "（字节）")     // 4字节
	fmt.Println("float64:       ", unsafe.Sizeof(float64(0)), "（字节）")     // 8字节
	fmt.Println("complex64:     ", unsafe.Sizeof(complex64(0)), "（字节）")   // 8字节 (2个float32)
	fmt.Println("complex128:    ", unsafe.Sizeof(complex128(0)), "（字节）")  // 16字节 (2个float64)
	fmt.Println("byte:          ", unsafe.Sizeof(byte(0)), "（字节）")        // 1字节 (uint8别名)
	fmt.Println("rune:          ", unsafe.Sizeof(rune(0)), "（字节）")        // 4字节 (int32别名)
	fmt.Println("string:        ", unsafe.Sizeof(""), "（字节）")             // 16字节 (指针+长度)
	fmt.Println("[]int:         ", unsafe.Sizeof([]int{}), "（字节）")        // 24字节 (指针+长度+容量)
	fmt.Println("map[int]int:   ", unsafe.Sizeof(map[int]int{}), "（字节）")  // 8字节 (仅哈希表指针)
	fmt.Println("func():        ", unsafe.Sizeof(func() {}), "（字节）")      // 8字节 (函数指针)
	fmt.Println("chan int:      ", unsafe.Sizeof(make(chan int)), "（字节）") // 8字节 (通道指针)
	fmt.Println("*int:          ", unsafe.Sizeof((*int)(nil)), "（字节）")    // 8字节 (任意指针)
}

func printSmallStruct() {
	type smallStruct struct {
		a int8
		b int8
	}

	var s smallStruct // 未初始化
	fmt.Println("字段类型大小：int8(1) + int8(1) = 2字节，无填充:", unsafe.Sizeof(s))
}

func printAlignStruct() {
	type AlignStruct struct {
		a int8
		b int64
		c int16
		d string
	}

	//1、核心误区纠正：对齐值≠类型大小（如 string 大小 16，但对齐值 8），结构体整体对齐值取内部字段的最大对齐值，而非最大大小。
	//2、计算逻辑：逐字段按 “偏移量是自身对齐值整数倍” 分配内存，最后校验总大小是最大对齐值的整数倍。
	//3、结果验证：AlignStruct 的最大对齐值是 8，逐字段计算后总大小 40（8 的 5 倍），因此最终大小为 40 字节。

	var s AlignStruct // 未初始化
	fmt.Println("打印各字段偏移量和总大小:")
	fmt.Printf("a偏移量: %d\n", unsafe.Offsetof(s.a))  // 0
	fmt.Printf("b偏移量: %d\n", unsafe.Offsetof(s.b))  // 8
	fmt.Printf("c偏移量: %d\n", unsafe.Offsetof(s.c))  // 16
	fmt.Printf("d偏移量: %d\n", unsafe.Offsetof(s.d))  // 24
	fmt.Printf("结构体总大小: %d\n", unsafe.Sizeof(s)) // 40
}

func printOptimizedStruct() {
	type AlignStruct struct {
		a int8  // 1字节
		b int64 // 8字节
		c int16 // 2字节
	}
	fmt.Println("优化前：字段未按大小从大到小排列:")
	fmt.Println(unsafe.Sizeof(AlignStruct{})) // 输出：24

	//优化前（AlignStruct：a→b→c）：总大小 24 字节
	//字段	偏移量	占用范围	占用大小	填充字节（原因）									累计大小
	//	a	0		0~0			1		7字节（1~7）→ 为了让 b 的偏移量是 8（int64 对齐值）	8
	//	b	8		8~15		8		0												16
	//	c	16		16~17		2		6字节（18~23）→ 为了让总大小是 8 的整数倍（18→24）	24
	//总填充字节：7+6=13 字节，实际有效数据仅 1+8+2=11 字节。

	type OptimizedStruct struct {
		b int64 // 8字节
		c int16 // 2字节
		a int8  // 1字节
	}
	fmt.Println("优化后：字段按大小从大到小排列:")
	fmt.Println(unsafe.Sizeof(OptimizedStruct{})) // 输出：16（而非24）

	//字段	偏移量	占用范围	占用大小	填充字节（原因）									累计大小
	//	b	 0	 	 0~7		8		0												8
	//	c	 8	 	 8~9		2		0（8 是 int16 对齐值 2 的整数倍）					10
	//	a	 10	 	 10~10		1		5 字节（11~15）→ 为了让总大小是 8 的整数倍（11→16）	16
	//总填充字节：5 字节，实际有效数据仍为 8+2+1=11 字节。

	//核心原因：减少 “碎片化填充”
	//调整字段顺序的本质是把大对齐值的字段放在前面，让小字段 “挤” 在大字段的对齐间隙里，避免小字段在前时，为了满足大字段的对齐要求产生大量填充。
}

func main() {
	printMemSize()
	printSmallStruct()
	printAlignStruct()
	printOptimizedStruct()
}
