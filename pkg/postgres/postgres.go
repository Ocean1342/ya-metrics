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
    id VARCHAR(255) NOT NULL,
    mtype VARCHAR(50) NOT NULL,
    delta BIGINT NULL,
    value DOUBLE PRECISION NULL,
    PRIMARY KEY (id, mtype))`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
