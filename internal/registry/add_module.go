package registry

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
)

type addModuleRepository interface {
	AddModule(Module) (Module, error)
	ByFQN(ModuleFQN) (m Module, err error)
}

type addModuleV1CommandValidator interface {
	RegisterStructLevelValidator(validator.StructLevelFunc, interface{})
	Validate(cmd interface{}) []ValidationError
	ModuleNameMustBeUnique(r uniqueModuleNameRepository, sl validator.StructLevel, fqn ModuleFQN)
}

type addModuleV1Command struct {
	DTO AddModuleV1DTO
}

type AddModuleV1Response struct {
	occurredAt       time.Time
	Status           RegistryHandlerStatus
	Module           Module
	ValidationErrors []ValidationError
}

func (r AddModuleV1Response) GetActionName() string {
	return "v1.modules.add"
}

func (r AddModuleV1Response) GetTimeOfOccurrence() time.Time {
	return r.occurredAt
}

func (r AddModuleV1Response) GetResponseStatus() RegistryHandlerStatus {
	return r.Status
}

func (r AddModuleV1Response) GetAuditMeta() map[string]interface{} {
	return map[string]interface{}{
		"module_id":         r.Module.Id,
		"validation_errors": r.ValidationErrors,
	}
}

type AddModuleV1DTO struct {
	Namespace string `validate:"required" json:"namespace"`
	Name      string `validate:"required" json:"name"`
	Provider  string `validate:"required" json:"provider"`
}

func (d AddModuleV1DTO) ToFQN() ModuleFQN {
	return ModuleFQN{
		Name:      d.Name,
		Namespace: d.Namespace,
		Provider:  d.Provider,
	}
}

func (dto AddModuleV1DTO) validate(r addModuleRepository, v addModuleV1CommandValidator) []ValidationError {
	v.RegisterStructLevelValidator(func(sl validator.StructLevel) {
		v.ModuleNameMustBeUnique(r, sl, ModuleFQN{
			Name:      dto.Name,
			Namespace: dto.Namespace,
			Provider:  dto.Provider,
		})
	}, AddModuleV1DTO{})

	return v.Validate(dto)
}

func (cmd addModuleV1Command) handle(r addModuleRepository, logger zerolog.Logger, v addModuleV1CommandValidator) (AddModuleV1Response, error) {
	occurred := time.Now().UTC()
	if errs := cmd.DTO.validate(r, v); len(errs) > 0 {
		return AddModuleV1Response{
			occurredAt:       occurred,
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	fqn := cmd.DTO.ToFQN()
	m, err := r.AddModule(Module{
		Id:        uuid.NewString(),
		Name:      fqn.Name,
		Namespace: fqn.Namespace,
		Provider:  fqn.Provider,
	})

	if err != nil {
		logger.Error().Err(err).Str("fqn", fqn.String()).Msg("failed to add module to store")

		return AddModuleV1Response{
			occurredAt: occurred,
			Status:     STATUS_INTERNAL_ERROR,
		}, err
	}

	return AddModuleV1Response{
		occurredAt: occurred,
		Status:     STATUS_CREATED,
		Module:     m,
	}, err
}
