package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/largezhou/gin_starter/app/logger"
	"go.uber.org/zap"
)

// Recovery panic 处理
func Recovery(formatter func(*gin.Context, any)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(ctx, "异常", zap.Any("app_error", err))
				formatter(ctx, err)
			}
		}()
		ctx.Next()
	}
}
