package registry

import (
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
)

type deleteModuleByFqnRepository interface {
	ByFQN(ModuleFQN) (m Module, err error)
	VersionsByModuleFQN(fqn ModuleFQN, chunkOpts ChunkingOptions) (m []ModuleVersion, err error)
	DeleteModule(mod Module) error
	DeleteVersionsForModule(mod Module) error
}

type deleteModuleByFqnV1CommandValidator interface {
	NoVersionsExistForModuleFQN(r noVersionsExistForModuleFQNRepository, sl validator.StructLevel, fqn ModuleFQN)
	RegisterStructLevelValidator(validator.StructLevelFunc, interface{})
	Validate(cmd interface{}) []ValidationError
	SetCustomMessageHandlers(handlers map[string]CustomValidationMessageHandler)
}

type deleteModuleV1ByFqnCommand struct {
	DTO DeleteModuleV1ByFqnDTO
}

type DeleteModuleV1ByFqnDTO struct {
	FQN            ModuleFQN `validate:"required"`
	DeleteVersions bool      `json:"delete_versions"`
}

func (dto DeleteModuleV1ByFqnDTO) validate(r deleteModuleByFqnRepository, v deleteModuleByFqnV1CommandValidator) []ValidationError {
	v.RegisterStructLevelValidator(func(sl validator.StructLevel) {
		if dto.DeleteVersions {
			return
		}

		v.NoVersionsExistForModuleFQN(r, sl, dto.FQN)
	}, DeleteModuleV1ByFqnDTO{})

	return v.Validate(dto)
}

func (cmd deleteModuleV1ByFqnCommand) handle(r deleteModuleByFqnRepository, logger zerolog.Logger, v deleteModuleByFqnV1CommandValidator) (DeleteModuleV1Response, error) {
	occurred := time.Now().UTC()

	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return DeleteModuleV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	m, err := r.ByFQN(cmd.DTO.FQN)

	if err != nil {
		if _, ok := err.(ErrResourceNotFound); ok {
			return DeleteModuleV1Response{
				occurredAt: occurred,
				Status:     STATUS_NOT_FOUND,
			}, nil
		}

		logger.Error().Err(err).Str("fqn", cmd.DTO.FQN.String()).Msg("failed to delete module from store")

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
