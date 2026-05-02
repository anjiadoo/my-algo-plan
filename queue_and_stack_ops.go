/*
 * ============================================================================
 *                    📘 栈与队列 · 核心记忆框架
 * ============================================================================
 * 【六大模式】
 *
 *   ① 基础栈（模拟）：用栈维护"层级/嵌套"结构
 *      适用：括号匹配、路径简化、表达式求值、字符串解码
 *
 *   ② 设计栈（辅助栈）：主栈 + 辅助栈同步维护额外信息
 *      适用：最小栈、最大频率栈
 *
 *   ③ 括号平衡（计数器代替栈）：用变量统计需求量，无需真正的栈
 *      适用：最少添加使括号有效、最少插入使括号平衡
 *
 *   ④ 单调栈：栈内元素保持单调性，用于"下一个更大/更小"问题
 *      适用：下一个更大元素、每日温度、柱状图最大矩形、移掉K位数字
 *
 *   ⑤ 循环队列/双端队列（设计题）：用数组 + 取模模拟环形结构
 *      适用：循环队列、循环双端队列
 *
 *   ⑥ 单调队列：队列内元素保持单调性，用于"滑动窗口最值"问题
 *      适用：滑动窗口最大值、绝对差不超过限制的最长子数组、最短子数组和≥K
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式一：基础栈 —— 模拟层级/嵌套】
 *
 *   核心思想：遇到"开始标记"入栈，遇到"结束标记"弹栈处理
 *
 *   simplifyPath（简化路径）：
 *     · 栈中存储的是路径的每一段文件夹名（不是字符！）
 *     · 按 "/" 分割后逐段处理：
 *       "" 或 "."  → 跳过（空段或当前目录）
 *       ".."       → 弹栈（返回上级，栈空则不操作）
 *       其他       → 入栈（进入子目录）
 *     · 最终用 "/" 拼接栈中所有段，前面加 "/"
 *
 *   isValid（有效括号）：
 *     · 左括号入栈，右括号检查栈顶是否匹配
 *     · 最终栈必须为空才是有效的
 *
 *   evalRPN（逆波兰表达式）：
 *     · 数字入栈，遇到操作符弹出两个操作数计算后压回
 *     · ⚠️ 易错点1：弹出顺序 a=倒数第二 b=栈顶，计算是 a op b
 *       stack[len-2] 是左操作数，stack[len-1] 是右操作数
 *       减法和除法不满足交换律，顺序错了结果就错
 *
 *   lengthLongestPath（文件最长路径）：
 *     · 栈模拟当前路径层级，\t 的数量 = 当前层级深度
 *     · ⚠️ 易错点2：平级或回退时需要弹栈直到 len(stack) == level
 *       每个新元素入栈前，先把"深度 >= 当前层级"的全部弹出
 *       这保证了栈始终是从根到当前节点的完整路径
 *     · 判断是否为文件：包含 "."（而非目录）
 *
 *   decodeString（字符串解码）：
 *     · 遇到 ']' 时分两步弹栈：
 *       Step1：弹出字符直到遇到 '['（得到待重复的字符串）
 *       Step2：继续弹出数字字符（得到重复次数）
 *     · 将 repeat(str, num) 的结果逐字符压回栈
 *     · ⚠️ 易错点3：数字可能是多位数（如 "12[a]"），必须循环弹出所有连续数字
 *     · ⚠️ 易错点4：弹出字符/数字时是逆序的，拼接时需要 string(ch) + str 前插
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式二：设计栈 —— 辅助栈同步维护额外信息】
 *
 *   MinStack（最小栈）：
 *     · 主栈 stack 正常存值
 *     · 辅助栈 minStack 长度与主栈相同，minStack[i] = stack[0..i] 的最小值
 *     · Push 时：minStack 压入 min(val, minStack栈顶)
 *     · Pop 时：两个栈同步弹出
 *     · GetMin：直接返回 minStack 栈顶
 *     · ⚠️ 易错点5：minStack 不是只存"出现过的最小值"，而是每层都存一个
 *       这样 Pop 时不需要额外判断，直接同步弹出即可
 *
 *   FreqStack（最大频率栈）：
 *     · 核心：两个映射 + 一个变量
 *       val2Freq  map[int]int    记录每个值的当前频率
 *       freq2Vals map[int][]int  记录每个频率对应的值列表（栈结构）
 *       maxFreq   int            当前最大频率
 *     · Push：freq[val]++，然后 freq2Vals[freq] 追加 val，更新 maxFreq
 *     · Pop：从 freq2Vals[maxFreq] 弹出栈顶，val2Freq[val]--
 *       如果 freq2Vals[maxFreq] 为空则 maxFreq--
 *     · ⚠️ 易错点6：同一个 val 会出现在多个频率层级中
 *       val=5 出现3次 → freq2Vals[1]、freq2Vals[2]、freq2Vals[3] 中都有 5
 *       这不是冗余！Pop 时只从最高频率层弹出，val2Freq 减1，完美对应
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式三：括号平衡 —— 计数器代替栈】
 *
 *   核心思想：括号匹配问题中，若只需统计"最少添加数"而不需要具体位置，
 *            用 needLeft/needRight 两个计数器即可，无需真正的栈。
 *
 *   minAddToMakeValid（最少添加使括号有效）：
 *     · needRight：当前未匹配的左括号数（= 还需要多少右括号）
 *     · needLeft：多余的右括号数（= 还需要多少左括号）
 *     · 遇到 '('：needRight++
 *     · 遇到 ')'：needRight--，若 < 0 则 needLeft++, needRight = 0
 *     · 答案：needLeft + needRight
 *
 *   minInsertions（最少插入使括号平衡，一个 '(' 对应两个 ')' ）：
 *     · ⚠️ 易错点7：一个左括号需要两个右括号匹配
 *       遇到 '(' 时 needRight += 2（不是+1）
 *     · ⚠️ 易错点8：needRight 为奇数时的修正
 *       遇到 '(' 时如果 needRight 是奇数，说明前面有一个落单的 ')'，
 *       需要插入一个 ')' 来配对，所以 needLeft++, needRight--
 *       这保证 needRight 始终是偶数（因为右括号总是成对消耗的）
 *     · 遇到 ')'：needRight--，若 == -1 则 needLeft++, needRight = 1
 *       （多出的 ')' 需要一个 '(' 来兜底，同时还差一个 ')' 凑成一对）
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式四：单调栈 —— "下一个更大/更小"问题族】
 *
 *   核心模板（从后往前遍历，单调递减栈 → 找下一个更大元素）：
 *     for i := n-1; i >= 0; i-- {
 *         for len(stack) > 0 && nums[i] >= stack[top] {
 *             stack.pop()   // 比当前矮的都没用了，弹掉
 *         }
 *         res[i] = stack为空 ? -1 : stack[top]   // 栈顶就是下一个更大
 *         stack.push(nums[i])
 *     }
 *
 *   ⚠️ 易错点9：栈中存的是值还是下标？
 *     · nextGreaterElement：存值（因为需要映射到另一个数组）
 *     · dailyTemperatures：存下标（因为需要计算距离 stack[top] - i）
 *     · canSeePersonsCount：存下标（需要计算可见人数）
 *     根据题目需要"值"还是"距离"来决定栈中存什么
 *
 *   ⚠️ 易错点10：环形数组的处理（nextGreaterElements）
 *     · 数组长度翻倍：遍历 2n-1 到 0，用 i%n 取实际下标
 *     · 这样每个元素都能"看到"它后面一圈的所有元素
 *
 *   canSeePersonsCount（队列中可见人数）：
 *     · ⚠️ 易错点11：被弹出的元素也是"可见的"！
 *       每弹出一个 → count++（那些矮个子被当前人看见了）
 *       最终如果栈非空 → 还能看见栈顶那个更高的人 → count + 1
 *     · 与普通单调栈不同：不仅关心栈顶，还关心弹出了几个
 *
 *   finalPrices（商品折扣）：
 *     · 本质是"下一个更小或相等元素"→ 单调递增栈
 *     · 弹栈条件：stack[top] > prices[i]（严格大于才弹）
 *     · 折扣 = 下一个 ≤ 自己的价格
 *
 *   removeKdigits（移掉K位数字）：
 *     · ⚠️ 易错点12：这里是从前往后遍历（不是从后往前）！
 *       维护单调递增栈：栈顶 > 当前字符 → 弹栈且 k--
 *       目的是让高位尽可能小
 *     · ⚠️ 易错点13：前导零处理
 *       如果栈为空且当前字符是 '0'，直接跳过（不入栈）
 *     · ⚠️ 易错点14：k 没用完的情况
 *       遍历结束后若 k > 0，此时栈一定是单调递增的，
 *       从栈顶删除 k 个（删最大的）
 *     · 栈为空返回 "0"
 *
 *   carFleet（车队）：
 *     · 关键洞察：按起始位置排序后，到达时间的单调递减序列数 = 车队数
 *     · 计算到达时间：time = (target - pos) / speed
 *     · 从后往前遍历（离终点最近的先看），如果 time[i] > maxTime
 *       说明形成了新车队，maxTime = time[i], 车队数++
 *     · ⚠️ 易错点15：到达时间更短的车会被前面更慢的车卡住，合并为一个车队
 *       所以只有严格递增的到达时间才产生新车队
 *
 *   findUnsortedSubarrayII（最短无序子数组，单调栈解法）：
 *     · 正向递增栈：弹出的元素下标的最小值 = 左边界
 *     · 反向递减栈：弹出的元素下标的最大值 = 右边界
 *     · ⚠️ 易错点16：两个栈都没有弹出过元素 → 数组本身有序，返回 0
 *
 *   largestRectangleArea（柱状图最大矩形）：
 *     · ⚠️ 易错点17：这是单调栈最经典也最难的题，模板不同于"下一个更大"
 *     · 使用单调递增栈（从前往后遍历），弹栈时机 = 找到了右边界
 *     · 弹出元素 h 时：
 *       右边界 = 当前 i（第一个比 h 矮的）
 *       左边界 = 弹出后的新栈顶（第一个比 h 矮的左边元素）
 *       宽度 = i - stack[top] - 1（开区间）
 *       面积 = h * 宽度
 *     · ⚠️ 易错点18：哨兵技巧 —— 在 heights 首尾各加一个 0
 *       头部的 0：保证栈底永远有元素，计算宽度时不会栈空
 *       尾部的 0：保证所有柱子最终都会被弹出（触发面积计算）
 *       不加哨兵就需要额外处理边界，代码更复杂
 *     · ⚠️ 易错点19：是先取高度再算宽度（先pop再用新栈顶算左边界）
 *       顺序：height = nums[stack.pop()] → width = i - stack[top] - 1
 *       如果先算宽度再pop，stack[top]还是自己，宽度就错了
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式五：循环队列/双端队列 —— 数组 + 取模】
 *
 *   核心公式：
 *     前进一步：(index + 1) % cap
 *     后退一步：(index - 1 + cap) % cap
 *
 *   MyCircularQueue（循环队列）：
 *     · first 指向队头，last 指向队尾的下一个空位
 *     · 用 size 变量区分"满"和"空"（而非浪费一个空间）
 *     · EnQueue：data[last] = val, last 前进
 *     · DeQueue：first 前进
 *     · Rear（取队尾）：data[(last - 1 + cap) % cap]
 *
 *   MyCircularDeque（循环双端队列）：
 *     · InsertFront：start 先后退一步，再写入
 *     · InsertLast：先写入 data[end]，end 再前进
 *     · ⚠️ 易错点20：InsertFront 是"先移后写"，InsertLast 是"先写后移"
 *       因为 start 指向第一个有效元素，end 指向最后一个有效元素的下一位
 *       前端插入需要在 start 前面开辟空间，后端插入直接写在 end 处
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【模式六：单调队列 —— 滑动窗口最值】
 *
 *   核心结构（MonotonicQueue）：
 *     · data[]：原始队列，维护窗口内所有元素
 *     · maxQ[]：单调递减队列，队头是当前窗口最大值
 *     · minQ[]：单调递增队列，队头是当前窗口最小值
 *
 *   push(elem)：
 *     · maxQ：从尾部弹出所有 < elem 的元素，再追加 elem
 *     · minQ：从尾部弹出所有 > elem 的元素，再追加 elem
 *   pop()：
 *     · 弹出 data 队头，如果等于 maxQ/minQ 的队头则同步弹出
 *
 *   ⚠️ 易错点21：pop 时的判等比较
 *     只有 data 队头 == maxQ/minQ 队头时才弹出辅助队列
 *     因为辅助队列中间的元素可能已经被 push 时弹掉了
 *     这里用"值相等"判断（适用于整数），如果有重复值且需要精确，应存下标
 *
 *   maxSlidingWindow（滑动窗口最大值）：
 *     · 窗口大小为 k，前 k-1 个元素只 push 不取值
 *     · 从第 k 个开始：push → getMax → pop（先加新元素，取最大，再移除最旧）
 *     · ⚠️ 易错点22：push/pop 的顺序影响窗口大小
 *       先 push 再 pop 保证窗口始终恰好 k 个元素时取最大值
 *
 *   longestSubarray（绝对差不超过限制的最长子数组）：
 *     · 滑动窗口 + 单调队列：窗口内 max - min <= limit
 *     · 右端 push 后，若 max - min > limit → 左端 pop 并 left++
 *     · ⚠️ 易错点23：窗口长度计算用 right - left（左闭右开）
 *       因为 right 在 push 后已经 ++了，所以区间是 [left, right)
 *
 *   shortestSubarray（和至少为K的最短子数组）：
 *     · ⚠️ 易错点24：有负数不能用普通滑动窗口！
 *       用前缀和 + 单调队列：preSum[right] - preSum[left] >= k
 *     · 单调队列维护前缀和的最小值（这样 preSum[right] - min 最大）
 *     · 内层 while：当 preSum[right] - queue.getMin() >= k 时，
 *       记录答案并 pop 左端（因为更短的子数组可能存在）
 *     · ⚠️ 易错点25：为什么pop不会错过答案？
 *       pop 出的 preSum[left] 对于之后更大的 right 也能满足，
 *       但 right - left 只会更大，不可能更优，所以可以安全丢弃
 *
 *   maxSubarraySumCircular（环形子数组最大和）：
 *     · 环形处理：前缀和数组长度 2n+1（原数组复制一遍）
 *     · 窗口大小限制为 n（子数组长度不超过原数组长度）
 *     · 对每个 preSum[i]，答案候选 = preSum[i] - 窗口内 min(preSum)
 *     · ⚠️ 易错点26：窗口大小维护
 *       当 window.size == n 时先 pop 再 push，保证窗口不超过 n
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【单调栈 vs 单调队列 —— 何时用哪个？】
 *
 *   单调栈：找"下一个更大/更小"，只关心一侧（通常从后往前）
 *     · 只需要栈顶操作，不需要从底部弹出
 *     · 每个元素入栈/出栈各一次 → O(n)
 *
 *   单调队列：维护"滑动窗口内的最值"，需要两端操作
 *     · 尾部：push 时弹出不符合单调性的元素
 *     · 头部：窗口滑出时弹出过期元素
 *     · 每个元素入队/出队各一次 → O(n)
 *
 *   判断标准：
 *     问 "某个元素右边（或左边）第一个比它大的是谁" → 单调栈
 *     问 "一个滑动窗口内的最大值/最小值"             → 单调队列
 *
 * ────────────────────────────────────────────────────────────────────────────
 * 【防错清单】—— 每次写完栈/队列题后对照检查
 *
 *     ✅ 基础栈：栈中存的是什么？（字符、路径段、数字、下标？）
 *     ✅ evalRPN：弹出顺序对了吗？a=倒数第二 是左操作数，b=栈顶 是右操作数？
 *     ✅ decodeString：数字是多位的情况处理了吗？拼接方向对了吗？
 *     ✅ MinStack：辅助栈是每层都存值（不是只存变化点）？
 *     ✅ FreqStack：Pop 后 freq2Vals 为空时 maxFreq-- 了吗？
 *     ✅ 括号平衡：needRight < 0 时的修正逻辑对了吗？
 *     ✅ 单调栈方向：找"下一个更大"→ 从后往前 + 递减栈？
 *     ✅ 单调栈存值/存下标：需要距离存下标，需要映射存值？
 *     ✅ 环形数组：遍历范围是 2n-1 到 0，下标用 i%n？
 *     ✅ removeKdigits：前导零跳过了吗？k没用完时从栈顶删了吗？
 *     ✅ 柱状图矩形：加了首尾哨兵0吗？先pop取高度再算宽度？
 *     ✅ 循环队列取模：后退用 (idx-1+cap)%cap（不是 idx-1%cap）？
 *     ✅ 单调队列pop：判断了队头 == 被移除元素才弹辅助队列？
 *     ✅ 滑动窗口大小：push/pop 顺序保证窗口恰好 k 个？
 *     ✅ shortestSubarray：用的是前缀和+单调队列（不是普通滑动窗口）？
 * ============================================================================
 */

