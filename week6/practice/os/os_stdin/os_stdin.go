package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	/*
		echo "ss -ntlp " > a.txt
		go run a.go < a.txt
	*/
	cmd := exec.Command("sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("run.err", err)
		return
	}
}