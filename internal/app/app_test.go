package app

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq" // To register the driver.
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupContainerConnection(t *testing.T) (*sql.DB, string, func()) {
	ctx := context.Background()
	re := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "./",
			Dockerfile: "TestDockerfile",
			KeepImage:  true,
		},

		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
			dsn := fmt.Sprintf("host=%s port=%s user=test_user password=test_password dbname=testdb sslmode=disable", host, port.Port())
			return dsn
		}).WithStartupTimeout(time.Second * 30),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: re,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("host=%s port=%s user=test_user password=test_password dbname=testdb sslmode=disable", host, port.Port())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	return db, dsn, func() {
		db.Close()
		container.Terminate(ctx)
	}
}

// Creates test database.
func createTestDb(name string, db *sql.DB) {
	q := fmt.Sprintf("create database %s;", name)
	_, err := db.Query(q)
	if err != nil {
		panic(err)
	}
}

// Gets all database names on server (?).
func getDatabases(db *sql.DB) []string {
	rows, err := db.Query("SELECT datname FROM pg_database;")
	if err != nil {
		panic(err)
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var result []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			panic(err)
		}
		result = append(result, name)
	}

	return result
}

// It Works test.
func TestRemoveDatabase(t *testing.T) {
	conn, dsn, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	dbname := "test_db_1"
	createTestDb(dbname, conn)

	require.Contains(t, getDatabases(conn), dbname)

	err := RemoveDatabases([]string{dbname}, dsn)
	require.NoError(t, err)

	require.NotContains(t, getDatabases(conn), dbname)
}

// Remove non-existing.

// Remove non-owned.
