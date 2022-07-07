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
		ve := err.(validator.ValidationErrors)
		response(ctx, apperror.InvalidParameter, (apperror.ValidationErrors{E: ve}).Error(), nil, nil)
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
	msg := fmt.Sprintf("%v", err)
	trace := string(debug.Stack())

	logger.Error(ctx, msg, zap.String("trace", trace))

	fields := gin.H{}

	if config.Config.App.Debug {
		fields["trace"] = trace
	} else {
		msg = http.StatusText(http.StatusInternalServerError)
	}

	response(ctx, apperror.UnknownErr, msg, nil, fields)
}

func ok(ctx *gin.Context, data any) {
	response(ctx, apperror.StatusOk, "", data, nil)
}

func okMsg(ctx *gin.Context, data any, msg string) {
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