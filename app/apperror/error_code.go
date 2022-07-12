package apperror

const (
	StatusOk         = 0
	InternalError    = 50000
	OperateFail      = 40000
	AuthFail         = 40001
	ResourceNotFound = 40004
	InvalidParameter = 40022
)

var ErrorCodeMap = map[int]string{
	StatusOk:         "成功",
	InternalError:    "服务器异常",
	OperateFail:      "操作失败",
	AuthFail:         "认证失败",
	ResourceNotFound: "资源不存在",
	InvalidParameter: "参数错误",
}

func GetMsg(code int) string {
	if msg, ok := ErrorCodeMap[code]; ok {
		return msg
	}

	return ErrorCodeMap[InternalError]
}
