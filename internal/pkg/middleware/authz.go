package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pachirode/go_blog_study/internal/pkg/core"
	"github.com/pachirode/go_blog_study/internal/pkg/errno"
	"github.com/pachirode/go_blog_study/internal/pkg/known"
	"github.com/pachirode/go_blog_study/internal/pkg/log"
)

type Auther interface {
	Authorize(sub, obj, act string) (bool, error)
}

func Authz(a Auther) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sub := ctx.GetString(known.Usernamekey)
		obj := ctx.Request.URL.Path
		act := ctx.Request.Method

		log.Debugw("Build authorize context", "sub", sub, "obj", obj, "act", act)
		if allowed, _ := a.Authorize(sub, obj, act); !allowed {
			core.WriteResponse(ctx, errno.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
	}
}
