package mdata

type Gauge interface {
	GetValue() float64
	GetName() string
}

type SimpleGauge struct {
	name  string
	value float64
}

func NewSimpleGauge(name string, value float64) *SimpleGauge {
	return &SimpleGauge{
		name:  name,
		value: value,
	}
}

func (g *SimpleGauge) GetValue() float64 {
	return g.value
}

func (g *SimpleGauge) GetName() string {
	return g.name
}
