package middlewares

import (
	"bytes"
	"go.uber.org/zap"
	"io"
	"net/http"
	"ya-metrics/internal/server/server"
	"ya-metrics/internal/server/server/handlers"
)

func HashableMiddleware(secretKey string, sugar *zap.SugaredLogger) server.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if secretKey == "" {
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			hash := r.Header.Get("HashSHA256")
			if hash == "none" {
				next.ServeHTTP(w, r)
				return
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(&handlers.CryptoResponseWriter{
				W:         w,
				SecretKey: secretKey,
				Body:      body,
			}, r)
		}
		return http.HandlerFunc(fn)
	}
}
