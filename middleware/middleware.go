package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Header inject default headers to HTTP response.
func Headers() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Keep-alive", "timeout=5, max=1000") // keep alive for max 100 connections
			w.Header().Add("Cache-Control", "max-age=600")      // 600 seconds for cache ttl
			next.ServeHTTP(w, req)
		})
	}
}
