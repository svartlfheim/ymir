package server

import (
	"time"

	"github.com/svartlfheim/ymir/internal/registry"
)

type AuditLogRepository interface {
	Save(action string, respStatus registry.RegistryHandlerStatus, occurred_at time.Time, meta map[string]interface{}) error
}
