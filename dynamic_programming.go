package main

import (
	"fmt"
	"math"
)

// 动态规划算法实现：
// 🌟技巧1：状态定义技巧 - DP 四步顺序：明确「状态」→ 明确「选择」→ 定义 dp 含义 → 明确 base case，顺序不能乱，状态决定 dp 数组维度
// 🌟技巧2：状态转移方程技巧 - dp[状态]=min/max(dp[状态-选择]+cost)，用「选择」连接当前状态和子状态；零钱兑换：dp(amount)=min(dp(amount-coin)+1)
// 🌟技巧3：memo初始化技巧 - 备忘录必须用不可能出现的特殊值（本题-666）区分「未计算」；0和-1均是合法返回值，用它们做哨兵会导致已算结果被覆盖
// 🌟技巧4：dp表正无穷技巧 - 求最小值时 dp table 用 amount+1 而非 math.MaxInt 表示「不可达」；math.MaxInt+1 会整数溢出，amount+1 永远大于任何合法答案且加1安全
// 🌟技巧5：subRes跳过技巧 - 自顶向下中，子问题返回-1（无解）时必须 continue 跳过，不能参与 min；1+(-1)=0 会被误判为最优解
// 🌟技巧6：状态遍历方向技巧 - 自底向上外层循环必须从小到大遍历 subAmount，保证 dp[i] 依赖的 dp[i-coin]（i-coin < i）已提前计算完毕
// 🌟技巧7：N叉树类比技巧 - 自顶向下本质是「加了 memo 剪枝的 N 叉树后序遍历」，每枚 coin 是一条分支；不加 memo 时 O(k^n)，加后变 O(k*n)

/**************************************************************
动态规划思维框架详解：

1. 明确「状态」:
	状态就是原问题中会变化的变量，通过状态的变化来描述问题的子结构。
	- 在零钱兑换问题中，状态就是目标金额amount
	- 状态转移：原问题amount的解可以通过子问题amount-coin的解推导出来

2. 明确「选择列表」:
	选择就是导致状态发生变化的行为，每个选择都会让状态转移到一个新的子问题。
	- 在零钱兑换中，选择就是选择不同面额的硬币coins
	- 每次选择一枚硬币，金额就会减少相应的面额

3. 定义dp数组/函数的含义:
	dp函数/数组表示的是在特定状态下的最优解。
	- dp(amount)表示凑出金额amount所需的最少硬币数
	- dp数组的设计要能够表达问题的解空间

4. 明确base case:
	base case是最小子问题的解，是递归或递推的终止条件。
	- 当amount=0时，不需要任何硬币，返回0
	- 当amount<0时，无解，返回-1
	- base case要确保所有可能的状态都能最终到达终止条件

**************************************************************

# 解法1：自顶向下递归的动态规划

def dp([状态1, 状态2, ...]):
    for 选择 in 所有可能的选择:
        # 此时的状态已经因为做了选择而改变
        result = 求最值(result, dp(状态1, 状态2, ...))
    return result

**************************************************************

# 解法2：自底向上迭代的动态规划

# 初始化 base case
dp[0][0][...] = base case
# 进行状态转移
for 状态1 in 状态1的所有取值：
    for 状态2 in 状态2的所有取值：
        for 选择 in 所有可能的选择:
            dp[状态1][状态2][...] = 求最值(选择1，选择2...)
**************************************************************/

//给你一个整数数组coins，表示不同面额的硬币(数量无限)；以及一个整数amount，表示总金额。
//计算并返回可以凑成总金额所需的「最少的硬币个数」。如果没有任何一种硬币组合能组成总金额，返回-1。

// 零钱兑换I - 自顶向下递归解法
func coinChange1(coins []int, amount int) int {
	// 1、明确「状态」 => amount
	// 2、明确「选择列表」 => coins
	// 3、明确「dp函数的定义」 => dp(amount)即答案
	// 4、明确「结束条件」 => amount==0 or amount<0

	memo := make([]int, amount+1)
	for i := 0; i < len(memo); i++ {
		memo[i] = -666
	}

	return dp(coins, amount, memo)
}

func dp(coins []int, amount int, memo []int) int {
	// base case
	if amount == 0 {
		return 0
	}
	if amount < 0 {
		return -1
	}

	// 查询备忘录
	if memo[amount] != -666 {
		return memo[amount]
	}

	// 类似N叉树的递归遍历
	res := math.MaxInt
	for _, coin := range coins {
		subRes := dp(coins, amount-coin, memo)
		if subRes == -1 {
			continue
		}
		res = minFunc(res, 1+subRes)
	}

	// 更新备忘录
	if res == math.MaxInt {
		memo[amount] = -1
	} else {
		memo[amount] = res
	}
	return memo[amount]
}

func minFunc(x, y int) int {
	if x > y {
		return y
	} else {
		return x
	}
}

