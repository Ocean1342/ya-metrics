package mdata

type Gauge interface {
	GetValue() float64
	GetName() string
	GetType() string
}

type SimpleGauge struct {
	name     string
	value    float64
	typeName string
}

// TODO: переписать на фабрику NewSimpleGauge ?
type GaugeFactory interface {
	NewGauge(name string, value float64) Gauge
}

func NewSimpleGauge(name string, value float64) *SimpleGauge {
	return &SimpleGauge{
		name:     name,
		value:    value,
		typeName: GAUGE,
	}
}

func (g *SimpleGauge) GetValue() float64 {
	return g.value
}

func (g *SimpleGauge) GetName() string {
	return g.name
}

func (g *SimpleGauge) GetType() string {
	return g.typeName
}
