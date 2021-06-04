package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/ymir/internal/registry"
)

type ModuleRegistryController struct {
	logger     zerolog.Logger
	moduleRepo registry.ModuleRepository
	cb         *registry.CommandBus
}

type ModuleVersionListVersionItem struct {
	Version string `json:"version"`
}

type ModuleVersionListItem struct {
	Source   string                         `json:"source"`
	Versions []ModuleVersionListVersionItem `json:"versions"`
}

type ModuleVersionList struct {
	Modules []ModuleVersionListItem `json:"modules"`
}

func (c *ModuleRegistryController) WellKnown(w http.ResponseWriter, r *http.Request) {
	cmd := registry.ServiceDiscoveryCommand{}
	resp := cmd.Handle()

	w.Header().Set("Content-Type", "application/json")

	if resp.Status != registry.STATUS_OKAY {
		c.logger.Error().Msg("unknown error occured")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	//nolint:errcheck
	json.NewEncoder(w).Encode(resp.Body)
}

func (c *ModuleRegistryController) ListModuleVersions(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ns := params["namespace"]
	name := params["name"]
	provider := params["provider"]

	res, err := c.cb.ListModuleVersionsV1ByFqn(registry.ListModuleVersionsByFqnV1DTO{
		FQN: registry.ModuleFQN{
			Name:      name,
			Namespace: ns,
			Provider:  provider,
		},
	})

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		c.logger.Error().Err(err).Str("action", "ModuleRegistry.ListModuleVersions").Msg("command failed")

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	switch res.Status {
	case registry.STATUS_OKAY:
		list := []ModuleVersionListVersionItem{}

		for _, v := range res.List {
			list = append(list, ModuleVersionListVersionItem{
				Version: v.Version,
			})
		}

		w.WriteHeader(http.StatusOK)

		//nolint:errcheck
		json.NewEncoder(w).Encode(list)
		return
	case registry.STATUS_NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		c.logger.Error().Str("status", string(res.Status)).Str("action", "Modules.GetModule").Msg("unhandle response")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *ModuleRegistryController) DownloadModule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ns := params["namespace"]
	name := params["name"]
	provider := params["provider"]
	version := params["version"]

	cmd := registry.DownloadModuleVersionV1Command{
		Namespace: ns,
		Name:      name,
		Provider:  provider,
		Version:   version,
	}

	resp, err := cmd.Handle(c.moduleRepo, registry.NewCommandValidator(c.logger), c.logger)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		c.logger.Error().Err(err).Msg("unknown error occurred")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if resp.Status == registry.STATUS_NOT_FOUND {
		c.logger.Debug().Str("version", version).Str("namespace", ns).Str("name", name).Str("provider", provider).Msg("module not found for download")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if resp.Status != registry.STATUS_OKAY {
		c.logger.Error().Err(err).Msg("unknown error occurred")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("X-Terraform-Get", resp.LocationURI)
	w.WriteHeader(http.StatusNoContent)

	//nolint:errcheck
	json.NewEncoder(w).Encode(params)
}

func (c *ModuleRegistryController) RegisterRoutes(r muxRouter) {
	r.HandleFunc("/.well-known/terraform.json", c.WellKnown)
	r.HandleFunc("/v1/modules/{namespace}/{name}/{provider}/versions", c.ListModuleVersions)
	r.HandleFunc("/v1/modules/{namespace}/{name}/{provider}/{version}/download", c.DownloadModule)
}

func NewModuleRegistryController(l zerolog.Logger, moduleRepo registry.ModuleRepository, cb *registry.CommandBus) *ModuleRegistryController {
	return &ModuleRegistryController{
		logger:     l,
		moduleRepo: moduleRepo,
		cb:         cb,
	}
}
