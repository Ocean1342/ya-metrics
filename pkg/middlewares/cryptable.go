package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"ya-metrics/internal/server/server"
)

func CryptoMiddleware(secretKey string, sugar *zap.SugaredLogger) server.Middleware {
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

			if hash == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			//посчитать хеш тела
			crypter := hmac.New(sha256.New, []byte(secretKey))
			//закодировать
			crypter.Write(body)
			countedHash := hex.EncodeToString(crypter.Sum(nil))
			if strings.EqualFold(countedHash, hash) {
				next.ServeHTTP(w, r)
			} else {
				sugar.Errorf("get wrong hash")
				w.WriteHeader(http.StatusBadRequest)
			}
		}
		return http.HandlerFunc(fn)
	}
}
