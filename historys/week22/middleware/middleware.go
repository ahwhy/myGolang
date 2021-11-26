package middleware

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []Middleware
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Use(m Middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *Router) Merge(h http.Handler) http.Handler {
	var mergedHandler = h

	// customizedHandler = logger(timeout(ratelimit(helloHandler)))
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}

	return mergedHandler
}

func TimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		// next handler
		next.ServeHTTP(wr, r)

		timeElapsed := time.Since(timeStart)
		log.Println(timeElapsed)
	})
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		log.Println("start")

		// next handler
		next.ServeHTTP(wr, r)

		log.Println("end")
	})
}
