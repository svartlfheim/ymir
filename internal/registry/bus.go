package registry

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type cliPrompter interface {
	Ask(string) (string, error)
}

type CommandBus struct {
	repo           ModuleRepository
	logger         zerolog.Logger
	buildValidator ValidatorBuilder
	prompter       cliPrompter
	fs             afero.Fs
}

type WithDependency func(*CommandBus)

func WithModuleRepo(r ModuleRepository) WithDependency {
	return func(cb *CommandBus) {
		cb.repo = r
	}
}

func WithLogger(l zerolog.Logger) WithDependency {
	return func(cb *CommandBus) {
		cb.logger = l
	}
}

func WithCommandValidatorBuilder(b ValidatorBuilder) WithDependency {
	return func(cb *CommandBus) {
		cb.buildValidator = b
	}
}

func WithPrompter(p cliPrompter) WithDependency {
	return func(cb *CommandBus) {
		cb.prompter = p
	}
}

func WithFS(fs afero.Fs) WithDependency {
	return func(cb *CommandBus) {
		cb.fs = fs
	}
}

func NewCommandBus(opts ...WithDependency) *CommandBus {
	cb := &CommandBus{}

	for _, opt := range opts {
		opt(cb)
	}

	return cb
}

func (cb *CommandBus) AddModuleV1FromCLI(filePath string) (AddModuleV1Response, error) {
	dto := AddModuleV1DTO{}

	if filePath != "" {
		b, err := afero.ReadFile(cb.fs, filePath)

		if err != nil {
			return AddModuleV1Response{}, ErrCouldNotReadFile{
				Path: filePath,
			}
		}

		// nolint: errcheck
		err = json.Unmarshal(b, &dto)

		if err != nil {
			return AddModuleV1Response{}, ErrCouldNotUnmarshalJSONToDTO{
				Path:    filePath,
				DTOName: "AddModuleV1DTO",
			}
		}
	} else {
		var name, ns, provider string
		var err error

		if provider, err = cb.prompter.Ask("Provider: "); err != nil {
			return AddModuleV1Response{}, ErrQuestionFailed{
				Question: "Provider: ",
			}
		}

		if ns, err = cb.prompter.Ask("Namespace: "); err != nil {
			return AddModuleV1Response{}, ErrQuestionFailed{
				Question: "Namespace: ",
			}
		}

		if name, err = cb.prompter.Ask("Name: "); err != nil {
			return AddModuleV1Response{}, ErrQuestionFailed{
				Question: "Name: ",
			}
		}

		dto.Name = name
		dto.Namespace = ns
		dto.Provider = provider
	}

	return cb.AddModuleV1FromDTO(dto)
}

func (cb *CommandBus) AddModuleV1FromDTO(dto AddModuleV1DTO) (AddModuleV1Response, error) {
	cmd := addModuleV1Command{
		DTO: dto,
	}

	v := cb.buildValidator(cb.logger)

	return cmd.handle(cb.repo, cb.logger, v)
}

func (cb *CommandBus) ListModulesV1FromCLI(p string, ns string) (ListModulesV1Response, error) {
	dto := ListModulesV1DTO{
		Provider:  p,
		Namespace: ns,
		ChunkOpts: ChunkingOptions{},
	}

	return cb.ListModulesV1FromDTO(dto)
}

func (cb *CommandBus) ListModulesV1FromDTO(dto ListModulesV1DTO) (ListModulesV1Response, error) {
	cmd := listModulesV1Command{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger)
}

func (cb *CommandBus) ShowModuleV1FromCLI(idOrFQN string) (ShowModuleV1Response, error) {
	if fqn, err := ParseModuleFQN(idOrFQN); err == nil {
		dto := ShowModuleV1ByFqnDTO{
			FQN: fqn,
		}
		return cb.ShowModuleV1ByFQN(dto)
	}

	if _, err := uuid.Parse(idOrFQN); err == nil {
		dto := ShowModuleV1DTO{
			Id: idOrFQN,
		}
		return cb.ShowModuleV1ByID(dto)
	}

	return ShowModuleV1Response{
		Status: STATUS_INVALID,
		ValidationErrors: []ValidationError{
			{
				Message: "id must be a uuid or an FQN formatted string (provider/namespace/name)",
				Rule:    "id_or_fqn",
				Field:   "id",
				Value:   idOrFQN,
			},
		},
	}, nil
}

func (cb *CommandBus) ShowModuleV1ByFQN(dto ShowModuleV1ByFqnDTO) (ShowModuleV1Response, error) {
	cmd := showModuleV1ByFqnCommand{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger, cb.buildValidator(cb.logger))
}

func (cb *CommandBus) ShowModuleV1ByID(dto ShowModuleV1DTO) (ShowModuleV1Response, error) {
	cmd := showModuleV1Command{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger, cb.buildValidator(cb.logger))
}

