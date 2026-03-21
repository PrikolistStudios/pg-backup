package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBackup(t *testing.T) {
	conn, _, config, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	dbname := "test_db_1"
	createTestDb(dbname, conn)

	dbs, _ := getDatabases(conn)
	require.Contains(t, dbs, dbname)

	err := BackupDatabases([]string{dbname}, config)
	require.NoError(t, err)

	// Check that dump was created.
	require.FileExists(t, dbname+".backup")
	_ = os.Remove(dbname + ".backup")
}

// Backup non-existing.
func TestBackupNonexisting(t *testing.T) {
	_, _, config, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	dbname := "nonexisting"

	err := BackupDatabases([]string{dbname}, config)
	require.Error(t, err)
	require.NoFileExists(t, dbname+".backup")
}
