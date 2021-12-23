package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/file/*filepath", http.Dir("./static/site_b"))
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("yyds({'message':'Chinese'})")) //返回一段json回调函数
	})
	http.ListenAndServe(":5657", router)
}

// go run http/validation/jsonp/site_b/main.go
// 浏览器中访问  http://localhost:5657/file/test.js ，js文件被当成普通的文本文件
