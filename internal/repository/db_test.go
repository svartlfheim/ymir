package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_wrapTransactionError(t *testing.T) {
	err := errors.New("fake error")

	subj := wrapTransactionError(err)

	assert.Equal(t, ErrDbTransaction{Wrapped: err}, subj)
}

func Test_wrapQueryError(t *testing.T) {
	err := errors.New("fake error")

	subj := wrapQueryError(err)

	assert.Equal(t, ErrDbQuery{Wrapped: err}, subj)
}

func Test_wrapHydrationError(t *testing.T) {
	err := errors.New("fake error")

	subj := wrapHydrationError("mytype", err)

	assert.Equal(t, ErrDbHydration{Wrapped: err, Type: "mytype"}, subj)
}