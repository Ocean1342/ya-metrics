package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"time"
	"ya-metrics/internal/agent/concurrencyagent"
)

var sugar *zap.SugaredLogger

func main() {
	initLogger()
	host := flag.String("a", "localhost:8080", "agent host")
	reportIntervalSec := flag.Int("r", 10, "report interval")
	pollIntervalSec := flag.Int("p", 2, "poll interval")
	secretKey := flag.String("k", "", "secret key")
	rateLimit := flag.Int("l", 5, "secret key")
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
	time.Sleep(11 * time.Second)
	cncrncyAgent := concurrencyagent.New(sugar, &http.Client{}, uint(*rateLimit))
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
