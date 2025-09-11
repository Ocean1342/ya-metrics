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
		gaugeStorage: srvrstrg.NewSimpleGaugeStorage(log),
		gaugeFactory: mdata.NewSimpleGauge,
		log:          log,
	}
}

func (s *GaugeDBStorage) Get(n string) mdata.Gauge {
	var name string
	var value float64

	row := s.db.QueryRow("SELECT id, value FROM metrics WHERE mtype=$1 AND id = $2", mdata.GAUGE, n)
	err := row.Scan(&name, &value)
	if err != nil {
		return nil
	}
	return s.gaugeFactory(name, value)
}

func (s *GaugeDBStorage) Set(m mdata.Gauge) error {
	_, err := s.db.Exec(
		"INSERT INTO metrics (id, mtype, delta, value) VALUES ($1,$2,$3,$4)"+
			"ON CONFLICT (id) DO UPDATE SET delta=EXCLUDED.delta, value=EXCLUDED.value",
		m.GetName(), m.GetType(), nil, m.GetValue(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *GaugeDBStorage) GetList() map[string]string {
	g := make(map[string]string)
	rows, err := s.db.Query("SELECT id,value FROM metrics WHERE mtype=$1", mdata.GAUGE)
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
	var countRows int
	row := s.db.QueryRow("SELECT count(*) FROM metrics WHERE mtype=$1", mdata.GAUGE)
	err := row.Scan(&countRows)
	if err != nil {
		s.log.Errorf("could not get count rows for counter data")
		return nil
	}
	md := make([]mdata.Metrics, countRows)
	rows, err := s.db.Query("SELECT id,value FROM metrics WHERE mtype=$1", mdata.GAUGE)
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
