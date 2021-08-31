package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var limitCh = make(chan struct{}, 100) // 最多并发处理100个请求

func getBoy(w http.ResponseWriter, r *http.Request) {
	time.Sleep(150 * time.Millisecond)
	w.Write([]byte("hi boy"))
}

func getGirl(w http.ResponseWriter, r *http.Request) {
	time.Sleep(150 * time.Millisecond)
	w.Write([]byte("hi girl"))
}

func timeMiddleWare(next http.Handler) http.Handler {
	// 通过HandlerFunc把一个func(rw http.ResponseWriter, r *http.Request)函数转为Handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		next.ServeHTTP(rw, r)
		timeElapsed := time.Since(begin)
		log.Printf("request %s use %d ms\n", r.URL.Path, timeElapsed.Milliseconds())
	})
}

func limitMiddleWare(next http.Handler) http.Handler {
	// 通过HandlerFunc返回一个Handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		limitCh <- struct{}{} // 并发度达到100时就会阻塞
		log.Printf("concurrence %d\n", len(limitCh))
		next.ServeHTTP(rw, r)
		<-limitCh
	})
}

/**
以下演示更优雅的中间件组织形式
*/
type middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []middleware
	mux             map[string]http.Handler // mux通常表示路由策略
}

func NewRouter() *Router {
	return &Router{
		middlewareChain: make([]middleware, 0, 10),
		mux:             make(map[string]http.Handler, 10),
	}
}

func (self *Router) Use(m middleware) {
	self.middlewareChain = append(self.middlewareChain, m)
}

func (self *Router) Add(path string, handler http.Handler) {
	var mergedHandler = handler
	for i := len(self.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = self.middlewareChain[i](mergedHandler) // 中间件层层嵌套
	}
	self.mux[path] = mergedHandler
}

func main() {
	// http.Handle("/", timeMiddleWare(limitMiddleWare(http.HandlerFunc(getBoy))))      // 中间层层嵌套
	// http.Handle("/home", timeMiddleWare(limitMiddleWare(http.HandlerFunc(getGirl)))) // 跟上面一行存在重复代码

	router := NewRouter()
	router.Use(limitMiddleWare)
	router.Use(timeMiddleWare)
	// 以下演示了2个路径(还可以更多)，每个路径都使用相同的middlewareChain
	router.Add("/", http.HandlerFunc(getBoy))
	router.Add("/home", http.HandlerFunc(getGirl))
	for path, handler := range router.mux {
		http.Handle(path, handler)
	}

	if err := http.ListenAndServe(":5656", nil); err != nil {
		fmt.Println(err)
	}
}
