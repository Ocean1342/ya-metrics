package mgen

import (
	"math/rand"
	"runtime"
	"ya-metrics/pkg/mdata"
)

func GenerateGaugeMetrics() []mdata.Gauge {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return []mdata.Gauge{
		mdata.NewSimpleGauge("Alloc", float64(memStats.Alloc)),
		mdata.NewSimpleGauge("BuckHashSys", float64(memStats.BuckHashSys)),
		mdata.NewSimpleGauge("Frees", float64(memStats.Frees)),
		mdata.NewSimpleGauge("GCCPUFraction", memStats.GCCPUFraction),
		mdata.NewSimpleGauge("GCSys", float64(memStats.GCSys)),
		mdata.NewSimpleGauge("HeapAlloc", float64(memStats.HeapAlloc)),
		mdata.NewSimpleGauge("HeapIdle", float64(memStats.HeapIdle)),
		mdata.NewSimpleGauge("HeapInuse", float64(memStats.HeapInuse)),
		mdata.NewSimpleGauge("HeapObjects", float64(memStats.HeapObjects)),
		mdata.NewSimpleGauge("HeapReleased", float64(memStats.HeapReleased)),
		mdata.NewSimpleGauge("HeapSys", float64(memStats.HeapSys)),
		mdata.NewSimpleGauge("LastGC", float64(memStats.LastGC)),
		mdata.NewSimpleGauge("Lookups", float64(memStats.Lookups)),
		mdata.NewSimpleGauge("MCacheInuse", float64(memStats.MCacheInuse)),
		mdata.NewSimpleGauge("MCacheSys", float64(memStats.MCacheSys)),
		mdata.NewSimpleGauge("MSpanInuse", float64(memStats.MSpanInuse)),
		mdata.NewSimpleGauge("MSpanSys", float64(memStats.MSpanSys)),
		mdata.NewSimpleGauge("NextGC", float64(memStats.NextGC)),
		mdata.NewSimpleGauge("NumForcedGC", float64(memStats.NumForcedGC)),
		mdata.NewSimpleGauge("NumGC", float64(memStats.NumGC)),
		mdata.NewSimpleGauge("OtherSys", float64(memStats.OtherSys)),
		mdata.NewSimpleGauge("PauseTotalNs", float64(memStats.PauseTotalNs)),
		mdata.NewSimpleGauge("StackInuse", float64(memStats.StackInuse)),
		mdata.NewSimpleGauge("StackSys", float64(memStats.StackSys)),
		mdata.NewSimpleGauge("TotalAlloc", float64(memStats.TotalAlloc)),
		mdata.NewSimpleGauge("RandomValue", rand.Float64()),
	}
}
