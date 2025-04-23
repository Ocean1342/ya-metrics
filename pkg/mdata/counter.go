package mdata

type Counter interface {
	GetValue() int64
	GetName() string
}
type SimpleCounter struct {
	Value int64
	Name  string
}

func NewSimpleCounter(name string, value int64) Counter {
	return &SimpleCounter{
		Value: value,
		Name:  name,
	}
}

func (s *SimpleCounter) GetValue() int64 {
	return s.Value
}
func (s *SimpleCounter) GetName() string {
	return s.Name
}
