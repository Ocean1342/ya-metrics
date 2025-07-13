package middlewares

import (
	"bytes"
	"compress/gzip"
	"go.uber.org/zap"
	"io"
	"net/http"
	"ya-metrics/internal/server/server"
)

func NewDecompressRequestMiddleware(log *zap.SugaredLogger) server.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Content-Encoding")
			if h == "gzip" {
				originalBody, err := io.ReadAll(r.Body)
				if err != nil {
					log.Error(err)
				}
				dec, err := gzip.NewReader(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(originalBody))
				if err != nil {
					log.Errorf("could not decode gzip:%s", err)
				} else {
					r.Body = dec
				}
				next.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}
