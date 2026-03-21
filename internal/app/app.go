package app

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ResolvePattern returns valid database names that fit pattern
func resolvePattern(pattern string) []string {
	panic("implement me")
}

func createConnection(config Config) (*sql.DB, error) {
	conn, err := sql.Open("postgres", getDsn(config))
	if err != nil {
		return nil, fmt.Errorf("create connection: %w", err)
	}

	// Verify that connection is successful.
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("create connection: %w", err)
	}

	return conn, nil
}

func getDsn(config Config) string {
	res := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)

	return res
}
