package stack_test

import (
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/week9/practice/data_structure/stack"
	"github.com/stretchr/testify/assert"
)

//
func TestStack(t *testing.T) {
	should := assert.New(t)

	s := stack.NewStack()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)

	// 判断返回 和 预期是否相等
	should.Equal(4, s.Pop())
	should.Equal(3, s.Pop())
	should.Equal(2, s.Pop())
	should.Equal(1, s.Pop())
}

func TestStackRich(t *testing.T) {
	s := stack.NewStackslice()
	s.Push(9)
	s.Push(1)
	s.Push(0)
	s.Push(2)

	t.Log(s.Empty())
	t.Log(s.Peak())
	t.Log(s.Len())
	t.Log(s.Pop())

	s.ForEach(func(v interface{}) {
		fmt.Println(v)
	})
}

func TestStackOrder(t *testing.T) {
	should := assert.New(t)

	s := stack.NewStackslice()
	s.Push(9)
	s.Push(1)
	s.Push(0)
	s.Push(2)

	// 被测试的函数
	s.Sort()

	should.Equal(9, s.Pop())
	should.Equal(2, s.Pop())
	should.Equal(1, s.Pop())
	should.Equal(0, s.Pop())
}

func BenchmarkStackPush(b *testing.B) {
	stack := stack.NewStacklist()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
}

func BenchmarkStackPop(b *testing.B) {
	stack := stack.NewStacklist()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Pop()
	}
}

func BenchmarkMyStackPush(b *testing.B) {
	stack := stack.NewStack()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
}

func BenchmarkMyStackPop(b *testing.B) {
	stack := stack.NewStack()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Pop()
	}
}

//go test -bench=. -benchmem
