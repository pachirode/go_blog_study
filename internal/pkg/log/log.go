package log

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pachirode/go_blog_study/internal/pkg/known"
)

// zapLogger Logger 接口的具体实现，底层封装了 zap.Logger
type zapLogger struct {
	z *zap.Logger
}

// Ensure zapLogger impletement Logger
var _ Logger = &zapLogger{}

var (
	mu sync.Mutex

	// 定义默认全局 Logger
	std = NewLogger(NewOptions())
)

// Init 初始化全局日志对象
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()

	std = NewLogger(opts)
}

// NewLogger 根据提供的 Options 参数创建一个自定义 zapLogger
func NewLogger(opts *Options) *zapLogger {
	// 没有提供配置选项，则使用默认
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	// 创建 encoder 配置，用于控制日志输出格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("2020-01-02 15:04:05.000"))
	}

	encoderConfig.EncodeDuration = func(d time.Duration, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	cfg := &zap.Config{
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Encoding:          opts.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       opts.OutputPaths,
		ErrorOutputPaths:  []string{"stderr"},
	}

	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{z: z}

	zap.RedirectStdLog(z)

	return logger
}

// Logger 定义日志接口
type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Sync()
}

func Debugw(msg string, keysAndValues ...interface{}) {
	std.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	std.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	std.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	std.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	std.Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	std.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Fatalw(msg, keysAndValues...)
}

// Sync 调用底层 zap.Logger 的 Sync 方法，将缓存日志刷新到磁盘文件中，需要再主程序退出前调用
func Sync() {
	std.Sync()
}

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}

// C 解析传入的 context，并尝试提取关注的键值，添加到日志中
func C(ctx context.Context) *zapLogger {
	return std.C(ctx)
}

func (l *zapLogger) C(ctx context.Context) *zapLogger {
	lc := l.clone()

	// 定义一个映射，关联 context 提取函数和日志字段名

	if requestID := ctx.Value(known.XRequestID); requestID != nil {
		lc.z = lc.z.With(zap.Any(known.XRequestID, requestID))
	}

	if userID := ctx.Value(known.XUserName); userID != nil {
		lc.z = lc.z.With(zap.Any(known.XUserName, userID))
	}

	return lc
}

func (l *zapLogger) clone() *zapLogger {
	lc := *l
	return &lc
}
