package apiserver

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	genericValidation "github.com/pachirode/pkg/validation"
	"google.golang.org/grpc"

	handler "github.com/pachirode/go_blog_study/internal/apiserver/handler/grpc"
	mw "github.com/pachirode/go_blog_study/internal/pkg/middleware/grpc"
	"github.com/pachirode/go_blog_study/internal/pkg/server"
	apiV1 "github.com/pachirode/go_blog_study/pkg/api/apiserver/v1/proto"
)

// grpcServer 定义一个 gRPC 服务器
type grpcServer struct {
	srv  server.Server
	stop func(context.Context)
}

// 确保 *grpcServer 实现了 server.Server 接口.
var _ server.Server = (*grpcServer)(nil)

// NewGRPCServerOr 创建一个 gRPC 或者 gRPC + grpc-gateway 混合服务器; 函数命令带有 Or 表示该函数可能存在多个功能
func (c *ServerConfig) NewGRPCServerOr() (server.Server, error) {
	// 配置服务器选项，设置拦截器
	serverOption := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			mw.RequestIDInterceptor(),
			selector.UnaryServerInterceptor(mw.AuthnInterceptor(c.retriever), NewAuthnWhiteListMatcher()),
			selector.UnaryServerInterceptor(mw.AuthzInterceptor(c.authz), NewAuthzWhiteListMatcher()),
			mw.DefaulterInterceptor(),
			mw.ValidatorInterceptor(genericValidation.NewValidator(c.val)),
		),
	}

	grpcSrv, err := server.NewGRPCServer(c.cfg.GRPCOptions, c.cfg.TLSOptions, serverOption, func(s grpc.ServiceRegistrar) {
		apiV1.RegisterMiniBlogServer(s, handler.NewHandler(c.biz))
	})

	if err != nil {
		return nil, err
	}

	if c.cfg.ServerMode == GRPCServerMode {
		return &grpcServer{srv: grpcSrv, stop: func(ctx context.Context) { grpcSrv.GracefulStop(ctx) }}, nil
	}

	go grpcSrv.RunOrDie()

	httpSrv, err := server.NewGRPCGatewayServer(c.cfg.HTTPOptions, c.cfg.GRPCOptions, c.cfg.TLSOptions, func(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
		return apiV1.RegisterMiniBlogHandler(context.Background(), mux, conn)
	})

	if err != nil {
		return nil, err
	}
	return &grpcServer{srv: httpSrv, stop: func(ctx context.Context) {
		grpcSrv.GracefulStop(ctx)
		httpSrv.GracefulStop(ctx)
	}}, nil
}

// RunOrDie 启动 gRPC 服务器或 HTTP 反向代理服务器，异常时退出.
func (s *grpcServer) RunOrDie() {
	s.srv.RunOrDie()
}

// GracefulStop 优雅停止 HTTP 和 gRPC 服务器.
func (s *grpcServer) GracefulStop(ctx context.Context) {
	s.stop(ctx)
}

// NewAuthnWhiteListMatcher 创建认证白名单匹配器.
func NewAuthnWhiteListMatcher() selector.Matcher {
	whitelist := map[string]struct{}{
		apiV1.MiniBlog_Healthz_FullMethodName:    {},
		apiV1.MiniBlog_CreateUser_FullMethodName: {},
		apiV1.MiniBlog_Login_FullMethodName:      {},
	}
	return selector.MatchFunc(func(ctx context.Context, call interceptors.CallMeta) bool {
		_, ok := whitelist[call.FullMethod()]
		return !ok
	})
}

// NewAuthzWhiteListMatcher 创建授权白名单匹配器.
func NewAuthzWhiteListMatcher() selector.Matcher {
	whitelist := map[string]struct{}{
		apiV1.MiniBlog_Healthz_FullMethodName:    {},
		apiV1.MiniBlog_CreateUser_FullMethodName: {},
		apiV1.MiniBlog_Login_FullMethodName:      {},
	}
	return selector.MatchFunc(func(ctx context.Context, call interceptors.CallMeta) bool {
		_, ok := whitelist[call.FullMethod()]
		return !ok
	})
}
