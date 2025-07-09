package miniblog

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/Miniblog/internal/miniblog/controller/v1/user"
	"github.com/marmotedu/Miniblog/internal/miniblog/store"
	"github.com/marmotedu/Miniblog/internal/pkg/core"
	"github.com/marmotedu/Miniblog/internal/pkg/errno"
	"github.com/marmotedu/Miniblog/internal/pkg/log"
	mw "github.com/marmotedu/Miniblog/internal/pkg/middleware"
	"github.com/marmotedu/Miniblog/pkg/auth"
)

func installRouters(g *gin.Engine) error {
	g.NoRoute(func(ctx *gin.Context) {
		core.WriteResponse(ctx, errno.ErrPageNotFound, nil)
	})

	g.GET("/test", func(ctx *gin.Context) {
		log.C(ctx).Infow("Test function called")

		core.WriteResponse(ctx, nil, map[string]string{"status": "ok"})
	})
	authz, err := auth.NewAuthz(store.S.DB())
	if err != nil {
		return nil
	}

	uc := user.New(store.S, authz)

	g.POST("/login", uc.Login)

	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
			userv1.PUT(":name/change-password", uc.ChangePassword)
			userv1.Use(mw.Authn(), mw.Authz(authz))
			userv1.GET(":name", uc.Get)
		}
	}

	return nil
}
