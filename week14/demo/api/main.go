package main

import (
	"github.com/ahwhy/myGolang/week14/demo/api/cmd"
)

func main() {
	cmd.Execute()
}

/*
go run main.go -f etc/demo.toml start
go build -ldflags "-s -w" -o demo-api main.go
go build -o demo-api -ldflags "-X gitee.com/infraboard/go-course/day14/demo/api/version.GIT_TAG='v0.0.1'" main.go
*/
