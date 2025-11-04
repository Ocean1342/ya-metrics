package middlewares

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"ya-metrics/internal/server/server"
	"ya-metrics/pkg/crypto"
)

func RSADecryptableMiddleware(privateCrypter *crypto.PrivateCrypter, sugar *zap.SugaredLogger) server.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Encrypted") != "true" {
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			decrBody, err := privateCrypter.DecryptHTTPRequest(r)
			if decrBody == nil {
				fmt.Println("test")
			}
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(decrBody))
		}
		return http.HandlerFunc(fn)
	}
}
