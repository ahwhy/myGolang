package main

import (
	"fmt"
)

func main() {
	nums := [...]int{
		0b1000100001100011,
		0b0110011100001101,
		0b0101010110011100,
		0b0110101100100010,
		0b0111101001111111,
		0b0100111000101101,
		0b0101011011111101,
		0b0111111010100010,
	}

	for i := 0; i < len(nums); i++ {
		fmt.Printf("二进制: %016b; Unicode: %U; 字符: %c.\n", nums[i], nums[i], nums[i])
	}
}