package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// 简化路径 https://leetcode.cn/problems/simplify-path/description/
func simplifyPath(path string) string {
	// 一定要明确栈中存储的是啥
	// 这里存的是文件夹组成路径
	var stack []string

	for _, str := range strings.Split(path, "/") {
		switch str {
		case "", ".":
			continue
		case "..":
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		default:
			stack = append(stack, str)
		}
	}
	if len(stack) == 0 {
		return "/"
	}
	return "/" + strings.Join(stack, "/")
}

// 有效的括号 https://leetcode.cn/problems/valid-parentheses/description/
func isValid(s string) bool {
	var stack []byte
	for _, ch := range s {
		switch ch {
		case '(', '{', '[':
			stack = append(stack, byte(ch))
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false
			}
			stack = stack[:len(stack)-1]
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				return false
			}
			stack = stack[:len(stack)-1]
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	if len(stack) > 0 {
		return false
	}
	return true
}

// 逆波兰表达式求值 https://leetcode.cn/problems/evaluate-reverse-polish-notation/
func evalRPN(tokens []string) int {
	var stack []int
	for _, token := range tokens {
		if strings.Contains("+-*/", token) {

			a := stack[len(stack)-2]
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "*":
				stack = append(stack, a*b)
			case "-":
				stack = append(stack, a-b)
			case "/":
				stack = append(stack, a/b)
			}
		} else {
			num, _ := strconv.Atoi(token)
			stack = append(stack, num)
		}
	}
	return stack[0]
}

