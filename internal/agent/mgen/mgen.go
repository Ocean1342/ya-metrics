package mgen

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
	"math/rand"
	"runtime"
	"sync"
	"time"
	"ya-metrics/pkg/mdata"
)

func GenerateGaugeMetrics(logger *zap.SugaredLogger) <-chan mdata.Gauge {
	ch := make(chan mdata.Gauge, 100)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		v, _ := mem.VirtualMemory()
		ch <- mdata.NewSimpleGauge("TotalMemory", float64(v.Total))
		ch <- mdata.NewSimpleGauge("FreeMemory", float64(v.Free))
		percent, err := cpu.Percent(time.Second, false)
		if err != nil {
			logger.Errorf("could not get CPUutilization1 err:%s", err)
			return
		}
		ch <- mdata.NewSimpleGauge("CPUutilization1", float64(percent[0]))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		gaugeMetrics := []mdata.Gauge{
			mdata.NewSimpleGauge("Alloc", float64(memStats.Alloc)),
			mdata.NewSimpleGauge("Mallocs", float64(memStats.Mallocs)),
			mdata.NewSimpleGauge("Sys", float64(memStats.Sys)),
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
		for _, g := range gaugeMetrics {
			ch <- g
		}
	}()
	wg.Wait()
	close(ch)
	return ch
}
