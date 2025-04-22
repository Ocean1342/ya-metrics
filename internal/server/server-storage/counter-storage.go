package server_storage

import (
	"strings"
	"ya-metrics/pkg/mdata"
)

type SimpleStorage struct {
	storage map[string]mdata.Counter
}

func NewSimpleCountStorage() *SimpleStorage {
	return &SimpleStorage{storage: make(map[string]mdata.Counter)}
}

func (s *SimpleStorage) Set(m mdata.Counter) error {
	s.storage[m.GetName()] = m
	//TODO: бизнес логика сохранения?
	return nil
}

func (s *SimpleStorage) Get(name string) (mdata.Counter, error) {
	return s.storage[strings.ToLower(name)], nil
}
