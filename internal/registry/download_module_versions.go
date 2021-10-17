package registry

import (
	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
)

type downloadModuleVersionRepository interface {
	ByFQN(ModuleFQN) (m Module, err error)
	VersionByFQN(fqn ModuleVersionFQN) (mv ModuleVersion, err error)
}

type downloadModuleV1CommandValidator interface {
	RegisterStructLevelValidator(validator.StructLevelFunc, interface{})
	Validate(cmd interface{}) []ValidationError
	ModuleMustExistByFQN(r moduleMustExistByFQNRepository, sl validator.StructLevel, fqn ModuleFQN)
	ModuleVersionMustBeUniqueForFQN(r moduleVersionMustBeUniqueForFQNRepository, sl validator.StructLevel, fqn ModuleVersionFQN)
}

type DownloadModuleVersionV1Command struct {
	Namespace string
	Name      string
	Provider  string
	Version   string
}

type HandleDownloadModuleVersionV1Response struct {
	Status           RegistryHandlerStatus
	LocationURI      string
	ValidationErrors []ValidationError
}

func (cmd DownloadModuleVersionV1Command) validate(r downloadModuleVersionRepository, v downloadModuleV1CommandValidator, logger zerolog.Logger) []ValidationError {
	v.RegisterStructLevelValidator(func(sl validator.StructLevel) {
		fqn := ModuleVersionFQN{
			ModuleFQN: ModuleFQN{
				Name:      cmd.Name,
				Namespace: cmd.Namespace,
				Provider:  cmd.Provider,
			},
			Version: cmd.Version,
		}
		v.ModuleMustExistByFQN(r, sl, fqn.ModuleFQN)
		v.ModuleVersionMustBeUniqueForFQN(r, sl, fqn)
	}, DownloadModuleVersionV1Command{})

	return v.Validate(cmd)
}

func (c DownloadModuleVersionV1Command) Handle(r downloadModuleVersionRepository, val downloadModuleV1CommandValidator, l zerolog.Logger) (HandleDownloadModuleVersionV1Response, error) {
	if errs := c.validate(r, val, l); len(errs) > 0 {
		return HandleDownloadModuleVersionV1Response{
			Status:           STATUS_INVALID,
			ValidationErrors: errs,
		}, nil
	}

	fqn := ModuleVersionFQN{
		ModuleFQN: ModuleFQN{
			Name:      c.Name,
			Namespace: c.Namespace,
			Provider:  c.Provider,
		},
		Version: c.Version,
	}

	version, err := r.VersionByFQN(fqn)

	if err != nil {
		// We won't bother checking for ResourceNotFound errors here
		// The validation already checked this should exist
		// If we got an error, it should be something like a query error
		l.Error().Err(err).Str("version", c.Version).Str("fqn", fqn.String()).Msg("failed to find module version")

		return HandleDownloadModuleVersionV1Response{
			Status: STATUS_INTERNAL_ERROR,
		}, err
	}

	return HandleDownloadModuleVersionV1Response{
		Status:      STATUS_OKAY,
		LocationURI: version.DownloadURL,
	}, nil
}
