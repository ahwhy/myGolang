package main

import (
	"log"
	"reflect"
)

func main() {
	var num float64 = 3.14
	log.Printf("[num原始值:%f]", num)

	// 通过reflect.ValueOf获取num中的value
	// 必须是指针才可以修改值
	pointer := reflect.ValueOf(&num)
	newValue := pointer.Elem()
	// 赋值
	newValue.SetFloat(5.6)
	log.Printf("[num新值:%f]", num)

	pointer = reflect.ValueOf(num)
	// reflect: call of reflect.Value.Elem on float64 Value
	// newValue = pointer.Elem()
}
