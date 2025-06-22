package srvrstrg

import "ya-metrics/pkg/mdata"

type CounterStorage interface {
	Set(m mdata.Counter) error
	Get(name string) (mdata.Counter, error)
	StorableStorage
}

type GaugeStorage interface {
	Get(n string) mdata.Gauge
	Set(m mdata.Gauge) error
	StorableStorage
}

type StorableStorage interface {
	GetList() map[string]string
	GetMetrics() []mdata.Metrics
	SetFrom(metrics []mdata.Metrics) error
}
