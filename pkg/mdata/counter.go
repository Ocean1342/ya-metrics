package mdata

type Counter interface {
	GetValue() int64
	GetName() string
	GetType() string
}
type SimpleCounter struct {
	value    int64
	name     string
	typeName string
}

func NewSimpleCounter(name string, value int64) Counter {
	return &SimpleCounter{
		value:    value,
		name:     name,
		typeName: COUNTER,
	}
}

func (s *SimpleCounter) GetValue() int64 {
	return s.value
}

func (s *SimpleCounter) GetName() string {
	return s.name
}

func (s *SimpleCounter) GetType() string { return s.typeName }
