package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/pkg/core"

	"github.com/pachirode/go_blog_study/internal/pkg/log"
	apiV1 "github.com/pachirode/go_blog_study/pkg/api/apiserver/v1/proto"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(ctx *gin.Context) {
	log.C(ctx.Request.Context()).Infow("Healthz handler is called", "method", "Healthz", "status", "healthy")
	core.WriteResponse(ctx, apiV1.HealthzResponse{
		Status:    apiV1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil)
}
