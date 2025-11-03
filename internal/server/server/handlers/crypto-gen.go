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
		h.log.Errorf("error generating private key: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	publicKeyString, err := publicKeyToString(&privateKey.PublicKey)
	if err != nil {
		h.log.Errorf("error generating public key: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	privateKeyString, err := privateKeyToString(privateKey)
	if err != nil {
		h.log.Errorf("error on convert private key to string: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
