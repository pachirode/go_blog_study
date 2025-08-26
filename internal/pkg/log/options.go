package log

import (
	"go.uber.org/zap/zapcore"
)

// Options 日志配置选项
type Options struct {
	// DisableCaller 是否禁止 caller 信息
	DisableCaller bool
	// DisableStacktrace 是否禁止 stacktrace
	DisableStacktrace bool
	// Level 日志级别
	// 可选：debug,info,warn,error,panic,fatal;默认info
	Level string
	// Format 日志格式
	// 可选：console,json;默认console
	Format string
	// OutputPaths 日志输出路径
	// 可选：stdout,stderr或者也可以指定文件路径;默认stdout
	OutputPaths []string
}

// NewOptions 创建并返回一个带有默认值的 Options 对象
func NewOptions() *Options {
	return &Options{
		DisableCaller:     false,
		DisableStacktrace: false,
		Level:             zapcore.InfoLevel.String(),
		Format:            "console",
		OutputPaths:       []string{"stdout"},
	}
}
