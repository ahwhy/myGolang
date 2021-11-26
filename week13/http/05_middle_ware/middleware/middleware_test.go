package middleware_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ahwhy/myGolang/week22/middleware"
)

func hello(wr http.ResponseWriter, r *http.Request) {
	wr.Write([]byte("hello"))
}

func TestHello(t *testing.T) {
	http.HandleFunc("/", hello)

	http.Handle("/", middleware.TimeMiddleware(http.HandlerFunc(hello)))
	// HandlerFunc 是一个类型, 把  hello这个函数，转换为 HandlerFunc类型, 使用int(a)是一个语法
	// HandlerFunc 实现了ServeHTTP方法, 这样的hello函数对象也就有了ServeHTTP方法
	// HandlerFunc 是个函数， 把他定义为一个Type, 然后给这个Type 绑定了一个函数
	// 通过这样的操作凡是 被转换为HandlerFunc类型的函数，都是一个http.Handler
	// 要完成 Type()的转换，函数签名必须一致

	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}

func TestMiddleware(t *testing.T) {
	r := middleware.NewRouter()
	r.Use(middleware.LogMiddleware)
	r.Use(middleware.TimeMiddleware)
	http.Handle("/", r.Merge(http.HandlerFunc(hello)))
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}
