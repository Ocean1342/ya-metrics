package mdata

import "fmt"

type Gauge interface {
	GetValue() float64
	GetName() string
	GetType() string
}

type GaugeMetric struct {
	name     string
	value    float64
	typeName string
}

// TODO: переписать на фабрику NewSimpleGauge ?
type GaugeFactory interface {
	NewGauge(name string, value float64) Gauge
}

func NewSimpleGauge(name string, value float64) *GaugeMetric {
	metric := GaugeMetric{
		name:     name,
		value:    value,
		typeName: GAUGE,
	}
	fmt.Println(metric)
	return &metric
}

func (g *GaugeMetric) GetValue() float64 {
	return g.value
}

func (g *GaugeMetric) GetName() string {
	return g.name
}

func (g *GaugeMetric) GetType() string {
	return g.typeName
}
