package repository

// import (
// 	"bytes"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/svartlfheim/ymir/internal/registry"
// 	ymirtest "github.com/svartlfheim/ymir/test"
// 	ymirtestdb "github.com/svartlfheim/ymir/test/db"
// 	ymirstubs "github.com/svartlfheim/ymir/test/stubs"
// )

// func Test_PostgresAuditLogs_Save(t *testing.T) {
// 	ymirtest.RunTestWithPreparedDb(ymirtestdb.PostgresDbOptions{}, t, func(tt *testing.T, dbCfg ymirtestdb.PostgresTestDb) {
// 		b := new(bytes.Buffer)
// 		l := ymirstubs.BuildZerologLogger(b)
// 		conn := ymirtestdb.DbConnFromPostgresDb(tt, dbCfg)

// 		repo := PostgresAuditLogs{
// 			db: conn,
// 			logger: l,
// 		}

// 		occurred, err := time.Parse(time.RFC3339, "2021-10-09T13:00:00Z")

// 		if err != nil {
// 			t.Log(err.Error())
// 			t.FailNow()
// 		}

// 		err = repo.Save("test-action", registry.STATUS_OKAY, occurred, map[string]interface{}{
// 			"myfield": "a value",
// 			"somelist": []string{
// 				"blah",
// 				"meh",
// 			},
// 		})

// 		assert.Nil(tt, err)

// 		rows, err := conn.Queryx(fmt.Sprintf("SELECT id, action, response_status, occurred_at, meta FROM %s", AuditLogsTableName))

// 		totalFound := 0
// 		for rows.Next() {
// 			totalFound++
// 			auditLog := postgresDbAuditLog{}
// 			err := rows.StructScan(&auditLog)
// 			if err != nil {
// 				panic(err)
// 			}

// 			parsedUuid, parseErr := uuid.Parse(auditLog.Id)

// 			assert.Nil(tt, parseErr)
// 			assert.IsType(tt, uuid.UUID{}, parsedUuid)
// 			assert.Equal(tt, auditLog.Action, "test-action")
// 			assert.Equal(tt, auditLog.ResponseStatus, string(registry.STATUS_OKAY))
// 			assert.Equal(tt, auditLog.OccurredAt, "2021-10-09T13:00:00Z")
// 			assert.Equal(tt, auditLog.Meta, "{\"myfield\": \"a value\", \"somelist\": [\"blah\", \"meh\"]}")
// 		}

// 		assert.Equal(tt, 1, totalFound)
// 	})
// }