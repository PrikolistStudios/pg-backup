package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateConnection(t *testing.T) {
	_, _, config, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	conn, err := CreateConnection(config)
	require.NoError(t, err)
	require.NotNil(t, conn)

	// Verify user.
	rows, _ := conn.Query("select current_user;")
	rows.Next()
	user := ""
	_ = rows.Scan(&user)
	require.Equal(t, config.User, user)

	// Verify database.
	rows, _ = conn.Query("select current_database();")
	rows.Next()
	db := ""
	_ = rows.Scan(&db)
	require.Equal(t, config.Database, db)
}

func TestBrokenConnection(t *testing.T) {
	_, _, config, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	config.Password = "wrong"
	_, err := CreateConnection(config)
	require.Error(t, err)
}
