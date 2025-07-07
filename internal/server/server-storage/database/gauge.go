package database

import (
	"database/sql"
	"go.uber.org/zap"
	"strconv"
	srvrstrg "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type GaugeDBStorage struct {
	db           *sql.DB
	gaugeStorage srvrstrg.GaugeStorage
	gaugeFactory func(name string, value float64) *mdata.GaugeMetric
	log          *zap.SugaredLogger
}

func NewGauge(db *sql.DB, log *zap.SugaredLogger) *GaugeDBStorage {
	return &GaugeDBStorage{db: db,
		gaugeStorage: srvrstrg.NewSimpleGaugeStorage(),
		gaugeFactory: mdata.NewSimpleGauge,
		log:          log,
	}
}

func (s *GaugeDBStorage) Get(n string) mdata.Gauge {
	var name string
	var value float64

	row := s.db.QueryRow("SELECT id,value FROM metrics WHERE mtype=$1 AND id = $2", mdata.GAUGE, n)
	err := row.Scan(&name, &value)
	if err != nil {
		return nil
	}
	return s.gaugeFactory(name, value)
}

func (s *GaugeDBStorage) Set(m mdata.Gauge) error {
	s.db.Exec(
		"INSERT INTO metrics (id, mtype, delta, value) VALUES ($1,$2,$3,$4)",
		m.GetName(), m.GetType(), nil, m.GetValue(),
	)
	return nil
}

func (s *GaugeDBStorage) GetList() map[string]string {
	g := make(map[string]string)
	rows, err := s.db.Query("SELECT name,value FROM metrics WHERE mtype=$1", mdata.GAUGE)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var value float64
		err = rows.Scan(name, value)
		if err != nil {
			return nil
		}
		g[name] = strconv.Itoa(int(value))
	}
	if err := rows.Err(); err != nil {
		s.log.Errorf("err on scan rows: %s", err)
	}
	return g
}
func (s *GaugeDBStorage) GetMetrics() []mdata.Metrics {
	md := make([]mdata.Metrics, 1_000_000)
	rows, err := s.db.Query("SELECT name,value FROM metrics WHERE mtype=$1", mdata.GAUGE)
	if err != nil {
		return nil
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var name string
		var value float64
		md[i] = mdata.Metrics{
			ID:    name,
			MType: mdata.GAUGE,
			Value: &value,
		}
		i++
	}
	if err := rows.Err(); err != nil {
		s.log.Errorf("err on scan rows: %s", err)
	}
	return md
}

func (s *GaugeDBStorage) SetFrom(metrics []mdata.Metrics) error {
	return s.gaugeStorage.SetFrom(metrics)
}
