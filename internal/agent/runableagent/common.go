package runableagent

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

func secretReqPrepare(secretKey string, req *http.Request) *http.Request {
	if secretKey == "" {
		return req
	}
	body, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	crypter := hmac.New(sha256.New, []byte(secretKey))
	crypter.Write(body)
	hash := crypter.Sum(nil)
	req.Header.Set("HashSHA256", hex.EncodeToString(hash))
	return req
}
