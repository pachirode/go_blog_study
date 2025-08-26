package gin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NoCache 是一个 Gin 中间件，用于禁止客户端缓存 HTTP 请求的返回结果.
func NoCache(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	ctx.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	ctx.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	ctx.Next()
}

// Cors 是一个 Gin 中间件，用于处理 CORS 请求.
func Cors(ctx *gin.Context) {
	// 处理预检请求
	if ctx.Request.Method == http.MethodOptions {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		ctx.Header("Allow", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Content-Type", "application/json")
		ctx.AbortWithStatus(http.StatusOK)
		return
	}
	ctx.Next() // 继续处理请求
}

// Secure 是一个 Gin 中间件，用于添加安全相关的 HTTP 头.
func Secure(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("X-Frame-Options", "DENY")
	ctx.Header("X-Content-Type-Options", "nosniff")
	ctx.Header("X-XSS-Protection", "1; mode=block")
	if ctx.Request.TLS != nil {
		ctx.Header("Strict-Transport-Security", "max-age=31536000")
	}
	ctx.Next()
}
