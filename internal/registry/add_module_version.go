package registry

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
)

type addModuleVersionRepository interface {
	AddVersion(ModuleVersion) (m ModuleVersion, err error)
	ById(id string) (m Module, err error)
	VersionByModuleAndValue(moduleId string, version string) (mv ModuleVersion, err error)
}

type addModuleVersionV1CommandValidator interface {
	RegisterStructLevelValidator(validator.StructLevelFunc, interface{})
	Validate(cmd interface{}) []ValidationError
	ModuleMustExistById(r moduleMustExistByIdRepository, sl validator.StructLevel, id string)
	ModuleVersionMustBeUniqueForId(r moduleVersionMustBeUniqueForIdRepository, sl validator.StructLevel, id string, version string)
}

type addModuleVersionV1Command struct {
	DTO AddModuleVersionV1DTO
}

type AddModuleVersionV1Response struct {
	occurredAt       time.Time
	Status           RegistryHandlerStatus
	ModuleVersion    ModuleVersion
	ValidationErrors []ValidationError
}

func (r AddModuleVersionV1Response) GetActionName() string {
	return "v1.modules.versions.add"
}

func (r AddModuleVersionV1Response) GetTimeOfOccurrence() time.Time {
	return r.occurredAt
}

func (r AddModuleVersionV1Response) GetResponseStatus() RegistryHandlerStatus {
	return r.Status
}

func (r AddModuleVersionV1Response) GetAuditMeta() map[string]interface{} {
	return map[string]interface{}{
		"module_version_id": r.ModuleVersion.Id,
		"validation_errors": r.ValidationErrors,
	}
}

type AddModuleVersionV1DTO struct {
	Version       string `json:"version" validate:"required,version"`
	ModuleId      string `json:"module_id" validate:"required,uuid"`
	Source        string `json:"source" validate:"required"`
	RepositoryURL string `json:"repository_url" validate:"required"`
}

func (dto AddModuleVersionV1DTO) validate(r addModuleVersionRepository, v addModuleVersionV1CommandValidator) []ValidationError {
	v.RegisterStructLevelValidator(func(sl validator.StructLevel) {
		v.ModuleMustExistById(r, sl, dto.ModuleId)
		v.ModuleVersionMustBeUniqueForId(r, sl, dto.ModuleId, dto.Version)
	}, AddModuleVersionV1DTO{})

	return v.Validate(dto)
}

func (cmd addModuleVersionV1Command) handle(r addModuleVersionRepository, logger zerolog.Logger, v addModuleVersionV1CommandValidator) (AddModuleVersionV1Response, error) {
	occurred := time.Now().UTC()

	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return AddModuleVersionV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	mv, err := r.AddVersion(ModuleVersion{
		Id:            uuid.NewString(),
		Version:       cmd.DTO.Version,
		ModuleId:      cmd.DTO.ModuleId,
		Source:        cmd.DTO.Source,
		RepositoryURL: cmd.DTO.RepositoryURL,
	})

	if err != nil {
		logger.Error().Err(err).Str("command", "add_module_version").Str("module_id", cmd.DTO.ModuleId).Str("version", cmd.DTO.Version).Msg("failed to add module version to store")

		return AddModuleVersionV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	return AddModuleVersionV1Response{
		occurredAt:    occurred,
		Status:        STATUS_CREATED,
		ModuleVersion: mv,
	}, err
}
