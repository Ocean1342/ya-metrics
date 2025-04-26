package srvrstrg

import "ya-metrics/pkg/mdata"

func NewSimpleGaugeStorage() GaugeStorage {
	return &SimpleGaugeStorage{
		storage: make(map[string]mdata.Gauge, 1_000_000),
	}
}

type SimpleGaugeStorage struct {
	storage map[string]mdata.Gauge
}

func (s *SimpleGaugeStorage) Get(n string) *mdata.Gauge {
	elem, ok := s.storage[n]
	if !ok {
		return nil
	}
	return &elem
}

func (s *SimpleGaugeStorage) Set(m mdata.Gauge) error {
	s.storage[m.GetName()] = m
	return nil
}
