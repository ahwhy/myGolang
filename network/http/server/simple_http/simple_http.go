package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// http协议具体内容
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request method: %s\n", r.Method)
	fmt.Printf("request host: %s\n", r.Host)
	fmt.Printf("request url: %s\n", r.URL)
	fmt.Printf("request proto: %s\n", r.Proto)

	fmt.Println("request header")
	for key, values := range r.Header {
		fmt.Printf("%s: %v\n", key, values)
	}
	fmt.Println()

	fmt.Println("request cookie")
	for _, cookie := range r.Cookies() {
		fmt.Printf("name=%s vaue=%s\n", cookie.Name, cookie.Value)
	}
	fmt.Println()

	fmt.Printf("request body: ")
	io.Copy(os.Stdout, r.Body) // 把r.Body流里的内容拷贝到os.Stdout流里
	fmt.Println()

	fmt.Fprint(w, "Hello Boy") // 把返回的内容写入http.ResponseWriter
	fmt.Printf("\n==========\n")
}

// TrunkedHandler 模拟http trunked过程
func TrunkedHandler(w http.ResponseWriter, r *http.Request) {
	flusher := w.(http.Flusher)
	// 两段消息，会自动触发 http trunked
	for i := 0; i < 2; i++ {
		fmt.Fprint(w, "This is Trunked\n")
		flusher.Flush()
		<-time.Tick(1 * time.Second)
	}
}

func main() {
	http.HandleFunc("/", HelloHandler) // 路由，请求要目录时去执行HelloHandler
	http.HandleFunc("/trunked", TrunkedHandler)
	if err := http.ListenAndServe(":5656", nil); err != nil { // ListenAndServe如果不发生error会一直阻塞；为每一个请求单独创建一个协程去处理
		panic(err)
	}
}
