package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/largezhou/gin_starter/app/app_error"
	"github.com/largezhou/gin_starter/app/config"
	"github.com/largezhou/gin_starter/app/middleware"
	"net/http"
	"runtime/debug"
)

func InitRouter(r *gin.Engine) {
	{
		g := r.Group("/api").Use(
			middleware.Recovery(failAny),
		)

		g.GET("/hello", func(ctx *gin.Context) {
			ok(ctx, "world", "")
		})
		g.GET("/error", func(ctx *gin.Context) {
			fail(ctx, errors.New("error"))
		})
	}
}

func failAny(ctx *gin.Context, err any) {
	realErr, ok := err.(error)
	if ok {
		fail(ctx, realErr)
	} else {
		handleDefaultError(ctx, err)
	}
}

func fail(ctx *gin.Context, err error) {
	switch {
	case errors.As(err, &validator.ValidationErrors{}):
		response(ctx, app_error.InvalidParameter, err.Error(), nil, nil)
	case errors.As(err, &app_error.Error{}):
		e := err.(app_error.Error)
		response(ctx, e.Code, e.Msg, nil, nil)
	default:
		handleDefaultError(ctx, err)
	}
}

func failWith(ctx *gin.Context, code int, msg string) {
	response(ctx, code, msg, nil, nil)
}

func handleDefaultError(ctx *gin.Context, err any) {
	var msg string
	fields := gin.H{}

	if config.Config.App.Debug {
		msg = fmt.Sprintf("%v", err)
		fields["trace"] = string(debug.Stack())
	} else {
		msg = http.StatusText(http.StatusInternalServerError)
	}

	response(ctx, app_error.UnknownErr, msg, nil, fields)
}

func ok(ctx *gin.Context, data any, msg string) {
	response(ctx, app_error.StatusOk, msg, data, nil)
}

func response(ctx *gin.Context, code int, msg string, data any, fields gin.H) {
	resp := gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	}

	for k, v := range fields {
		resp[k] = v
	}

	ctx.JSON(http.StatusOK, resp)
}
