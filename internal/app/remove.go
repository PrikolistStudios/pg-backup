package app

import (
	"database/sql"
	"fmt"
	"os"
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

func RemoveDatabases(patterns []string, config Config) error {
	// Connect.
	conn, err := createConnection(config)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: unable to connect to postgres database")
		return fmt.Errorf("remove databases: %w", err)
	}

	defer func(conn *sql.DB) {
		_ = conn.Close()
	}(conn)

	// Accumulate errors
	acc := NewErrDatabaseRemoval()

	// Loop through patterns and resolve each one
	for _, pattern := range patterns {
		err = removeDatabase(pattern, config.ForceRemove, conn)
		if err != nil {
			acc.Err = append(acc.Err, err)
			acc.Tables = append(acc.Tables, pattern)
		}
	}

	if len(acc.Err) > 0 {
		return acc
	}

	return nil
}
