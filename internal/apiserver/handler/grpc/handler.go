package grpc

import (
	"github.com/pachirode/go_blog_study/internal/apiserver/biz"
	apiV1 "github.com/pachirode/go_blog_study/pkg/api/apiserver/v1/proto"
)

// Handler 负责处理博客模块的请求.
type Handler struct {
	apiV1.UnimplementedMiniBlogServer

	biz biz.IBiz
}

// NewHandler 创建一个新的 Handler 实例.
func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{
		biz: biz,
	}
}
