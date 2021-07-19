package web

import (
	"fmt"
	"net/http"

	"github.com/ahwhy/myGolang/week6/homework/simple-http-probe/probe"
	"github.com/gin-gonic/gin"
)

func HttpProbe(c *gin.Context) {
	// 解析传过来的host
	host := c.Query("host")
	ishttps := c.Query("is_https")

	// 校验入参
	if host == "" {
		c.String(http.StatusBadRequest, "empty host")
		return
	}

	schema := "http"
	if ishttps == "1" {
		schema = "https"
	}

	url := fmt.Sprintf("%s://%s", schema, host)
	res := probe.DoHttpProbe(url)
	c.String(http.StatusOK, res)
}
