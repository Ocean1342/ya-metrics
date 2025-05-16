package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
	"ya-metrics/internal/agent/runagent"
)

func main() {
	host := flag.String("a", "localhost:8080", "agent host")
	reportIntervalSec := flag.Int("r", 10, "report interval")
	pollIntervalSec := flag.Int("p", 2, "poll interval")
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

	srvrAddr := fmt.Sprintf("http://%s", *host)
	timeToWork := time.Duration(120) * time.Second
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeToWork))
	defer cancel()
	j := runagent.JsonAgent{}
	a := runagent.SimpleAgent{}
	for {
		select {
		case <-ctx.Done():
			fmt.Println("shutting down")
			return
		default:
			a.Run(srvrAddr, int64(*pollIntervalSec), *reportIntervalSec)
			j.Run(srvrAddr, int64(*pollIntervalSec), *reportIntervalSec)
		}
	}
}
