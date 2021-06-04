package ymir

import "fmt"

type ErrNoArgAtIndex struct {
	index int
}

func (e ErrNoArgAtIndex) Error() string {
	return fmt.Sprintf("not arg was supplied at index: %d", e.index)
}
