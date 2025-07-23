package srvrstrg

import (
	"go.uber.org/zap"
	"strconv"
	"sync"
	"ya-metrics/pkg/mdata"
)

func NewSimpleGaugeStorage(log *zap.SugaredLogger) GaugeStorage {
	return &SimpleGaugeStorage{
		mu:      sync.RWMutex{},
		storage: make(map[string]mdata.Gauge, 1_000_000),
		log:     log,
	}
}

type SimpleGaugeStorage struct {
	mu      sync.RWMutex
	storage map[string]mdata.Gauge
	log     *zap.SugaredLogger
}

func (s *SimpleGaugeStorage) Get(n string) mdata.Gauge {
	s.mu.RLock()
	defer s.mu.RUnlock()
	elem, ok := s.storage[n]
	if !ok {
		return nil
	}
	return elem
}

func (s *SimpleGaugeStorage) Set(m mdata.Gauge) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage[m.GetName()] = m
	return nil
}

func (s *SimpleGaugeStorage) GetList() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make(map[string]string, len(s.storage))
	for k, v := range s.storage {
		res[k] = strconv.Itoa(int(v.GetValue()))
	}
	return res
}
func (s *SimpleGaugeStorage) GetMetrics() []mdata.Metrics {
	md := make([]mdata.Metrics, len(s.storage))
	i := 0
	s.mu.RLock()
	defer s.mu.RUnlock()
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

func (s *SimpleGaugeStorage) SetFrom(metrics []mdata.Metrics) error {
	factory := mdata.NewSimpleGauge
	for _, m := range metrics {
		if m.MType != mdata.GAUGE {
			continue
		}

		if m.Value == nil {
			s.log.Errorf("received nil gauge:%s", m.ID)
			continue
		}
		value := *m.Value
		err := s.Set(factory(m.ID, value))
		if err != nil {
			return err
		}
	}
	return nil
}
