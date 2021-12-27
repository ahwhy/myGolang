package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ahwhy/myGolang/network/http/server/middleware/mw"
)

func hello(wr http.ResponseWriter, r *http.Request) {
	time.Sleep(150 * time.Millisecond)
	wr.Write([]byte("hello"))
}

func getBoy(w http.ResponseWriter, r *http.Request) {
	time.Sleep(150 * time.Millisecond)
	w.Write([]byte("hi boy"))
}

func getGirl(w http.ResponseWriter, r *http.Request) {
	time.Sleep(150 * time.Millisecond)
	w.Write([]byte("hi girl"))
}

func main() {
	// 初始的中间件
	// HandlerFunc 是一个类型, 把  hello这个函数，转换为 HandlerFunc类型, 同int(a)是一个语法
	// HandlerFunc 实现了ServeHTTP方法, 这样的hello函数对象也就有了ServeHTTP方法
	// HandlerFunc 是个函数， 把他定义为一个Type, 然后给这个Type 绑定了一个函数
	// 通过这样的操作凡是 被转换为HandlerFunc类型的函数，都是一个http.Handler
	// 要完成 Type()的转换，函数签名必须一致
	http.Handle("/init", mw.TimeMiddleWare(mw.LogMiddleware(mw.LimitMiddleWare(http.HandlerFunc(hello)))))

	// 不包含限流的中间件
	r1 := mw.NewRouter()
	r1.Use(mw.TimeMiddleWare)
	r1.Use(mw.LogMiddleware)

	http.Handle("/", r1.Merge(http.HandlerFunc(hello)))

	// 包含限流的中间件
	r2 := mw.NewRouter()
	r2.Use(mw.TimeMiddleWare)
	r2.Use(mw.LogMiddleware)
	r2.Use(mw.LimitMiddleWare)
	r2.Add("/boy", http.HandlerFunc(getBoy))
	r2.Add("/girl", http.HandlerFunc(getGirl))
	for path, handler := range r2.Mux {
		http.Handle(path, handler)
	}

	if err := http.ListenAndServe(":5656", nil); err != nil {
		log.Println(err)
	}
}
