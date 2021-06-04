package registry

import (
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
)

type addModuleVersionByFqnRepository interface {
	AddVersion(ModuleVersion) (m ModuleVersion, err error)
	ById(id string) (m Module, err error)
	ByFQN(ModuleFQN) (m Module, err error)
	VersionByFQN(fqn ModuleVersionFQN) (mv ModuleVersion, err error)
	VersionByModuleAndValue(moduleId string, version string) (mv ModuleVersion, err error)
}

type addModuleVersionByFqnV1CommandValidator interface {
	RegisterStructLevelValidator(validator.StructLevelFunc, interface{})
	Validate(cmd interface{}) []ValidationError
	ModuleMustExistByFQN(r moduleMustExistByFQNRepository, sl validator.StructLevel, fqn ModuleFQN)
	ModuleMustExistById(r moduleMustExistByIdRepository, sl validator.StructLevel, id string)
	ModuleVersionMustBeUniqueForFQN(r moduleVersionMustBeUniqueForFQNRepository, sl validator.StructLevel, fqn ModuleVersionFQN)
	ModuleVersionMustBeUniqueForId(r moduleVersionMustBeUniqueForIdRepository, sl validator.StructLevel, id string, version string)
}

type addModuleVersionV1ByModuleFqnCommand struct {
	DTO AddModuleVersionV1ByModuleFqnDTO
}

type AddModuleVersionV1ByModuleFqnDTO struct {
	Version       string    `json:"version" validate:"required,version"`
	ModuleFQN     ModuleFQN `json:"required"`
	Source        string    `json:"source" validate:"required"`
	RepositoryURL string    `json:"repository_url" validate:"required"`
}

func (dto AddModuleVersionV1ByModuleFqnDTO) validate(r addModuleVersionByFqnRepository, v addModuleVersionByFqnV1CommandValidator) []ValidationError {
	v.RegisterStructLevelValidator(func(sl validator.StructLevel) {
		v.ModuleMustExistByFQN(r, sl, dto.ModuleFQN)
		v.ModuleVersionMustBeUniqueForFQN(r, sl, ModuleVersionFQN{
			ModuleFQN: dto.ModuleFQN,
			Version:   dto.Version,
		})
	}, AddModuleVersionV1DTO{})

	return v.Validate(dto)
}

func (cmd addModuleVersionV1ByModuleFqnCommand) handle(r addModuleVersionByFqnRepository, logger zerolog.Logger, v addModuleVersionByFqnV1CommandValidator) (AddModuleVersionV1Response, error) {
	occurred := time.Now().UTC()

	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return AddModuleVersionV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	m, err := r.ByFQN(cmd.DTO.ModuleFQN)

	if err != nil {
		logger.Error().Err(err).Str("command", "add_module_version").Str("fqn", cmd.DTO.ModuleFQN.String()).Msg("failed find module")

		return AddModuleVersionV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	byIdCmd := addModuleVersionV1Command{
		DTO: AddModuleVersionV1DTO{
			Version:       cmd.DTO.Version,
			Source:        cmd.DTO.Source,
			RepositoryURL: cmd.DTO.RepositoryURL,
			ModuleId:      m.Id,
		},
	}

	return byIdCmd.handle(r, logger, v)
}
