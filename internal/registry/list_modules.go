package registry

import (
	"time"

	"github.com/rs/zerolog"
)

type listModulesRepository interface {
	All(chunkOpts ChunkingOptions, filters ModuleFilters) ([]Module, error)
}

type ListModulesV1DTO struct {
	Namespace string
	Provider  string
	ChunkOpts ChunkingOptions
}

type listModulesV1Command struct {
	DTO ListModulesV1DTO
}

type ListModulesV1Response struct {
	occurredAt time.Time
	Status     RegistryHandlerStatus
	List       []Module
}

func (r ListModulesV1Response) GetActionName() string {
	return "v1.modules.list"
}

func (r ListModulesV1Response) GetTimeOfOccurrence() time.Time {
	return r.occurredAt
}

func (r ListModulesV1Response) GetResponseStatus() RegistryHandlerStatus {
	return r.Status
}

func (r ListModulesV1Response) GetAuditMeta() map[string]interface{} {
	return map[string]interface{}{
		"total": len(r.List),
	}
}

func (cmd listModulesV1Command) handle(r listModulesRepository, l zerolog.Logger) (ListModulesV1Response, error) {
	occurred := time.Now().UTC()

	// empty chunking opts as not used by repo yet
	modules, err := r.All(cmd.DTO.ChunkOpts, ModuleFilters{
		Provider:  cmd.DTO.Provider,
		Namespace: cmd.DTO.Namespace,
	})

	if err != nil {
		l.Error().Str("provider-filter", cmd.DTO.Provider).Str("namespace-filter", cmd.DTO.Namespace).Err(err).Msg("error listing modules")

		return ListModulesV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	return ListModulesV1Response{
		occurredAt: occurred,
		Status:     STATUS_OKAY,
		List:       modules,
	}, nil
}
