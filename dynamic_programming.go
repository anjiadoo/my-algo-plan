package main

import (
	"fmt"
	"math"
)

//给你一个整数数组coins，表示不同面额的硬币(数量无限)；以及一个整数amount，表示总金额。
//计算并返回可以凑成总金额所需的「最少的硬币个数」。如果没有任何一种硬币组合能组成总金额，返回-1。

// 零钱兑换 - 自顶向下递归解法
func coinChange1(coins []int, amount int) int {
	// 1、状态 => amount
	// 2、选择列表 => coins
	// 3、明确dp函数的定义 => dp(amount)即答案
	// 4、结束条件 => amount==0 or amount<0

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
	//
	return -1
}

func main() {
	coins := []int{1, 2, 5}
	fmt.Println(coinChange1(coins, 11))
}
