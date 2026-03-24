package app

import (
	"database/sql"
	"fmt"
)

func removeDatabase(name string, force bool, conn *sql.DB) error {
	q := ""
	if force {
		q = fmt.Sprintf("drop database %s with (force);", name)
	} else {
		q = fmt.Sprintf("drop database %s;", name)
	}
	_, err := conn.Query(q)
	return err
}

func NewRemoveAction(force bool, conn *sql.DB) DatabaseAction {
	return func(name string) error {
		return removeDatabase(name, force, conn)
	}
}
