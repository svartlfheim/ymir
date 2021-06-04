package registry

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
)

const requiredTag string = "required"
const requiredWithoutTag string = "required_without"
const uniqueModulePathTag string = "unique_module_path"
const moduleMustExistTag string = "module_must_exist"
const uniqueModuleVersionTag string = "unique_module_version"
const noVersionsExistForIdTag string = "no_versions_exist_for_id"
const noVersionsExistForModuleFQNTag string = "no_versions_exist_for_module_fqn"
const uuidTag string = "uuid"
const versionTag string = "version"

type ValidatorBuilder func(l zerolog.Logger) CommandValidator

type CustomValidationMessageHandler func(e validator.FieldError) (string, error)

type commandValidator struct {
	validate              *validator.Validate
	logger                zerolog.Logger
	customMessageHandlers map[string]CustomValidationMessageHandler
}

type ValidationError struct {
	Message string
	Rule    string
	Field   string
	Value   interface{}
}

type CommandValidator interface {
	ModuleNameMustBeUnique(r uniqueModuleNameRepository, sl validator.StructLevel, fqn ModuleFQN)
	ModuleVersionMustBeUniqueForFQN(r moduleVersionMustBeUniqueForFQNRepository, sl validator.StructLevel, fqn ModuleVersionFQN)
	ModuleVersionMustBeUniqueForId(r moduleVersionMustBeUniqueForIdRepository, sl validator.StructLevel, id string, version string)
	ModuleMustExistById(r moduleMustExistByIdRepository, sl validator.StructLevel, id string)
	ModuleMustExistByFQN(r moduleMustExistByFQNRepository, sl validator.StructLevel, fqn ModuleFQN)
	NoVersionsExistForModuleFQN(r noVersionsExistForModuleFQNRepository, sl validator.StructLevel, fqn ModuleFQN)
	NoVersionsExistForModuleId(r noVersionsExistForIdRepository, sl validator.StructLevel, id string)
	SetCustomMessageHandlers(handlers map[string]CustomValidationMessageHandler)
	RegisterStructLevelValidator(f validator.StructLevelFunc, t interface{})
	Validate(cmd interface{}) []ValidationError
}

type uniqueModuleNameRepository interface {
	ByFQN(ModuleFQN) (m Module, err error)
}

func (v *commandValidator) ModuleNameMustBeUnique(r uniqueModuleNameRepository, sl validator.StructLevel, fqn ModuleFQN) {
	_, err := r.ByFQN(fqn)

	if err == nil {
		v.logger.Error().Err(err).Msg("validation error")
		sl.ReportError(fqn.String(), "ns/name/provider", "CompositeField", uniqueModulePathTag, fqn.String())
		return
	}

	if _, ok := err.(ErrResourceNotFound); !ok {
		v.logger.Error().Err(err).Msg("unexpected repository error during validation")
	}
}

type moduleVersionMustBeUniqueForFQNRepository interface {
	VersionByFQN(fqn ModuleVersionFQN) (mv ModuleVersion, err error)
}

func (v *commandValidator) ModuleVersionMustBeUniqueForFQN(r moduleVersionMustBeUniqueForFQNRepository, sl validator.StructLevel, fqn ModuleVersionFQN) {
	_, err := r.VersionByFQN(fqn)

	if _, ok := err.(ErrResourceNotFound); !ok {
		sl.ReportError(fqn.String(), "ns/name/provider@version", "CompositeField", uniqueModuleVersionTag, fqn.String())
	}

	if err != nil {
		v.logger.Error().Err(err).Msg("error during ModuleVersionMustBeUniqueForFQN validation query")
	}
}

type moduleVersionMustBeUniqueForIdRepository interface {
	VersionByModuleAndValue(moduleId string, version string) (mv ModuleVersion, err error)
}

func (v *commandValidator) ModuleVersionMustBeUniqueForId(r moduleVersionMustBeUniqueForIdRepository, sl validator.StructLevel, id string, version string) {
	_, err := r.VersionByModuleAndValue(id, version)

	if _, ok := err.(ErrResourceNotFound); !ok {
		sl.ReportError(id, "uuid@version", "CompositeField", uniqueModuleVersionTag, id)
	}

	if err != nil {
		v.logger.Error().Err(err).Msg("error during ModuleVersionMustBeUniqueForFQN validation query")
	}
}

type moduleMustExistByFQNRepository interface {
	ByFQN(ModuleFQN) (m Module, err error)
}

func (v *commandValidator) ModuleMustExistByFQN(r moduleMustExistByFQNRepository, sl validator.StructLevel, fqn ModuleFQN) {
	_, err := r.ByFQN(fqn)

	if _, ok := err.(ErrResourceNotFound); ok {
		sl.ReportError(fqn.String(), "ns/name/provider", "CompositeField", moduleMustExistTag, fqn.String())
	}
}

type moduleMustExistByIdRepository interface {
	ById(id string) (m Module, err error)
}

func (v *commandValidator) ModuleMustExistById(r moduleMustExistByIdRepository, sl validator.StructLevel, id string) {
	_, err := r.ById(id)

	if _, ok := err.(ErrResourceNotFound); ok {
		sl.ReportError(id, "uuid", "CompositeField", moduleMustExistTag, id)
	}
}

