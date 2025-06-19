package srvrstrg

import (
	"strconv"
	"ya-metrics/pkg/mdata"
)

func NewSimpleGaugeStorage() GaugeStorage {
	return &SimpleGaugeStorage{
		storage: make(map[string]mdata.Gauge, 1_000_000),
	}
}

type SimpleGaugeStorage struct {
	storage map[string]mdata.Gauge
}

func (s *SimpleGaugeStorage) Get(n string) mdata.Gauge {
	elem, ok := s.storage[n]
	if !ok {
		return nil
	}
	return elem
}

func (s *SimpleGaugeStorage) Set(m mdata.Gauge) error {
	s.storage[m.GetName()] = m
	return nil
}

func (s *SimpleGaugeStorage) GetList() map[string]string {
	res := make(map[string]string, len(s.storage))
	for k, v := range s.storage {
		res[k] = strconv.Itoa(int(v.GetValue()))
	}
	return res
}
func (s *SimpleGaugeStorage) GetMetrics() []mdata.Metrics {
	md := make([]mdata.Metrics, len(s.storage))
	i := 0
	for _, gauge := range s.storage {
		value := gauge.GetValue()
		md[i] = mdata.Metrics{
			ID:    gauge.GetName(),
			MType: mdata.GAUGE,
			Value: &value,
		}
		i++
	}
	return md
}
