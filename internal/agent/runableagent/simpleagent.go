package runableagent

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
	"ya-metrics/internal/agent/mgen"
	"ya-metrics/pkg/mdata"
)

type SimpleAgent struct {
	SecretKey string
	SendCh    chan *http.Request
	Logger    *zap.SugaredLogger
}

func (s *SimpleAgent) SendMetrics(srvrAddr string, pCount int64, reportIntervalSec int) {
	buffer := bytes.NewBuffer([]byte(""))
	for _, m := range mgen.GenerateGaugeMetrics() {
		url := s.prepareURLGauge(srvrAddr, m.GetType(), m.GetName(), m.GetValue())
		req, err := s.requestPrepare(url, http.MethodPost, buffer)
		if err != nil {
			s.Logger.Error(err)
			continue
		}
		s.sendRequest(req)
	}
	pCount++
	c := mdata.NewSimpleCounter("PollCount", pCount)
	url := s.prepareURLCounter(srvrAddr, c.GetType(), c.GetName(), c.GetValue())
	req, err := s.requestPrepare(url, http.MethodPost, buffer)

	if err != nil {
		s.Logger.Error(err)
	}
	s.sendRequest(req)
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
		s.Logger.Error(err)
		panic(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	return req, nil
}

func (s *SimpleAgent) sendRequest(req *http.Request) {
	secretReqPrepare(s.SecretKey, req)
	s.SendCh <- req
}
