package registry

import (
	"time"

	"github.com/rs/zerolog"
)

type AuditableAction interface {
	GetActionName() string
	GetTimeOfOccurrence() time.Time
	GetResponseStatus() RegistryHandlerStatus
	GetAuditMeta() map[string]interface{}
}

type auditLogsRepo interface {
	Save(action string, respStatus RegistryHandlerStatus, occurred_at time.Time, meta map[string]interface{}) error
}

type Auditor struct {
	repo   auditLogsRepo
	logger zerolog.Logger
}

func (a *Auditor) Record(action AuditableAction) {

	err := a.repo.Save(action.GetActionName(), action.GetResponseStatus(), action.GetTimeOfOccurrence(), action.GetAuditMeta())

	if err != nil {
		a.logger.Error().Err(err).Msg("error during audit log save process")
	}
}

func NewAuditor(r auditLogsRepo, l zerolog.Logger) *Auditor {
	return &Auditor{
		repo:   r,
		logger: l,
	}
}
