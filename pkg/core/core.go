package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/Miniblog/pkg/errno"
)

type ErrResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteResponse(ctx *gin.Context, err error, data interface{}) {
	if err != nil {
		http, code, msg := errno.Derrcode(err)
		ctx.JSON(http, ErrResponse{
			Code:    code,
			Message: msg,
		})

		return
	}

	ctx.JSON(http.StatusOK, data)
}
