package data

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testSuite *TestSuite

func TestMain(m *testing.M) {
	var err error
	testSuite, err = NewTestSuiteFromLocal()
	if err != nil {
		panic("failed to load test config: " + err.Error())
	}

	if err := testSuite.Setup(); err != nil {
		panic("failed to setup test suite: " + err.Error())
	}

	code := m.Run()

	if err := testSuite.TearDown(); err != nil {
		panic("failed to teardown test suite: " + err.Error())
	}

	os.Exit(code)
}

func TestTestSuite_MySQL(t *testing.T) {
	db := testSuite.DB()
	require.NotNil(t, db)

	// Test raw query
	var result int
	err := db.Raw("SELECT 1").Scan(&result).Error
	require.NoError(t, err)
	require.Equal(t, 1, result)
}

func TestTestSuite_Redis(t *testing.T) {
	rdb := testSuite.Redis()
	require.NotNil(t, rdb)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test set and get
	err := rdb.Set(ctx, "test_key", "test_value", time.Minute).Err()
	require.NoError(t, err)

	val, err := rdb.Get(ctx, "test_key").Result()
	require.NoError(t, err)
	require.Equal(t, "test_value", val)
}

func TestTestSuite_ClearRedis(t *testing.T) {
	rdb := testSuite.Redis()
	require.NotNil(t, rdb)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set some data
	err := rdb.Set(ctx, "clear_test_key", "value", time.Minute).Err()
	require.NoError(t, err)

	// Clear Redis
	err = testSuite.ClearRedis()
	require.NoError(t, err)

	// Verify data is cleared
	_, err = rdb.Get(ctx, "clear_test_key").Result()
	require.Error(t, err)
}
