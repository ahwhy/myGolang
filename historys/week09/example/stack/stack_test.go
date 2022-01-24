package stack_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/ahwhy/myGolang/historys/week09/example/stack"
	"github.com/stretchr/testify/assert"
)

var StackMap = make(map[string]stack.Stool)

var (
	Stacklist  = stack.NewStacklist()
	Stack      = stack.NewStack()
	Stackslice = stack.NewStackslice()
)

func init() {
	StackMap["Stacklist"] = Stacklist
	StackMap["Stack"] = Stack
	StackMap["Stackslice"] = Stackslice
}

func TestAllStack(t *testing.T) {
	should := assert.New(t)

	for n, s := range StackMap {
		fmt.Printf("This is %s\n", n)
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
}

func TestAllStackRich(t *testing.T) {
	for n, s := range StackMap {
		fmt.Printf("This is %s\n", n)
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
}

func TestAllStackOrder(t *testing.T) {
	for n, s := range StackMap {
		fmt.Printf("This is %s\n", n)

		for i := 0; i < 10; i++ {
			s.Push(i)
		}

		// 被测试的函数
		s.Sort()

		s.ForEach(func(v interface{}) {
			fmt.Println(v)
		})
	}
}

const (
	MAX_RAND_LIMIT = 10000
)

func TestAllStackGo(t *testing.T) {
	group := &sync.WaitGroup{}
	group.Add(3)
	go func() {
		for _, s := range StackMap {
			for i := 0; i < 100; i++ {
				num := rand.Intn(MAX_RAND_LIMIT)
				s.Push(num)
			}
		}
		group.Done()
	}()

	go func() {
		for _, s := range StackMap {
			for i := 0; i < 100; i++ {
				num := rand.Intn(MAX_RAND_LIMIT)
				s.Push(num)
			}
		}
		group.Done()
	}()

	go func() {
		for _, s := range StackMap {
			for i := 0; i < 100; i++ {
				num := rand.Intn(MAX_RAND_LIMIT)
				s.Push(num)
			}
		}
		group.Done()
	}()
	group.Wait()

	for n, s := range StackMap {
		fmt.Printf("This is %s\n", n)
		s.Sort()

		s.ForEach(func(v interface{}) {
			fmt.Println(v)
		})
	}
}
