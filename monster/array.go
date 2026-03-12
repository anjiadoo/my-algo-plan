package main

import (
	"fmt"
	"strings"
)

// 删除有序数组中的重复项 https://leetcode.cn/problems/remove-duplicates-from-sorted-array/description/
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	slow, fast := 0, 0
	for fast < len(nums) {
		if nums[slow] != nums[fast] {
			slow++
			nums[slow] = nums[fast] // 维护nums[0..slow]无重复
		}
		fast++
	}
	return slow + 1
}

// 删除有序数组中的重复项II https://leetcode.cn/problems/remove-duplicates-from-sorted-array-ii/description/
func removeDuplicates2(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	slow, fast, count := 0, 0, 0
	for fast < len(nums) {
		// 1, 1, 1, 2, 2, 2, 3, 4, 5, 5, 6
		// s
		// f
		// 维护nums[0..slow]无重复
		// slow++在fast++前面，所以需要判断 slow < fast
		if nums[slow] != nums[fast] {
			slow++
			nums[slow] = nums[fast]
		} else if slow < fast && count < 2 {
			slow++
			nums[slow] = nums[fast]
		}
		fast++
		count++
		if fast < len(nums) && nums[fast] != nums[fast-1] {
			count = 0
		}
	}
	return slow + 1
}

// 移除指定元素 https://leetcode.cn/problems/remove-element/description/
func removeElement(nums []int, val int) int {
	slow := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != val {
			nums[slow] = nums[fast] // 维护nums[0..slow)无val
			slow++
		}
	}
	return slow
}

// 移动零 https://leetcode.cn/problems/move-zeroes/
func moveZeroes(nums []int) {
	slow := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != 0 {
			nums[slow] = nums[fast] // 维护nums[0..slow)无val
			slow++
		}
	}
	for ; slow < len(nums); slow++ {
		nums[slow] = 0
	}
}

// 两数之和II - 输入有序数组 https://leetcode.cn/problems/two-sum-ii-input-array-is-sorted/
func twoSum(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1
	for left < right {
		sum := numbers[left] + numbers[right]
		if sum == target {
			return []int{left + 1, right + 1}
		} else if sum > target {
			right--
		} else {
			left++
		}
	}
	return []int{-1, -1}
}

// 最长回文子串 https://leetcode.cn/problems/longest-palindromic-substring/
func longestPalindrome(s string) string {
	// 从中心向两端扩散的双指针技巧
	maxLongestStr := ""

	palindrome := func(s string, l, r int) string {
		// 防止索引越界
		for l >= 0 && r < len(s) && s[l] == s[r] {
			l--
			r++
		}
		return s[l+1 : r]
	}

	for i := 0; i < len(s); i++ {
		str1 := palindrome(s, i, i)
		if len(str1) > len(maxLongestStr) {
			maxLongestStr = str1
		}

		str2 := palindrome(s, i, i+1)
		if len(str2) > len(maxLongestStr) {
			maxLongestStr = str2
		}
	}
	return maxLongestStr
}

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
// 副对角线 i + j == n-1 → 遍历 i + j < n-1，即 j < n-1-i（左上半区）
// 代码写 j < n-i 多包含了对角线上的点（自己和自己交换，无害）

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

// 验证回文串 https://leetcode.cn/problems/valid-palindrome/
func isPalindrome1(s string) bool {
	sb := strings.Builder{}
	for i := 0; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' {
			sb.WriteByte(s[i])
		} else if s[i] >= 'A' && s[i] <= 'Z' {
			sb.WriteByte(s[i] + 32)
		} else if s[i] >= '0' && s[i] <= '9' {
			sb.WriteByte(s[i])
		}
	}

	s = sb.String()

	left, right := 0, len(s)-1
	for left < right {
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}
	return true
}

// 颜色分类 https://leetcode.cn/problems/sort-colors/
func sortColors(nums []int) []int {
	// 双指针p0/p2=>[0,p0)存放0，(p2,len(nums)-1]存放2
	p0, p2 := 0, len(nums)-1
	p := 0

	for p <= p2 {
		if nums[p] == 0 {
			nums[p], nums[p0] = nums[p0], nums[p]
			p0++
		} else if nums[p] == 2 {
			nums[p], nums[p2] = nums[p2], nums[p]
			p2--
		} else {
			p++
		}
		if p < p0 {
			p = p0
		}
	}
	return nums
}

// 合并两个有序数组 https://leetcode.cn/problems/merge-sorted-array/description/
func merge1(nums1 []int, m int, nums2 []int, n int) {
	//输入：nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3
	//输出：[1,2,2,3,5,6]
	//解释：需要合并 [1,2,3] 和 [2,5,6] 。
	//合并结果是 [1,2,2,3,5,6] ，其中斜体加粗标注的为 nums1 中的元素。

}

func main() {

	fmt.Println(sortColors([]int{2, 0, 2, 1, 1, 0}))
	fmt.Println(sortColors([]int{2, 0, 1}))
	fmt.Println(sortColors([]int{0, 2, 1, 2, 1, 0, 2, 0, 1}))

	//fmt.Println(isPalindrome1("A man, a plan, a canal: Panama"))
	//fmt.Println(isPalindrome1("race a car"))
	//fmt.Println(isPalindrome1("0P"))

	//fmt.Println(generateMatrix(3))
	//fmt.Println(generateMatrix(1))
	//fmt.Println(generateMatrix(4))

	//matrix1 := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}
	//matrix2 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	//fmt.Println(spiralOrder(matrix2)) // [1 2 3 6 9 8 7 4 5]
	//fmt.Println(spiralOrder(matrix1)) // [1 2 3 4 8 12 11 10 9 5 6 7]

	//matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	//rotate2(matrix)
	//fmt.Println(matrix)

	//fmt.Println(longestPalindrome("babad"))
	//fmt.Println(longestPalindrome("cbbd"))

	//fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
	//fmt.Println(twoSum([]int{2, 3, 4}, 6))
	//fmt.Println(twoSum([]int{-1, 0}, -1))

	//nums := []int{1, 0, 2, 0, 3, 0, 4, 5}
	//moveZeroes(nums)
	//fmt.Println(nums)

	//fmt.Println(removeElement([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 5))

	//nums := []int{1, 1, 1, 2, 2, 2, 3, 4, 5, 5, 6}
	//fmt.Println(removeDuplicates2(nums))
	//fmt.Println(nums)
}
