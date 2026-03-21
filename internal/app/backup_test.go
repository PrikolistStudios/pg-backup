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

	require.Contains(t, getDatabases(conn), dbname)

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

// Backup non-owned
func TestBackupNonOwned(t *testing.T) {
	db, _, config, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	// Create a different user.
	_, err := db.Query(`CREATE ROLE different_user WITH
	LOGIN
	CREATEDB
	CONNECTION LIMIT -1
	PASSWORD '12345678';`)
	require.NoError(t, err)

	// Connect as a new user.
	config.User = "different_user"
	config.Password = "12345678"

	// Should not be owned by this user.
	dbname := "postgres"

	err = RemoveDatabases([]string{dbname}, config)
	require.Error(t, err)
	require.NoFileExists(t, dbname+".backup")
}
