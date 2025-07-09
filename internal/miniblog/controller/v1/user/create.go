package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/Miniblog/internal/pkg/core"
	"github.com/marmotedu/Miniblog/internal/pkg/errno"
	"github.com/marmotedu/Miniblog/internal/pkg/log"
	v1 "github.com/marmotedu/Miniblog/pkg/api/miniblog/v1"
)

const defaultMethod = "(GET)|(POST)|(PUT)|(DELETE)"

func (ctrl *UserController) Create(ctx *gin.Context) {
	log.C(ctx).Infow("Create user function called")

	var r v1.CreateUserRequest

	if err := ctx.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(ctx, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(ctx, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	if err := ctrl.b.Users().Create(ctx, &r); err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	if _, err := ctrl.a.AddNamedPolicy("p", r.Username, "/v1/users/"+r.Username, defaultMethod); err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	core.WriteResponse(ctx, nil, nil)
}
