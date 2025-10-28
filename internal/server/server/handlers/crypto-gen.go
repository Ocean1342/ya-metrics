package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
)

type CryptoKeysResponse struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func (h *Handler) CryptoGen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	publicKeyString, _ := publicKeyToString(&privateKey.PublicKey)
	privateKeyString, _ := privateKeyToString(privateKey)
	data, err := json.Marshal(CryptoKeysResponse{PrivateKey: privateKeyString, PublicKey: publicKeyString})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func privateKeyToString(key *rsa.PrivateKey) (string, error) {
	// Преобразуем приватный ключ в DER-формат
	privKeyBytes := x509.MarshalPKCS1PrivateKey(key)

	// Кодируем в PEM
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	return string(privKeyPEM), nil
}

func publicKeyToString(key *rsa.PublicKey) (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}

	// Кодируем в PEM
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), nil
}
