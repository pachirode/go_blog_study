package miniblog

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/Miniblog/internal/miniblog/controller/v1/user"
	"github.com/marmotedu/Miniblog/internal/miniblog/store"
	"github.com/marmotedu/Miniblog/internal/pkg/core"
	"github.com/marmotedu/Miniblog/internal/pkg/errno"
	"github.com/marmotedu/Miniblog/internal/pkg/log"
)

func installRouters(g *gin.Engine) error {
	g.NoRoute(func(ctx *gin.Context) {
		core.WriteResponse(ctx, errno.ErrPageNotFound, nil)
	})

	g.GET("/test", func(ctx *gin.Context) {
		log.C(ctx).Infow("Test function called")

		core.WriteResponse(ctx, nil, map[string]string{"status": "ok"})
	})

	uc := user.New(store.S)

	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
		}
	}

	return nil
}
