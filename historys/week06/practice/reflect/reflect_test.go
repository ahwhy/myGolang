package main

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func Test_Interface(test *testing.T) {
	var s interface{} = true
	switch s.(type) {
	case string:
		fmt.Println("s.type=string")
	case int:
		fmt.Println("s.type=int")
	case bool:
		fmt.Println("s.type=bool")
	default:
		fmt.Println("未知的类型")
	}
}

func Test_Reflect1(test *testing.T) {
	var s interface{} = 12345.87987
	// TypeOf会返回模板的对象
	reflectType := reflect.TypeOf(s)
	reflectValue := reflect.ValueOf(s)

	log.Printf("[typeof:%v]", reflectType)
	log.Printf("[valueof:%v]", reflectValue)
}