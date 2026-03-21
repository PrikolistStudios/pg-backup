package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func backupDatabase(name string, config Config) error {
	config.Database = name
	cmd := exec.Command("pg_dump", getDsn(config))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("backup database: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("backup database: %w", err)
	}

	out, err := io.ReadAll(stdout)
	if err != nil {
		return fmt.Errorf("backup database: %w", err)
	}

	exitErr := cmd.Wait()
	if exitErr != nil {
		return fmt.Errorf("backup database: %w", exitErr)
	}

	err = os.WriteFile(name+".backup", out, 0644)
	if err != nil {
		return fmt.Errorf("backup database: %w", err)
	}

	return nil
}

// BackupDatabases When using pgdump the configuration is kept except the database name.
func BackupDatabases(names []string, config Config) error {
	// Accumulate errors
	acc := NewErrAccumulatedErrors()

	for _, name := range names {
		err := backupDatabase(name, config)
		if err != nil {
			acc.Err = append(acc.Err, ErrBackup)
			acc.Items = append(acc.Items, name)
		}
	}

	if len(acc.Err) > 0 {
		return acc
	}

	return nil
}
