package database

import (
	"database/sql"
	"go.uber.org/zap"
	"strconv"
	srvrstrg "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type CounterDBStorage struct {
	db             *sql.DB
	counterStorage srvrstrg.CounterStorage
	log            *zap.SugaredLogger
}

func NewCounter(db *sql.DB, log *zap.SugaredLogger) *CounterDBStorage {
	return &CounterDBStorage{db: db, counterStorage: srvrstrg.NewSimpleCountStorage(mdata.NewSimpleCounter), log: log}
}

func (s *CounterDBStorage) Set(m mdata.Counter) error {
	s.db.Exec(
		"INSERT INTO metrics (id, mtype, delta, value) VALUES ($1,$2,$3,$4)",
		m.GetName(), m.GetType(), m.GetValue(), nil,
	)
	return nil
}

func (s *CounterDBStorage) Get(n string) (mdata.Counter, error) {
	var g mdata.SimpleCounter
	row := s.db.QueryRow("SELECT * FROM metrics WHERE mtype=$1 AND id = $2", mdata.COUNTER, n)
	err := row.Scan(&g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *CounterDBStorage) GetList() map[string]string {
	g := make(map[string]string)
	rows, err := s.db.Query("SELECT name,value FROM metrics WHERE mtype=$1", mdata.COUNTER)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var value float64
		err = rows.Scan(&name, &value)
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
func (s *CounterDBStorage) GetMetrics() []mdata.Metrics {
	md := make([]mdata.Metrics, 1_000_000)
	rows, err := s.db.Query("SELECT name,value FROM metrics WHERE mtype=$1", mdata.COUNTER)
	if err != nil {
		return nil
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var name string
		var delta int64
		err = rows.Scan(&name, &delta)
		if err != nil {
			//TODO: логер?
			continue
		}
		md[i] = mdata.Metrics{
			ID:    name,
			MType: mdata.COUNTER,
			Delta: &delta,
		}
		i++
	}
	if err := rows.Err(); err != nil {
		s.log.Errorf("err on scan rows: %s", err)
	}
	return md
}

func (s *CounterDBStorage) SetFrom(metrics []mdata.Metrics) error {
	return s.counterStorage.SetFrom(metrics)
}