// 文件的最长绝对路径 https://leetcode.cn/problems/longest-absolute-file-path/
func lengthLongestPath(input string) int {
	var stack []string
	maxLength := 0

	for _, part := range strings.Split(input, "\n") {
		level := strings.LastIndexByte(part, '\t') + 1

		// 平级需要出栈
		for len(stack) > level {
			stack = stack[:len(stack)-1]
		}

		// 子级需要入栈
		stack = append(stack, part[level:])

		// 如果是文件，更新最长路径
		if strings.Contains(part, ".") {
			fullPath := strings.Join(stack, "/")
			maxLength = max(maxLength, len(fullPath))
		}
	}
	return maxLength
}

// 字符串解码 https://leetcode.cn/problems/decode-string/description/
func decodeString(s string) string {
	var stack []byte
	for _, ch := range s {
		if ch == ']' {
			str := ""
			for len(stack) > 0 && stack[len(stack)-1] != '[' {
				str = string(stack[len(stack)-1]) + str
				stack = stack[:len(stack)-1]
			}
			// 去掉'['
			stack = stack[:len(stack)-1]

			sNum := ""
			for len(stack) > 0 && stack[len(stack)-1] >= '0' && stack[len(stack)-1] <= '9' {
				sNum = string(stack[len(stack)-1]) + sNum
				stack = stack[:len(stack)-1]
			}

			num, _ := strconv.Atoi(sNum)
			subStr := strings.Repeat(str, num)
			stack = append(stack, []byte(subStr)...)
		} else {
			stack = append(stack, byte(ch))
		}
	}
	return string(stack)
}

