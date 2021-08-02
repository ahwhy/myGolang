package list_test

import (
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/week9/practice/data_structure/list"
)

func TestListBasic(t *testing.T) {
	l := list.NewIntList(1)
	l.AddNode(list.NewIntNode(2))
	l.AddNode(list.NewIntNode(3))
	l.AddNode(list.NewIntNode(4))
	l.Traverse(func(n *list.Node) {
		if n.Next != nil {
			fmt.Printf("%d --> ", n.Value)
		} else {
			fmt.Printf("%d", n.Value)
		}
	})
}

func PrintNode(n *list.Node) {
	if n.Next != nil {
		fmt.Printf("%d --> ", n.Value)
	} else {
		fmt.Printf("%d", n.Value)
	}
}

func TestListRich(t *testing.T) {
	l := list.NewIntList(1)
	n2, n3, n4 := list.NewIntNode(2), list.NewIntNode(3), list.NewIntNode(4)
	l.AddNode(n2)
	l.AddNode(n3)
	l.AddNode(n4)
	l.Traverse(PrintNode)

	// 测试插入
	l.InsertAfter(n2, list.NewIntNode(20))
	l.Traverse(PrintNode)
}

func TestListWithPre(t *testing.T) {
	l := list.NewIntList(1)
	n2, n3, n4 := list.NewIntNode(2), list.NewIntNode(3), list.NewIntNode(4)
	l.AddNode(n2)
	l.AddNode(n3)
	l.AddNode(n4)
	l.Traverse(PrintNode)

	// 测试插入
	l.InsertBefore(n2, list.NewIntNode(20))
	l.Traverse(PrintNode)

	// 测试删除
	l.Remove(n3)
	l.Traverse(PrintNode)
}

func TestListRing(t *testing.T) {
	l := list.NewIntList(1)
	n2, n3, n4 := list.NewIntNode(2), list.NewIntNode(3), list.NewIntNode(4)
	l.AddNode(n2)
	l.AddNode(n3)
	l.AddNode(n4)
	l.Traverse(PrintNode)

	// 测试插入
	l.ChangeToRing()
	l.Traverse(PrintNode)
	fmt.Println()
	l.InsertAfter(n3, list.NewIntNode(100))
	l.Traverse(PrintNode)
}
