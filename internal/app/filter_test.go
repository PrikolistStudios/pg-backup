package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGlobWorks(t *testing.T) {
	conn, _, _, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	db_names := []string{"test_db_1", "test_db_2", "test_db_3"}
	for _, name := range db_names {
		createTestDb(name, conn)
	}

	filtered, err := FilterPatterns([]string{"test_db_*"}, conn)
	require.NoError(t, err)
	require.ElementsMatch(t, db_names, filtered)
}

func TestBrokenGlob(t *testing.T) {
	conn, _, _, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	filtered, err := FilterPatterns([]string{"[ some broken glob ["}, conn)
	require.Len(t, filtered, 0)
	require.Error(t, err)
	//require.ElementsMatch(t, db_names, filtered)
}

func TestNoMatches(t *testing.T) {
	conn, _, _, closeFunc := setupContainerConnection(t)
	defer closeFunc()

	filtered, err := FilterPatterns([]string{"test_db_*"}, conn)
	require.Len(t, filtered, 0)
	require.ErrorIs(t, err, ErrNoMatch)
}
