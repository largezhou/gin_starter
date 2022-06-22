package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 简单跨域中间件
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