// MinStack 最小栈 https://leetcode.cn/problems/min-stack/
type MinStack struct {
	stack    []int // 原始栈
	minStack []int // 存最小元素栈长度与stack相等
}

func (s *MinStack) Push(val int) {
	s.stack = append(s.stack, val)
	if len(s.minStack) == 0 || val < s.minStack[len(s.minStack)-1] {
		s.minStack = append(s.minStack, val)
	} else {
		s.minStack = append(s.minStack, s.minStack[len(s.minStack)-1])
	}
}

func (s *MinStack) Pop() {
	s.stack = s.stack[:len(s.stack)-1]
	s.minStack = s.minStack[:len(s.minStack)-1]
}

func (s *MinStack) Top() int {
	return s.stack[len(s.stack)-1]
}

func (s *MinStack) GetMin() int {
	return s.minStack[len(s.minStack)-1]
}

// FreqStack 最大频率栈 https://leetcode.cn/problems/maximum-frequency-stack/description/
type FreqStack struct {
	maxFreq   int           // 记录FreqStack中元素的最大频率
	val2Freq  map[int]int   // 记录FreqStack中每个val对应的出现频率
	freq2Vals map[int][]int // 记录频率freq对应的val列表
}

func Constructor() FreqStack {
	return FreqStack{
		val2Freq:  make(map[int]int),
		freq2Vals: make(map[int][]int),
	}
}

func (s *FreqStack) Push(val int) {
	freq := s.val2Freq[val] + 1
	s.val2Freq[val] = freq

	s.freq2Vals[freq] = append(s.freq2Vals[freq], val)
	if freq > s.maxFreq {
		s.maxFreq = freq
	}
}

func (s *FreqStack) Pop() int {
	vals := s.freq2Vals[s.maxFreq]
	res := vals[len(vals)-1]
	vals = vals[:len(vals)-1]
	s.freq2Vals[s.maxFreq] = vals

	s.val2Freq[res]--
	if len(s.freq2Vals[s.maxFreq]) == 0 {
		delete(s.freq2Vals, s.maxFreq)
		s.maxFreq--
	}
	return res
}

