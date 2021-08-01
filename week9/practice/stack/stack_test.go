package stack_test

import (
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/week9/practice/stack"
	"github.com/stretchr/testify/assert"
)

//  1
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
	s := stack.NewStack()
	s.Push(9)
	s.Push(1)
	s.Push(0)
	s.Push(2)

	t.Log(s.IsEmpty())
	t.Log(s.Peek())
	t.Log(s.Len())
	t.Log(s.Pop())

	s.ForEach(func(item stack.Item) {
		fmt.Println(item)
	})
}

func TestStackOrder(t *testing.T) {
	should := assert.New(t)

	s := stack.NewStack()
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
