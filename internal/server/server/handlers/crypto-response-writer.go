package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
)

type CryptoResponseWriter struct {
	W         http.ResponseWriter
	SecretKey string
	Body      []byte
}

func (w *CryptoResponseWriter) Write(b []byte) (int, error) {
	return w.W.Write(b)
}

func (w *CryptoResponseWriter) Header() http.Header {
	return w.W.Header()
}

func (w *CryptoResponseWriter) WriteHeader(statusCode int) {
	w.W.Header().Set("Content-Type", "application/json")
	if w.SecretKey != "" {
		crypter := hmac.New(sha256.New, []byte(w.SecretKey))
		crypter.Write(w.Body)
		countedHash := hex.EncodeToString(crypter.Sum(nil))
		w.Header().Set("HashSHA256", countedHash)
		fmt.Println(w.Header())
	}
	w.W.WriteHeader(statusCode)
}
