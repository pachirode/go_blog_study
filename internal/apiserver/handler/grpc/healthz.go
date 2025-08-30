package grpc

import (
	"context"
	"time"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/pachirode/go_blog_study/internal/pkg/log"
	apiV1 "github.com/pachirode/go_blog_study/pkg/api/apiserver/v1/proto"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(ctx context.Context, rq *emptypb.Empty) (*apiV1.HealthzResponse, error) {
	log.C(ctx).Infow("Healthz handler is called", "method", "Healthz", "status", "healthy")
	return &apiV1.HealthzResponse{
		Status:    apiV1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
