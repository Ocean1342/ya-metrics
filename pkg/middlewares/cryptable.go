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
			hash := r.Header.Get("HashSHA256")
			if hash == "" {
				w.WriteHeader(http.StatusBadRequest)
				http.Error(w, "empty hash header", http.StatusBadRequest)
				return
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				http.Error(w, "could not read body", http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			//посчитать хеш тела
			crypter := hmac.New(sha256.New, []byte(secretKey))
			//закодировать
			crypter.Write(body)
			countedHash := hex.EncodeToString(crypter.Sum(nil))

			if strings.ToLower(countedHash) == strings.ToLower(hash) {
				next.ServeHTTP(w, r)
				return
			}
			sugar.Errorf("get wrong hash")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "not same hash", http.StatusBadRequest)
		}
		return http.HandlerFunc(fn)
	}
}
