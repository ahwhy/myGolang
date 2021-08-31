package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HostMap map[string]http.Handler

// 作为http.Handler必须实现ServeHTTP接口；HostMap首先是个map，其次它还具有了ServeHTTP的功能
func (hm HostMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, exists := hm[r.Host]; exists {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}

func main() {
	bookRouter := httprouter.New()
	bookRouter.POST("/", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		rw.Write([]byte("read book"))
	})
	
	foodRouter := httprouter.New()
	foodRouter.POST("/", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		rw.Write([]byte("eat food"))
	})

	/**
	- 不同的二级域名，对应不同的Router
	- 需要在/etc/hosts里加入两行	
		127.0.0.1 book.dianshang
		127.0.0.1 food.dianshang
	*/
	hm := make(HostMap)
	hm["book.dianshang:5656"] = bookRouter
	hm["food.dianshang:5656"] = foodRouter
	http.ListenAndServe(":5656", hm)
}
