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
		Code:  CommonError,
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
	return fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag() + param)
}
