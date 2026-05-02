/*
 * ============================================================================
 *                    📘 矩阵操作 · 核心记忆框架
 * ============================================================================
 * 【两大核心思想】
 *
 *   ① 坐标变换思想：矩阵操作本质是"点的坐标映射"
 *      · 找到不动点（对称轴/不动线），确定映射规则
 *      · 推导: 新坐标 = f(旧坐标)
 *
 *   ② 降维思想：把二维问题转化为一维问题
 *      · 一维下标 index = i * n + j（行优先展平）
 *      · 反算坐标：i = index / n,  j = index % n
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式一：矩阵旋转】
 *
 *   口诀：「对角线翻转 + 逐行翻转」= 旋转
 *
 *   顺时针旋转90°:  先沿【主对角线】翻转，再逐行翻转（左右镜像）
 *   逆时针旋转90°:  先沿【副对角线】翻转，再逐行翻转（左右镜像）
 *
 *   ⚠️ 易错点1：两种旋转都需要"逐行翻转"，区别只在对角线的选择！
 *
 *   主对角线翻转（坐标互换）:
 *     · 不动点：i == j
 *     · 映射：(i, j) ↔ (j, i)
 *     · 遍历范围：j 从 i 开始（j >= i），只取右上半区，防止重复交换
 *       for i := 0; i < n; i++ { for j := i; j < n; j++ { swap } }
 *
 *   副对角线翻转:
 *     · 不动点：i + j == n - 1
 *     · 映射：(i, j) ↔ (n-1-j, n-1-i)
 *     · 遍历范围：j < n-i（即 i+j < n），只取左上半区，防止重复交换
 *       for i := 0; i < n; i++ { for j := 0; j < n-i; j++ { swap } }
 *
 *   ⚠️ 易错点2：遍历范围写错导致重复交换（翻回去了）
 *      主对角线：j 从 i 开始（不是从 0）
 *      副对角线：j < n-i（不是 j < n）
 *
 *   现场推导法（忘了不怕）:
 *     设 (i,j) 旋转后到 (i',j')，以顺时针90°为例：
 *     想象一个3×3矩阵，(0,0)顺时针转到(0,2)，得出：
 *     i' = j,  j' = n-1-i
 *     验证：(0,0)→(0,2)✅  (0,2)→(2,2)✅  (2,2)→(2,0)✅  (2,0)→(0,0)✅
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式二：螺旋遍历/生成】
 *
 *   口诀：四个边界，顺序收缩，走完一边就缩一边
 *     upperBound, lowerBound, leftBound, rightBound
 *     顺序：→ 上边  ↓ 右边  ← 下边  ↑ 左边（顺时针）
 *     走完一边，对应边界 +1 或 -1
 *
 *   模板骨架:
 *     for 未走完 {
 *         if upper <= lower  { 走上边(←→); upper++ }
 *         if left  <= right  { 走右边(↓);  right--  }
 *         if upper <= lower  { 走下边(→←); lower--  }
 *         if left  <= right  { 走左边(↑);  left++   }
 *     }
 *
 *   ⚠️ 易错点3：走下边/走左边前，必须再次判断边界！
 *      原因：走完上边后 upper++，此时 upper 可能已经 > lower（矩阵只有一行时）
 *            若不判断，下边会重复遍历已走过的行
 *      关键：四个方向共用2个条件，但每次走完一边，边界已变，下一步同条件的
 *            判断结果可能不同，所以 if 不能改成 else if！
 *
 *   ⚠️ 易错点4：终止条件区分
 *      spiralOrder（读取）：for len(res) < m*n
 *      generateMatrix（生成）：for k <= n*n
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式三：对角线分组】
 *
 *   口诀：同一对角线上的元素，坐标之差 (i-j) 相同
 *
 *   核心：用 map[i-j][]int 把同一主对角线的元素聚合
 *         （副对角线用 i+j 作为 key）
 *
 *   分组 → 排序 → 回填（按原来的遍历顺序弹出）
 *
 *   ⚠️ 易错点5：diagonalSort 排序方向陷阱
 *      题目要求从左上到右下升序，但回填时是从 i=0,j=0 开始的
 *      → 每次从 slice 尾部弹出（取 nums[len-1]），所以排序要降序
 *        sort.Slice(nums, func(i,j int) bool { return nums[i] > nums[j] })
 *      如果正序排序再从头弹出，逻辑也对，但弹出操作变成移除头部，效率低
 *
 *   ⚠️ 易错点6：回填时 mapDiagonal[key] 要及时更新（缩短 slice）
 *      mat[i][j] = nums[len(nums)-1]
 *      mapDiagonal[key] = nums[:len(nums)-1]   // 必须写！否则下一格取到同一值
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式四：二维数组降维（循环移位）】
 *
 *   口诀：展平 → 旋转数组三步法 → 原地完成
 *
 *   坐标互转:
 *     index → (i, j)：  i = index / cols,  j = index % cols
 *     (i, j) → index：  index = i * cols + j
 *
 *   三步翻转法（循环右移k位）:
 *     ① 翻转 [0, mn-k-1]
 *     ② 翻转 [mn-k, mn-1]
 *     ③ 整体翻转 [0, mn-1]
 *
 *   ⚠️ 易错点7：k 需要先对 mn 取模！
 *      k = k % mn
 *      原因：k 可能比 mn 大，移位 mn 圈等于没移
 *
 *   ⚠️ 易错点8：翻转区间边界 mn-k-1 和 mn-k（差1，不要写错）
 *      前段：[0, mn-k-1]（长度 mn-k）
 *      后段：[mn-k, mn-1]（长度 k）
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式五：转置矩阵】
 *
 *   口诀：行列互换，新矩阵维度反转 (m×n → n×m)
 *     newMatrix[j][i] = matrix[i][j]
 *
 *   ⚠️ 易错点9：非方阵必须新建矩阵！
 *      m != n 时无法原地转置，必须 make([][]int, n)，内层 make([]int, m)
 *      只有方阵（m==n）才能原地转置（但一般直接新建更安全）
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次写完矩阵题后对照检查
 *
 *     ✅ 旋转题：遍历范围是否只取了半区（防止重复交换翻回去）？
 *     ✅ 螺旋题：走下边/走左边前是否再次判断了边界条件？
 *     ✅ 对角线题：排序方向与弹出方向是否配对正确？弹出后是否更新slice？
 *     ✅ 移位题：是否对 mn 取模？翻转区间端点差1是否写对？
 *     ✅ 转置题：是否新建了正确维度（n行m列）的矩阵？
 * ============================================================================
 */

