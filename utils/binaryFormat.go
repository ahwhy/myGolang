package utils

import (
	"math"
	"strings"
)

// BinaryFormat 输出一个int32对应的二进制表示
func BinaryFormat(n int32) string {
	a := uint32(n)
	sb := strings.Builder{}
	c := uint32(math.Pow(2, 31)) // 最高位上是1，其他位全是0  10000000000000000000000000000000

	for i := 0; i < 32; i++ {
		if a&c != 0 { // 判断n的当前位上是否为1
			sb.WriteString("1")
		} else {
			sb.WriteString("0")
		}

		c >>= 1 // "1"往右移一位
	}

	return sb.String()
}
