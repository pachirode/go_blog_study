package log

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestMain(m *testing.M) {
	opts := &Options{
		Level:             "debug",
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            "json",
		OutputPaths:       []string{"stdout"},
	}

	Init(opts)

	m.Run()
}

func TestLogger(t *testing.T) {
	opts := &Options{
		Level:             "debug",
		Format:            "json",
		DisableCaller:     false,
		DisableStacktrace: false,
		OutputPaths:       []string{"stdout"},
	}

	Init(opts)

	Debugw("debug message", "key1", "value1", "key2", "value2")
	Infow("info message", "key1", "value1", "key2", "value2")
	Warnw("warn message", "key1", "value1", "key2", "value2")
	Errorw("error message", "key1", "value1", "something error")

	Sync()
}

func TestLoggerMethods(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() {
		Debugw("debug message")
		Infow("info message")
		Warnw("warn message")
		Errorw("error message")
	}, "Log method should not cause a panic in this test")

	assert.Panics(func() {
		Panicw("panic message")
	}, "Panic method should cause a panic in this test")
}

func TestLoggerInitialization(t *testing.T) {
	opts := NewOptions()
	logger := NewLogger(opts)

	assert.NotNil(t, logger)
	assert.IsType(t, &zapLogger{}, logger)
}

func TestSync(t *testing.T) {
	assert.NotPanics(t, func() {
		Sync()
	}, "Sync should not panic")
}

func BenchmarkZapLogger(b *testing.B) {
	// 使用 zap.NewNop() 模拟 logger
	logger := &zapLogger{z: zap.NewNop()}

	ctx := context.WithValue(context.Background(), "requestID", "12345")
	ctx = context.WithValue(ctx, "userID", "user1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.C(ctx).Debugw("debug message", "key1", "value1", "key2", "value2")
		logger.C(ctx).Infow("info message", "key1", "value1", "key2", "value2")
		logger.C(ctx).Warnw("warn message", "key1", "value1", "key2", "value2")
		logger.C(ctx).Errorw("error message", "key1", "value1", "something error")
	}
}
