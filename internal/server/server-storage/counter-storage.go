package srvrstrg

import (
	"strconv"
	"strings"
	"ya-metrics/pkg/mdata"
)

type SimpleCountStorage struct {
	storage map[string]mdata.Counter
	factory func(name string, value int64) mdata.Counter
}

func NewSimpleCountStorage(factory func(name string, value int64) mdata.Counter) *SimpleCountStorage {
	return &SimpleCountStorage{storage: make(map[string]mdata.Counter), factory: factory}
}

func (s *SimpleCountStorage) Set(m mdata.Counter) error {
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
	return s.storage[strings.ToLower(name)], nil
}

func (s *SimpleCountStorage) GetList() map[string]string {
	res := make(map[string]string, len(s.storage))
	for k, v := range s.storage {
		res[k] = strconv.Itoa(int(v.GetValue()))
	}
	return res
}
