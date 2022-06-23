package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/largezhou/gin_starter/app/appconst"
)

func SetTraceId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader(appconst.TraceIdHeaderKey)
		if traceId == "" {
			traceId = uuid.NewString()
		}
		ctx.Set(appconst.TraceIdKey, traceId)
		ctx.Header(appconst.TraceIdHeaderKey, traceId)

		ctx.Next()
	}
}