// 使括号有效的最少添加 https://leetcode.cn/problems/minimum-add-to-make-parentheses-valid/
func minAddToMakeValid(s string) int {
	needRight := 0 // 对右括号的需求
	needLeft := 0  // 对左括号的需求

	for _, ch := range s {
		if ch == '(' {
			needRight++
		}
		if ch == ')' {
			needRight--
			if needRight < 0 {
				needRight = 0
				needLeft++
			}
		}
	}
	return needLeft + needRight
}

// 平衡括号字符串的最少插入次数 https://leetcode.cn/problems/minimum-insertions-to-balance-a-parentheses-string/description/
func minInsertions(s string) int {
	needRight := 0 // 对右括号的需求
	needLeft := 0  // 对左括号的需求

	for _, ch := range s {
		if ch == '(' {
			needRight += 2

			// 难点：当遇到左括号时
			// 若对右括号的需求量为奇数
			// 需要插入 1 个右括号
			// 一个左括号对应两个右括号
			if needRight%2 == 1 {
				needLeft++
				needRight--
			}
		}
		if ch == ')' {
			needRight--
			if needRight == -1 {
				needLeft++
				needRight = 1
			}
		}
	}
	return needLeft + needRight
}

// RecentCounter 最近的请求次数 https://leetcode.cn/problems/number-of-recent-calls/
type RecentCounter struct {
	queue []int
}

func (q *RecentCounter) Ping(t int) int {
	q.queue = append(q.queue, t)
	for len(q.queue) > 0 && q.queue[0] < t-3000 {
		q.queue = q.queue[1:]
	}
	return len(q.queue)
}

// MyCircularQueue 设计循环队列 https://leetcode.cn/problems/design-circular-queue/
type MyCircularQueue struct {
	first int
	last  int
	size  int
	cap   int
	data  []int
}

func NewMyCircularQueue(k int) MyCircularQueue {
	return MyCircularQueue{
		cap:  k,
		data: make([]int, k),
	}
}

func (m *MyCircularQueue) EnQueue(value int) bool {
	if m.IsFull() {
		return false
	}
	m.data[m.last] = value
	m.last = (m.last + 1 + m.cap) % m.cap
	m.size++
	return true
}

func (m *MyCircularQueue) DeQueue() bool {
	if m.IsEmpty() {
		return false
	}
	m.data[m.first] = -1
	m.first = (m.first + 1 + m.cap) % m.cap
	m.size--
	return true
}

func (m *MyCircularQueue) Front() int {
	if m.IsEmpty() {
		return -1
	}
	return m.data[m.first]
}

func (m *MyCircularQueue) Rear() int {
	if m.IsEmpty() {
		return -1
	}
	return m.data[(m.last-1+m.cap)%m.cap]
}

func (m *MyCircularQueue) IsEmpty() bool {
	return m.size == 0
}

func (m *MyCircularQueue) IsFull() bool {
	return m.size == m.cap
}

// MyCircularDeque 设计循环双端队列 https://leetcode.cn/problems/design-circular-deque/
type MyCircularDeque struct {
	start int
	end   int
	size  int
	cap   int
	data  []int
}

func NewMyCircularDeque(k int) MyCircularDeque {
	return MyCircularDeque{
		cap:  k,
		data: make([]int, k),
	}
}

func (m *MyCircularDeque) InsertFront(value int) bool {
	if m.IsFull() {
		return false
	}
	m.start = (m.start - 1 + m.cap) % m.cap
	m.data[m.start] = value
	m.size++
	return true
}

func (m *MyCircularDeque) InsertLast(value int) bool {
	if m.IsFull() {
		return false
	}
	m.data[m.end] = value
	m.end = (m.end + 1 + m.cap) % m.cap
	m.size++
	return true
}

func (m *MyCircularDeque) DeleteFront() bool {
	if m.IsEmpty() {
		return false
	}
	m.data[m.start] = -1
	m.start = (m.start + 1 + m.cap) % m.cap
	m.size--
	return true
}

func (m *MyCircularDeque) DeleteLast() bool {
	if m.IsEmpty() {
		return false
	}
	m.end = (m.end - 1 + m.cap) % m.cap
	m.size--
	return true
}

func (m *MyCircularDeque) GetFront() int {
	if m.IsEmpty() {
		return -1
	}
	return m.data[m.start]
}

func (m *MyCircularDeque) GetRear() int {
	if m.IsEmpty() {
		return -1
	}
	return m.data[(m.end-1+m.cap)%m.cap]
}

func (m *MyCircularDeque) IsEmpty() bool {
	return m.size == 0
}

func (m *MyCircularDeque) IsFull() bool {
	return m.size == m.cap
}

// 下一个更大元素I https://leetcode.cn/problems/next-greater-element-i/description/
func nextGreaterElement(nums1 []int, nums2 []int) []int {
	mapNextGreater := map[int]int{}

	// 从后往前遍历，栈里存的是当前元素的下一个更大元素
	var stack []int

	for i := len(nums2) - 1; i >= 0; i-- {
		for len(stack) > 0 && nums2[i] >= stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			mapNextGreater[nums2[i]] = -1
		} else {
			mapNextGreater[nums2[i]] = stack[len(stack)-1]
		}
		stack = append(stack, nums2[i])
	}

	var res []int
	for _, num := range nums1 {
		res = append(res, mapNextGreater[num])
	}
	return res
}

