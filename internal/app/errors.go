package app

import (
	"errors"
	"fmt"
	"strings"
)

var ErrBackup = errors.New("database backup error")

type ErrDatabaseAction struct {
	Err    []error
	Tables []string
}

func (e ErrDatabaseAction) Error() string {
	return fmt.Sprintf("failed to perform action on databases: %s", strings.Join(e.Tables, ","))
}

func (e ErrDatabaseAction) Unwrap() []error {
	return e.Err
}

func NewErrDatabaseRemoval() ErrDatabaseAction {
	return ErrDatabaseAction{
		Err:    make([]error, 0),
		Tables: make([]string, 0),
	}
}
