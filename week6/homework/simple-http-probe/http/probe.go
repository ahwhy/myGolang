package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/ahwhy/myGolang/week6/homework/simple-http-probe/probe"
)

func HttpProbe(c *gin.Context) {
	// 解析传过来的host
	host := c.Query("host")
	isHttps := c.Query("is_https")
	// validate 校验入参
	if host == "" {
		c.String(http.StatusBadRequest, "empty host")
		return
	}
	schema := "http"
	if isHttps == "1" {
		schema = "https"
	}
	url := fmt.Sprintf("%s://%s", schema, host)
	res := probe.DoHttpProbe(url)
	c.String(http.StatusOK, res)

}
