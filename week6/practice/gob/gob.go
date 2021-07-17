package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

type User struct {
	Id   int
	Name string
}

func main() {
	enusers := []User{
		{1, "aa"},
		{2, "bb"},
	}
	// 注册
	gob.Register(User{})
	// 编码
	file, err := os.Create("users.gob")
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder := gob.NewEncoder(file)
	fmt.Println(encoder.Encode(enusers))
	fmt.Println(enusers)
	file.Close()
	// 解码
	file, err = os.Open("users.gob")
	if err != nil {
		return
	}
	decoder := gob.NewDecoder(file)
	var deusers []User
	// deusers := make([]User,0,2)
	fmt.Println(decoder.Decode(&deusers))
	fmt.Println(deusers)
	file.Close()
}
