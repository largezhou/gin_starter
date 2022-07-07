package middleware

import (
	"github.com/gin-gonic/gin"
)

var setToContextFuncList []func(*gin.Context)

func RegisterSetToContextFunc(f func(ctx *gin.Context)) {
	setToContextFuncList = append(setToContextFuncList, f)
}

// SetToContext 把一些值设置到 gin 的 context 中
func SetToContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, f := range setToContextFuncList {
			f(ctx)
		}

		ctx.Next()
	}
}
