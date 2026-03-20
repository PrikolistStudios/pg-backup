package app

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// ResolvePattern returns valid database names that fit pattern
func resolvePattern(pattern string) []string {
	panic("implement me")
}

// Removes specified database
func removeDatabase(name string, conn *sql.DB) error {
	q := fmt.Sprintf("drop database %s;", name)
	_, err := conn.Query(q)
	return err
}

func createConnection(url string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.Wrap(err, "create connection: ")
	}

	// Verify that connection is successful.
	err = conn.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "create connection: ")
	}

	return conn, nil
}

func RemoveDatabases(patterns []string, url string) error {
	// Connect.
	conn, err := createConnection(url)
	if err != nil {
		return errors.Wrap(err, "Remove Databases: ")
	}

	defer func(conn *sql.DB) {
		_ = conn.Close()
	}(conn)

	// Loop through patterns and resolve each one
	for _, pattern := range patterns {
		err = removeDatabase(pattern, conn)
		if err != nil {
			return errors.Wrap(err, "Remove Databases: ")
		}
	}

	return nil
}

func BackupDatabases(patterns []string) error {
	panic("implement me")
}
