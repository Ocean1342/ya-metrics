package srvrstrg

import "ya-metrics/pkg/mdata"

type CounterStorage interface {
	Set(m mdata.Counter) error
	Get(name string) (mdata.Counter, error)
	Listable
}

type GaugeStorage interface {
	Get(n string) mdata.Gauge
	Set(m mdata.Gauge) error
	Listable
}

type Listable interface {
	GetList() map[string]string
}