// 每日温度 https://leetcode.cn/problems/daily-temperatures/
func dailyTemperatures(temperatures []int) []int {
	res := make([]int, len(temperatures))
	var stack []int

	for i := len(temperatures) - 1; i >= 0; i-- {
		for len(stack) > 0 && temperatures[i] >= temperatures[stack[len(stack)-1]] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			res[i] = 0
		} else {
			res[i] = stack[len(stack)-1] - i
		}
		stack = append(stack, i)
	}
	return res
}

// 下一个更大元素II https://leetcode.cn/problems/next-greater-element-ii/
func nextGreaterElements(nums []int) []int {
	result := make([]int, len(nums))
	n := len(nums)
	var stack []int

	// 数组长度加倍模拟环形数组
	for i := 2*n - 1; i >= 0; i-- {
		for len(stack) > 0 && nums[i%n] >= stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			result[i%n] = -1
		} else {
			result[i%n] = stack[len(stack)-1]
		}
		stack = append(stack, nums[i%n])
	}

	return result
}

type ListNode1 struct {
	Val  int
	Next *ListNode1
}

func NewMyListNode_(nums []int) *ListNode1 {
	dummy := &ListNode1{}
	p := dummy
	for i := 0; i < len(nums); i++ {
		p.Next = &ListNode1{Val: nums[i]}
		p = p.Next
	}
	return dummy.Next
}

// 链表中的下一个更大节点 https://leetcode.cn/problems/next-greater-node-in-linked-list/description/
func nextLargerNodes(head *ListNode1) []int {
	var nums []int
	for p := head; p != nil; p = p.Next {
		nums = append(nums, p.Val)
	}

	res := make([]int, len(nums))
	var stack []int

	for i := len(nums) - 1; i >= 0; i-- {
		for len(stack) > 0 && stack[len(stack)-1] <= nums[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			res[i] = 0
		} else {
			res[i] = stack[len(stack)-1]
		}
		stack = append(stack, nums[i])
	}
	return res
}

// 队列中可以看到的人数 https://leetcode.cn/problems/number-of-visible-people-in-a-queue/description/
func canSeePersonsCount(heights []int) []int {
	res := make([]int, len(heights))
	var stack []int

	for i := len(heights) - 1; i >= 0; i-- {
		count := 0
		for len(stack) > 0 && heights[stack[len(stack)-1]] <= heights[i] {
			stack = stack[:len(stack)-1]
			count++ //⚠️记录被挤掉的人
		}
		if len(stack) == 0 {
			res[i] = count
		} else {
			res[i] = count + 1 //⚠️除了可以看到“被挤掉的人”外，比自己高的那个人也能看到
		}
		stack = append(stack, i)
	}

	return res
}

// 商品折扣后的最终价格 https://leetcode.cn/problems/final-prices-with-a-special-discount-in-a-shop/description/
func finalPrices(prices []int) []int {
	res := make([]int, len(prices))
	var stack []int

	for i := len(prices) - 1; i >= 0; i-- {
		for len(stack) > 0 && stack[len(stack)-1] > prices[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			res[i] = prices[i]
		} else {
			res[i] = prices[i] - stack[len(stack)-1]
		}
		stack = append(stack, prices[i])
	}
	return res
}

// 移掉K位数字 https://leetcode.cn/problems/remove-k-digits/description/
func removeKdigits(num string, k int) string {
	var stack []rune
	for _, ch := range num {
		// 上一个更小或相等
		for len(stack) > 0 && stack[len(stack)-1] > (ch) && k > 0 {
			stack = stack[:len(stack)-1]
			k--
		}

		// 防止 0 作为数字的开头
		if len(stack) == 0 && ch == '0' {
			continue
		}
		stack = append(stack, ch)
	}

	// ⚠️k没用完，栈一定是单调递增的，上面的for循环
	for k > 0 && len(stack) > 0 {
		stack = stack[:len(stack)-1]
		k--
	}

	if len(stack) == 0 {
		return "0"
	}
	return string(stack)
}

// 车队 https://leetcode.cn/problems/car-fleet/
func carFleet(target int, position []int, speed []int) int {
	// 关键：按起始位置排序后，到达时间快的车会被后面到达时间慢的车卡住
	// 所以应该是计算单调递减序列

	type Pair struct {
		pos   int
		speed int
	}

	var pairs []Pair
	for i := 0; i < len(position); i++ {
		pairs = append(pairs, Pair{pos: position[i], speed: speed[i]})
	}

	// 按起始位置排序
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].pos < pairs[j].pos
	})

	// 计算到达时间
	var times []float64
	for _, pair := range pairs {
		times = append(times, float64(target-pair.pos)/float64(pair.speed))
	}

	// 单调递减栈大小就是答案
	//var stack []int
	//for i := len(times) - 1; i >= 0; i-- {
	//	for len(stack) > 0 && times[i] <= stack[len(stack)-1] {
	//		stack = stack[:len(stack)-1]
	//	}
	//	stack = append(stack, times[i])
	//}

	// 避免使用栈模拟，倒序遍历取递增序列就是答案
	var maxTime float64
	res := 0
	for i := len(times) - 1; i >= 0; i-- {
		if times[i] > maxTime {
			res++
			maxTime = times[i]
		}
	}
	return res
}

// 最短无序连续子数组 https://leetcode.cn/problems/shortest-unsorted-continuous-subarray/
func findUnsortedSubarray(nums []int) int {
	temp := make([]int, len(nums))
	copy(temp, nums)
	sort.Ints(temp)

	left := 0
	for i := 0; i < len(nums); i++ {
		if temp[left] != nums[i] {
			break
		}
		left++
	}

	right := len(nums) - 1
	for i := len(nums) - 1; i >= 0; i-- {
		if temp[right] != nums[i] {
			break
		}
		right--
	}

	if left == len(nums) && right == -1 {
		return 0
	}
	return right - left + 1
}

