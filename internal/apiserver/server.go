package apiserver

import (
	"github.com/pachirode/go_blog_study/internal/pkg/server"
	"time"

	genericOptions "github.com/pachirode/pkg/options"
)

const (
	// GRPCServerMode 定义 gRPC 服务模式.
	GRPCServerMode = "grpc"
	// GRPCGatewayServerMode 定义 gRPC + HTTP 服务模式.
	GRPCGatewayServerMode = "grpc-gateway"
	// GinServerMode 定义 Gin 服务模式.
	GinServerMode = "gin"
)

// Config 配置结构体，用于存储应用相关的配置.
// 不用 viper.Get，是因为这种方式能更加清晰的知道应用提供了哪些配置项.
type Config struct {
	ServerMode        string
	JWTKey            string
	Expiration        time.Duration
	EnableMemoryStore bool
	TLSOptions        *genericOptions.TLSOptions
	HTTPOptions       *genericOptions.HTTPOptions
	GRPCOptions       *genericOptions.GRPCOptions
	MySQLOptions      *genericOptions.MySQLOptions
}

// UnionServer type UnionServer struct {
type UnionServer struct {
	srv server.Server
}
