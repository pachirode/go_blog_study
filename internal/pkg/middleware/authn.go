package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pachirode/go_blog_study/internal/pkg/core"
	"github.com/pachirode/go_blog_study/internal/pkg/errno"
	"github.com/pachirode/go_blog_study/internal/pkg/known"
	"github.com/pachirode/go_blog_study/pkg/token"
)

func Authn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, err := token.ParseRequest(ctx)
		if err != nil {
			core.WriteResponse(ctx, errno.ErrTokenInvalid, nil)
			ctx.Abort()

			return
		}

		ctx.Set(known.Usernamekey, username)
		ctx.Next()
	}
}
