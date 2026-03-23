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

	action := NewRemoveAction(false, conn)
	err := PerformDatabasesAction([]string{dbname}, action)
	require.NoError(t, err)

	dbs, _ = getDatabases(conn)
	require.NotContains(t, dbs, dbname)
}

// Remove non-existing.
func TestRemoveNonexisting(t *testing.T) {
	conn, _, _, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	dbname := "nonexisting"

	action := NewRemoveAction(false, conn)
	err := PerformDatabasesAction([]string{dbname}, action)
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

	action := NewRemoveAction(false, conn)
	err = PerformDatabasesAction([]string{dbname}, action)
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
	require.NoError(t, err)
	action := NewRemoveAction(true, newConn)
	err = PerformDatabasesAction([]string{dbname}, action)
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
	require.NoError(t, err)
	action := NewRemoveAction(false, newConn)
	err = PerformDatabasesAction([]string{dbname}, action)
	require.Error(t, err)

	pqerr := &pq.Error{}
	require.ErrorAs(t, err, &pqerr)
	require.Equal(t, pqerr.Code, pqerror.InsufficientPrivilege)
}
