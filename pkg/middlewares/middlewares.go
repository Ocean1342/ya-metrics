package middlewares

import (
	"go.uber.org/zap"
	"ya-metrics/internal/server/server"
)

func InitMiddlewares(sugar *zap.SugaredLogger) []server.Middleware {
	return []server.Middleware{
		NewLogResponseMiddleware(sugar),
		NewCompressResponseMiddleware(),
		NewLogRequestMiddleware(sugar),
		NewDecompressRequestMiddleware(),
	}
}
