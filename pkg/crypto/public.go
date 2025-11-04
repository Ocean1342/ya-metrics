package crypto

import (
	"bytes"
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

type PublicCrypter struct {
	publicKey *rsa.PublicKey
	log       *zap.SugaredLogger
	Enabled   bool
}

func NewPublicCrypter(publicKeyFilePath string, log *zap.SugaredLogger) (*PublicCrypter, error) {
	key, err := loadPublicKeyFromFile(publicKeyFilePath)
	if err != nil {
		return &PublicCrypter{Enabled: false, log: log}, err
	}
	return &PublicCrypter{
		publicKey: key,
		log:       log,
		Enabled:   true,
	}, nil
}

func loadPublicKeyFromFile(filename string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	switch block.Type {
	case "PUBLIC KEY":
		// PKIX формат
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return pub.(*rsa.PublicKey), nil
	case "RSA PUBLIC KEY":
		// PKCS1 формат
		return x509.ParsePKCS1PublicKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}
}

func (p *PublicCrypter) CryptRequest(req *http.Request) *http.Request {
	// Читаем тело запроса
	body, err := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		p.log.Errorf("could not crypt request body. err: %v", err)
		return req
	}
	defer req.Body.Close()

	// Если тело пустое, ничего не шифруем
	if len(body) == 0 {
		req.Body = io.NopCloser(bytes.NewReader(body))
		return req
	}

	// Шифруем тело запроса
	encryptedBody, err := encryptData(p.publicKey, body)
	if err != nil {
		p.log.Errorf("could not crypt request body. err: %v", err)
		return req
	}

	// Устанавливаем зашифрованное тело обратно в запрос
	req.Body = io.NopCloser(bytes.NewReader(encryptedBody))
	req.ContentLength = int64(len(encryptedBody))
	req.Header.Set("X-Encrypted", "true")
	return req
}

func encryptData(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	encrypted, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		data,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("RSA encryption failed: %w", err)
	}

	return encrypted, nil
}
