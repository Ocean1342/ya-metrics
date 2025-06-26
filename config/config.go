package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port             int               `json:"port"`
	Host             string            `json:"host"`
	HostString       string            `json:"host_str"`
	PermStoreOptions *PermStoreOptions `json:"perm_store_options"`
}

type PermStoreOptions struct {
	StoreInterval   int64  `env:"STORE_INTERVAL" default:"300"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" default:"perm_storage.local.json"`
	RestoreOnStart  bool   `env:"RESTORE" default:"false"`
}

func New() *Config {
	//TODO: envы должны быть в приоритете
	hostStr := flag.String("a", "localhost:8080", "server address")
	storeInterval := flag.Int64("i", 300, "server address")
	fileStoragePath := flag.String("f", "./perm_storage.json", "server address")
	restoreOnStart := flag.Bool("r", false, "restore storage from file")
	flag.Parse()
	if os.Getenv("ADDRESS") != "" {
		*hostStr = os.Getenv("ADDRESS")
	}
	if os.Getenv("STORE_INTERVAL") != "" {
		v, err := strconv.Atoi(os.Getenv("STORE_INTERVAL"))
		if err != nil {
			panic(err)
		}
		*storeInterval = int64(v)
	}
	if os.Getenv("FILE_STORAGE_PATH") != "" {
		*fileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	}
	restoreEnv := os.Getenv("RESTORE")
	if restoreEnv != "" {
		switch strings.ToLower(restoreEnv) {
		case "true":
			*restoreOnStart = true
		case "false":
			*restoreOnStart = false
		default:
			panic(fmt.Sprintf("invalid RESTORE env value: %s", restoreEnv))
		}
		*restoreOnStart = true
	}

	return &Config{
		Port:       8080,
		Host:       "localhost",
		HostString: *hostStr,
		PermStoreOptions: &PermStoreOptions{
			StoreInterval:   *storeInterval,
			FileStoragePath: *fileStoragePath,
			RestoreOnStart:  *restoreOnStart,
		},
	}
}
