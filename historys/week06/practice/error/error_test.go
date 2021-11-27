package error_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func validateArgs(name string) (ok bool, err error) {
	if strings.HasPrefix(name, "mysql") {
		return true, nil
	} else {
		return false, errors.New("name must startwith mysql")
	}
}

func Test_ErrorNew(t *testing.T) {
	s1 := "mysql-abc"
	s2 := "redis-abc"
	_, err := validateArgs(s1)
	if err != nil {
		fmt.Println("[s1 validate 失败]", err)
	}
	_, err = validateArgs(s2)
	if err != nil {
		fmt.Println("[s2 validate 失败]", err)
	}
}

func Test_ErrorOS(t *testing.T) {
	file, err := os.Stat("test.txt")
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			log.Printf("PathError")
		case *os.LinkError:
			log.Printf("LinkError")
		case *os.SyscallError:
			log.Printf("SyscallError")
		default:
			log.Printf("unknow error")
		}
	} else {
		fmt.Println(file)
	}
}

type MyError struct {
	err error
	msg string // 自定义的error字符串
}

func (e *MyError) Error() string {
	return e.err.Error() + e.msg
}

func Test_Customize(t *testing.T) {
	err := errors.New("原始的错误 ")
	newErr := MyError{
		err: err,
		msg: "自定义的错误",
	}
	// fmt.Println(newErr.Error())
	fmt.Printf("%T %v", newErr, newErr)
}

func Test_Wrap(t *testing.T) {
	e := errors.New("原始的错误")
	w := fmt.Errorf("Wrap了一个新的错误: %w", e)
	fmt.Println(w)
}
