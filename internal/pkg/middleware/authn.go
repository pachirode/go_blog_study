package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/Miniblog/internal/pkg/core"
	"github.com/marmotedu/Miniblog/internal/pkg/errno"
	"github.com/marmotedu/Miniblog/internal/pkg/known"
	"github.com/marmotedu/Miniblog/pkg/token"
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
