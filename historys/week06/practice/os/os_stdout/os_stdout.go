package main

import (
	"fmt"
	"os"
)

func main() {
	os.Stdout.Write([]byte("aa,bb"))
	fmt.Println("aa,bb")
}
