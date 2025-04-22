package server_storage

import "ya-metrics/pkg/mdata"

type CounterStorage interface {
	Set(m mdata.Counter) error
	Get(name string) (mdata.Counter, error)
}

type GaugeStorage interface {
	Get(n string) *mdata.Gauge
	Set(m mdata.Gauge) error
}
