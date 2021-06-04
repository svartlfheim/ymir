package gopoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	s := "blah"

	assert.IsType(t, &s, ToString(s))
	assert.Equal(t, &s, ToString(s))
}

func TestInt(t *testing.T) {
	i := 173

	assert.IsType(t, &i, ToInt(i))
	assert.Equal(t, &i, ToInt(i))
}

func TestBool(t *testing.T) {
	b := false

	assert.IsType(t, &b, ToBool(b))
	assert.Equal(t, &b, ToBool(b))
}
