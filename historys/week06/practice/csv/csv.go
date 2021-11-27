package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type User struct {
	Id   int
	Name string
}

func main() {
	wusers := []User{
		{1, "aa"},
		{2, "bb"},
	}
	// 写入
	file, err := os.Create("users.csv")
	if err != nil {
		return
	}
	writer := csv.NewWriter(file)
	for _, user := range wusers {
		writer.Write([]string{strconv.Itoa(user.Id), user.Name})
	}
	writer.Flush()
	file.Close()
	// 读取
	file, err = os.Open("users.csv")
	if err != nil {
		return
	}
	var rusers []User
	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		id, _ := strconv.Atoi(line[0])
		rusers = append(rusers, User{id, line[1]})
	}
	fmt.Println(rusers)
}
