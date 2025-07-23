package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
	"ya-metrics/internal/agent/concurrencyagent"
)

var sugar *zap.SugaredLogger

func main() {
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
	//TODO: костыль, чтобы дать время серверу подняться
	//time.Sleep(11 * time.Second)
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3                   // Максимальное количество попыток
	retryClient.RetryWaitMin = 1 * time.Second // Минимальное время ожидания между попытками
	retryClient.RetryWaitMax = 5 * time.Second
	cncrncyAgent := concurrencyagent.New(sugar, retryClient, uint(*rateLimit))
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
