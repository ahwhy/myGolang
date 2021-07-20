package main

import (
	"math/rand"
	"testing"
	"time"
)

// 制定大的cap的切片
func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

// 动态扩容的slice
func generateDynamic(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func BenchmarkGenerateWithCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generateWithCap(100000)
	}
}
func BenchmarkGenerateDynamic(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generateDynamic(100000)
	}
}

func benchmarkGenerate(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		generateDynamic(i)
	}
}

func BenchmarkGenerateDynamic1000(b *testing.B)     { benchmarkGenerate(1000, b) }
func BenchmarkGenerateDynamic10000(b *testing.B)    { benchmarkGenerate(10000, b) }
func BenchmarkGenerateDynamic100000(b *testing.B)   { benchmarkGenerate(100000, b) }
func BenchmarkGenerateDynamic1000000(b *testing.B)  { benchmarkGenerate(1000000, b) }
func BenchmarkGenerateDynamic10000000(b *testing.B) { benchmarkGenerate(10000000, b) }