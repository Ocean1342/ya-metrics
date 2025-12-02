package middlewares

import (
	"go.uber.org/zap"
	"net/http"
	"ya-metrics/internal/server/server"
	"ya-metrics/pkg/netcmprr"
)

func TrustableMiddleware(sugar *zap.SugaredLogger, trustedSubnet string) server.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet == "" {
				next.ServeHTTP(w, r)
				return
			}
			reqIP := r.Header.Get("X-Real-Ip")
			if reqIP == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			isTrusted, err := netcmprr.IsTrustedSubnet(trustedSubnet, reqIP)
			if err != nil {
				sugar.Errorf("error on count trusted subnet: %s", err)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if !isTrusted {
				sugar.Errorf("request from untrusted subnet: %s", reqIP)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