type noVersionsExistForIdRepository interface {
	VersionsByModule(moduleId string, chunkOpts ChunkingOptions) (m []ModuleVersion, err error)
}

func (v *commandValidator) NoVersionsExistForModuleId(r noVersionsExistForIdRepository, sl validator.StructLevel, id string) {
	res, err := r.VersionsByModule(id, ChunkingOptions{})

	if err != nil {
		v.logger.Error().Err(err).Msg("unexpected repository error during validation")
	}

	if err != nil || len(res) > 0 {
		sl.ReportError(id, "id", "CompositeField", noVersionsExistForIdTag, id)
	}
}

type noVersionsExistForModuleFQNRepository interface {
	VersionsByModuleFQN(fqn ModuleFQN, chunkOpts ChunkingOptions) (m []ModuleVersion, err error)
}

func (v *commandValidator) NoVersionsExistForModuleFQN(r noVersionsExistForModuleFQNRepository, sl validator.StructLevel, fqn ModuleFQN) {
	res, err := r.VersionsByModuleFQN(fqn, ChunkingOptions{})

	if err != nil {
		v.logger.Error().Err(err).Msg("unexpected repository error during validation")
	}

	if err != nil || len(res) > 0 {
		sl.ReportError(fqn.String(), "id", "CompositeField", noVersionsExistForModuleFQNTag, fqn.String())
	}
}

func (v *commandValidator) RegisterStructLevelValidator(f validator.StructLevelFunc, t interface{}) {
	v.validate.RegisterStructValidation(f, t)
}

func (v *commandValidator) SetCustomMessageHandlers(handlers map[string]CustomValidationMessageHandler) {
	v.customMessageHandlers = handlers
}

func (v *commandValidator) Validate(cmd interface{}) []ValidationError {
	errs := []ValidationError{}
	err := v.validate.Struct(cmd)

	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			v.logger.Fatal().Err(err).Msg("error from validator")
		}

		for _, fieldErr := range err.(validator.ValidationErrors) {
			var msg string
			var msgErr error
			h, ok := v.customMessageHandlers[fieldErr.Tag()]

			if ok {
				msg, msgErr = h(fieldErr)
			}

			if _, wasRejectedByHandler := msgErr.(ErrRejectedByMessageHandler); !ok || wasRejectedByHandler {
				msg, msgErr = messageForFieldError(fieldErr)
			}

			if msgErr != nil {
				v.logger.Fatal().Err(msgErr).Str("error-tag", fieldErr.Tag()).Msg("could not generate message for field rule")
			}

			errs = append(errs, ValidationError{
				Message: msg,
				Field:   fieldErr.Field(),
				Rule:    fieldErr.Tag(),
				Value:   fieldErr.Param(),
			})
		}
	}

	return errs
}

func messageForFieldError(e validator.FieldError) (string, error) {
	switch e.Tag() {
	case requiredTag:
		return "is required", nil
	case uniqueModulePathTag:
		return "a module with this name already exists", nil
	case uniqueModuleVersionTag:
		return "version already exists for this module", nil
	case moduleMustExistTag:
		return "the module does not exist", nil
	case requiredWithoutTag:
		oneOf := strings.Split(e.Param(), " ")
		oneOf = append(oneOf, e.Field())
		return fmt.Sprintf("one of [%s], must be supplied", strings.Join(oneOf, ",")), nil
	case uuidTag:
		return "must be a valid uuid", nil
	case versionTag:
		return "must be semver (^[0-9]+\\.[0-9]+\\.[0-9]+$) or prefixed with 'dev-'", nil
	default:
		return "", errors.New("type not implemented")
	}
}

func uuidValidator(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())

	return err == nil
}

func versionValidator(fl validator.FieldLevel) bool {
	val := fl.Field().String()

	// We'll let anything if it is a dev version
	if strings.HasPrefix(val, "dev-") {
		return true
	}

	matched, err := regexp.MatchString(`^[0-9]+\.[0-9]+\.[0-9]+$`, val)

	return matched && err == nil
}

func buildRequiredModuleVersionRuleMessage(e validator.FieldError) (string, error) {
	switch e.StructField() {
	case "Id", "ModuleName", "ModuleVersion", "ModuleNamespace", "ModuleProvider":
		return "id or (ModuleName, ModuleNamespace, ModuleVersion, and ModuleProvider) must be supplied", nil
	default:
		return "", ErrRejectedByMessageHandler{
			Tag:   e.Tag(),
			Field: e.StructField(),
		}
	}
}

func NewCommandValidator(l zerolog.Logger) CommandValidator {
	v := validator.New()
	err := v.RegisterValidation("uuid", uuidValidator)

	if err != nil {
		l.Error().Err(err).Msg("failed to register uuid validator")
	}

	err = v.RegisterValidation("version", versionValidator)
	if err != nil {
		l.Error().Err(err).Msg("failed to register version validator")
	}

	return &commandValidator{
		validate: v,
		logger:   l,
	}
}