// 零钱兑换I - 自底向上递推解法
func coinChange2(coins []int, amount int) int {
	// 1、确定 dp table 的定义
	// 2、确定状态
	// 3、确定选择列表
	// 4、base case

	// 数组大小为 amount + 1，初始值也为 amount + 1
	dpTable := make([]int, amount+1)
	for i := 1; i < len(dpTable); i++ {
		dpTable[i] = amount + 1
	}

	// base case
	dpTable[0] = 0

	// 外层 for 循环在遍历所有状态的所有取值
	for subAmount := 0; subAmount <= amount; subAmount++ {
		// 内层 for 循环在求所有选择的最小值
		for _, coin := range coins {
			// 子问题无解，跳过
			if subAmount-coin < 0 {
				continue
			}
			dpTable[subAmount] = min(dpTable[subAmount], 1+dpTable[subAmount-coin])
		}
	}

	if dpTable[amount] == amount+1 {
		return -1
	}

	return dpTable[amount]
}

// 分割等和子集 https://leetcode.cn/problems/partition-equal-subset-sum/description/
func canPartition(nums []int) bool {
	// 题意转化为背包问题就是：
	// 给一个容量为 sum/2 的背包 和 n 个物品，每个物品的重量为 nums[i]
	// 问，是否存在一种装法可以刚好装满背包？

	// 明确「状态」和「选择」
	// 状态就是「可选择的物品」和「背包容量」，选择就是「装进背包」和「不装进背包」

	// 定义 dp table
	// dp[i][j] = x 表示：只使用nums中的前 i 个物品，当背包的容量为 j 时，
	// 若 x=true 表示刚好凑满背包，x = false 表示不能刚好凑满背包。

	sum := 0
	for _, num := range nums {
		sum += num
	}

	// 和为奇数时，不可能划分成两个和相等的集合
	if sum%2 != 0 {
		return false
	}

	sum = sum / 2
	n := len(nums)
	dp := make([][]bool, n+1)
	for i := 0; i < len(dp); i++ {
		dp[i] = make([]bool, sum+1)
	}

	// base case：容量为0相当于背包满了，没有物品选择相当于凑不满(默认)
	for i := 0; i < len(dp); i++ {
		dp[i][0] = true
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= sum; j++ {
			// 这里的判断是前提，能不能装得下，如果能才有「选」或「不选」
			if j-nums[i-1] >= 0 {
				// 能不能刚好凑满 = 不加入背包的结果 || 加入背包的结果
				dp[i][j] = dp[i-1][j] || dp[i-1][j-nums[i-1]]
			} else {
				// 背包容量不足，不能装入第 i 个物品，继承 i-1 结果
				dp[i][j] = dp[i-1][j]
			}
		}
	}
	return dp[n][sum]
}

// 零钱兑换II https://leetcode.cn/problems/coin-change-ii/
func change(amount int, coins []int) int {
	// 定义 dp table
	// dp[i][j] = x 表示：如果只使用coins中的前 i 硬币面值，若想凑出金额 j, 有 dp[i][j] 种筹法。

	n := len(coins)
	dp := make([][]int, len(coins)+1)
	for i := 0; i < len(dp); i++ {
		dp[i] = make([]int, amount+1)
	}

	// base case: 背包容量为0，不操作就直接满了，所以组合数为1
	for i := 0; i < len(dp); i++ {
		dp[i][0] = 1
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= amount; j++ {
			// 这里的判断是前提，能不能装得下，如果能才有「选」或「不选」
			if j-coins[i-1] >= 0 {
				// 总共有多少种凑法 = 不选择(i-1) + 选择(i，因为硬币可以重复使用，不能重复使用的话就是上面的「i-1」形式了)
				dp[i][j] = dp[i-1][j] + dp[i][j-coins[i-1]]
			} else {
				// 背包容量不足，不能装入第 i 个物品，继承 i-1 结果
				dp[i][j] = dp[i-1][j]
			}
		}
	}
	return dp[n][amount]
}

func main() {

	fmt.Println(canPartition([]int{1, 5, 11, 5}))
	fmt.Println(canPartition([]int{1, 2, 3, 5}))
	fmt.Println(canPartition([]int{1, 2, 5}))

	//coins := []int{1, 2, 5}
	//fmt.Println(change(5, coins))
	//coins = []int{3}
	//fmt.Println(change(2, coins))

	//coins := []int{1, 2, 5}
	//fmt.Println(coinChange1(coins, 11))
	//fmt.Println(coinChange2(coins, 11))
	//
	//coins = []int{2}
	//fmt.Println(coinChange1(coins, 3))
	//fmt.Println(coinChange2(coins, 3))
	//
	//coins = []int{1}
	//fmt.Println(coinChange1(coins, 0))
	//fmt.Println(coinChange2(coins, 0))
}
