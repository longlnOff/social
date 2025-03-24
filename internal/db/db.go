package db

import (
	"database/sql"
	"time"
)

func New(connection string, maxOpenConns int, maxIdleConns int, maxIdleTime int) (*sql.DB, error) {
	db, err := sql.Open("postgres", connection)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Minute)

}
