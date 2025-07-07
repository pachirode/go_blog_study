package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/Miniblog/internal/pkg/core"
	"github.com/marmotedu/Miniblog/internal/pkg/errno"
	"github.com/marmotedu/Miniblog/internal/pkg/log"
	v1 "github.com/marmotedu/Miniblog/pkg/api/miniblog/v1"
)

func (ctrl *UserController) ChangePassword(ctx *gin.Context) {
	log.C(ctx).Infow("Change user password called")

	var r v1.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(ctx, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(ctx, errno.ErrInvalidParameter, nil)

		return
	}

	if err := ctrl.b.Users().ChangePassword(ctx, ctx.Param("name"), &r); err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	core.WriteResponse(ctx, nil, nil)
}
