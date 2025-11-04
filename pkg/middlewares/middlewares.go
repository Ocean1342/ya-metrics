package middlewares

import (
	"go.uber.org/zap"
	"ya-metrics/config"
	"ya-metrics/internal/server/server"
	"ya-metrics/pkg/crypto"
)

func InitMiddlewares(cfg *config.Config, sugar *zap.SugaredLogger, crypter *crypto.PrivateCrypter) []server.Middleware {
	return []server.Middleware{
		RSADecryptableMiddleware(crypter, sugar),
		CryptoMiddleware(cfg.SecretKey, sugar),
		NewLogResponseMiddleware(sugar),
		NewCompressResponseMiddleware(),
		NewLogRequestMiddleware(sugar),
		HashableMiddleware(cfg.SecretKey, sugar),
		NewDecompressRequestMiddleware(sugar),
	}
}
