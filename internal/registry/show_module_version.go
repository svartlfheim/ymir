package registry

import (
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
)

type showModuleVersionRepository interface {
	VersionByFQN(ModuleVersionFQN) (m ModuleVersion, err error)
	VersionById(string) (m ModuleVersion, err error)
}

type showModuleVersionV1CommandValidator interface {
	RegisterStructLevelValidator(validator.StructLevelFunc, interface{})
	Validate(cmd interface{}) []ValidationError
	SetCustomMessageHandlers(handlers map[string]CustomValidationMessageHandler)
}

type showModuleVersionV1Command struct {
	DTO ShowModuleVersionV1DTO
}

type ShowModuleVersionV1Response struct {
	occurredAt       time.Time
	Status           RegistryHandlerStatus
	ModuleVersion    ModuleVersion
	ValidationErrors []ValidationError
}

func (r ShowModuleVersionV1Response) GetActionName() string {
	return "v1.modules.versions.show"
}

func (r ShowModuleVersionV1Response) GetTimeOfOccurrence() time.Time {
	return r.occurredAt
}

func (r ShowModuleVersionV1Response) GetResponseStatus() RegistryHandlerStatus {
	return r.Status
}

func (r ShowModuleVersionV1Response) GetAuditMeta() map[string]interface{} {
	return map[string]interface{}{
		"module_version_id": r.ModuleVersion.Id,
		"validation_errors": r.ValidationErrors,
	}
}

type ShowModuleVersionV1DTO struct {
	Id string `validate:"required,uuid"`
}

func (dto ShowModuleVersionV1DTO) validate(r showModuleVersionRepository, v showModuleVersionV1CommandValidator) []ValidationError {
	return v.Validate(dto)
}

func (cmd showModuleVersionV1Command) handle(r showModuleVersionRepository, logger zerolog.Logger, v showModuleVersionV1CommandValidator) (ShowModuleVersionV1Response, error) {
	occurred := time.Now().UTC()

	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return ShowModuleVersionV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	mv, err := r.VersionById(cmd.DTO.Id)

	if err != nil {
		if _, ok := err.(ErrResourceNotFound); ok {
			return ShowModuleVersionV1Response{
				occurredAt: occurred,
				Status:     STATUS_NOT_FOUND,
			}, nil
		}

		logger.Error().Err(err).Str("id", cmd.DTO.Id).Msg("failed to add module to store")

		return ShowModuleVersionV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	return ShowModuleVersionV1Response{
		occurredAt:    occurred,
		Status:        STATUS_OKAY,
		ModuleVersion: mv,
	}, err
}
