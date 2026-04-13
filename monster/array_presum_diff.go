package main

import "fmt"

// NumArray 一维区域和检索-数组不可变（一维前缀和） https://leetcode.cn/problems/range-sum-query-immutable/description/
type NumArray struct {
	PreSum []int
}

func NewNumArray(nums []int) NumArray {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] += preSum[i-1] + nums[i-1]
	}
	return NumArray{
		PreSum: preSum,
	}
}

func (n *NumArray) SumRange(left int, right int) int {
	return n.PreSum[right+1] - n.PreSum[left]
}

// NumMatrix 二维区域和检索-矩阵不可变（二维前缀和） https://leetcode.cn/problems/range-sum-query-2d-immutable/
type NumMatrix struct {
	PreSum [][]int
}

func NewNumMatrix(matrix [][]int) NumMatrix {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return NumMatrix{}
	}

	preSum := make([][]int, len(matrix)+1)
	for i := range preSum {
		preSum[i] = make([]int, len(matrix[0])+1)
	}

	for i := 1; i < len(preSum); i++ {
		for j := 1; j < len(preSum[0]); j++ {
			// [ 3, 0, 1, 4, 2 ]
			// [ 5, 6, 3, 2, 1 ]
			// [ 1, 2, 0, 1, 5 ]
			// [ 4, 1, 0, 1, 7 ]
			// [ 1, 0, 3, 0, 5 ]
			preSum[i][j] = preSum[i-1][j] + preSum[i][j-1] - preSum[i-1][j-1] + matrix[i-1][j-1]
		}
	}

	return NumMatrix{PreSum: preSum}
}

func (n *NumMatrix) SumRegion(x1, y1 int, x2, y2 int) int {
	return n.PreSum[x2+1][y2+1] - n.PreSum[x1][y2+1] - n.PreSum[x2+1][y1] + n.PreSum[x1][y1]
}

// 矩阵区域和 https://leetcode.cn/problems/matrix-block-sum/description/
func matrixBlockSum(mat [][]int, k int) [][]int {
	preSum := make([][]int, len(mat)+1)
	for i := range preSum {
		preSum[i] = make([]int, len(mat[0])+1)
	}

	for i := 1; i < len(preSum); i++ {
		for j := 1; j < len(preSum[0]); j++ {
			preSum[i][j] = preSum[i-1][j] + preSum[i][j-1] - preSum[i-1][j-1] + mat[i-1][j-1]
		}
	}

	answer := make([][]int, len(mat))
	for i := range answer {
		answer[i] = make([]int, len(mat[0]))
	}

	for i := 0; i < len(answer); i++ {
		for j := 0; j < len(answer[0]); j++ {
			// i - k <= r <= i + k, j - k <= c <= j + k 且 (r, c) 在矩阵内。
			x1 := max(i-k, 0)
			y1 := max(j-k, 0)
			x2 := min(i+k, len(answer)-1)
			y2 := min(j+k, len(answer[0])-1)
			answer[i][j] = preSum[x2+1][y2+1] - preSum[x1][y2+1] - preSum[x2+1][y1] + preSum[x1][y1]
		}
	}
	return answer
}

// 寻找数组的中心下标 https://leetcode.cn/problems/find-pivot-index/
func pivotIndex(nums []int) int {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] += preSum[i-1] + nums[i-1]
	}

	for i := 1; i < len(preSum); i++ {
		// 判断：左侧之和 与 右侧之和
		if preSum[i-1] == preSum[len(nums)]-preSum[i] {
			return i - 1
		}
	}
	return -1
}

