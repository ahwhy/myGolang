package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var limitCh = make(chan struct{}, 100) // 最多并发处理100个请求

func timeMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		begin := time.Now()
		ctx.Next() // 执行业务逻辑
		timeElapsed := time.Since(begin)
		log.Printf("request %s use %d ms\n", ctx.Request.URL.Path, timeElapsed.Milliseconds())
	}
}

func limitMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitCh <- struct{}{} // 并发度达到100时就会阻塞
		log.Printf("concurrence %d\n", len(limitCh))
		ctx.Next() // 执行业务逻辑
		<-limitCh
	}
}

func main() {
	engine := gin.Default()
	engine.Use(timeMiddleWare()) // 全局MiddleWare
	engine.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hi boy")
	})
	engine.GET("/girl", limitMiddleWare(), func(ctx *gin.Context) { // 局部MiddleWare
		ctx.String(http.StatusOK, "hi girl")
	})
	engine.Run(":5656")
}
