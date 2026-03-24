package app

import (
	"database/sql"
	"fmt"

	"github.com/gobwas/glob"
)

var ErrNoMatch = fmt.Errorf("pattern did not have any matches")

func getDatabases(conn *sql.DB) ([]string, error) {
	rows, err := conn.Query("SELECT datname FROM pg_database;")
	if err != nil {
		return nil, fmt.Errorf("getting databases: %w", err)
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var result []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("getting databases: %w", err)
		}
		result = append(result, name)
	}

	return result, nil
}

// FilterPatterns
// Accepts array of glob patterns and returns array of existing databases with names matching.
func FilterPatterns(patterns []string, conn *sql.DB) ([]string, error) {
	// Get databases.
	databases, err := getDatabases(conn)
	if err != nil {
		return nil, fmt.Errorf("filtering patterns: %w", err)
	}

	// Try to apply each pattern to each database.
	acc := NewErrAccumulatedErrors()
	set := make(map[string]struct{})
	for _, pattern := range patterns {
		gb, err := glob.Compile(pattern)
		if err != nil {
			acc.Err = append(acc.Err, err)
			acc.Items = append(acc.Items, pattern)
			continue
		}

		matched := false
		for _, db := range databases {
			if gb.Match(db) {
				set[db] = struct{}{}
				matched = true
			}
		}

		if !matched {
			acc.Items = append(acc.Items, pattern)
			acc.Err = append(acc.Err, ErrNoMatch)
		}
	}

	// Create result.
	result := make([]string, 0, len(set))
	for db := range set {
		result = append(result, db)
	}

	err = nil
	if len(acc.Err) > 0 {
		err = acc
	}

	return result, err
}
