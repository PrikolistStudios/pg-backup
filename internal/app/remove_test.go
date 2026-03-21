package app

import (
	"testing"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/stretchr/testify/require"
)

// It Works test.
func TestRemoveDatabase(t *testing.T) {
	conn, _, _, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	dbname := "test_db_1"
	createTestDb(dbname, conn)

	dbs, _ := getDatabases(conn)
	require.Contains(t, dbs, dbname)

	err := RemoveDatabases([]string{dbname}, false, conn)
	require.NoError(t, err)

	dbs, _ = getDatabases(conn)
	require.NotContains(t, dbs, dbname)
}

// Remove non-existing.
func TestRemoveNonexisting(t *testing.T) {
	conn, _, _, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	dbname := "nonexisting"

	err := RemoveDatabases([]string{dbname}, false, conn)
	require.Error(t, err)

	pqerr := &pq.Error{}
	require.ErrorAs(t, err, &pqerr)
	require.Equal(t, pqerr.Code, pqerror.InvalidCatalogName)
}

// Remove non-owned.
func TestRemoveNoForce(t *testing.T) {
	conn, _, config, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	// Create a different user.
	_, err := conn.Query(`CREATE ROLE different_superuser WITH
  SUPERUSER
	LOGIN
	CREATEDB
	CONNECTION LIMIT -1
	PASSWORD '12345678';`)
	require.NoError(t, err)

	// Connect as a new user.
	config.User = "different_superuser"
	config.Password = "12345678"

	// Try remove the database with a connection.
	dbname := config.Database

	err = RemoveDatabases([]string{dbname}, false, conn)
	require.Error(t, err)

	pqerr := &pq.Error{}
	require.ErrorAs(t, err, &pqerr)
	require.Equal(t, pqerr.Code, pqerror.ObjectInUse)
}

func TestRemoveWithForce(t *testing.T) {
	conn, _, config, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	// Create a different user.
	_, err := conn.Query(`CREATE ROLE different_superuser WITH
  SUPERUSER
	LOGIN
	CREATEDB
	CONNECTION LIMIT -1
	PASSWORD '12345678';`)
	require.NoError(t, err)

	// Connect as a new user.
	config.User = "different_superuser"
	config.Password = "12345678"
	config.ForceRemove = true

	// Try to remove the database with a connection.
	dbname := config.Database

	// Connect to different database while removing the other.
	config.Database = "postgres"
	newConn, err := CreateConnection(config)
	err = RemoveDatabases([]string{dbname}, true, newConn)
	require.NoError(t, err)

	dbs, _ := getDatabases(newConn)
	require.NotContains(t, dbs, dbname)
}

func TestRemoveNonOwned(t *testing.T) {
	conn, _, config, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	// Create a different user.
	_, err := conn.Query(`CREATE ROLE different_user WITH 
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
	newConn, err := CreateConnection(config)
	err = RemoveDatabases([]string{dbname}, false, newConn)
	require.Error(t, err)

	pqerr := &pq.Error{}
	require.ErrorAs(t, err, &pqerr)
	require.Equal(t, pqerr.Code, pqerror.InsufficientPrivilege)
}
