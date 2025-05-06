package mdata

import "strings"

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

type AvailableMetricsTypes interface {
	Isset(typeName string) bool
}

type YaMetricsTypes struct {
	list map[int]string
}

// TODO: вынести определение метрик на уровень cmd
func InitMetrics() AvailableMetricsTypes {
	list := make(map[int]string, 2)
	list[0] = GAUGE
	list[1] = COUNTER
	return &YaMetricsTypes{list: list}
}
func (ym YaMetricsTypes) Isset(typeName string) bool {
	for _, m := range ym.list {
		if strings.EqualFold(m, typeName) {
			return true
		}
	}
	return false
}
