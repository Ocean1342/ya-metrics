package middlewares

import (
	"go.uber.org/zap"
	"ya-metrics/config"
	"ya-metrics/internal/server/server"
)

func InitMiddlewares(cfg *config.Config, sugar *zap.SugaredLogger) []server.Middleware {
	return []server.Middleware{
		CryptoMiddleware(cfg.SecretKey, sugar),
		NewLogResponseMiddleware(sugar),
		NewCompressResponseMiddleware(),
		NewLogRequestMiddleware(sugar),
		NewDecompressRequestMiddleware(),
		HashableMiddleware(cfg.SecretKey, sugar),
	}
}
