package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func BuildAuditLogsForPostgres(conn *sqlx.DB, logger zerolog.Logger) (*PostgresAuditLogs) {
	return &PostgresAuditLogs{
		db:     conn,
		logger: logger,
	}
}