package main

import (
	"sort"
)

// 将二维矩阵原地顺时针旋转90度 https://leetcode.cn/problems/rotate-image/description/
func rotate(matrix [][]int) {
	// 主对角线翻转（左上→右下）
	// 不动点特征：主对角线上的点满足 i == j
	// 映射规则：只需要把 i 和 j 互换
	// 遍历范围：j 从 i 开始（即 j >= i），只遍历主对角线右上方

	// 沿「左上→右下」翻转（对角线）
	n := len(matrix)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}
	reverse := func(nums []int) {
		i, j := 0, len(nums)-1
		for i < j {
			nums[i], nums[j] = nums[j], nums[i]
			i++
			j--
		}
	}
	// 然后反转二维矩阵的每一行
	for i := 0; i < len(matrix); i++ {
		reverse(matrix[i])
	}
}

// 推导验证法（考场救命）
// 如果考场上忘了公式，只要记住不动点条件就能现场推：
// 副对角线不动点：i + j = n - 1，即 j = n - 1 - i
// 翻转本质是沿这条线做对称，新旧两点到这条线的"距离"相等
// 设 (i, j) 映射到 (i', j')，则需要满足：
// i' + j = n - 1 → i' = n - 1 - j ✅
// i + j' = n - 1 → j' = n - 1 - i ✅
// 这就推出了 (i, j) ↔ (n-1-j, n-1-i)，3 秒推导完毕。
//
// 遍历范围的记法
// 也很简单：遍历不动线的一侧，避免重复交换：
// 主对角线 i == j → 遍历 j >= i（右上半区）
// 副对角线 i + j == n-1 → 遍历 i + j <= n-1，即 j <= n-1-i（左上半区）
// 代码写 j < n-i（与j <= n-1-i等价） 多包含了对角线上的点（自己和自己交换，无害）

