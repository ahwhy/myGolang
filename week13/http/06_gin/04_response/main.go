package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func text(engine *gin.Engine) {
	engine.GET("/user/text", func(c *gin.Context) {
		c.String(http.StatusOK, "hi boy") // response Content-Type:text/plain
	})
}

func json1(engine *gin.Engine) {
	engine.GET("/user/json1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"name": "zcy", "addr": "bj"}) // response Content-Type:application/json
	})
}

func json2(engine *gin.Engine) {
	var stu struct { // 匿名结构体
		Name string
		Addr string
	}
	stu.Name = "zcy"
	stu.Addr = "bj"
	engine.GET("/user/json2", func(c *gin.Context) {
		c.JSON(http.StatusOK, stu) // response Content-Type:application/json
	})
}

func jsonp(engine *gin.Engine) {
	var stu struct {
		Name string
		Addr string
	}
	stu.Name = "zcy"
	stu.Addr = "bj"
	engine.GET("/user/jsonp", func(ctx *gin.Context) {
		// 如果请求参数里有callback=xxx，则response Content-Type为application/javascript，否则response Content-Type为application/json
		ctx.JSONP(http.StatusOK, stu)
	})
}

func xml(engine *gin.Engine) {
	var stu struct {
		Name string
		Addr string
	}
	stu.Name = "zcy"
	stu.Addr = "bj"
	engine.GET("/user/xml", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"name": "zcy", "addr": "bj"}) // response Content-Type:application/xml
	})
}

func yaml(engine *gin.Engine) {
	var stu struct {
		Name string
		Addr string
	}
	stu.Name = "zcy"
	stu.Addr = "bj"
	engine.GET("/user/yaml", func(c *gin.Context) {
		c.YAML(http.StatusOK, stu)
	})
}

func html(engine *gin.Engine) {
	engine.LoadHTMLFiles("static/template.html")
	engine.GET("/user/html", func(c *gin.Context) {
		// 通过json往前端页面上传值
		c.HTML(http.StatusOK, "template.html", gin.H{"title": "用户信息", "name": "zcy", "addr": "bj"})
	})
}

func redirect(engine *gin.Engine) {
	engine.GET("/not_exists", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://localhost:5656/user/html")
	})
}

func main() {
	engine := gin.Default()
	text(engine)     // http://localhost:5656/user/text
	json1(engine)    // http://localhost:5656/user/json1
	json2(engine)    // http://localhost:5656/user/json2
	jsonp(engine)    // http://localhost:5656/user/jsonp?callback=yyds
	xml(engine)      // http://localhost:5656/user/xml
	yaml(engine)     // http://localhost:5656/user/yaml
	html(engine)     // http://localhost:5656/user/html
	redirect(engine) // http://localhost:5656/not_exists
	engine.Run(":5656")
}
