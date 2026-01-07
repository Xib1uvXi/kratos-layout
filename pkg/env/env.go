package env

import (
	"os"
	"strconv"
	"time"
)

// Get returns the value of the environment variable.
// Returns empty string if the variable is not set.
func Get(key string) string {
	return os.Getenv(key)
}

// GetOrDefault returns the value of the environment variable.
// If the variable is not set, it returns the default value.
func GetOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Lookup returns the value of the environment variable and whether it exists.
func Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

// GetInt returns the value of the environment variable as an integer.
// Returns an error if the variable is not set or cannot be parsed.
func GetInt(key string) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, &NotSetError{Key: key}
	}
	return strconv.Atoi(value)
}

// GetIntOrDefault returns the value of the environment variable as an integer.
// If the variable is not set or cannot be parsed, it returns the default value.
func GetIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetInt64 returns the value of the environment variable as an int64.
// Returns an error if the variable is not set or cannot be parsed.
func GetInt64(key string) (int64, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, &NotSetError{Key: key}
	}
	return strconv.ParseInt(value, 10, 64)
}

// GetInt64OrDefault returns the value of the environment variable as an int64.
// If the variable is not set or cannot be parsed, it returns the default value.
func GetInt64OrDefault(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetFloat64 returns the value of the environment variable as a float64.
// Returns an error if the variable is not set or cannot be parsed.
func GetFloat64(key string) (float64, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, &NotSetError{Key: key}
	}
	return strconv.ParseFloat(value, 64)
}

// GetFloat64OrDefault returns the value of the environment variable as a float64.
// If the variable is not set or cannot be parsed, it returns the default value.
func GetFloat64OrDefault(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetBool returns the value of the environment variable as a boolean.
// Accepts: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
// Returns an error if the variable is not set or cannot be parsed.
func GetBool(key string) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return false, &NotSetError{Key: key}
	}
	return strconv.ParseBool(value)
}

// GetBoolOrDefault returns the value of the environment variable as a boolean.
// If the variable is not set or cannot be parsed, it returns the default value.
func GetBoolOrDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetDuration returns the value of the environment variable as a time.Duration.
// The value should be a valid duration string (e.g., "1h", "30m", "10s").
// Returns an error if the variable is not set or cannot be parsed.
func GetDuration(key string) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, &NotSetError{Key: key}
	}
	return time.ParseDuration(value)
}

// GetDurationOrDefault returns the value of the environment variable as a time.Duration.
// If the variable is not set or cannot be parsed, it returns the default value.
func GetDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// NotSetError is returned when an environment variable is not set.
type NotSetError struct {
	Key string
}

func (e *NotSetError) Error() string {
	return "environment variable not set: " + e.Key
}
