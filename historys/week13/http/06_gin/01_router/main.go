package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func boy(c *gin.Context) { // 使用者所需要的东西全都封装在了gin.Context里面，包括http.Request和ResponseWriter
	c.String(http.StatusOK, "hi boy") // 通过gin.Context.String返回一个text/plain类型的正文
}

func girl(c *gin.Context) {
	c.String(http.StatusOK, "hi girl")
}

func main() {
	// gin.SetMode(gin.ReleaseMode)        // 发布模式，默认是Debug模式

	engine := gin.Default() // 默认的engine已自带了Logger和Recovery两个中间件
	engine.GET("/", boy)
	engine.POST("/", girl)

	//路由分组
	oldVersion := engine.Group("/v1")
	oldVersion.GET("/student", boy) // http://localhost:5656/v1/student
	oldVersion.GET("/teacher", boy) // http://localhost:5656/v1/teacher

	newVersion := engine.Group("/v2")
	newVersion.GET("/student", girl) // http://localhost:5656/v2/student
	newVersion.GET("/teacher", girl) // http://localhost:5656/v2/teacher

	engine.Run(":5656")
}