func findUnsortedSubarrayII(nums []int) int {
	left, right := len(nums), -1

	// 递增栈，弹出的元素都是乱序元素
	var incrStack []int
	for i := 0; i < len(nums); i++ {
		for len(incrStack) > 0 && nums[incrStack[len(incrStack)-1]] > nums[i] {
			left = min(left, incrStack[len(incrStack)-1])
			incrStack = incrStack[:len(incrStack)-1]
		}
		incrStack = append(incrStack, i)
	}

	// 递减栈，弹出的元素都是乱序元素
	var decrStack []int
	for i := len(nums) - 1; i >= 0; i-- {
		for len(decrStack) > 0 && nums[decrStack[len(decrStack)-1]] < nums[i] {
			right = max(right, decrStack[len(decrStack)-1])
			decrStack = decrStack[:len(decrStack)-1]
		}
		decrStack = append(decrStack, i)
	}

	// 单调栈没有弹出任何元素，说明nums本来就是有序的
	if left == len(nums) && right == -1 {
		return 0
	}
	return right - left + 1
}

// 柱状图中最大的矩形 https://leetcode.cn/problems/largest-rectangle-in-histogram/
func largestRectangleArea(heights []int) int {
	nums := make([]int, len(heights)+2)
	for i := 0; i < len(heights); i++ {
		nums[i+1] = heights[i]
	}

	// 求最大，那么就需要穷举，思路是：单调递增栈
	// 找 heights[i] 左侧第一个更小元素的索引 left
	// 找 heights[i] 右侧第一个更小元素的索引 right
	// 单栈需要 pop 元素的时候，就找到了右侧第一个更小

	var stack []int
	maxArea := 0

	for i := 0; i < len(nums); i++ {
		for len(stack) > 0 && nums[stack[len(stack)-1]] > nums[i] {
			// 先取高度，先 pop
			height := nums[stack[len(stack)-1]]
			stack = stack[:len(stack)-1]

			// 再算宽度：用pop后的新栈顶作为左边界（开区间）
			width := i - stack[len(stack)-1] - 1
			maxArea = max(maxArea, width*height)
		}
		stack = append(stack, i)
	}
	return maxArea
}

// 滑动窗口最大值 https://leetcode.cn/problems/sliding-window-maximum/
func maxSlidingWindow(nums []int, k int) []int {
	mq := NewMonotonicQueue()
	var res []int

	for i := 0; i < len(nums); i++ {
		if i < k-1 {
			mq.push(nums[i])
		} else {
			mq.push(nums[i])
			res = append(res, mq.getMax())
			mq.pop()
		}
	}
	return res
}

// MonotonicQueue 单调队列
type MonotonicQueue struct {
	data []int //原始队列
	maxQ []int //单调递减，维护队列最大值
	minQ []int //单调递增，维护队列最小值
}

func NewMonotonicQueue() *MonotonicQueue {
	return &MonotonicQueue{}
}

func (mq *MonotonicQueue) getMax() int {
	if len(mq.maxQ) == 0 {
		return -1
	}
	return mq.maxQ[0]
}

func (mq *MonotonicQueue) getMin() int {
	if len(mq.minQ) == 0 {
		return -1
	}
	return mq.minQ[0]
}

func (mq *MonotonicQueue) size() int {
	return len(mq.data)
}

func (mq *MonotonicQueue) pop() int {
	if len(mq.data) == 0 {
		return -1
	}
	elem := mq.data[0]
	mq.data = mq.data[1:]
	if elem == mq.getMax() {
		mq.maxQ = mq.maxQ[1:]
	}
	if elem == mq.getMin() {
		mq.minQ = mq.minQ[1:]
	}
	return elem
}

func (mq *MonotonicQueue) push(elem int) {
	mq.data = append(mq.data, elem)
	for len(mq.maxQ) > 0 && mq.maxQ[len(mq.maxQ)-1] < elem {
		mq.maxQ = mq.maxQ[:len(mq.maxQ)-1]
	}
	mq.maxQ = append(mq.maxQ, elem)

	for len(mq.minQ) > 0 && mq.minQ[len(mq.minQ)-1] > elem {
		mq.minQ = mq.minQ[:len(mq.minQ)-1]
	}
	mq.minQ = append(mq.minQ, elem)
}

// 绝对差不超过限制的最长连续子数组 https://leetcode.cn/problems/longest-continuous-subarray-with-absolute-diff-less-than-or-equal-to-limit/description/
func longestSubarray(nums []int, limit int) int {
	mq := NewMonotonicQueue()
	left, right := 0, 0
	res := 0

	for right < len(nums) {
		mq.push(nums[right])
		right++
		for mq.getMax()-mq.getMin() > limit {
			mq.pop()
			left++
		}
		// 左闭右开，因为right++了
		res = max(res, right-left)
	}
	return res
}

