package api

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/largezhou/gin_starter/app/apperror"
	"github.com/largezhou/gin_starter/app/config"
	"github.com/largezhou/gin_starter/app/middleware"
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
		response(ctx, apperror.InvalidParameter, err.Error(), nil, nil)
	case errors.As(err, &apperror.Error{}):
		e := err.(apperror.Error)
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

	response(ctx, apperror.UnknownErr, msg, nil, fields)
}

func ok(ctx *gin.Context, data any, msg string) {
	response(ctx, apperror.StatusOk, msg, data, nil)
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
