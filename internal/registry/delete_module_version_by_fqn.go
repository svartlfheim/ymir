package registry

import (
	"time"

	"github.com/rs/zerolog"
)

type deleteModuleVersionByFqnRepository interface {
	VersionByFQN(ModuleVersionFQN) (m ModuleVersion, err error)
	DeleteModuleVersion(ModuleVersion) error
}

type deleteModuleVersionByFqnV1CommandValidator interface {
	Validate(cmd interface{}) []ValidationError
}

type DeleteModuleVersionByFqnV1DTO struct {
	FQN ModuleVersionFQN `validate:"required"`
}

type deleteModuleVersionByFqnV1Command struct {
	DTO DeleteModuleVersionByFqnV1DTO
}

func (dto DeleteModuleVersionByFqnV1DTO) validate(r deleteModuleVersionByFqnRepository, v deleteModuleVersionByFqnV1CommandValidator) []ValidationError {
	return v.Validate(dto)
}

func (cmd deleteModuleVersionByFqnV1Command) handle(r deleteModuleVersionByFqnRepository, logger zerolog.Logger, v deleteModuleVersionByFqnV1CommandValidator) (DeleteModuleVersionV1Response, error) {
	occurred := time.Now().UTC()
	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return DeleteModuleVersionV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	mv, err := r.VersionByFQN(cmd.DTO.FQN)

	if err != nil {
		if _, ok := err.(ErrResourceNotFound); ok {
			return DeleteModuleVersionV1Response{
				occurredAt: occurred,
				Status:     STATUS_NOT_FOUND,
			}, nil
		}

		logger.Error().Err(err).Str("fqn", cmd.DTO.FQN.String()).Msg("failed to find module version")

		return DeleteModuleVersionV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	err = r.DeleteModuleVersion(mv)

	if err != nil {
		logger.Error().Err(err).Str("fqn", cmd.DTO.FQN.String()).Msg("failed to delete module version")

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
