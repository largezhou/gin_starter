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
	"github.com/largezhou/gin_starter/app/logger"
	"go.uber.org/zap"
)

func fail(ctx *gin.Context, err any) {
	switch err.(type) {
	case error:
		failError(ctx, err.(error))
	case string:
		failWith(ctx, apperror.OperateFail, err.(string))
	default:
		handleDefaultError(ctx, err)
	}
}

func failError(ctx *gin.Context, err error) {
	switch {
	case errors.As(err, &validator.ValidationErrors{}):
		ve := err.(validator.ValidationErrors)
		failWith(ctx, apperror.InvalidParameter, (apperror.ValidationErrors{E: ve}).Error())
	case errors.As(err, &apperror.Error{}):
		e := err.(apperror.Error)
		failWith(ctx, e.Code, e.Msg)
	default:
		handleDefaultError(ctx, err)
	}
}

func failWith(ctx *gin.Context, code int, msg string) {
	response(ctx, code, msg, nil, nil)
}

func handleDefaultError(ctx *gin.Context, err any) {
	msg := fmt.Sprintf("%v", err)
	trace := string(debug.Stack())

	logger.Error(ctx, msg, zap.String("trace", trace))

	fields := gin.H{}

	if config.Config.App.Debug {
		fields["trace"] = trace
	} else {
		msg = apperror.GetMsg(apperror.InternalError)
	}

	response(ctx, apperror.InternalError, msg, nil, fields)
}

func ok(ctx *gin.Context, data any) {
	response(ctx, apperror.StatusOk, "", data, nil)
}

func okMsg(ctx *gin.Context, msg string, data any) {
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
