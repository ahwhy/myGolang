package sort_test

import (
	"math/rand"
	"testing"

	"github.com/ahwhy/myGolang/week9/practice/sort"
	"github.com/ahwhy/myGolang/week9/practice/stack"
)

const (
	MAX_RAND_LIMIT = 10000
)

// 准备数据
func generateRandomArray(arrayLen int) []int {
	var a []int
	for i := 0; i < arrayLen; i++ {
		a = append(a, rand.Intn(MAX_RAND_LIMIT))
	}
	return a
}

func benchmarkBubbleSort(i int, b *testing.B) {
	a := generateRandomArray(i)
	sort.BubbleSort(a)
}

func BenchmarkBubbleSort100(b *testing.B) {
	benchmarkBubbleSort(100, b)
}

func BenchmarkBubbleSort1000(b *testing.B) {
	benchmarkBubbleSort(1000, b)
}

func BenchmarkBubbleSort10000(b *testing.B) {
	benchmarkBubbleSort(10000, b)
}

func benchmarkSelectSort(i int, b *testing.B) {
	a := generateRandomArray(i)
	sort.SelectSort(a)
}

func BenchmarkSelectSort100(b *testing.B) {
	benchmarkSelectSort(100, b)
}

func BenchmarkSelectSort1000(b *testing.B) {
	benchmarkSelectSort(1000, b)
}

func BenchmarkSelectSort10000(b *testing.B) {
	benchmarkSelectSort(10000, b)
}

func benchmarkInsertSort(i int, b *testing.B) {
	a := generateRandomArray(i)
	s := stack.NewNumberStack(a)
	s.Sort()
}

func BenchmarkInsertSort100(b *testing.B) {
	benchmarkInsertSort(100, b)
}

func BenchmarkInsertSort1000(b *testing.B) {
	benchmarkInsertSort(1000, b)
}

func BenchmarkInsertSort10000(b *testing.B) {
	benchmarkInsertSort(10000, b)
}

func benchmarkBuildInSort(i int, b *testing.B) {
	a := generateRandomArray(i)
	sort.BuildInSort(a)
}

func BenchmarkBuildInSort100(b *testing.B) {
	benchmarkBuildInSort(100, b)
}

func BenchmarkBuildInSort1000(b *testing.B) {
	benchmarkBuildInSort(1000, b)
}

func BenchmarkBuildInSort10000(b *testing.B) {
	benchmarkBuildInSort(10000, b)
}
