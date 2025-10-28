package config

import (
	"encoding/json"
	"os"
)

type FromFIleConfig struct {
	HostStr          string `json:"host_str"`
	StoreInterval    int64  `json:"store_interval"`
	FileStoragePath  string `json:"file_storage_path"`
	RestoreOnStart   bool   `json:"restore_on_start"`
	ProfileEnabled   bool   `json:"profile_enabled"`
	SecretKey        string `json:"secret_key"`
	DBURL            string `json:"db_url"`
	CryptoPrivateKey string `json:"crypto_private_key"`
}

func ParseFromFile(filePath string) (*FromFIleConfig, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config FromFIleConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
