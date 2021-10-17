package repository

import "fmt"

type ErrDriverNotImplemented struct {
	Driver string
}

func (e ErrDriverNotImplemented) Error() string {
	if e.Driver == "" {
		return "storage driver was not set"
	}

	return fmt.Sprintf("driver: '%s' is not implemented", e.Driver)
}

type ErrDbTransaction struct {
	Wrapped error
}

func (e ErrDbTransaction) Error() string {
	return fmt.Sprintf("error during database transaction: %s", e.Wrapped.Error())
}

type ErrDbQuery struct {
	Wrapped error
}

func (e ErrDbQuery) Error() string {
	return fmt.Sprintf("error during database query: %s", e.Wrapped.Error())
}

type ErrDbHydration struct {
	Type    string
	Wrapped error
}

func (e ErrDbHydration) Error() string {
	return fmt.Sprintf("error during database hydration for type %s: %s", e.Type, e.Wrapped.Error())
}
