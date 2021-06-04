package repository

import (
	"bytes"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	ymirstubs "github.com/svartlfheim/ymir/test/stubs"
)

func TestBuildAuditLogsForPostgres(t *testing.T) {
	db := &sqlx.DB{}
	b := new(bytes.Buffer)
	l := ymirstubs.BuildZerologLogger(b)

	repo := BuildAuditLogsForPostgres(db, l)

	assert.IsType(t, &PostgresAuditLogs{}, repo)
}