package app

import (
	"database/sql"
	"fmt"
)

// Removes specified database
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

func RemoveDatabases(names []string, force bool, conn *sql.DB) error {
	// Accumulate errors
	acc := NewErrAccumulatedErrors()

	// Loop through patterns and resolve each one
	for _, name := range names {
		err := removeDatabase(name, force, conn)
		if err != nil {
			acc.Err = append(acc.Err, err)
			acc.Items = append(acc.Items, name)
		}
	}

	if len(acc.Err) > 0 {
		return acc
	}

	return nil
}
