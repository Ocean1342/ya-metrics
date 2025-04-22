package server_storage

import "ya-metrics/pkg/mdata"

func NewSimpleGaugeStorage() GaugeStorage {
	return &SimpleGaugeStorage{}
}

type SimpleGaugeStorage struct {
}

func (s *SimpleGaugeStorage) Get(n string) *mdata.Gauge {
	return nil
}

func (s *SimpleGaugeStorage) Set(m mdata.Gauge) error {
	return nil
}
