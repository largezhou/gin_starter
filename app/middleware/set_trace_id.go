package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/largezhou/gin_starter/app/app_const"
)

func SetTraceId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader(app_const.TraceIdHeaderKey)
		if traceId == "" {
			traceId = uuid.NewString()
		}
		ctx.Set(app_const.TraceIdKey, traceId)
		ctx.Header(app_const.TraceIdHeaderKey, traceId)

		ctx.Next()
	}
}
