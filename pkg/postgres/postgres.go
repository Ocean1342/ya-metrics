package postgres

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"time"
)

var ErrConnection = errors.New("txt")

// var retryTimes = []int{1, 3, 5}
var retryTimes = []int{1, 1}

func New(url string, log *zap.SugaredLogger) (*sql.DB, error) {
	var res *sql.DB
	var err error
	for _, sleep := range retryTimes {
		res, err = repeatableNew(url)
		if err != nil {
			if errors.Is(err, ErrConnection) {
				log.Infof("detected repeatable error:`%s`. sleep for:%d", err, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				continue
			}
		}
		return res, nil
	}
	if err != nil {
		return nil, err
	}
	return res, err
}

func repeatableNew(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	//
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS metrics (
    id VARCHAR(255) NOT NULL UNIQUE,
    mtype VARCHAR(50) NOT NULL,
    delta BIGINT NULL,
    value DOUBLE PRECISION NULL)`)
	if err != nil {

		var pgErr *pgconn.ConnectError
		//TODO:почему с As() работает, а с Is нет?
		if errors.As(err, &pgErr) {
			err = ErrConnection
		}
		return nil, err
	}

	return db, nil
}
