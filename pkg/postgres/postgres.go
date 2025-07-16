package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"time"
)

var ErrRetrybleConnection = errors.New("error retryable")

var retries = 3

func New(url string, log *zap.SugaredLogger) (*sql.DB, error) {
	var res *sql.DB
	var err error
	for i := 1; i <= retries; i++ {
		sleep := 2*i - 1
		res, err = repeatableNew(url)
		if err != nil {
			if errors.Is(err, ErrRetrybleConnection) {
				log.Infof("detected retryble error:`%s`. sleep for:%d", err, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				continue
			}
		}
		if res == nil {
			log.Infof("could not start db. sleep for:%d", sleep)
			time.Sleep(time.Duration(sleep) * time.Second)
			continue
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
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS metrics (
    id VARCHAR(255) NOT NULL UNIQUE,
    mtype VARCHAR(50) NOT NULL,
    delta BIGINT NULL,
    value DOUBLE PRECISION NULL)`)
	if err != nil {
		var connectError *pgconn.ConnectError
		var pgError *pgconn.PgError
		if errors.As(err, &connectError) {
			return nil, fmt.Errorf("retraible error:%w", ErrRetrybleConnection)
		}
		if errors.As(err, &pgError) {
			if pgerrcode.IsConnectionException(pgError.Code) {
				return nil, fmt.Errorf("retraible error:%w", pgError)
			}
		}
		return nil, err
	}

	return db, nil
}