// 除了自身以外数组的乘积 https://leetcode.cn/problems/product-of-array-except-self/
func productExceptSelf(nums []int64) []int64 {
	n := len(nums)

	// prefix[i] 表示num[0...i]的乘积
	prefix := make([]int64, n)
	prefix[0] = nums[0]
	for i := 1; i < n; i++ {
		prefix[i] = prefix[i-1] * nums[i]
	}

	// prefix[i] 表示num[i...n-1]的乘积
	suffix := make([]int64, n)
	suffix[n-1] = nums[n-1]
	for i := n - 2; i >= 0; i-- {
		suffix[i] = suffix[i+1] * nums[i]
	}

	answer := make([]int64, n)
	answer[0] = suffix[1]
	answer[n-1] = prefix[n-2]

	for i := 1; i <= n-2; i++ {
		answer[i] = prefix[i-1] * suffix[i+1]
	}
	return answer
}

// ProductOfNumbers 最后K个数的乘积 https://leetcode.cn/problems/product-of-the-last-k-numbers/
type ProductOfNumbers struct {
	preProduct []int
}

func NewProductOfNumbers() ProductOfNumbers {
	return ProductOfNumbers{preProduct: []int{1}}
}

func (p *ProductOfNumbers) Add(num int) {
	if num == 0 {
		p.preProduct = []int{1}
		return
	}
	n := len(p.preProduct)
	p.preProduct = append(p.preProduct, p.preProduct[n-1]*num)
}

func (p *ProductOfNumbers) GetProduct(k int) int {
	n := len(p.preProduct)
	// 不足k个元素，是因为最后k个元素存在0，非0情况下:k+1=n
	if k > n-1 {
		return 0
	}
	return p.preProduct[n-1] / p.preProduct[n-1-k]
}

// 连续数组 https://leetcode.cn/problems/contiguous-array/
func findMaxLength(nums []int) int {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		if nums[i-1] == 0 {
			preSum[i] = preSum[i-1] - 1
		} else {
			preSum[i] = preSum[i-1] + 1
		}
	}

	res := 0

	// 子数组和为零，说明preSum[i]==preSum[j]
	mapSumToIdx := make(map[int]int)
	for i := 0; i < len(preSum); i++ {
		if idx, ok := mapSumToIdx[preSum[i]]; !ok {
			mapSumToIdx[preSum[i]] = i
		} else {
			res = max(res, i-idx)
		}
	}
	return res
}

// 连续的子数组和 https://leetcode.cn/problems/continuous-subarray-sum/
func checkSubarraySum(nums []int, k int) bool {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	// 子数组元素总和为k的倍数,说明(preSum[j]-preSum[i])%k==0
	// 推导：preSum[j]%k == preSum[i]%k

	mapSumToIdx := make(map[int]int)
	for i := 0; i < len(preSum); i++ {
		if idx, ok := mapSumToIdx[preSum[i]%k]; ok {
			if i-idx >= 2 {
				return true
			}
		} else {
			mapSumToIdx[preSum[i]%k] = i
		}
	}
	return false
}

// 和为K的子数组 https://leetcode.cn/problems/subarray-sum-equals-k/
func subarraySum(nums []int, k int) int {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	res := 0

	// 该数组中和为k的子数组的个数: preSum[i]-preSum[j] = k
	// preSum[j] = preSum[i]-k
	// preSum[i]-k 是否存在，存在多少次就加多少

	mapSumToCounter := make(map[int]int)
	for i := 0; i < len(preSum); i++ {
		if val, ok := mapSumToCounter[preSum[i]-k]; ok {
			res += val
		}
		mapSumToCounter[preSum[i]]++
	}
	return res
}

