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
	rateLimit := flag.Int("l", 1, "secret key")
	flag.Parse()
	if os.Getenv("ADDRESS") != "" {
		*host = os.Getenv("ADDRESS")
	}
	if os.Getenv("REPORT_INTERVAL") != "" {
		envValReportIntervalSec, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL"))
		if err != nil {
			panic(err)
		}
		*reportIntervalSec = envValReportIntervalSec
	}
	if os.Getenv("POLL_INTERVAL") != "" {
		valEnvPollIntervalSec, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
		if err != nil {
			panic(err)
		}
		*pollIntervalSec = valEnvPollIntervalSec
	}
	if os.Getenv("KEY") != "" {
		*secretKey = os.Getenv("KEY")
	}
	if os.Getenv("RATE_LIMIT") != "" {
		valRateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
		if err != nil {
			panic(err)
		}
		*rateLimit = valRateLimit
	}

	srvrAddr := fmt.Sprintf("http://%s", *host)
	timeToWork := time.Duration(180) * time.Second
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeToWork))
	defer cancel()

	cncrncyAgent := concurrencyagent.New(sugar, initClient(), uint(*rateLimit))
	cncrncyAgent.Run(ctx, srvrAddr, int64(*pollIntervalSec), *reportIntervalSec, *secretKey)
	for range ctx.Done() {
		sugar.Info("client shutting down")
		return
	}
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