// 将二维矩阵原地逆时针旋转90度
func rotate2(matrix [][]int) {
	// 副对角线翻转（左下→右上）
	// 不动点特征：副对角线上的点满足 i + j == n - 1
	// 映射规则：用n-1做"镜子"，每个坐标都用n-1减去对方
	// 遍历范围：j < n - i，只遍历副对角线左上方

	// 沿「左下→右上」翻转（副对角线）
	n := len(matrix)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i; j++ {
			// (i, j) ↔ (n-1-j, n-1-i)
			matrix[i][j], matrix[n-1-j][n-1-i] = matrix[n-1-j][n-1-i], matrix[i][j]
		}
	}
	reverse := func(nums []int) {
		i, j := 0, len(nums)-1
		for i < j {
			nums[i], nums[j] = nums[j], nums[i]
			i++
			j--
		}
	}
	// 然后反转二维矩阵的每一行
	for i := 0; i < len(matrix); i++ {
		reverse(matrix[i])
	}
}

// 螺旋矩阵 https://leetcode.cn/problems/spiral-matrix/description/
func spiralOrder(matrix [][]int) []int {
	m, n := len(matrix), len(matrix[0])
	upperBound := 0
	rightBound := n - 1
	lowerBound := m - 1
	leftBound := 0

	var res []int

	for len(res) < m*n {

		// if upperBound <= lowerBound：防止只剩一行时，上边和下边重复遍历同一行
		// if leftBound  <= rightBound：防止只剩一列时，右边和左边重复遍历同一列

		if upperBound <= lowerBound {
			for i := leftBound; i <= rightBound; i++ {
				res = append(res, matrix[upperBound][i])
			}
			upperBound++
		}
		if leftBound <= rightBound {
			for i := upperBound; i <= lowerBound; i++ {
				res = append(res, matrix[i][rightBound])
			}
			rightBound--
		}
		if upperBound <= lowerBound {
			for i := rightBound; i >= leftBound; i-- {
				res = append(res, matrix[lowerBound][i])
			}
			lowerBound--
		}
		if leftBound <= rightBound {
			for i := lowerBound; i >= upperBound; i-- {
				res = append(res, matrix[i][leftBound])
			}
			leftBound++
		}
	}
	return res
}

// 生成螺旋矩阵 https://leetcode.cn/problems/spiral-matrix-ii/description/
func generateMatrix(n int) [][]int {
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		array := make([]int, n)
		matrix[i] = array
	}

	upperBound := 0
	rightBound := n - 1
	lowerBound := n - 1
	leftBound := 0

	// if upperBound <= lowerBound：防止只剩一行时，上边和下边重复遍历同一行
	// if leftBound  <= rightBound：防止只剩一列时，右边和左边重复遍历同一列

	k := 1
	for k <= n*n {
		if upperBound <= lowerBound {
			for i := leftBound; i <= rightBound; i++ {
				matrix[upperBound][i] = k
				k++
			}
			upperBound++
		}
		if leftBound <= rightBound {
			for i := upperBound; i <= lowerBound; i++ {
				matrix[i][rightBound] = k
				k++
			}
			rightBound--
		}
		if upperBound <= lowerBound {
			for i := rightBound; i >= leftBound; i-- {
				matrix[lowerBound][i] = k
				k++
			}
			lowerBound--
		}
		if leftBound <= rightBound {
			for i := lowerBound; i >= upperBound; i-- {
				matrix[i][leftBound] = k
				k++
			}
			leftBound++
		}
	}
	return matrix
}

