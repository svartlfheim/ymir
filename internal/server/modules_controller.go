package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/ymir/internal/registry"
)

type ModulesController struct {
	logger  zerolog.Logger
	cb      *registry.CommandBus
	auditor requestAuditor
}

type requestAuditor interface {
	Record(action registry.AuditableAction)
}

func (c *ModulesController) ListModules(w http.ResponseWriter, r *http.Request) {
	res, err := c.cb.ListModulesV1FromDTO(registry.ListModulesV1DTO{})

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.ListModules").Msg("command failed")

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	go c.auditor.Record(res)

	switch res.Status {
	case registry.STATUS_OKAY:
		handleResourceResponse(res.List, http.StatusOK, w)
		return
	default:
		c.logger.Error().Str("status", string(res.Status)).Str("action", "Modules.ListModules").Msg("unhandled response")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *ModulesController) PostModule(w http.ResponseWriter, r *http.Request) {

	dto := registry.AddModuleV1DTO{}
	err := json.NewDecoder(r.Body).Decode(&dto)

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.PostModule").Msg("failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	res, err := c.cb.AddModuleV1FromDTO(dto)

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.PostModule").Msg("command failed")

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	go c.auditor.Record(res)

	switch res.Status {
	case registry.STATUS_INVALID:
		handleValidationErrorsResponse(res.ValidationErrors, http.StatusUnprocessableEntity, w)
		return
	case registry.STATUS_CREATED:
		handleResourceResponse(res.Module, http.StatusCreated, w)
		return
	default:
		c.logger.Error().Str("status", string(res.Status)).Str("action", "Modules.PostModule").Msg("unhandled response")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *ModulesController) GetModule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	res, err := c.cb.ShowModuleV1ByID(registry.ShowModuleV1DTO{
		Id: id,
	})

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.GetModule").Msg("command failed")

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	go c.auditor.Record(res)

	switch res.Status {
	case registry.STATUS_INVALID:
		handleValidationErrorsResponse(res.ValidationErrors, http.StatusBadRequest, w)
		return
	case registry.STATUS_OKAY:
		handleResourceResponse(res.Module, http.StatusOK, w)
		return
	case registry.STATUS_NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		c.logger.Error().Str("status", string(res.Status)).Str("action", "Modules.GetModule").Msg("unhandled response")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *ModulesController) DeleteModule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	dto := registry.DeleteModuleV1DTO{}
	err := json.NewDecoder(r.Body).Decode(&dto)

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.DeleteModule").Msg("failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.logger.Info().Bool("should", dto.DeleteVersions).Msg("delete versions")
	dto.Id = id

	res, err := c.cb.DeleteModuleV1ById(dto)

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.DeleteModule").Msg("command failed")

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	go c.auditor.Record(res)

	switch res.Status {
	case registry.STATUS_INVALID:
		handleValidationErrorsResponse(res.ValidationErrors, http.StatusUnprocessableEntity, w)
		return
	case registry.STATUS_OKAY:
		handleResourceResponse(res.Module, http.StatusOK, w)
		return
	case registry.STATUS_NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		c.logger.Error().Str("status", string(res.Status)).Str("action", "Modules.GetModule").Msg("unhandled response")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *ModulesController) ListModuleVersions(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	moduleId := params["module_id"]

	dto := registry.ListModuleVersionsV1DTO{
		ModuleId: moduleId,
	}

	res, err := c.cb.ListModuleVersionsV1ById(dto)

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.ListModuleVersions").Msg("command failed")

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	go c.auditor.Record(res)

	switch res.Status {
	case registry.STATUS_OKAY:
		handleResourceResponse(res.List, http.StatusOK, w)
		return
	case registry.STATUS_NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		c.logger.Error().Str("status", string(res.Status)).Str("action", "Modules.GetModule").Msg("unhandled response")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *ModulesController) CreateModuleVersion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	moduleId := params["module_id"]

	dto := registry.AddModuleVersionV1DTO{}
	err := json.NewDecoder(r.Body).Decode(&dto)

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.ListModuleVersions").Msg("failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	dto.ModuleId = moduleId

	res, err := c.cb.AddModuleVersionV1ForModuleId(dto)

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.ListModuleVersions").Msg("command failed")

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	go c.auditor.Record(res)

	switch res.Status {
	case registry.STATUS_INVALID:
		handleValidationErrorsResponse(res.ValidationErrors, http.StatusUnprocessableEntity, w)
		return
	case registry.STATUS_CREATED:
		handleResourceResponse(res.ModuleVersion, http.StatusCreated, w)
		return
	default:
		c.logger.Error().Str("status", string(res.Status)).Str("action", "Modules.ListModuleVersions").Msg("unhandled response")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *ModulesController) GetModuleVersion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	res, err := c.cb.ShowModuleVersionV1ById(registry.ShowModuleVersionV1DTO{
		Id: id,
	})

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.GetModuleVersion").Msg("command failed")

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	go c.auditor.Record(res)

	switch res.Status {
	case registry.STATUS_INVALID:
		handleValidationErrorsResponse(res.ValidationErrors, http.StatusBadRequest, w)
		return
	case registry.STATUS_OKAY:
		handleResourceResponse(res.ModuleVersion, http.StatusOK, w)
		return
	case registry.STATUS_NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		c.logger.Error().Str("status", string(res.Status)).Str("action", "Modules.GetModuleVersion").Msg("unhandled response")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *ModulesController) DeleteModuleVersion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	res, err := c.cb.DeleteModuleVersionV1ById(registry.DeleteModuleVersionV1DTO{
		Id: id,
	})

	if err != nil {
		c.logger.Error().Err(err).Str("action", "Modules.DeleteModuleVersion").Msg("command failed")

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	go c.auditor.Record(res)

	switch res.Status {
	case registry.STATUS_INVALID:
		handleValidationErrorsResponse(res.ValidationErrors, http.StatusBadRequest, w)
		return
	case registry.STATUS_OKAY:
		handleResourceResponse(res.ModuleVersion, http.StatusOK, w)
		return
	case registry.STATUS_NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		c.logger.Error().Str("status", string(res.Status)).Str("action", "Modules.DeleteModuleVersion").Msg("unhandled response")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *ModulesController) RegisterRoutes(r muxRouter) {
	api := r.PathPrefix("/api").Subrouter()
	api.Use(apiMiddleware)

	api.HandleFunc("/v1/modules", c.ListModules).Methods("GET")
	api.HandleFunc("/v1/modules", c.PostModule).Methods("POST")
	api.HandleFunc("/v1/modules/{id}", c.GetModule).Methods("GET")
	api.HandleFunc("/v1/modules/{id}", c.DeleteModule).Methods("DELETE")

	api.HandleFunc("/v1/modules/{module_id}/versions", c.ListModuleVersions).Methods("GET")
	api.HandleFunc("/v1/modules/{module_id}/versions", c.CreateModuleVersion).Methods("POST")

	api.HandleFunc("/v1/module-versions/{id}", c.GetModuleVersion).Methods("GET")
	api.HandleFunc("/v1/module-versions/{id}", c.DeleteModuleVersion).Methods("DELETE")
}

func NewModulesController(l zerolog.Logger, cb *registry.CommandBus, a requestAuditor) *ModulesController {
	return &ModulesController{
		logger:  l,
		cb:      cb,
		auditor: a,
	}
}