func (cb *CommandBus) DeleteModuleV1FromCLI(idOrFQN string, deleteVersions bool, force bool) (DeleteModuleV1Response, error) {
	fqn, fqnParseErr := ParseModuleFQN(idOrFQN)
	_, uuidParseErr := uuid.Parse(idOrFQN)

	if fqnParseErr != nil && uuidParseErr != nil {
		return DeleteModuleV1Response{
			Status: STATUS_INVALID,
			ValidationErrors: []ValidationError{
				{
					Message: "id must be a uuid or an FQN formatted string (provider/namespace/name)",
					Rule:    "id_or_fqn",
					Field:   "id",
					Value:   idOrFQN,
				},
			},
		}, nil
	}

	if !force {
		isSure, err := cb.prompter.Ask("Are you sure? [y/n]")

		if err != nil || !strings.EqualFold("y", isSure) {
			return DeleteModuleV1Response{}, ErrFailedToConfirmAction{
				Action: "delete_module:" + idOrFQN,
			}
		}
	}

	if fqnParseErr == nil {
		dto := DeleteModuleV1ByFqnDTO{
			FQN:            fqn,
			DeleteVersions: deleteVersions,
		}

		return cb.DeleteModuleV1ByFqn(dto)
	}

	dto := DeleteModuleV1DTO{
		Id:             idOrFQN,
		DeleteVersions: deleteVersions,
	}

	return cb.DeleteModuleV1ById(dto)
}

func (cb *CommandBus) DeleteModuleV1ByFqn(dto DeleteModuleV1ByFqnDTO) (DeleteModuleV1Response, error) {
	cmd := deleteModuleV1ByFqnCommand{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger, cb.buildValidator(cb.logger))
}

func (cb *CommandBus) DeleteModuleV1ById(dto DeleteModuleV1DTO) (DeleteModuleV1Response, error) {
	cmd := deleteModuleV1Command{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger, cb.buildValidator(cb.logger))
}

func (cb *CommandBus) AddModuleVersionV1FromCLI(filePath string) (AddModuleVersionV1Response, error) {

	if filePath != "" {
		b, err := afero.ReadFile(cb.fs, filePath)

		if err != nil {
			return AddModuleVersionV1Response{}, ErrCouldNotReadFile{
				Path: filePath,
			}
		}

		dto := AddModuleVersionV1DTO{}
		// nolint: errcheck
		err = json.Unmarshal(b, &dto)

		if err != nil {
			return AddModuleVersionV1Response{}, ErrCouldNotUnmarshalJSONToDTO{
				Path:    filePath,
				DTOName: "AddModuleV1DTO",
			}
		}

		return cb.AddModuleVersionV1ForModuleId(dto)
	}

	var version, source, repoUrl, idOrFQN string
	var err error

	if idOrFQN, err = cb.prompter.Ask("Module (id or FQN): "); err != nil {
		return AddModuleVersionV1Response{}, ErrQuestionFailed{
			Question: "Module (id or FQN): ",
		}
	}

	if version, err = cb.prompter.Ask("Version: "); err != nil {
		return AddModuleVersionV1Response{}, ErrQuestionFailed{
			Question: "Version: ",
		}
	}

	if source, err = cb.prompter.Ask("Source: "); err != nil {
		return AddModuleVersionV1Response{}, ErrQuestionFailed{
			Question: "Source: ",
		}
	}

	if repoUrl, err = cb.prompter.Ask("Repository URL: "); err != nil {
		return AddModuleVersionV1Response{}, ErrQuestionFailed{
			Question: "Repository URL: ",
		}
	}

	fqn, fqnParseErr := ParseModuleFQN(idOrFQN)
	_, uuidParseErr := uuid.Parse(idOrFQN)

	if fqnParseErr != nil && uuidParseErr != nil {
		return AddModuleVersionV1Response{
			Status: STATUS_INVALID,
			ValidationErrors: []ValidationError{
				{
					Message: "module id must be a uuid or an FQN formatted string (provider/namespace/name)",
					Rule:    "module_id_or_fqn",
					Field:   "module_id",
					Value:   idOrFQN,
				},
			},
		}, nil
	}

	if fqnParseErr == nil {
		dto := AddModuleVersionV1ByModuleFqnDTO{
			Version:       version,
			Source:        source,
			RepositoryURL: repoUrl,
			ModuleFQN:     fqn,
		}

		return cb.AddModuleVersionV1ForModuleFqn(dto)
	}

	dto := AddModuleVersionV1DTO{
		Version:       version,
		Source:        source,
		RepositoryURL: repoUrl,
		ModuleId:      idOrFQN,
	}

	return cb.AddModuleVersionV1ForModuleId(dto)
}

func (cb *CommandBus) AddModuleVersionV1ForModuleFqn(dto AddModuleVersionV1ByModuleFqnDTO) (AddModuleVersionV1Response, error) {
	cmd := addModuleVersionV1ByModuleFqnCommand{
		DTO: dto,
	}

	v := cb.buildValidator(cb.logger)

	return cmd.handle(cb.repo, cb.logger, v)
}

func (cb *CommandBus) AddModuleVersionV1ForModuleId(dto AddModuleVersionV1DTO) (AddModuleVersionV1Response, error) {
	cmd := addModuleVersionV1Command{
		DTO: dto,
	}

	v := cb.buildValidator(cb.logger)

	return cmd.handle(cb.repo, cb.logger, v)
}

