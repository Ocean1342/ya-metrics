package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
	"ya-metrics/internal/server/server"
)

func NewCompressResponseMiddleware() server.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			b := validateAcceptEncoding(r.Header.Values("Accept-Encoding")) && validateContentType(r.Header.Values("Content-Type"))
			if b {
				nw := NewCompressableResponseWriter(w)
				next.ServeHTTP(nw, r)
				//TODO: ?
				defer nw.Close()
			} else {
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}

func validateAcceptEncoding(values []string) bool {
	for _, v := range values {
		if strings.Contains(v, "gzip") {
			return true
		}
	}
	return false
}

func validateContentType(values []string) bool {
	for _, v := range values {
		if strings.EqualFold(v, "text/html") || strings.EqualFold(v, "application/json") {
			return true
		}
	}
	return false
}

type CompressableResponseWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressableResponseWriter(w http.ResponseWriter) *CompressableResponseWriter {
	return &CompressableResponseWriter{
		w,
		gzip.NewWriter(w),
	}
}

func (c *CompressableResponseWriter) Write(b []byte) (int, error) {
	return c.zw.Write(b)
}

func (c *CompressableResponseWriter) Header() http.Header {
	return c.w.Header()
}

func (c *CompressableResponseWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

func (c *CompressableResponseWriter) Close() error {
	return c.zw.Close()
}
