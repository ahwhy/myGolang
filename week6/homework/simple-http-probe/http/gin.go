package http

import (
	"net/http"

	"github.com/ahwhy/myGolang/week6/homework/simple-http-probe/config"
	"github.com/gin-gonic/gin"
)

func StartGin(c *config.Config) {
	// 初始化gin 实例
	r := gin.Default()
	// 绑定路由
	Routes(r)
	r.Run(c.HttpListenAddr)

}

// 添加路由的函数
//  /probe/http?host=baidu.com&is_https=1
func Routes(r *gin.Engine) {
	// api group贡献前缀path
	api := r.Group("/api")
	api.GET("/probe/http", HttpProbe)
	api.GET("/v1", func(c *gin.Context) {
		c.String(http.StatusOK, "你好我是 http prober")
	})
}