// 表现良好的最长时间段 https://leetcode.cn/problems/longest-well-performing-interval/
func longestWPI(hours []int) int {
	preSum := make([]int, len(hours)+1)
	for i := 1; i < len(preSum); i++ {
		if hours[i-1] > 8 {
			preSum[i] = preSum[i-1] + 1
		} else {
			preSum[i] = preSum[i-1] - 1
		}
	}

	// preSum重要特性：每次变化只能是±1（因为每个小时要么+1要么-1），这意味着：
	// 1、前缀和是连续变化的：从 0 开始，每次只能增加或减少 1
	// 2、数值的出现顺序：较小的前缀和数值（如 -4、-5）通常出现在较大的前缀和数值（如 -3、-2）之后。

	res := 0
	mapSumToIdx := make(map[int]int)

	for i := 0; i < len(preSum); i++ {
		// 求最长子数组，那么记录第一次出现就好
		if _, exists := mapSumToIdx[preSum[i]]; !exists {
			mapSumToIdx[preSum[i]] = i
		}

		if preSum[i] > 0 {
			// preSum[i] 为正，说明 hours[0..i-1] 都是「表现良好的时间段」
			res = max(res, i)
		} else {
			// 求子数组和 >0 ，即：preSum[i] - preSum[j] > 0 表示区间 [j, i-1] 的和 > 0
			// preSum[j] < preSum[i] => preSum[j] <= preSum[i]-1
			// preSum[j] 最大值就是：preSum[i]-1，因为preSum的单调性(递减)，所以负数越大越靠前
			// 所以查找 mapSumToIdx[preSum[i]-1] 就是在查找最小的 j
			if j, found := mapSumToIdx[preSum[i]-1]; found {
				res = max(res, i-j)
			}
		}
	}
	return res
}

// 和可被K整除的子数组 https://leetcode.cn/problems/subarray-sums-divisible-by-k/
func subarraysDivByK(nums []int, k int) int {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	// (preSum[i] - preSum[j])%k == 0
	// preSum[i]%k - preSum[j]%k == 0
	// preSum[i]%k == preSum[j]%k

	res := 0
	mapCounter := make(map[int]int)

	for i := 0; i < len(preSum); i++ {

		// 数学上模运算的结果总是在 [0, m-1] 范围内，例如：
		// -1 mod 5 = 4（因为 -1 = (-1)×5 + 4）
		// -3 mod 5 = 2（因为 -3 = (-1)×5 + 2）
		// 但大多数编程语言的 % 运算符不同与此不同，Go中：
		// -1 % 5 = -1（而不是数学上的4）=> -1+5=4
		// -3 % 5 = -3（而不是数学上的2）=> -3+5=2

		// 子数组和可能为负，把负余数调整为正余数
		remainder := preSum[i] % k
		if remainder < 0 {
			remainder += k
		}

		if num, ok := mapCounter[remainder]; !ok {
			mapCounter[remainder] = 1
		} else {
			res += num
			mapCounter[remainder] = num + 1
		}
	}
	return res
}

// DiffArray 差分数组
type DiffArray struct {
	diff []int
}

func NewDiffArray(nums []int) *DiffArray {
	diff := make([]int, len(nums))
	diff[0] = nums[0]
	for i := 1; i < len(nums); i++ {
		diff[i] = nums[i] - nums[i-1]
	}
	return &DiffArray{diff: diff}
}

func (d *DiffArray) Add(i, j, val int) {
	d.diff[i] += val
	if j+1 < len(d.diff) {
		d.diff[j+1] -= val
	}
}
func (d *DiffArray) Result() []int {
	res := make([]int, len(d.diff))
	res[0] = d.diff[0]

	for i := 1; i < len(d.diff); i++ {
		res[i] = res[i-1] + d.diff[i]
	}
	return res
}

// 航班预订统计 https://leetcode.cn/problems/corporate-flight-bookings/
func corpFlightBookings(bookings [][]int, n int) []int {
	diff := make([]int, n)
	for _, booking := range bookings {
		start := booking[0]
		end := booking[1]
		seats := booking[2]

		diff[start-1] += seats
		if end < n {
			diff[end] -= seats
		}
	}
	res := make([]int, n)
	res[0] = diff[0]
	for i := 1; i < n; i++ {
		res[i] = res[i-1] + diff[i]
	}
	return res
}

