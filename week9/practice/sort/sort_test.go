package sort_test

import (
	"testing"

	"gitee.com/infraboard/go-course/day9/sort"
	"github.com/stretchr/testify/assert"
)

func TestBubbleSort(t *testing.T) {
	should := assert.New(t)

	raw := []int{3, 6, 4, 2, 11, 10, 5}
	target := sort.BubbleSort(raw)

	should.Equal([]int{2, 3, 4, 5, 6, 10, 11}, target)
}

func TestSelectSort(t *testing.T) {
	should := assert.New(t)

	raw := []int{3, 6, 4, 2, 11, 10, 5}
	target := sort.SelectSort(raw)

	should.Equal([]int{2, 3, 4, 5, 6, 10, 11}, target)
}

func TestBuildInSort(t *testing.T) {
	should := assert.New(t)

	raw := []int{3, 6, 4, 2, 11, 10, 5}
	target := sort.BuildInSort(raw)

	should.Equal([]int{2, 3, 4, 5, 6, 10, 11}, target)
}
