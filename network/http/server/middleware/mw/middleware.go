package mw

import (
	"log"
	"net/http"
	"time"
)

func TimeMiddleWare(next http.Handler) http.Handler {
	// 通过HandlerFunc把一个func(rw http.ResponseWriter, r *http.Request)函数转为Handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		// next handler
		next.ServeHTTP(rw, r)
		timeElapsed := time.Since(timeStart)
		log.Printf("request %s use %d ms\n", r.URL.Path, timeElapsed.Milliseconds())
	})
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		log.Println("start")
		// next handler
		next.ServeHTTP(wr, r)
		log.Println("end")
	})
}

var limitCh = make(chan struct{}, 100) // 最多并发处理100个请求

func LimitMiddleWare(next http.Handler) http.Handler {
	// 通过HandlerFunc返回一个Handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		limitCh <- struct{}{} // 并发度达到100时就会阻塞
		log.Printf("concurrence %d\n", len(limitCh))
		// next handler
		next.ServeHTTP(rw, r)
		<-limitCh
	})
}
