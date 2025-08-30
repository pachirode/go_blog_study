package apiserver

import (
	"time"

	"github.com/pachirode/go_blog_study/internal/apiserver/biz"
	"github.com/pachirode/go_blog_study/internal/apiserver/pkg/validation"
	mw "github.com/pachirode/go_blog_study/internal/pkg/middleware/gin"
	"github.com/pachirode/go_blog_study/internal/pkg/server"
	"github.com/pachirode/pkg/authz"
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

// UnionServer 联合服务器，根据 ServerMode 决定要启动的服务器类型
// 可选服务器类型
//   - GRPC
//   - HTTP 反向代理
//   - grpc-gateway 反向代理，根据 TLS 决定启动 HTTP 还是 HTTPS
//   - Gin
//   - 根据是否启动 TLS 来判断启动 HTTP 还是 HTTPS
type UnionServer struct {
	srv server.Server
}

// ServerConfig 服务器依赖和配置项目
type ServerConfig struct {
	cfg       *Config
	biz       biz.IBiz
	val       *validation.Validator
	retriever mw.UserRetriever
	authz     *authz.Authz
}
