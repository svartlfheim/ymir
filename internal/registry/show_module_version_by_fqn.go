package registry

import (
	"time"

	"github.com/rs/zerolog"
)

type showModuleVersionByFqnRepository interface {
	VersionByFQN(ModuleVersionFQN) (m ModuleVersion, err error)
}

type showModuleVersionByFqnV1CommandValidator interface {
	Validate(cmd interface{}) []ValidationError
}

type ShowModuleVersionByFqnV1DTO struct {
	FQN ModuleVersionFQN `validate:"required"`
}

type showModuleVersionByFqnV1Command struct {
	DTO ShowModuleVersionByFqnV1DTO
}

func (dto ShowModuleVersionByFqnV1DTO) validate(r showModuleVersionByFqnRepository, v showModuleVersionByFqnV1CommandValidator) []ValidationError {
	return v.Validate(dto)
}

func (cmd showModuleVersionByFqnV1Command) handle(r showModuleVersionByFqnRepository, logger zerolog.Logger, v showModuleVersionByFqnV1CommandValidator) (ShowModuleVersionV1Response, error) {
	occurred := time.Now().UTC()

	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return ShowModuleVersionV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	mv, err := r.VersionByFQN(cmd.DTO.FQN)

	if err != nil {
		if _, ok := err.(ErrResourceNotFound); ok {
			return ShowModuleVersionV1Response{
				occurredAt: occurred,
				Status:     STATUS_NOT_FOUND,
			}, nil
		}

		logger.Error().Err(err).Str("fqn", cmd.DTO.FQN.String()).Msg("failed to add module to store")

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
