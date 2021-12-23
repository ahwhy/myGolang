package main

import (
	"flag"
	"fmt"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { //只允许访问根路径
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" { //只允许GET请求
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	port := flag.String("port", "5657", "http service port") //如果命令行不指定port参数，则默认为5657
	flag.Parse()                                             //解析命令行输入的port参数
	hub := NewHub()
	go hub.Run()
	//注册每种请求对应的处理函数
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		ServeWs(hub, rw, r)
	})
	if err := http.ListenAndServe(":"+*port, nil); err != nil { //如果启动成功，该行会一直阻塞，hub.run()会一直运行
		fmt.Printf("start http service error: %s\n", err)
	}
}

//go run chat_room/*.go --port 5657
