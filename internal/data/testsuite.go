package data

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/pkg/orm"
)

// LocalConfigFileName is the name of the local config file for testing.
const LocalConfigFileName = ".local.config.yaml"

// TestSuite provides test utilities for data layer testing.
// It manages MySQL and Redis connections and provides cleanup methods.
type TestSuite struct {
	conf    *conf.Data
	db      orm.DB
	utilDB  orm.DBUtil
	redis   *redis.Client
	cleanup []func()
}

// NewTestSuiteFromLocal creates a new test suite loading config from configs/.local.config.yaml
func NewTestSuiteFromLocal() (*TestSuite, error) {
	configPath, err := findLocalConfig()
	if err != nil {
		return nil, err
	}

	c := config.New(
		config.WithSource(
			file.NewSource(configPath),
		),
	)

	if err := c.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, fmt.Errorf("failed to scan config: %w", err)
	}

	if err := c.Close(); err != nil {
		return nil, fmt.Errorf("failed to close config: %w", err)
	}

	return &TestSuite{
		conf:    bc.Data,
		cleanup: make([]func(), 0),
	}, nil
}

// findLocalConfig searches for the local config file.
func findLocalConfig() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// Navigate from internal/data to project root
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	configPath := filepath.Join(projectRoot, "configs", LocalConfigFileName)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", fmt.Errorf("local config file not found: %s\nPlease create configs/%s for local testing", configPath, LocalConfigFileName)
	}

	return configPath, nil
}

// Setup initializes all database connections.
// Call this in TestMain or at the beginning of your test.
func (ts *TestSuite) Setup() error {
	if err := ts.validateConfig(); err != nil {
		return err
	}

	if err := ts.setupMySQL(); err != nil {
		return fmt.Errorf("setup mysql failed: %w", err)
	}

	if err := ts.setupRedis(); err != nil {
		return fmt.Errorf("setup redis failed: %w", err)
	}

	return nil
}

// TearDown cleans up all data and closes connections.
// Call this at the end of your test or in TestMain.
func (ts *TestSuite) TearDown() error {
	var errs []error

	if err := ts.ClearAll(); err != nil {
		errs = append(errs, fmt.Errorf("clear all data failed: %w", err))
	}

	for _, fn := range ts.cleanup {
		fn()
	}

	if len(errs) > 0 {
		return fmt.Errorf("teardown errors: %v", errs)
	}
	return nil
}

// ClearAll clears both MySQL and Redis data.
func (ts *TestSuite) ClearAll() error {
	if err := ts.ClearMySQL(); err != nil {
		return err
	}
	if err := ts.ClearRedis(); err != nil {
		return err
	}
	return nil
}

// ClearMySQL clears all data from MySQL tables.
func (ts *TestSuite) ClearMySQL() error {
	if ts.db == nil {
		return nil
	}
	return ts.db.ClearAllData()
}

// ClearRedis clears all data from Redis.
func (ts *TestSuite) ClearRedis() error {
	if ts.redis == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return ts.redis.FlushDB(ctx).Err()
}

// DB returns the gorm.DB instance for test queries.
func (ts *TestSuite) DB() *gorm.DB {
	if ts.db == nil {
		return nil
	}
	return ts.db.GetDB()
}

// Redis returns the redis.Client instance for test operations.
func (ts *TestSuite) Redis() *redis.Client {
	return ts.redis
}

func (ts *TestSuite) validateConfig() error {
	if ts.conf == nil {
		return fmt.Errorf("config is nil")
	}

	dbName := ts.conf.Database.DbName
	if !strings.Contains(dbName, "test") && !strings.Contains(dbName, "dev") {
		return fmt.Errorf("database name must contain 'test' or 'dev' for safety, got: %s", dbName)
	}

	return nil
}

func (ts *TestSuite) setupMySQL() error {
	db := ts.conf.Database
	dbConfig := &orm.DBConfig{
		Username:        db.Username,
		Password:        db.Password,
		Host:            db.Host,
		Port:            fmt.Sprintf("%d", db.Port),
		DBName:          db.DbName,
		DBCharset:       db.DbCharset,
		MaxIdleConns:    int(db.MaxIdleConns),
		MaxOpenConns:    int(db.MaxOpenConns),
		ConnMaxLifetime: db.ConnMaxLifetime.AsDuration(),
		ConnMaxIdleTime: db.ConnMaxIdleTime.AsDuration(),
	}

	// Create util DB for database management
	utilDB, err := orm.MakeDBUtil(dbConfig)
	if err != nil {
		return fmt.Errorf("create util db failed: %w", err)
	}
	ts.utilDB = utilDB
	ts.cleanup = append(ts.cleanup, func() { utilDB.Close() })

	// Create database if not exists
	if createErr := utilDB.CreateDB(); createErr != nil {
		return fmt.Errorf("create database failed: %w", createErr)
	}

	// Create main DB connection
	ormDB, err := orm.MakeDB(dbConfig)
	if err != nil {
		return fmt.Errorf("create db failed: %w", err)
	}
	ts.db = ormDB
	ts.cleanup = append(ts.cleanup, func() { ormDB.Close() })

	return nil
}

func (ts *TestSuite) setupRedis() error {
	if ts.conf.Redis == nil || ts.conf.Redis.Addr == "" {
		return nil // Redis is optional
	}

	redisCfg := ts.conf.Redis
	ts.redis = redis.NewClient(&redis.Options{
		Addr:         redisCfg.Addr,
		Password:     redisCfg.Password,
		DB:           int(redisCfg.Db),
		DialTimeout:  redisCfg.DialTimeout.AsDuration(),
		ReadTimeout:  redisCfg.ReadTimeout.AsDuration(),
		WriteTimeout: redisCfg.WriteTimeout.AsDuration(),
	})
	ts.cleanup = append(ts.cleanup, func() { ts.redis.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ts.redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis failed: %w", err)
	}

	return nil
}
