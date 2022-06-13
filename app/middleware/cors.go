package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 简单跨域中间件
func Cors(r gin.IRoutes) gin.HandlerFunc {
	corsFunc := func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
	r.Use(corsFunc).OPTIONS("/*ignore", func(ctx *gin.Context) {
		return
	})

	return corsFunc
}
