package env

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testKey = "TEST_ENV_VAR"

func TestGet(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	assert.Equal(t, "", Get(testKey))

	// Set
	os.Setenv(testKey, "hello")
	assert.Equal(t, "hello", Get(testKey))

	os.Unsetenv(testKey)
}

func TestGetOrDefault(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set, returns default
	assert.Equal(t, "default", GetOrDefault(testKey, "default"))

	// Set, returns value
	os.Setenv(testKey, "value")
	assert.Equal(t, "value", GetOrDefault(testKey, "default"))

	os.Unsetenv(testKey)
}

func TestLookup(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	value, exists := Lookup(testKey)
	assert.False(t, exists)
	assert.Equal(t, "", value)

	// Set
	os.Setenv(testKey, "value")
	value, exists = Lookup(testKey)
	assert.True(t, exists)
	assert.Equal(t, "value", value)

	// Set to empty string
	os.Setenv(testKey, "")
	value, exists = Lookup(testKey)
	assert.True(t, exists)
	assert.Equal(t, "", value)

	os.Unsetenv(testKey)
}

func TestGetInt(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	_, err := GetInt(testKey)
	require.Error(t, err)
	assert.IsType(t, &NotSetError{}, err)

	// Valid int
	os.Setenv(testKey, "42")
	val, err := GetInt(testKey)
	require.NoError(t, err)
	assert.Equal(t, 42, val)

	// Invalid int
	os.Setenv(testKey, "not-a-number")
	_, err = GetInt(testKey)
	require.Error(t, err)

	os.Unsetenv(testKey)
}

func TestGetIntOrDefault(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	assert.Equal(t, 100, GetIntOrDefault(testKey, 100))

	// Valid int
	os.Setenv(testKey, "42")
	assert.Equal(t, 42, GetIntOrDefault(testKey, 100))

	// Invalid int, returns default
	os.Setenv(testKey, "invalid")
	assert.Equal(t, 100, GetIntOrDefault(testKey, 100))

	os.Unsetenv(testKey)
}

func TestGetInt64(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	_, err := GetInt64(testKey)
	require.Error(t, err)

	// Valid int64
	os.Setenv(testKey, "9223372036854775807")
	val, err := GetInt64(testKey)
	require.NoError(t, err)
	assert.Equal(t, int64(9223372036854775807), val)

	os.Unsetenv(testKey)
}

func TestGetInt64OrDefault(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	assert.Equal(t, int64(1000), GetInt64OrDefault(testKey, 1000))

	// Valid int64
	os.Setenv(testKey, "9999")
	assert.Equal(t, int64(9999), GetInt64OrDefault(testKey, 1000))

	os.Unsetenv(testKey)
}

func TestGetFloat64(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	_, err := GetFloat64(testKey)
	require.Error(t, err)

	// Valid float
	os.Setenv(testKey, "3.14")
	val, err := GetFloat64(testKey)
	require.NoError(t, err)
	assert.InDelta(t, 3.14, val, 0.001)

	os.Unsetenv(testKey)
}

func TestGetFloat64OrDefault(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	assert.InDelta(t, 1.5, GetFloat64OrDefault(testKey, 1.5), 0.001)

	// Valid float
	os.Setenv(testKey, "2.5")
	assert.InDelta(t, 2.5, GetFloat64OrDefault(testKey, 1.5), 0.001)

	// Invalid float, returns default
	os.Setenv(testKey, "invalid")
	assert.InDelta(t, 1.5, GetFloat64OrDefault(testKey, 1.5), 0.001)

	os.Unsetenv(testKey)
}

func TestGetBool(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	_, err := GetBool(testKey)
	require.Error(t, err)

	// True values
	for _, v := range []string{"1", "t", "T", "true", "TRUE", "True"} {
		os.Setenv(testKey, v)
		val, err := GetBool(testKey)
		require.NoError(t, err)
		assert.True(t, val, "expected true for %q", v)
	}

	// False values
	for _, v := range []string{"0", "f", "F", "false", "FALSE", "False"} {
		os.Setenv(testKey, v)
		val, err := GetBool(testKey)
		require.NoError(t, err)
		assert.False(t, val, "expected false for %q", v)
	}

	os.Unsetenv(testKey)
}

func TestGetBoolOrDefault(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	assert.True(t, GetBoolOrDefault(testKey, true))
	assert.False(t, GetBoolOrDefault(testKey, false))

	// Valid bool
	os.Setenv(testKey, "true")
	assert.True(t, GetBoolOrDefault(testKey, false))

	os.Setenv(testKey, "false")
	assert.False(t, GetBoolOrDefault(testKey, true))

	// Invalid bool, returns default
	os.Setenv(testKey, "invalid")
	assert.True(t, GetBoolOrDefault(testKey, true))

	os.Unsetenv(testKey)
}

func TestGetDuration(t *testing.T) {
	os.Unsetenv(testKey)

	// Not set
	_, err := GetDuration(testKey)
	require.Error(t, err)

	// Valid duration
	os.Setenv(testKey, "1h30m")
	val, err := GetDuration(testKey)
	require.NoError(t, err)
	assert.Equal(t, 90*time.Minute, val)

	// Another valid duration
	os.Setenv(testKey, "30s")
	val, err = GetDuration(testKey)
	require.NoError(t, err)
	assert.Equal(t, 30*time.Second, val)

	os.Unsetenv(testKey)
}

func TestGetDurationOrDefault(t *testing.T) {
	os.Unsetenv(testKey)

	defaultDuration := 5 * time.Minute

	// Not set
	assert.Equal(t, defaultDuration, GetDurationOrDefault(testKey, defaultDuration))

	// Valid duration
	os.Setenv(testKey, "10m")
	assert.Equal(t, 10*time.Minute, GetDurationOrDefault(testKey, defaultDuration))

	// Invalid duration, returns default
	os.Setenv(testKey, "invalid")
	assert.Equal(t, defaultDuration, GetDurationOrDefault(testKey, defaultDuration))

	os.Unsetenv(testKey)
}

func TestNotSetError(t *testing.T) {
	err := &NotSetError{Key: "MY_VAR"}
	assert.Equal(t, "environment variable not set: MY_VAR", err.Error())
}
