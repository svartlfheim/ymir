package registry

import (
	"time"

	"github.com/rs/zerolog"
)

type listModuleVersionsByFqnRepository interface {
	ByFQN(ModuleFQN) (m Module, err error)
	VersionsByModule(moduleId string, chunkOpts ChunkingOptions) (m []ModuleVersion, err error)
}

type ListModuleVersionsByFqnV1DTO struct {
	FQN ModuleFQN `validate:"required"`
}

type listModuleVersionsByFqnV1Command struct {
	DTO ListModuleVersionsByFqnV1DTO
}

func (cmd listModuleVersionsByFqnV1Command) handle(r listModuleVersionsByFqnRepository, l zerolog.Logger) (ListModuleVersionsV1Response, error) {
	module, err := r.ByFQN(cmd.DTO.FQN)

	occurred := time.Now().UTC()

	if _, ok := err.(ErrResourceNotFound); ok {
		return ListModuleVersionsV1Response{
			occurredAt: occurred,
			Status:     STATUS_NOT_FOUND,
		}, nil
	}

	if err != nil {
		l.Error().Err(err).Str("fqn", cmd.DTO.FQN.String()).Msg("error finding module")

		return ListModuleVersionsV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	var moduleVersions []ModuleVersion

	moduleVersions, err = r.VersionsByModule(module.Id, ChunkingOptions{})

	if err != nil {
		// We won't bother checking for ResourceNotFound errors here
		// The validation already checked this should exist
		// If we got an error, it should be something like a query error
		l.Error().Err(err).Str("fqn", cmd.DTO.FQN.String()).Msg("error listing module versions")

		return ListModuleVersionsV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	return ListModuleVersionsV1Response{
		occurredAt: occurred,
		Status:     STATUS_OKAY,
		List:       moduleVersions,
	}, nil
}
