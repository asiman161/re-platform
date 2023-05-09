package bootstrap

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func InitDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "can't init DB")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "can't ping DB")
	}

	return db, nil
}
