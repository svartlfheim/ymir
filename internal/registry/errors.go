package registry

import "fmt"

type ErrResourceNotFound struct {
	Type string
	URI  string
}

func (e ErrResourceNotFound) Error() string {
	return fmt.Sprintf("resource of type '%s', with URI: %s could not be found", e.Type, e.URI)
}

type ErrCouldNotReadFile struct {
	Path string
}

func (e ErrCouldNotReadFile) Error() string {
	return fmt.Sprintf("file could not be read: %s", e.Path)
}

type ErrCouldNotUnmarshalJSONToDTO struct {
	Path    string
	DTOName string
}

func (e ErrCouldNotUnmarshalJSONToDTO) Error() string {
	return fmt.Sprintf("file %s could not be unmarshaled to DTO %s", e.Path, e.DTOName)
}

type ErrQuestionFailed struct {
	Question string
}

func (e ErrQuestionFailed) Error() string {
	return fmt.Sprintf("error during question: %s", e.Question)
}

type ErrCouldNotParseModuleFQN struct {
	Value   string
	Message string
}

func (e ErrCouldNotParseModuleFQN) Error() string {
	return fmt.Sprintf("could not parse '%s' as a ModuleFQN: %s", e.Value, e.Message)
}

type ErrCouldNotParseModuleVersionFQN struct {
	Value   string
	Message string
}

func (e ErrCouldNotParseModuleVersionFQN) Error() string {
	return fmt.Sprintf("could not parse '%s' as a ModuleVersionFQN: %s", e.Value, e.Message)
}

type ErrRejectedByMessageHandler struct {
	Tag   string
	Field string
}

func (e ErrRejectedByMessageHandler) Error() string {
	return fmt.Sprintf("rule '%s' for field '%s' cannot be generated by this handler", e.Tag, e.Field)
}

type ErrFailedToConfirmAction struct {
	Action string
}

func (e ErrFailedToConfirmAction) Error() string {
	return fmt.Sprintf("failed to confirm action: %s", e.Action)
}