func (cb *CommandBus) ListModuleVersionsV1FromCLI(idOrFQN string) (ListModuleVersionsV1Response, error) {
	fqn, fqnParseErr := ParseModuleFQN(idOrFQN)
	_, uuidParseErr := uuid.Parse(idOrFQN)

	if fqnParseErr != nil && uuidParseErr != nil {
		return ListModuleVersionsV1Response{}, ErrCouldNotParseModuleFQN{
			Value:   idOrFQN,
			Message: "id must be a uuid or an FQN formatted string (provider/namespace/name)",
		}
	}

	if fqnParseErr == nil {
		dto := ListModuleVersionsByFqnV1DTO{
			FQN: fqn,
		}

		return cb.ListModuleVersionsV1ByFqn(dto)
	}

	dto := ListModuleVersionsV1DTO{
		ModuleId: idOrFQN,
	}

	return cb.ListModuleVersionsV1ById(dto)
}

func (cb *CommandBus) ListModuleVersionsV1ByFqn(dto ListModuleVersionsByFqnV1DTO) (ListModuleVersionsV1Response, error) {
	cmd := listModuleVersionsByFqnV1Command{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger)
}

func (cb *CommandBus) ListModuleVersionsV1ById(dto ListModuleVersionsV1DTO) (ListModuleVersionsV1Response, error) {
	cmd := listModuleVersionsV1Command{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger)
}

func (cb *CommandBus) ShowModuleVersionV1FromCLI(idOrFQN string) (ShowModuleVersionV1Response, error) {
	fqn, fqnParseErr := ParseModuleVersionFQN(idOrFQN)
	_, uuidParseErr := uuid.Parse(idOrFQN)

	if fqnParseErr != nil && uuidParseErr != nil {
		return ShowModuleVersionV1Response{
			Status: STATUS_INVALID,
			ValidationErrors: []ValidationError{
				{
					Message: "id must be a uuid or an FQN formatted string (provider/namespace/name)",
					Rule:    "id_or_fqn",
					Field:   "id",
					Value:   idOrFQN,
				},
			},
		}, nil
	}

	if fqnParseErr == nil {
		dto := ShowModuleVersionByFqnV1DTO{
			FQN: fqn,
		}

		return cb.ShowModuleVersionV1ByFqn(dto)
	}

	dto := ShowModuleVersionV1DTO{
		Id: idOrFQN,
	}

	return cb.ShowModuleVersionV1ById(dto)
}

func (cb *CommandBus) ShowModuleVersionV1ByFqn(dto ShowModuleVersionByFqnV1DTO) (ShowModuleVersionV1Response, error) {
	cmd := showModuleVersionByFqnV1Command{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger, cb.buildValidator(cb.logger))
}

func (cb *CommandBus) ShowModuleVersionV1ById(dto ShowModuleVersionV1DTO) (ShowModuleVersionV1Response, error) {
	cmd := showModuleVersionV1Command{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger, cb.buildValidator(cb.logger))
}

func (cb *CommandBus) DeleteModuleVersionV1FromCLI(idOrFQN string, force bool) (DeleteModuleVersionV1Response, error) {
	fqn, fqnParseErr := ParseModuleVersionFQN(idOrFQN)
	_, uuidParseErr := uuid.Parse(idOrFQN)

	if fqnParseErr != nil && uuidParseErr != nil {
		return DeleteModuleVersionV1Response{
			Status: STATUS_INVALID,
			ValidationErrors: []ValidationError{
				{
					Message: "id must be a uuid or an FQN formatted string (provider/namespace/name)",
					Rule:    "id_or_fqn",
					Field:   "id",
					Value:   idOrFQN,
				},
			},
		}, nil
	}

	if !force {
		isSure, err := cb.prompter.Ask("Are you sure? [y/n]")

		if err != nil || !strings.EqualFold("y", isSure) {
			return DeleteModuleVersionV1Response{}, ErrFailedToConfirmAction{
				Action: "delete_module_version:" + idOrFQN,
			}
		}
	}

	if fqnParseErr == nil {
		dto := DeleteModuleVersionByFqnV1DTO{
			FQN: fqn,
		}

		return cb.DeleteModuleVersionV1ByFqn(dto)
	}

	dto := DeleteModuleVersionV1DTO{
		Id: idOrFQN,
	}

	return cb.DeleteModuleVersionV1ById(dto)
}

func (cb *CommandBus) DeleteModuleVersionV1ByFqn(dto DeleteModuleVersionByFqnV1DTO) (DeleteModuleVersionV1Response, error) {
	cmd := deleteModuleVersionByFqnV1Command{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger, cb.buildValidator(cb.logger))
}

func (cb *CommandBus) DeleteModuleVersionV1ById(dto DeleteModuleVersionV1DTO) (DeleteModuleVersionV1Response, error) {
	cmd := deleteModuleVersionV1Command{
		DTO: dto,
	}

	return cmd.handle(cb.repo, cb.logger, cb.buildValidator(cb.logger))
}
