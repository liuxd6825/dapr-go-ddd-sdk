package intutils

import "strings"

// base36 定义了目标进制的字符集。
// 它的长度就是进制的基数（这里是36）。
const base36 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

/*
3位	ZZZ	        46,655
4位	ZZZZ	 1,679,615
5位	ZZZZZ	60,466,175
*/

// IntToBase36
// @Description: 将一个10进制转换为36进制字符串。
// @param n
// @return string
func IntToBase36(n int64) string {
	// 0 是一个特殊情况，直接返回字符集中的第一个字符 '0'。
	if n == 0 {
		return string(base36[0])
	}

	// 获取基数，即字符集的长度。
	base := int64(len(base36))
	var isNegative bool

	// 处理负数：记录符号，然后对绝对值进行计算。
	if n < 0 {
		isNegative = true
		n = -n
	}

	// 使用 strings.Builder 来高效地构建字符串，避免在循环中产生大量字符串副本。
	var sb strings.Builder

	// 核心算法：除基取余
	for n > 0 {
		// 计算余数
		remainder := n % base
		// 将余数对应的字符写入 builder
		sb.WriteByte(base36[remainder])
		// 更新 n 为商
		n = n / base
	}

	// 如果原始数字是负数，在反转后的字符串前添加负号。
	if isNegative {
		sb.WriteString("-")
	}

	// 因为我们得到的是反向的字符串，所以需要将其反转。
	return reverseString(sb.String())
}

// reverseString 是一个辅助函数，用于反转字符串。
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
