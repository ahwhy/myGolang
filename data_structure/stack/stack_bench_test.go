package stack_test

import (
	"testing"

	"github.com/ahwhy/myGolang/data_structure/stack"
)

func BenchmarkStacklistPush(b *testing.B) {
	stack := stack.NewStacklist()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
}

func BenchmarkStacklistPop(b *testing.B) {
	stack := stack.NewStacklist()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Pop()
	}
}

func BenchmarkStackPush(b *testing.B) {
	stack := stack.NewStack()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
}

func BenchmarkStackPop(b *testing.B) {
	stack := stack.NewStack()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Pop()
	}
}

func BenchmarkStackslicePush(b *testing.B) {
	stack := stack.NewStackslice()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
}

func BenchmarkStackslicePop(b *testing.B) {
	stack := stack.NewStackslice()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Pop()
	}
}

//go test -bench=. -benchmem