// 拼车 https://leetcode.cn/problems/car-pooling/description/
func carPooling(trips [][]int, capacity int) bool {
	maxSite := 1001

	diff := make([]int, maxSite)
	for _, trip := range trips {
		val := trip[0]
		i := trip[1] // 第i站上车
		j := trip[2] // 第j站下车

		// 更新区间应为：diff[i,j-1]
		diff[i] += val
		if j < len(diff)-1 {
			diff[j] -= val
		}
	}

	res := make([]int, maxSite)
	res[0] = diff[0]
	// 起点就超载了
	if res[0] > capacity {
		return false
	}

	for i := 1; i < maxSite; i++ {
		res[i] = res[i-1] + diff[i]
		if res[i] > capacity {
			return false
		}
	}
	return true
}

func main() {

	fmt.Println(carPooling([][]int{{2, 1, 5}, {3, 3, 7}}, 4))
	fmt.Println(carPooling([][]int{{2, 1, 5}, {3, 3, 7}}, 5))
	fmt.Println(carPooling([][]int{{9, 0, 1}, {3, 3, 7}}, 4))

	//fmt.Println(corpFlightBookings([][]int{{1, 2, 10}, {2, 3, 20}, {2, 5, 25}}, 5))
	//fmt.Println(corpFlightBookings([][]int{{1, 2, 10}, {2, 2, 15}}, 2))

	//da := NewDiffArray([]int{0, 0, 0, 0, 0, 0, 0, 0})
	//da.Add(1, 3, 1)
	//da.Add(4, 6, -1)
	//da.Add(2, 5, 10)
	//fmt.Println(da.Result())

	//fmt.Println(subarraysDivByK([]int{-1, 2, 9}, 2))
	//fmt.Println(subarraysDivByK([]int{4, 5, 0, -2, -3, 1}, 5))

	//fmt.Println(longestWPI([]int{8, 10, 6, 16, 5}))
	//fmt.Println(longestWPI([]int{6, 6, 6}))
	//fmt.Println(longestWPI([]int{9, 9, 6, 0, 6, 6, 9}))

	//fmt.Println(subarraySum([]int{1, -1, 0}, 0))
	//fmt.Println(subarraySum([]int{1, 2, 3}, 3))
	//fmt.Println(subarraySum([]int{1, 9, 2, 8, 3, 7, 4, 6, 5}, 10))

	//fmt.Println(checkSubarraySum([]int{5, 0, 0, 0}, 3))
	//fmt.Println(checkSubarraySum([]int{0}, 1))
	//fmt.Println(checkSubarraySum([]int{23, 2, 4, 6, 7}, 6))
	//fmt.Println(checkSubarraySum([]int{23, 2, 6, 4, 7}, 13))

	//fmt.Println(findMaxLength([]int{0, 1}))
	//fmt.Println(findMaxLength([]int{0, 1, 0}))
	//fmt.Println(findMaxLength([]int{0, 1, 1, 1, 1, 1, 0, 0, 0}))

	//pp := NewProductOfNumbers()
	//pp.Add(4)
	//pp.Add(3)
	//pp.Add(2)
	//pp.Add(1)
	//fmt.Println(pp.GetProduct(1))
	//fmt.Println(pp.GetProduct(2))
	//fmt.Println(pp.GetProduct(3))
	//fmt.Println(pp.GetProduct(4))

	//fmt.Println(productExceptSelf([]int64{1, 2, 3, 4}))
	//fmt.Println(productExceptSelf([]int64{-1, 1, 0, -3, 3}))

	//fmt.Println(pivotIndex([]int{1, 7, 3, 6, 5, 6}))
	//fmt.Println(pivotIndex([]int{2, 1, -1}))

	//// 输入：mat = [[1,2,3],[4,5,6],[7,8,9]], k = 1
	//// 输出：[[12,21,16],[27,45,33],[24,39,28]]
	//fmt.Println(matrixBlockSum([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, 1))
	//
	////输入：mat = [[1,2,3],[4,5,6],[7,8,9]], k = 2
	////输出：[[45,45,45],[45,45,45],[45,45,45]]
	//fmt.Println(matrixBlockSum([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, 2))
}
