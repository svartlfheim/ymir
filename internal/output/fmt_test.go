package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
COLOUR output

The fmt struct should change the colour of output in the CLI.
This is difficult to test here, as the colour codes don't seem to appear in the buffered output

Maybe look into it later, but it's a pretty minor issue
*/

func Test_Info(t *testing.T) {
	t.Run("basic info text", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Info("blah")

		assert.Equal(t, "blah", buff.String())
	})

	t.Run("basic infoln", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Infoln("blah")

		assert.Equal(t, "blah\n", buff.String())
	})

	t.Run("basic infof", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Infof("this is int: %d, and this is string: %s", 45, "someval")

		assert.Equal(t, "this is int: 45, and this is string: someval", buff.String())
	})
}

func Test_Warn(t *testing.T) {
	t.Run("basic warn text", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Warn("blah")

		assert.Equal(t, "blah", buff.String())
	})

	t.Run("basic warnln", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Warnln("blah")

		assert.Equal(t, "blah\n", buff.String())
	})

	t.Run("basic warnf", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Warnf("this is int: %d, and this is string: %s", 45, "someval")

		assert.Equal(t, "this is int: 45, and this is string: someval", buff.String())
	})
}

func Test_Success(t *testing.T) {
	t.Run("success warn text", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Success("blah")

		assert.Equal(t, "blah", buff.String())
	})

	t.Run("success warnln", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Successln("blah")

		assert.Equal(t, "blah\n", buff.String())
	})

	t.Run("success warnf", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Successf("this is int: %d, and this is string: %s", 45, "someval")

		assert.Equal(t, "this is int: 45, and this is string: someval", buff.String())
	})
}

func Test_Error(t *testing.T) {
	t.Run("error warn text", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Error("blah")

		assert.Equal(t, "blah", buff.String())
	})

	t.Run("error warnln", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Errorln("blah")

		assert.Equal(t, "blah\n", buff.String())
	})

	t.Run("error warnf", func(tt *testing.T) {
		buff := new(bytes.Buffer)
		fmt := NewFmt(buff)

		fmt.Errorf("this is int: %d, and this is string: %s", 45, "someval")

		assert.Equal(t, "this is int: 45, and this is string: someval", buff.String())
	})
}
