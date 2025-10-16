package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
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

func (p *PrivateCrypter) DecryptHTTPRequest(req *http.Request) ([]byte, error) {
	encryptedBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	defer req.Body.Close()
	if len(encryptedBody) == 0 {
		return nil, nil
	}
	var decryptedData []byte
	decryptedData, err = decryptData(p.privateKey, encryptedBody)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}
	return decryptedData, nil
}

// decryptData расшифровывает данные с использованием RSA-OAEP
func decryptData(privateKey *rsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	decrypted, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privateKey,
		encryptedData,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("RSA decryption failed: %w", err)
	}

	return decrypted, nil
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
