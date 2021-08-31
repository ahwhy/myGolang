package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func handle(method string, w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Printf("request method: %s\n", r.Method)
	fmt.Printf("request body: ")
	io.Copy(os.Stdout, r.Body) // 把r.Body流里的内容拷贝到os.Stdout流里
	fmt.Println()
	w.Write([]byte("Hi boy, you request " + method))
}

func get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handle("get", w, r, params)
}

func post(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handle("post", w, r, params)
}

func head(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handle("head", w, r, params)
}

func options(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handle("options", w, r, params)
}

func put(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handle("put", w, r, params)
}

func patch(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handle("patch", w, r, params)
}

func delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handle("delete", w, r, params)
}

func panic(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var arr []int
	_ = arr[1] // 数组越界panic
}

func main2() {
	router := httprouter.New()
	router.GET("/", get)
	router.POST("/", post)
	router.HEAD("/", head)
	router.OPTIONS("/", options)
	router.PUT("/", put)
	router.PATCH("/", patch)
	router.DELETE("/", delete)
	// router没有提供CONNECT和TRACE

	// *只能有一个，且必须放path的末尾；catch-all routes are only allowed at the end of the path
	router.POST("/user/:name/:type/*addr", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Printf("name:%s, type:%s, addr:%s\n", p.ByName("name"), p.ByName("type"), p.ByName("addr"))
	})

	// 必须以/*filepath结尾，因为要获取访问的路径信息
	// 在浏览器中访问：http://127.0.0.1:5656/file/home.html
	// 或 http://127.0.0.1:5656/file/readme.md
	router.ServeFiles("/file/*filepath", http.Dir("./"))

	// 通过recover捕获panic
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		w.WriteHeader(http.StatusInternalServerError) // 设置response status
		fmt.Fprintf(w, "error:%s", err)               // 线上环境不要把原始错误信息返回给前端
	}
	router.GET("/panic", panic)

	// router实现了ServerHTTP接口，所以它是一种http.Handler
	http.ListenAndServe(":5656", router)
}
