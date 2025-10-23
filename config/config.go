package config

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	filecfg "ya-metrics/internal/server/config"
)

type Config struct {
	Port             int               `json:"port"`
	Host             string            `json:"host"`
	HostString       string            `json:"host_str"`
	PermStoreOptions *PermStoreOptions `json:"perm_store_options"`
	DBURL            string            `json:"db_url"`
	SecretKey        string            `json:"secret_key"`
	ProfilingEnabled bool              `json:"profiling_enabled"`
	CryptoKey        string            `json:"crypto_key"`
}

type PermStoreOptions struct {
	StoreInterval   int64  `env:"STORE_INTERVAL" default:"300"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" default:"perm_storage.local.json"`
	RestoreOnStart  bool   `env:"RESTORE" default:"false"`
}

func New(log *zap.SugaredLogger) *Config {
	hostStr := flag.String("a", "localhost:8080", "server address")
	storeInterval := flag.Int64("i", 300, "server address")
	fileStoragePath := flag.String("f", "./perm_storage.json", "server address")
	restoreOnStart := flag.Bool("r", false, "restore storage from file")
	profileEnabled := flag.Bool("pprof", true, "enable profiling")
	secretKey := flag.String("k", "", "secret key")
	//dbDefaultString := "host=localhost port=5432 user=ya password=ya dbname=ya sslmode=disable"
	dbURL := flag.String("d", "", "server address")
	cryptoPrivateKey := flag.String("crypto-key", "", "crypto private key")
	cfgFilePath := flag.String("config", "", "crypto public key")
	flag.Parse()
	if os.Getenv("CONFIG") != "" {
		*cfgFilePath = os.Getenv("CONFIG")
	}
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
	if os.Getenv("DATABASE_DSN") != "" {
		*dbURL = os.Getenv("DATABASE_DSN")
	}
	if os.Getenv("KEY") != "" {
		*secretKey = os.Getenv("KEY")
	}
	if os.Getenv("CRYPTO_KEY") != "" {
		*cryptoPrivateKey = os.Getenv("CRYPTO_KEY")
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

	profileEnabledEnv := os.Getenv("PROFILING_ENABLED")
	if profileEnabledEnv != "" {
		switch strings.ToLower(profileEnabledEnv) {
		case "true":
			*profileEnabled = true
		case "false":
			*profileEnabled = false
		default:
			panic(fmt.Sprintf("invalid profileEnabled env value: %s", profileEnabledEnv))
		}
		*profileEnabled = true
	}

	if *cfgFilePath != "" {
		cfg, err := filecfg.ParseFromFile(*cfgFilePath)
		if err != nil {
			log.Errorf("wrong config file path: %s", *cfgFilePath)
		} else {
			if *hostStr == "" {
				*hostStr = cfg.HostStr
			}
			if *storeInterval == 0 {
				*storeInterval = cfg.StoreInterval
			}
			if *fileStoragePath == "" {
				*fileStoragePath = cfg.FileStoragePath
			}
			if restoreOnStart == nil {
				restoreOnStart = &cfg.RestoreOnStart
			}
			if profileEnabled == nil {
				profileEnabled = &cfg.ProfileEnabled
			}
			if *secretKey == "" {
				*secretKey = cfg.SecretKey
			}
			if *dbURL == "" {
				*dbURL = cfg.DBURL
			}
			if *cryptoPrivateKey == "" {
				*cryptoPrivateKey = cfg.CryptoPrivateKey
			}
		}
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
		DBURL:            *dbURL,
		SecretKey:        *secretKey,
		ProfilingEnabled: *profileEnabled,
		CryptoKey:        *cryptoPrivateKey,
	}
}
