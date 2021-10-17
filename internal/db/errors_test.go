package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DriverNotImplemented(t *testing.T) {
	err := ErrDbDriverNotImplemented{
		Driver: "somedriver",
	}

	assert.Equal(t, "database driver is not implemented: somedriver", err.Error())
}
