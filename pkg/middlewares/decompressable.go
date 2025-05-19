package middlewares

import (
	"compress/gzip"
	"net/http"
	"ya-metrics/internal/server/server"
)

func NewDecodeableRequestMiddleware() server.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Content-Encoding")
			if h == "gzip" {
				dec, _ := gzip.NewReader(r.Body)
				r.Body = dec
				defer r.Body.Close()
				next.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}
