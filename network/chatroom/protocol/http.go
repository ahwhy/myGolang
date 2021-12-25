package protocol

import (
	"fmt"
	"net/http"

	"github.com/ahwhy/myGolang/network/chatroom/module"
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
	http.ServeFile(w, r, "../browser/home.html")
}

func Run(port *string) {
	hub := module.NewHub()
	go hub.Run()
	//注册每种请求对应的处理函数
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		ServeWs(hub, rw, r)
	})
	if err := http.ListenAndServe(":"+*port, nil); err != nil { // 如果启动成功，该行会一直阻塞，hub.run()会一直运行
		fmt.Printf("start http service error: %s\n", err)
	}
}
