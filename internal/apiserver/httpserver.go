package apiserver

import (
	"github.com/pachirode/go_blog_study/internal/pkg/server"
)

// ginServer 定义一个使用 Gin 框架开发的 HTTP 服务
type ginServer struct {
	srv server.Server
}