// 将矩阵按对角线排序 https://leetcode.cn/problems/sort-the-matrix-diagonally/
func diagonalSort(mat [][]int) [][]int {
	// 重点：在同一个对角线上的元素，其横纵坐标之差是相同的

	// 思路：把同一对角线上的元素放到一起，排序，然后再回填到矩阵中
	mapDiagonal := make(map[int][]int)
	for i := 0; i < len(mat); i++ {
		for j := 0; j < len(mat[0]); j++ {
			diagonal := i - j
			mapDiagonal[diagonal] = append(mapDiagonal[diagonal], mat[i][j])
		}
	}

	for _, nums := range mapDiagonal {
		sort.Slice(nums, func(i, j int) bool {
			return nums[i] > nums[j]
		})
	}

	for i := 0; i < len(mat); i++ {
		for j := 0; j < len(mat[0]); j++ {
			diagonal := i - j
			nums := mapDiagonal[diagonal]
			mat[i][j] = nums[len(nums)-1]
			mapDiagonal[diagonal] = nums[:len(nums)-1]
		}
	}
	return mat
}

// 二维网格迁移 https://leetcode.cn/problems/shift-2d-grid/
func shiftGrid(grid [][]int, k int) [][]int {
	// 思路：把二维数组抽象成一维数组
	// 然后，分别对「前m*n-k个子数组」、「后k个子数组」翻转，然后再整体翻转。

	get := func(grid [][]int, index int) int {
		i := index / len(grid[0])
		j := index % len(grid[0])
		return grid[i][j]
	}

	set := func(grid [][]int, index, val int) {
		i := index / len(grid[0])
		j := index % len(grid[0])
		grid[i][j] = val
	}

	reverse := func(grid [][]int, start, end int) {
		for start < end {
			left := get(grid, start)
			right := get(grid, end)
			set(grid, start, right)
			set(grid, end, left)
			start++
			end--
		}
	}

	mn := len(grid) * len(grid[0])
	k = k % mn

	reverse(grid, 0, mn-k-1)
	reverse(grid, mn-k, mn-1)
	reverse(grid, 0, mn-1)

	return grid
}

// 转置矩阵 https://leetcode.cn/problems/transpose-matrix/description/
func transpose(matrix [][]int) [][]int {
	m, n := len(matrix), len(matrix[0])

	// 要new新的矩阵，因为行列不一定相等
	newMatrix := make([][]int, n)
	for i := 0; i < n; i++ {
		newMatrix[i] = make([]int, m)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			newMatrix[j][i] = matrix[i][j]
		}
	}
	return newMatrix
}

func main() {
	//matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	//fmt.Println(transpose(matrix))
	//
	//grid1 := [][]int{{3, 8, 1, 9}, {19, 7, 2, 5}, {4, 6, 11, 10}, {12, 0, 21, 13}}
	//grid2 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	//grid3 := [][]int{{1}, {2}, {3}, {4}, {7}, {6}, {5}}
	//fmt.Println(shiftGrid(grid1, 4))
	//fmt.Println(shiftGrid(grid2, 1))
	//fmt.Println(shiftGrid(grid3, 22))
	//
	//mat := [][]int{{3, 3, 1, 1}, {2, 2, 1, 2}, {1, 1, 1, 2}}
	//fmt.Println(diagonalSort(mat))
	//
	//fmt.Println(generateMatrix(3))
	//fmt.Println(generateMatrix(1))
	//fmt.Println(generateMatrix(4))
	//
	//matrix1 := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}
	//matrix2 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	//fmt.Println(spiralOrder(matrix2)) // [1 2 3 6 9 8 7 4 5]
	//fmt.Println(spiralOrder(matrix1)) // [1 2 3 4 8 12 11 10 9 5 6 7]
	//
	//matrix = [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	//rotate2(matrix)
	//fmt.Println(matrix)
}
