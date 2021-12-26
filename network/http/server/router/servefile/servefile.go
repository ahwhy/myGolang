package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/file/*filepath", http.Dir("./static"))
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("yyds({'message':'Chinese'})")) // 返回一段json回调函数
	})
	http.ListenAndServe(":5657", router)
}

// go build -o ./servefile servefile.go
// 浏览器中访问  http://localhost:5656/file/xss.html
