package utils

// 计算字符串长度
func CalculateString(line string) int {
	sum := 0
	for _, c := range line {
		sum += int(c)
	}

	return sum
}
