package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/ymir/internal/registry"
)

type postgresDbAuditLog struct {
	Id             string `db:"id"`
	Action         string `db:"action"`
	ResponseStatus string `db:"response_status"`
	OccurredAt     string `db:"occurred_at"`
	Meta           string `db:"meta"` //it's JSONB
}

type PostgresAuditLogs struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

func (s *PostgresAuditLogs) startTransaction() (*sqlx.Tx, error) {
	tx, err := s.db.Beginx()

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, ErrDbTransaction{
			Wrapped: err,
		}
	}

	return tx, nil
}

func (s *PostgresAuditLogs) Save(action string, respStatus registry.RegistryHandlerStatus, occurred_at time.Time, meta map[string]interface{}) error {
	id := uuid.New().String()

	jsonMetaVal, err := json.Marshal(&meta)

	if err != nil {
		return err
	}

	aL := postgresDbAuditLog{
		Id:             id,
		Action:         action,
		ResponseStatus: string(respStatus),
		OccurredAt:     occurred_at.UTC().Format(time.RFC3339),
		Meta:           string(jsonMetaVal),
	}

	tx, err := s.startTransaction()

	if err != nil {
		return err
	}

	insert := fmt.Sprintf(`
INSERT INTO %s (id, action, response_status, occurred_at, meta) VALUES (:id, :action, :response_status, :occurred_at, :meta);`,
		AuditLogsTableName)

	_, err = tx.NamedExec(insert, aL)

	if err != nil {
		return wrapTransactionError(err)
	}

	if err := tx.Commit(); err != nil {
		rollbackErr := tx.Rollback()

		if rollbackErr != nil {
			return wrapTransactionError(rollbackErr)
		}

		return wrapTransactionError(err)
	}

	return nil
}
