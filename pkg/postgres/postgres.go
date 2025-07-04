package postgres

import "database/sql"

func New(url string) (*sql.DB, error) {
	//defer db.Close
	db, err := sql.Open("pgx", url)
	if err != nil {
		panic(err)
	}
	return db, nil
}
