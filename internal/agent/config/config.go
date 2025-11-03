package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type AgentConfig struct {
	Host              string `json:"host"`
	ReportIntervalSec int    `json:"report_interval_sec"`
	PollIntervalSec   int    `json:"poll_interval_sec"`
	SecretKey         string `json:"secret_key"`
	RateLimit         int    `json:"rate_limit"`
	CryptoPublicKey   string `json:"crypto_public_key"`
}

func ParseFromFile(filePath string) (*AgentConfig, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var config AgentConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return &config, nil
}
