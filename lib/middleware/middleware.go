package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/largezhou/gin_starter/lib/config"
	"github.com/largezhou/gin_starter/lib/logger"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
	"time"
)

// Logger 记录请求，添加 requestId
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.NewString()
		c.Set("requestId", requestId)
		logger.Logger = logger.Logger.With(zap.String("requestId", requestId))

		start := time.Now()

		logger.Info(
			"request",
			zap.String("clientIp", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
		)

		c.Next()

		logger.Info(
			"response",
			zap.Duration("cost", time.Now().Sub(start)),
			zap.Int("code", c.Writer.Status()),
		)
	}
}

// Recovery panic 处理
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("异常", zap.Any("error", err))

				resp := gin.H{}

				if config.Config.App.Debug {
					stack := debug.Stack()

					resp["message"] = fmt.Sprintf("%v", err)
					resp["trace"] = string(stack)
				} else {
					resp["message"] = http.StatusText(http.StatusInternalServerError)
				}

				c.JSON(http.StatusInternalServerError, resp)
			}
		}()
		c.Next()
	}
}
