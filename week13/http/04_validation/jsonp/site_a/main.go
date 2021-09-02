package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/file/*filepath", http.Dir("./static/site_a"))
	http.ListenAndServe(":5656", router)
}

// go run http/validation/jsonp/site_a/main.go
// 浏览器中访问  http://localhost:5656/file/cross_site.html
