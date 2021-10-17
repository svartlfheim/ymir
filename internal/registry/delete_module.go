package registry

import (
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
)

type deleteModuleRepository interface {
	ByFQN(ModuleFQN) (m Module, err error)
	ById(string) (m Module, err error)
	VersionsByModuleFQN(fqn ModuleFQN, chunkOpts ChunkingOptions) (m []ModuleVersion, err error)
	VersionsByModule(moduleId string, chunkOpts ChunkingOptions) (m []ModuleVersion, err error)
	DeleteModule(mod Module) error
	DeleteVersionsForModule(mod Module) error
}

type deleteModuleV1CommandValidator interface {
	NoVersionsExistForModuleId(r noVersionsExistForIdRepository, sl validator.StructLevel, id string)
	NoVersionsExistForModuleFQN(r noVersionsExistForModuleFQNRepository, sl validator.StructLevel, fqn ModuleFQN)
	RegisterStructLevelValidator(validator.StructLevelFunc, interface{})
	Validate(cmd interface{}) []ValidationError
	SetCustomMessageHandlers(handlers map[string]CustomValidationMessageHandler)
}

type deleteModuleV1Command struct {
	DTO DeleteModuleV1DTO
}

type DeleteModuleV1Response struct {
	occurredAt       time.Time
	attemptedFor     string
	Status           RegistryHandlerStatus
	Module           Module
	ValidationErrors []ValidationError
}

func (r DeleteModuleV1Response) GetActionName() string {
	return "v1.modules.delete"
}

func (r DeleteModuleV1Response) GetTimeOfOccurrence() time.Time {
	return r.occurredAt
}

func (r DeleteModuleV1Response) GetResponseStatus() RegistryHandlerStatus {
	return r.Status
}

func (r DeleteModuleV1Response) GetAuditMeta() map[string]interface{} {
	return map[string]interface{}{
		"module_id":         r.Module.Id,
		"validation_errors": r.ValidationErrors,
	}
}

type DeleteModuleV1DTO struct {
	Id             string `validate:"required,uuid"`
	DeleteVersions bool   `json:"delete_versions"`
}

func (dto DeleteModuleV1DTO) validate(r deleteModuleRepository, v deleteModuleV1CommandValidator) []ValidationError {
	v.RegisterStructLevelValidator(func(sl validator.StructLevel) {
		if dto.DeleteVersions {
			return
		}

		v.NoVersionsExistForModuleId(r, sl, dto.Id)
	}, AddModuleV1DTO{})

	return v.Validate(dto)
}

func (cmd deleteModuleV1Command) handle(r deleteModuleRepository, logger zerolog.Logger, v deleteModuleV1CommandValidator) (DeleteModuleV1Response, error) {
	occurred := time.Now().UTC()

	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return DeleteModuleV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	m, err := r.ById(cmd.DTO.Id)

	if err != nil {
		if _, ok := err.(ErrResourceNotFound); ok {
			return DeleteModuleV1Response{
				occurredAt: occurred,
				Status:     STATUS_NOT_FOUND,
			}, nil
		}

		logger.Error().Err(err).Str("id", cmd.DTO.Id).Msg("failed to add module to store")

		return DeleteModuleV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	if cmd.DTO.DeleteVersions {
		if err := r.DeleteVersionsForModule(m); err != nil {
			return DeleteModuleV1Response{
				occurredAt: occurred,
				Status:     STATUS_INTERNAL_ERROR,
			}, err
		}
	}

	if err := r.DeleteModule(m); err != nil {
		return DeleteModuleV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	return DeleteModuleV1Response{
		occurredAt: occurred,
		Status:     STATUS_OKAY,
		Module:     m,
	}, err
}
