package registry

import (
	"time"

	"github.com/rs/zerolog"
)

type listModuleVersionsRepository interface {
	ByFQN(ModuleFQN) (m Module, err error)
	ById(id string) (m Module, err error)
	VersionsByModule(moduleId string, chunkOpts ChunkingOptions) (m []ModuleVersion, err error)
}

type ListModuleVersionsV1DTO struct {
	ModuleId string `validate:"required,uuid"`
}

type listModuleVersionsV1Command struct {
	DTO ListModuleVersionsV1DTO
}

type ListModuleVersionsV1Response struct {
	attemptedFor     string
	occurredAt       time.Time
	Status           RegistryHandlerStatus
	ValidationErrors []ValidationError
	List             []ModuleVersion
}

func (r ListModuleVersionsV1Response) GetActionName() string {
	return "v1.modules.versions.list"
}

func (r ListModuleVersionsV1Response) GetTimeOfOccurrence() time.Time {
	return r.occurredAt
}

func (r ListModuleVersionsV1Response) GetResponseStatus() RegistryHandlerStatus {
	return r.Status
}

func (r ListModuleVersionsV1Response) GetAuditMeta() map[string]interface{} {
	return map[string]interface{}{
		"total":             len(r.List),
		"module_id":         r.attemptedFor,
		"validation_errors": r.ValidationErrors,
	}
}

func (cmd listModuleVersionsV1Command) handle(r listModuleVersionsRepository, l zerolog.Logger) (ListModuleVersionsV1Response, error) {
	module, err := r.ById(cmd.DTO.ModuleId)

	occurred := time.Now().UTC()

	if _, ok := err.(ErrResourceNotFound); ok {
		return ListModuleVersionsV1Response{
			occurredAt: occurred,
			Status:     STATUS_NOT_FOUND,
		}, nil
	}

	if err != nil {
		l.Error().Err(err).Str("id", cmd.DTO.ModuleId).Msg("error finding module")

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
		l.Error().Err(err).Str("id", cmd.DTO.ModuleId).Msg("error listing module versions")

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
