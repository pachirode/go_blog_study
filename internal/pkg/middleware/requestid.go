package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/marmotedu/Miniblog/internal/pkg/known"
)

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.Request.Header.Get(known.RequestUUID)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Set(known.RequestUUID, requestID)
		ctx.Writer.Header().Set(known.RequestUUID, requestID)
		ctx.Next()
	}
}
