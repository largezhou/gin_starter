package apperror

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	Code  int
	Msg   string
	Cause error
}

func (e Error) Error() string {
	msg := ""

	if e.Msg != "" {
		msg = e.Msg
	} else {
		msg = errorCodeMap[e.Code]
	}

	if e.Cause != nil {
		msg += ": " + e.Cause.Error()
	}

	return msg
}

func (e Error) SetCode(code int) Error {
	e.Code = code
	return e
}

func (e Error) SetMsg(msg string) Error {
	e.Msg = msg
	return e
}

func (e Error) SetCause(err error) Error {
	e.Cause = err
	return e
}

func New(msg string) Error {
	return Error{
		Code:  OperateFail,
		Msg:   msg,
		Cause: nil,
	}
}

// ValidationErrors 包装一下 validator.ValidationErrors
// 暂时没有处理 i18n 的情况
type ValidationErrors struct {
	E validator.ValidationErrors
}

func (ve ValidationErrors) Error() string {
	if len(ve.E) == 0 {
		return ""
	}
	fe := ve.E[0]
	param := fe.Param()
	if param != "" {
		param = ":" + param
	}
	return fmt.Sprintf("字段 [ %s ] 验证规则 [ %s ] 失败", fe.Field(), fe.Tag() + param)
}
