package runagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"ya-metrics/internal/agent/mgen"
	"ya-metrics/pkg/mdata"
)

type JSONAgent struct{}

func (s *JSONAgent) Run(srvrAddr string, pCount int64, reportIntervalSec int) {
	func() {
		url := s.prepareURL(srvrAddr)
		for _, m := range mgen.GenerateGaugeMetrics() {
			req, err := s.gaugeRequestPrepare(m, url, http.MethodPost)
			if err != nil {
				fmt.Println(err)
			}
			resp := s.sendRequest(req)
			if resp == nil {
				fmt.Println("response is nil")
				continue
			}
			defer resp.Body.Close()
			err = s.responseAnalyze(resp)
			if err != nil {
				fmt.Println(err)
			}
		}
		pCount++
		c := mdata.NewSimpleCounter("PollCount", pCount)
		req, err := s.counterRequestPrepare(c, url, http.MethodPost)
		defer req.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
		resp := s.sendRequest(req)
		if resp == nil {
			fmt.Println("response is nil")
			return
		}
		defer resp.Body.Close()

		err = s.responseAnalyze(resp)
		if err != nil {
			fmt.Println(err)
		}
		//sleep
		time.Sleep(time.Second * time.Duration(reportIntervalSec))
	}()
}

func (s *JSONAgent) prepareURL(base string) string {
	return fmt.Sprintf("%s/update/", base)
}

func (s *JSONAgent) counterRequestPrepare(c mdata.Counter, url string, method string) (*http.Request, error) {
	value := c.GetValue()
	metric := mdata.Metrics{ID: c.GetName(), MType: c.GetType(), Delta: &value}
	j, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(j))
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (s *JSONAgent) gaugeRequestPrepare(g mdata.Gauge, url string, method string) (*http.Request, error) {
	value := g.GetValue()
	data, err := json.Marshal(mdata.Metrics{ID: g.GetName(), MType: g.GetType(), Value: &value})
	b := bytes.NewBuffer(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (s *JSONAgent) sendRequest(req *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	time.Sleep(1000 * time.Microsecond)
	return resp
}

func (s *JSONAgent) responseAnalyze(resp *http.Response) error {
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