// 和至少为K的最短子数组 https://leetcode.cn/problems/shortest-subarray-with-sum-at-least-k/description/
func shortestSubarray(nums []int, k int) int {
	preSum := make([]int, len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	window := NewMonotonicQueue()
	left, right := 0, 0
	res := math.MaxInt

	for right < len(preSum) {
		window.push(preSum[right])
		right++

		for right < len(preSum) && len(window.data) > 0 &&
			preSum[right]-window.getMin() >= k {
			res = min(res, right-left)
			window.pop()
			left++
		}
	}
	if res == math.MaxInt {
		return -1
	}
	return res
}

// 环形子数组的最大和 https://leetcode.cn/problems/maximum-sum-circular-subarray/
func maxSubarraySumCircular(nums []int) int {
	// 求最大要穷举，遍历前缀和数组，每次把新进元素与
	// 窗口内最小元素求和，取最大
	// 注意：窗口大小为不超过nums长度

	preSum := make([]int, 2*len(nums)+1)
	for i := 1; i < len(preSum); i++ {
		idx := (i - 1) % len(nums)
		preSum[i] = preSum[i-1] + nums[idx]
	}

	window := NewMonotonicQueue()
	window.push(preSum[0])
	res := math.MinInt

	for i := 1; i < len(preSum); i++ {
		res = max(res, preSum[i]-window.getMin())

		// 维护窗口的大小为nums数组的大小
		if len(window.data) == len(nums) {
			window.pop()
		}
		window.push(preSum[i])
	}
	return res
}

func main() {

	fmt.Println(maxSubarraySumCircular([]int{1, -2, 3, -2}))
	fmt.Println(maxSubarraySumCircular([]int{5, -3, 5}))
	fmt.Println(maxSubarraySumCircular([]int{3, -2, 2, -3}))

	//fmt.Println(shortestSubarray([]int{1}, 1))
	//fmt.Println(shortestSubarray([]int{1, 2}, 4))
	//fmt.Println(shortestSubarray([]int{2, -1, 2}, 3))

	//fmt.Println(longestSubarray([]int{8, 2, 4, 7}, 4))
	//fmt.Println(longestSubarray([]int{10, 1, 2, 4, 7, 2}, 5))
	//fmt.Println(longestSubarray([]int{4, 2, 2, 2, 4, 4, 2, 2}, 0))

	//fmt.Println(maxSlidingWindow([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3))
	//fmt.Println(maxSlidingWindow([]int{1, 3, 1, 2, 0, 5}, 3))

	//fmt.Println(largestRectangleArea([]int{2, 1, 5, 6, 2, 3}))
	//fmt.Println(largestRectangleArea([]int{2, 1, 2}))

	//fmt.Println(findUnsortedSubarrayII([]int{2, 6, 4, 8, 10, 9, 15}))
	//fmt.Println(findUnsortedSubarrayII([]int{1, 2, 3, 4, 5}))

	//fmt.Println(carFleet(10, []int{6, 8}, []int{3, 2}))
	//fmt.Println(carFleet(12, []int{10, 8, 0, 5, 3}, []int{2, 4, 1, 1, 3}))

	//fmt.Println(removeKdigits("1432219", 3))
	//fmt.Println(removeKdigits("10200", 1))

	//fmt.Println(finalPrices([]int{10, 1, 1, 6}))
	//fmt.Println(finalPrices([]int{8, 4, 6, 2, 3}))
	//fmt.Println(finalPrices([]int{1, 2, 3, 4, 5}))

	//fmt.Println(canSeePersonsCount([]int{10, 6, 8, 5, 11, 9}))
	//fmt.Println(canSeePersonsCount([]int{5, 1, 2, 3, 10}))

	//fmt.Println(nextLargerNodes(NewMyListNode([]int{2, 1, 5})))
	//fmt.Println(nextLargerNodes(NewMyListNode([]int{2, 7, 4, 3, 5})))

	//fmt.Println(nextGreaterElements([]int{1, 2, 1}))
	//fmt.Println(nextGreaterElements([]int{1, 2, 3, 4, 3}))

	//fmt.Println(dailyTemperatures([]int{73, 74, 75, 71, 69, 72, 76, 73}))
	//fmt.Println(dailyTemperatures([]int{30, 40, 50, 60}))

	//fmt.Println(nextGreaterElement([]int{4, 1, 2}, []int{1, 3, 4, 2}))
	//fmt.Println(nextGreaterElement([]int{2, 4}, []int{1, 2, 3, 4}))

	//fmt.Println(minInsertions("(()))(()))()())))"))
	//fmt.Println(minInsertions("))())("))
	//fmt.Println(minInsertions("(((((("))

	//fmt.Println(minAddToMakeValid("()))(("))
	//fmt.Println(minAddToMakeValid("((("))
	//fmt.Println(minAddToMakeValid("()()()"))

	//fmt.Println(decodeString("3[a2[c]]"))
	//fmt.Println(decodeString("2[abc]3[cd]ef"))

	//fmt.Println(lengthLongestPath("a\n\tb1\n\t\tf1.txt\n\taaaaa\n\t\tf2.txt"))
	//fmt.Println(lengthLongestPath("a"))
	//fmt.Println(lengthLongestPath("dir\n\tsubdir1\n\t\tfile1.ext\n\t\tsubsubdir1\n\tsubdir2\n\t\tsubsubdir2\n\t\t\tfile2.ext"))

	//fmt.Println(evalRPN([]string{"2", "1", "+", "3", "*"}))
	//fmt.Println(evalRPN([]string{"4", "13", "5", "/", "+"}))
	//fmt.Println(evalRPN([]string{"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"}))

	//fmt.Println(isValid("()"))
	//fmt.Println(isValid("()[]{}"))
	//fmt.Println(isValid("([])"))
	//fmt.Println(isValid("([)]"))

	//fmt.Println(simplifyPath("/home/"))
	//fmt.Println(simplifyPath("/home//foo/"))
	//fmt.Println(simplifyPath("/home/user/Documents/../Pictures"))
	//fmt.Println(simplifyPath("/../"))
}
