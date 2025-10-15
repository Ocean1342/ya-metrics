package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go.uber.org/zap"
	"os"
)

type PrivateCrypter struct {
	privateKey *rsa.PrivateKey
	log        *zap.SugaredLogger
}

func NewPrivateCrypter(privateKeyFilePath string, log *zap.SugaredLogger) (*PrivateCrypter, error) {
	key, err := loadPrivateKeyFromFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}
	return &PrivateCrypter{
		privateKey: key,
		log:        log,
	}, nil
}

func loadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		// PKCS1 формат
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		// PKCS8 формат
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return key.(*rsa.PrivateKey), nil
	default:
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}
}
