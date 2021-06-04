package output

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

type Fmt struct {
	writer io.Writer
}

func green() *color.Color {
	return color.New(color.FgGreen)
}

func yellow() *color.Color {
	return color.New(color.FgYellow)
}

func red() *color.Color {
	return color.New(color.FgRed)
}

func (f *Fmt) Infoln(a ...interface{}) (int, error) {
	return fmt.Fprintln(f.writer, a...)
}

func (f *Fmt) Infof(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(f.writer, format, a...)
}

func (f *Fmt) Info(a ...interface{}) (int, error) {
	return fmt.Fprint(f.writer, a...)
}

func (f *Fmt) Warnln(a ...interface{}) (int, error) {
	return yellow().Fprintln(f.writer, a...)
}

func (f *Fmt) Warnf(format string, a ...interface{}) (int, error) {
	return yellow().Fprintf(f.writer, format, a...)
}

func (f *Fmt) Warn(a ...interface{}) (int, error) {
	return yellow().Fprint(f.writer, a...)
}

func (f *Fmt) Successln(a ...interface{}) (int, error) {
	return green().Fprintln(f.writer, a...)
}

func (f *Fmt) Successf(format string, a ...interface{}) (int, error) {
	return green().Fprintf(f.writer, format, a...)
}

func (f *Fmt) Success(a ...interface{}) (int, error) {
	return green().Fprint(f.writer, a...)
}

func (f *Fmt) Errorln(a ...interface{}) (int, error) {
	return red().Fprintln(f.writer, a...)
}

func (f *Fmt) Errorf(format string, a ...interface{}) (int, error) {
	return red().Fprintf(f.writer, format, a...)
}

func (f *Fmt) Error(a ...interface{}) (int, error) {
	return red().Fprint(f.writer, a...)
}

func NewFmt(w io.Writer) *Fmt {
	return &Fmt{
		writer: w,
	}
}
