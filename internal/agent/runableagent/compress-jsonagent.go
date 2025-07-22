package runableagent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
	"ya-metrics/internal/agent/mgen"
	"ya-metrics/pkg/mdata"
)

type CompressJSONAgent struct {
	SecretKey string
	SendCh    chan *http.Request
	Logger    *zap.SugaredLogger
}

func (s *CompressJSONAgent) SendMetrics(srvrAddr string, pCount int64, reportIntervalSec int) {
	url := s.prepareURL(srvrAddr)
	for _, m := range mgen.GenerateGaugeMetrics() {
		req, err := s.gaugeRequestPrepare(m, url, http.MethodPost)
		if err != nil {
			s.Logger.Error(err)
		}
		s.sendRequest(req)
	}
	pCount++
	req, err := s.counterRequestPrepare(mdata.NewSimpleCounter("PollCount", pCount), url, http.MethodPost)
	if err != nil {
		s.Logger.Error(err)
	}
	s.sendRequest(req)

	//sleep
	time.Sleep(time.Second * time.Duration(reportIntervalSec))
}

func (s *CompressJSONAgent) prepareURL(base string) string {
	return fmt.Sprintf("%s/update/", base)
}

func (s *CompressJSONAgent) compressData(in []byte) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	w := gzip.NewWriter(b)
	_, err := w.Write(in)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (s *CompressJSONAgent) counterRequestPrepare(c mdata.Counter, url string, method string) (*http.Request, error) {
	value := c.GetValue()
	metric := mdata.Metrics{ID: c.GetName(), MType: c.GetType(), Delta: &value}
	j, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}
	compressed, err := s.compressData(j)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(compressed))
	if err != nil {
		s.Logger.Error("Error creating request:", err)
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")
	return req, nil
}

func (s *CompressJSONAgent) gaugeRequestPrepare(g mdata.Gauge, url string, method string) (*http.Request, error) {
	value := g.GetValue()
	data, err := json.Marshal(mdata.Metrics{ID: g.GetName(), MType: g.GetType(), Value: &value})
	if err != nil {
		return nil, err
	}
	compressed, err := s.compressData(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(compressed))
	if err != nil {
		s.Logger.Error("Error creating request:", err)
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")
	return req, nil
}

func (s *CompressJSONAgent) sendRequest(req *http.Request) {
	secretReqPrepare(s.SecretKey, req)
	s.SendCh <- req
}
