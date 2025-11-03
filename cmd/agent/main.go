package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"time"
	"ya-metrics/internal/agent/concurrencyagent"
	"ya-metrics/internal/agent/config"
	"ya-metrics/pkg/crypto"
)

var (
	sugar        *zap.SugaredLogger
	buildVersion string = "N\\A"
	buildDate    string = "N\\A"
	buildCommit  string = "N\\A"
)

func main() {
	fmt.Printf("Build version:=%s, Build date=%s Build commit=%s\n", buildVersion, buildDate, buildCommit)
	initLogger()
	host := flag.String("a", "localhost:8080", "agent host")
	reportIntervalSec := flag.Int("r", 5, "report interval")
	pollIntervalSec := flag.Int("p", 5, "poll interval")
	secretKey := flag.String("k", "", "secret key")
	rateLimit := flag.Int("l", 1, "rate limit")
	cryptoPublicKey := flag.String("crypto-key", "", "crypto public key")
	cfgFilePath := flag.String("config", "", "crypto public key")
	flag.Parse()

	if v, ok := os.LookupEnv("CONFIG"); ok {
		*cfgFilePath = v
	}
	if v, ok := os.LookupEnv("ADDRESS"); ok {
		*host = v
	}
	if v, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		envValReportIntervalSec, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		*reportIntervalSec = envValReportIntervalSec
	}
	if v, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		valEnvPollIntervalSec, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		*pollIntervalSec = valEnvPollIntervalSec
	}
	if v, ok := os.LookupEnv("KEY"); ok {
		*secretKey = v
	}
	if v, ok := os.LookupEnv("RATE_LIMIT"); ok {
		valRateLimit, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		*rateLimit = valRateLimit
	}

	if v, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		*cryptoPublicKey = v
	}

	if *cfgFilePath != "" {
		cfg, err := config.ParseFromFile(*cfgFilePath)
		if err != nil {
			sugar.Errorf("wrong config file path: %s", *cfgFilePath)
		} else {
			if *host == "" {
				*host = cfg.Host
			}
			if *reportIntervalSec == 0 {
				*reportIntervalSec = cfg.ReportIntervalSec
			}
			if *pollIntervalSec == 0 {
				*pollIntervalSec = cfg.PollIntervalSec
			}
			if *secretKey == "" {
				*secretKey = cfg.SecretKey
			}
			if *rateLimit == 0 {
				*rateLimit = cfg.RateLimit
			}
			if *cryptoPublicKey == "" {
				*cryptoPublicKey = cfg.CryptoPublicKey
			}
		}
	}

	srvrAddr := fmt.Sprintf("http://%s", *host)
	timeToWork := time.Duration(180) * time.Second
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeToWork))
	defer cancel()

	publicCrypter, err := crypto.NewPublicCrypter(*cryptoPublicKey, sugar)
	if err != nil {
		sugar.Errorf("could not create crypter. Abort server init.")
		cancel()
	}
	cncrncyAgent := concurrencyagent.New(sugar, initClient(), uint(*rateLimit), publicCrypter)
	cncrncyAgent.Run(ctx, srvrAddr, int64(*pollIntervalSec), *reportIntervalSec, *secretKey)
	//graceful shutdown
	<-ctx.Done()
	sugar.Info("client shutting down")
	return

}

func initLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("could not start logger")
	}
	defer logger.Sync()
	sugar = logger.Sugar()
}

func initClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second
	return retryClient.StandardClient()
}
