package mw

import "net/http"

type middleware func(http.Handler) http.Handler

func NewRouter() *Router {
	return &Router{
		MiddlewareChain: make([]middleware, 0, 10),
		Mux:             make(map[string]http.Handler, 10),
	}
}

type Router struct {
	MiddlewareChain []middleware
	Mux             map[string]http.Handler // mux通常表示路由策略
}

func (r *Router) Use(m middleware) {
	r.MiddlewareChain = append(r.MiddlewareChain, m)
}

func (r *Router) Merge(hander http.Handler) http.Handler {
	var mergedHandler = hander
	for i := len(r.MiddlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.MiddlewareChain[i](mergedHandler) // 中间件层层嵌套
	}

	return mergedHandler
}

func (r *Router) Add(path string, hander http.Handler) {
	var mergedHandler = hander
	for i := len(r.MiddlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.MiddlewareChain[i](mergedHandler) // 中间件层层嵌套
	}
	r.Mux[path] = mergedHandler
}
