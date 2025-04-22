package mdata

type Counter interface {
	GetValue() int64
	GetName() string
}
type SimpleCounter struct {
	value int64
	name  string
}

func InitSimpleCounter(name string, value int64) *SimpleCounter {
	return &SimpleCounter{
		value: value,
		name:  name,
	}
}

func (s *SimpleCounter) GetValue() int64 {
	return s.value
}
func (s *SimpleCounter) GetName() string {
	return s.name
}
