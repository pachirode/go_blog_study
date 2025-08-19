package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pachirode/go_blog_study/internal/pkg/core"
	"github.com/pachirode/go_blog_study/internal/pkg/errno"
	"github.com/pachirode/go_blog_study/internal/pkg/log"
	v1 "github.com/pachirode/go_blog_study/pkg/api/miniblog/v1"
)

func (ctrl *UserController) Login(ctx *gin.Context) {
	log.C(ctx).Infow("Login function called")

	var r v1.LoginRequest
	if err := ctx.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(ctx, errno.ErrBind, nil)

		return
	}

	resp, err := ctrl.b.Users().Login(ctx, &r)
	if err != nil {
		core.WriteResponse(ctx, err, nil)

		return
	}

	core.WriteResponse(ctx, nil, resp)
}
