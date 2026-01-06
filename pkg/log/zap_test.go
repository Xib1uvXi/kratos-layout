package log

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapLogger(t *testing.T) {
	eConfig := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	logger := NewZapLogger(
		zapcore.NewConsoleEncoder(eConfig),
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)
	zlog := log.NewHelper(logger)
	zlog.Infow("name", "kratos", "from", "opensource")
	zlog.Infow("name", "kratos", "from")

	// zap stdout/stderr Sync bugs in OSX, see https://github.com/uber-go/zap/issues/370
	_ = logger.Sync()
}

func TestInitDefaultLogger(t *testing.T) {
	logger := InitDefaultLogger(zapcore.DebugLevel)
	logger.Log(log.LevelDebug, "name", "kratos", "from", "opensource")
}

func TestInitJSONLogger(t *testing.T) {
	logger := InitJSONLogger(zapcore.DebugLevel)
	// Use log.Helper to test caller correctly, as CallerSkip(2) is designed for it
	helper := log.NewHelper(logger)
	helper.Info("test json logger")
}
