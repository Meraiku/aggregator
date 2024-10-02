package sql

import (
	"database/sql"

	"github.com/uptrace/bun/driver/pgdriver"
)

func ConnectPostgres(dsn string) (*Queries, error) {

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	db := New(sqlDB)
	return db, nil
}
