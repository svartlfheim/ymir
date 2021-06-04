package registry

import (
	"time"

	"github.com/rs/zerolog"
)

type showModuleByFqnRepository interface {
	ByFQN(ModuleFQN) (m Module, err error)
}

type showModuleByFqnV1CommandValidator interface {
	Validate(cmd interface{}) []ValidationError
}

type showModuleV1ByFqnCommand struct {
	DTO ShowModuleV1ByFqnDTO
}

type ShowModuleV1ByFqnDTO struct {
	FQN ModuleFQN `validate:"required"`
}

func (dto ShowModuleV1ByFqnDTO) validate(r showModuleByFqnRepository, v showModuleByFqnV1CommandValidator) []ValidationError {
	return v.Validate(dto)
}

func (cmd showModuleV1ByFqnCommand) handle(r showModuleByFqnRepository, logger zerolog.Logger, v showModuleByFqnV1CommandValidator) (ShowModuleV1Response, error) {
	occurred := time.Now().UTC()
	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return ShowModuleV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	m, err := r.ByFQN(cmd.DTO.FQN)

	if err != nil {
		if _, ok := err.(ErrResourceNotFound); ok {
			return ShowModuleV1Response{
				Status: STATUS_NOT_FOUND,
			}, nil
		}

		logger.Error().Err(err).Str("fqn", cmd.DTO.FQN.String()).Msg("failed to add module to store")

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
