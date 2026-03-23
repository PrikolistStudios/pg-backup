package app

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq" // To register the driver.
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupContainerConnection(t *testing.T) (*sql.DB, string, Config, func()) {
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

	config := Config{
		Host:     host,
		Port:     port.Port(),
		Database: "testdb",
		User:     "test_user",
		Password: "test_password",
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
	return db, dsn, config, func() {
		_ = db.Close()
		_ = container.Terminate(ctx)
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
