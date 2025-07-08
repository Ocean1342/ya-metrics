package postgres

import "database/sql"

func New(url string) (*sql.DB, error) {
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
		return nil, err
	}

	return db, nil
}
