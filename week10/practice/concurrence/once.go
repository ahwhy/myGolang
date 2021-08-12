package main

import (
	"fmt"
	"sync"
	"time"
)

var globalResource map[string]string = make(map[string]string)
var loadResourceOnce sync.Once

func LoadResource() {
	loadResourceOnce.Do(func() {
		fmt.Println("load global resource")
		globalResource["1"] = "A"
		globalResource["2"] = "B"
	})
}

type Singleton struct {
	Name string
}

var singleton *Singleton
var singletonOnce sync.Once

func GetSingletonInstance() *Singleton {
	singletonOnce.Do(func() {
		fmt.Println("init Singleton")
		singleton = &Singleton{Name: "Tom"}
	})
	return singleton
}

func main() {
	go LoadResource()
	go LoadResource()
	inst1 := GetSingletonInstance()
	inst2 := GetSingletonInstance()
	fmt.Printf("inst1 address %v\n", []*Singleton{inst1})
	fmt.Printf("inst2 address %v\n", []*Singleton{inst2})
	time.Sleep(100 * time.Millisecond)
}
