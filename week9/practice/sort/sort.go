package sort

import (
	"sort"
)

func BubbleSort(numbers []int) []int {
	for range numbers {
		for j := 0; j < len(numbers)-1; j++ {
			// 当前值 numbers[j], 后一个值是多少 numbers[j+1]
			// fmt.Printf("数据: 当前: %d, 比对: %d\n", numbers[j], numbers[j+1])

			// 比较大小，交互位置, 大数底部
			if numbers[j] > numbers[j+1] {
				numbers[j], numbers[j+1] = numbers[j+1], numbers[j]
			}
		}

		// fmt.Printf("第%d趟冒泡: %v\n", i+1, numbers)
	}
	return numbers
}

func SelectSort(numbers []int) []int {
	for i := range numbers {
		for j := i + 1; j < len(numbers); j++ {
			// fmt.Printf("数据 --> 当前数据: %d, 比对的数据: %d\n", numbers[i], numbers[j])
			// 交换位置
			if numbers[i] > numbers[j] {
				numbers[i], numbers[j] = numbers[j], numbers[i]
				// fmt.Printf("交换 --> 当前数据: %d, 比对数据: %d\n", numbers[i], numbers[j])
			}
		}

		// fmt.Printf("%d趟: %v\n", i+1, numbers)
	}

	return numbers
}

type IntSlice []int

func (s IntSlice) Len() int { return len(s) }

func (s IntSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s IntSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func BuildInSort(number []int) []int {
	sort.Sort(IntSlice(number))
	return number
}
