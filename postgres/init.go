package postgres

import (
	"github.com/jmoiron/sqlx"
)

func GetPostgresConnection() (*sqlx.DB, error) {
	//psqlInfo := "postgres://achekslime:iZ2jH3uxeqkR@ep-silent-hat-022245.eu-central-1.aws.neon.tech/flixray"

	connStr := "user=achekslime password=iZ2jH3uxeqkR dbname=flixray host=ep-silent-hat-022245.eu-central-1.aws.neon.tech sslmode=require"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
