package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"ya-metrics/internal/agent/mgen"
	"ya-metrics/pkg/mdata"
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

	for {
		select {
		case <-ctx.Done():
			fmt.Println("shutting down")
			return
		default:
			run(srvrAddr, int64(*pollIntervalSec), *reportIntervalSec)
		}

	}
}

func run(srvrAddr string, pCount int64, reportIntervalSec int) {
	func() {
		buffer := bytes.NewBuffer([]byte(""))
		for _, m := range mgen.GenerateGaugeMetrics() {
			url := prepareURLGauge(srvrAddr, m.GetType(), m.GetName(), m.GetValue())
			req, err := requestPrepare(url, http.MethodPost, buffer)
			if err != nil {
				fmt.Println(err)
			}
			resp := sendRequest(req)
			if resp == nil {
				fmt.Println("response is nil")
				continue
			}
			defer resp.Body.Close()
			err = responseAnalyze(resp)
			if err != nil {
				fmt.Println(err)
			}
		}
		pCount++
		c := mdata.NewSimpleCounter("PollCount", pCount)
		url := prepareURLCounter(srvrAddr, c.GetType(), c.GetName(), c.GetValue())
		req, err := requestPrepare(url, http.MethodPost, buffer)
		//TODO: добавить обработку ошибок
		if err != nil {
			fmt.Println(err)
		}
		resp := sendRequest(req)
		if resp == nil {
			fmt.Println("response is nil")
			return
		}
		defer resp.Body.Close()
		err = responseAnalyze(resp)
		if err != nil {
			fmt.Println(err)
		}
		//sleep
		time.Sleep(time.Second * time.Duration(reportIntervalSec))
	}()
}

func prepareURLCounter(base, typeName, name string, value int64) string {
	return base + fmt.Sprintf("/update/%s/%s/%d", typeName, name, value)
}

func prepareURLGauge(base, typeName, name string, value float64) string {
	return base + fmt.Sprintf("/update/%s/%s/%v", typeName, name, value)
}

func requestPrepare(url string, method string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	return req, nil
}

func sendRequest(req *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return resp
}

// TODO: доделать анализ ответа
func responseAnalyze(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("nil response")
	}
	// Читаем ответ
	defer resp.Body.Close()
	_, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return fmt.Errorf("error reading response")
	}

	fmt.Println("Response Status:", resp.Status)
	return nil
}
