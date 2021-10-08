package main

import (
	"fmt"
)

func main() {
	var nums []int = []int{101, 103, 107, 102, 106}
	nums = append(nums, 123, 12, 1244, 345, 363455)

	bubble_Sort(nums...)
	fmt.Println(nums)
}

func bubble_Sort(nums ...int) {
	for i := range nums {
		for j := i; j < len(nums); j++ {
		// for j, _ := range nums{
			if nums[i] < nums[j] {
				nums[i], nums[j] = nums[j], nums[i]
			}
		}
	}
}
