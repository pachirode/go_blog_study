package user

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/Miniblog/internal/pkg/core"
	"github.com/marmotedu/Miniblog/internal/pkg/log"
)

func (ctrl *UserController) Get(ctx *gin.Context) {
	log.C(ctx).Infow("Get user function called.")

	userInfo, err := ctrl.b.Users().Get(ctx, ctx.Param("name"))
	if err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	core.WriteResponse(ctx, nil, userInfo)
}
