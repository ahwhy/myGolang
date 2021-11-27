package web

import (
	"net/http"

	"github.com/ahwhy/myGolang/historys/week06/example/simple-http-probe/config"
	"github.com/gin-gonic/gin"
)

func StartGin(c *config.Config) {
	// 初始化 gin实例
	r := gin.Default()

	// 测试
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "done",
		})
	})

	// 绑定路由
	Routers(r)

	// 启动 gin实例
	r.Run(c.HttpListenAddr)
}

/*
- 添加路由的函数
- url: /api/probe/http?host=baidu.com&is_https=1
*/
func Routers(r *gin.Engine) {
	// group前缀path -> /api
	api := r.Group("/api")
	api.GET("/probe/http", HttpProbe)
	api.GET("/v1", func(c *gin.Context) {
		c.String(http.StatusOK, "This is api v1.")
	})
}
