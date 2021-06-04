package registry

import (
	"time"

	"github.com/rs/zerolog"
)

type showModuleRepository interface {
	ById(string) (m Module, err error)
}

type showModuleV1CommandValidator interface {
	Validate(cmd interface{}) []ValidationError
}

type showModuleV1Command struct {
	DTO ShowModuleV1DTO
}

type ShowModuleV1Response struct {
	attemptedSearch  string
	occurredAt       time.Time
	Status           RegistryHandlerStatus
	Module           Module
	ValidationErrors []ValidationError
}

func (r ShowModuleV1Response) GetActionName() string {
	return "v1.modules.show"
}

func (r ShowModuleV1Response) GetTimeOfOccurrence() time.Time {
	return r.occurredAt
}

func (r ShowModuleV1Response) GetResponseStatus() RegistryHandlerStatus {
	return r.Status
}

func (r ShowModuleV1Response) GetAuditMeta() map[string]interface{} {
	return map[string]interface{}{
		"module_id":         r.Module.Id,
		"validation_errors": r.ValidationErrors,
		"searched_for":      r.attemptedSearch,
	}
}

type ShowModuleV1DTO struct {
	Id string `validate:"required,uuid"`
}

func (dto ShowModuleV1DTO) validate(r showModuleRepository, v showModuleV1CommandValidator) []ValidationError {
	return v.Validate(dto)
}

func (cmd showModuleV1Command) handle(r showModuleRepository, logger zerolog.Logger, v showModuleV1CommandValidator) (ShowModuleV1Response, error) {
	occurred := time.Now().UTC()
	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return ShowModuleV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	m, err := r.ById(cmd.DTO.Id)

	if err != nil {
		if _, ok := err.(ErrResourceNotFound); ok {
			return ShowModuleV1Response{
				Status: STATUS_NOT_FOUND,
			}, nil
		}

		logger.Error().Err(err).Str("id", cmd.DTO.Id).Msg("failed to add module to store")

		return ShowModuleV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	return ShowModuleV1Response{
		occurredAt: occurred,
		Status:     STATUS_OKAY,
		Module:     m,
	}, err
}
