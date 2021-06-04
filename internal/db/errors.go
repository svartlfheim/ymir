package db

import "fmt"

type ErrDbDriverNotImplemented struct {
	Driver string
}

func (e ErrDbDriverNotImplemented) Error() string {
	return fmt.Sprintf("database driver is not implemented: %s", e.Driver)
}
