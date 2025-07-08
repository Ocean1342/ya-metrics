package runableagent

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
	"ya-metrics/internal/agent/mgen"
	"ya-metrics/pkg/mdata"
)

type SimpleAgent struct{}

func (s *SimpleAgent) SendMetrics(srvrAddr string, pCount int64, reportIntervalSec int) {
	buffer := bytes.NewBuffer([]byte(""))
	for _, m := range mgen.GenerateGaugeMetrics() {
		url := s.prepareURLGauge(srvrAddr, m.GetType(), m.GetName(), m.GetValue())
		req, err := s.requestPrepare(url, http.MethodPost, buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		resp := s.sendRequest(req)
		if resp == nil {
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
	url := s.prepareURLCounter(srvrAddr, c.GetType(), c.GetName(), c.GetValue())
	req, err := s.requestPrepare(url, http.MethodPost, buffer)
	//TODO: добавить обработку ошибок
	if err != nil {
		fmt.Println(err)
	}
	resp := s.sendRequest(req)
	if resp == nil {
		return
	}
	defer resp.Body.Close()
	err = s.responseAnalyze(resp)
	if err != nil {
		fmt.Println(err)
	}
	//sleep
	time.Sleep(time.Second * time.Duration(reportIntervalSec))
}

func (s *SimpleAgent) prepareURLCounter(base, typeName, name string, value int64) string {
	return base + fmt.Sprintf("/update/%s/%s/%d", typeName, name, value)
}

func (s *SimpleAgent) prepareURLGauge(base, typeName, name string, value float64) string {
	return base + fmt.Sprintf("/update/%s/%s/%v", typeName, name, value)
}

func (s *SimpleAgent) requestPrepare(url string, method string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	return req, nil
}

func (s *SimpleAgent) sendRequest(req *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	time.Sleep(1000 * time.Microsecond)
	return resp
}

// TODO: доделать анализ ответа
func (s *SimpleAgent) responseAnalyze(resp *http.Response) error {
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

	return nil
}
