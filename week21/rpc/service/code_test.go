package service_test

import (
	"fmt"
	"testing"

	"gitee.com/infraboard/go-course/day21/rpc/service"
)

type TestStruct struct {
	Name  string
	Value string
}

func TestGobCode(t *testing.T) {
	t1 := &TestStruct{"name", "value"}
	resp, err := service.GobEncode(t1)
	fmt.Println(resp, err)

	t2 := &TestStruct{}
	service.GobDecode(resp, t2)
	fmt.Println(t2, err)
}
