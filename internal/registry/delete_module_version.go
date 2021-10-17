package registry

import (
	"time"

	"github.com/rs/zerolog"
)

type deleteModuleVersionRepository interface {
	VersionById(string) (m ModuleVersion, err error)
	DeleteModuleVersion(ModuleVersion) error
}

type deleteModuleVersionV1CommandValidator interface {
	Validate(cmd interface{}) []ValidationError
}

type deleteModuleVersionV1Command struct {
	DTO DeleteModuleVersionV1DTO
}

type DeleteModuleVersionV1Response struct {
	occurredAt       time.Time
	Status           RegistryHandlerStatus
	ModuleVersion    ModuleVersion
	ValidationErrors []ValidationError
}

func (r DeleteModuleVersionV1Response) GetActionName() string {
	return "v1.modules.versions.delete"
}

func (r DeleteModuleVersionV1Response) GetTimeOfOccurrence() time.Time {
	return r.occurredAt
}

func (r DeleteModuleVersionV1Response) GetResponseStatus() RegistryHandlerStatus {
	return r.Status
}

func (r DeleteModuleVersionV1Response) GetAuditMeta() map[string]interface{} {
	return map[string]interface{}{
		"module_version_id": r.ModuleVersion.Id,
		"validation_errors": r.ValidationErrors,
	}
}

type DeleteModuleVersionV1DTO struct {
	Id string `validate:"required,uuid"`
}

func (dto DeleteModuleVersionV1DTO) validate(r deleteModuleVersionRepository, v deleteModuleVersionV1CommandValidator) []ValidationError {
	return v.Validate(dto)
}

func (cmd deleteModuleVersionV1Command) handle(r deleteModuleVersionRepository, logger zerolog.Logger, v deleteModuleVersionV1CommandValidator) (DeleteModuleVersionV1Response, error) {
	occurred := time.Now().UTC()
	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return DeleteModuleVersionV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	mv, err := r.VersionById(cmd.DTO.Id)

	if err != nil {
		if _, ok := err.(ErrResourceNotFound); ok {
			return DeleteModuleVersionV1Response{
				occurredAt: occurred,
				Status:     STATUS_NOT_FOUND,
			}, nil
		}

		logger.Error().Err(err).Str("id", cmd.DTO.Id).Msg("failed to find module version")

		return DeleteModuleVersionV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	err = r.DeleteModuleVersion(mv)

	if err != nil {
		logger.Error().Err(err).Str("id", cmd.DTO.Id).Msg("failed to delete module version")

		return DeleteModuleVersionV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	return DeleteModuleVersionV1Response{
		occurredAt:    occurred,
		Status:        STATUS_OKAY,
		ModuleVersion: mv,
	}, err
}
