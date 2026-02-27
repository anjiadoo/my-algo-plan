package main

import (
	"fmt"
	"math"
)

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

// 零钱兑换 - 自顶向下递归解法
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

// 零钱兑换 - 自底向上递推解法
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

func main() {
	coins := []int{1, 2, 5}
	fmt.Println(coinChange1(coins, 11))
	fmt.Println(coinChange2(coins, 11))

	coins = []int{2}
	fmt.Println(coinChange1(coins, 3))
	fmt.Println(coinChange2(coins, 3))

	coins = []int{1}
	fmt.Println(coinChange1(coins, 0))
	fmt.Println(coinChange2(coins, 0))
}
