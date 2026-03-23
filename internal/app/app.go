package app

import (
	_ "github.com/lib/pq"
)

type DatabaseAction func(name string) error

// PerformDatabasesAction performs given action on each database and accumulates returned errors for each database.
func PerformDatabasesAction(names []string, action DatabaseAction) error {
	// Accumulate errors
	acc := NewErrAccumulatedErrors()

	errChan := make(chan error, len(names))

	// Launch each action concurrently.
	for _, name := range names {
		go func() {
			errChan <- action(name)
		}()
	}

	// Gather errors
	for _, name := range names {
		err := <-errChan
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
