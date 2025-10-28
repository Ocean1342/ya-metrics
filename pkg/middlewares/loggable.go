package middlewares

import (
	"go.uber.org/zap"
	"net/http"
	"time"
	"ya-metrics/internal/server/server"
)

func NewLogRequestMiddleware(sugar *zap.SugaredLogger) server.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
			next.ServeHTTP(w, r)
			sugar.Infof(
				"Request: Url: %s. Method: %s TimeExecute: %d ms",
				r.URL.Path,
				r.Method,
				time.Since(t).Microseconds(),
			)
		}
		return http.HandlerFunc(fn)
	}
}

func NewLogResponseMiddleware(sugar *zap.SugaredLogger) server.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			lr := loggingResponseWriter{
				l: sugar,
				w: w,
				r: responseData{size: 0, status: 0},
			}
			next.ServeHTTP(&lr, r)
			sugar.Infof("Response: Size:%d b, Status:%d", lr.r.size, lr.r.status)
		}
		return http.HandlerFunc(fn)
	}
}

type responseData struct {
	size   int
	status int
}

type loggingResponseWriter struct {
	l *zap.SugaredLogger
	w http.ResponseWriter
	r responseData
}

func (lr *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lr.w.Write(b)
	lr.r.size += size
	return size, err
}

func (lr *loggingResponseWriter) Header() http.Header {
	return lr.w.Header()
}

func (lr *loggingResponseWriter) WriteHeader(statusCode int) {
	lr.r.status = statusCode
	lr.w.WriteHeader(statusCode)
}
