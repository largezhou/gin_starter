package app_error

type Error struct {
	Code  int
	Msg   string
	Cause error
}

func (e Error) Error() string {
	if e.Msg != "" {
		return e.Msg
	} else {
		return errorCodeMap[e.Code]
	}
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
