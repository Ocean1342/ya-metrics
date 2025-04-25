package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
	"ya-metrics/internal/agent/mgen"
	"ya-metrics/pkg/mdata"
)

func main() {
	//TODO: вынести в конфиги
	var pCount int64
	srvrAddr := "https://localhost:8080"
	reportIntervalSec := 10

	for {
		buffer := bytes.NewBuffer([]byte(""))
		for _, m := range mgen.GenerateGaugeMetrics() {
			url := prepareUrlGauge(srvrAddr, m.GetType(), m.GetName(), m.GetValue())
			req, err := requestPrepare(url, http.MethodPost, buffer)

			if err != nil {
				fmt.Println(err)
			}

			resp := sendRequest(req)
			err = responseAnalyze(resp)
			if err != nil {
				fmt.Println(err)
			}
		}
		pCount++
		c := mdata.NewSimpleCounter("PollCount", pCount)
		url := prepareUrlCounter(srvrAddr, c.GetType(), c.GetName(), c.GetValue())
		req, err := requestPrepare(url, http.MethodPost, buffer)
		//TODO: добавить обработку ошибок
		if err != nil {
			fmt.Println(err)
		}
		resp := sendRequest(req)
		err = responseAnalyze(resp)
		if err != nil {
			fmt.Println(err)
		}
		//sleep
		time.Sleep(time.Second * time.Duration(reportIntervalSec))
	}

}

func prepareUrlCounter(base, typeName, name string, value int64) string {
	return base + fmt.Sprintf("/update/%s/%s/%d", typeName, name, value)
}

func prepareUrlGauge(base, typeName, name string, value float64) string {
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
	// Читаем ответ
	if resp == nil {
		return fmt.Errorf("nil response")
	}
	_, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return fmt.Errorf("error reading response")
	}

	fmt.Println("Response Status:", resp.Status)
	return nil
}
