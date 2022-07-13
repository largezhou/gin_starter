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
	"github.com/largezhou/gin_starter/app/trans"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// apiRecover 中间件的错误处理函数
func apiRecover(ctx *gin.Context, err any) {
	if s, ok := err.(string); ok {
		err = errors.New(s)
	}

	fail(ctx, err)
}

func fail(ctx *gin.Context, err any) {
	switch err.(type) {
	case error:
		failError(ctx, err.(error))
	case string:
		failWith(ctx, apperror.OperateFail, err.(string))
	case int:
		code := err.(int)
		if msg, ok := apperror.ErrorCodeMap[code]; ok {
			failWith(ctx, code, msg)
		} else {
			handleDefaultError(ctx, err)
		}
	default:
		handleDefaultError(ctx, err)
	}
}

// failError 单独处理各种特定的错误
func failError(ctx *gin.Context, err error) {
	switch {
	case errors.As(err, &validator.ValidationErrors{}):
		ve := err.(validator.ValidationErrors)
		msg := "参数校验失败"
		if len(ve) >= 0 {
			msg = ve[0].Translate(trans.Translator)
		}
		failWith(ctx, apperror.InvalidParameter, msg)
	case errors.As(err, &apperror.Error{}):
		e := err.(apperror.Error)
		failWith(ctx, e.Code, e.Msg)
	case errors.Is(err, gorm.ErrRecordNotFound):
		fail(ctx, apperror.ResourceNotFound)
	default:
		handleDefaultError(ctx, err)
	}
}

func failWith(ctx *gin.Context, code int, msg string) {
	response(ctx, code, msg, nil, nil)
}

// handleDefaultError 处理非预期的错误，会记录日志和堆栈
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
