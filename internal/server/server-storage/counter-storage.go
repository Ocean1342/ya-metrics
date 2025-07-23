package srvrstrg

import (
	"strconv"
	"sync"
	"ya-metrics/pkg/mdata"
)

type SimpleCountStorage struct {
	mu      sync.RWMutex
	storage map[string]mdata.Counter
	factory func(name string, value int64) mdata.Counter
}

func NewSimpleCountStorage(factory func(name string, value int64) mdata.Counter) *SimpleCountStorage {
	return &SimpleCountStorage{mu: sync.RWMutex{}, storage: make(map[string]mdata.Counter), factory: factory}
}

func (s *SimpleCountStorage) Set(m mdata.Counter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	el, ok := s.storage[m.GetName()]
	if !ok {
		s.storage[m.GetName()] = m
		return nil
	}
	newVal := el.GetValue() + m.GetValue()
	s.storage[m.GetName()] = s.factory(m.GetName(), newVal)
	return nil
}

func (s *SimpleCountStorage) Get(name string) (mdata.Counter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val := s.storage[name]
	return val, nil
}

func (s *SimpleCountStorage) GetList() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make(map[string]string, len(s.storage))
	for k, v := range s.storage {
		res[k] = strconv.Itoa(int(v.GetValue()))
	}
	return res
}

func (s *SimpleCountStorage) GetMetrics() []mdata.Metrics {
	md := make([]mdata.Metrics, len(s.storage))
	i := 0
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, counter := range s.storage {
		value := counter.GetValue()
		md[i] = mdata.Metrics{
			ID:    counter.GetName(),
			MType: counter.GetType(),
			Delta: &value,
		}
		i++
	}
	return md
}

func (s *SimpleCountStorage) SetFrom(metrics []mdata.Metrics) error {
	factory := mdata.NewSimpleCounter
	for _, m := range metrics {
		if m.MType != mdata.COUNTER {
			continue
		}
		err := s.Set(factory(m.ID, *m.Delta))
		if err != nil {
			return err
		}
	}
	return nil
}
