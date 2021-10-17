package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ErrDriverNotImplemented(t *testing.T) {
	err := ErrDriverNotImplemented{
		Driver: "blah",
	}

	assert.Equal(t, "driver: 'blah' is not implemented", err.Error())

	err = ErrDriverNotImplemented{}

	assert.Equal(t, "storage driver was not set", err.Error())
}

func Test_ErrDbTransaction(t *testing.T) {
	wrapped := errors.New("wrapped error")
	err := ErrDbTransaction{
		Wrapped: wrapped,
	}

	assert.Equal(t, "error during database transaction: wrapped error", err.Error())
}

func Test_ErrDbQuery(t *testing.T) {
	wrapped := errors.New("wrapped error")
	err := ErrDbQuery{
		Wrapped: wrapped,
	}

	assert.Equal(t, "error during database query: wrapped error", err.Error())
}


func Test_ErrDbHydration(t *testing.T) {
	wrapped := errors.New("wrapped error")
	err := ErrDbHydration{
		Wrapped: wrapped,
		Type: "MyType",
	}

	assert.Equal(t, "error during database hydration for type MyType: wrapped error", err.Error())
}
