package config

import (
	"encoding/json"
	"os"
)

type AgentConfig struct {
	HostStr          string `json:"host_str"`
	StoreInterval    int64  `json:"store_interval"`
	FileStoragePath  string `json:"file_storage_path"`
	RestoreOnStart   bool   `json:"restore_on_start"`
	ProfileEnabled   bool   `json:"profile_enabled"`
	SecretKey        string `json:"secret_key"`
	DbURL            string `json:"db_url"`
	CryptoPrivateKey string `json:"crypto_private_key"`
}

func ParseFromFile(filePath string) (*AgentConfig, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config AgentConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
