package app

import (
	_ "github.com/lib/pq"
)

type DatabaseAction func(name string) error

// PerformDatabasesAction performs given action on each database and accumulates returned errors for each database.
func PerformDatabasesAction(names []string, action DatabaseAction) error {
	// Accumulate errors
	acc := NewErrAccumulatedErrors()

	for _, name := range names {
		err := action(name)
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
