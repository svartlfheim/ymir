package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func BuildModulesForPostgres(conn *sqlx.DB, logger zerolog.Logger) (*PostgresModules) {
	return &PostgresModules{
		db:     conn,
		logger: logger,
	}
}
