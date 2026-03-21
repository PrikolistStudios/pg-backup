package app

import (
	"errors"
	"fmt"
	"strings"
)

var ErrBackup = errors.New("database backup error")

type ErrAccumulatedErrors struct {
	Err   []error
	Items []string
}

func (e ErrAccumulatedErrors) Error() string {
	var sb strings.Builder
	for i, err := range e.Err {
		sb.WriteString(fmt.Sprintf("%s: %s\n", e.Items[i], err))
	}
	return sb.String()
}

func (e ErrAccumulatedErrors) Unwrap() []error {
	return e.Err
}

func NewErrAccumulatedErrors() ErrAccumulatedErrors {
	return ErrAccumulatedErrors{
		Err:   make([]error, 0),
		Items: make([]string, 0),
	}
}
