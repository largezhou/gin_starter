package app_error

const (
	StatusOk         = 0
	UnknownErr       = 50000
	CommonError      = 40000
	AuthFail         = 40001
	ResourceNotFound = 40004
	InvalidParameter = 40022
)

var errorCodeMap = map[int]string{
	StatusOk:         "成功",
	UnknownErr:       "未知错误",
	CommonError:      "出错了",
	AuthFail:         "认证失败",
	ResourceNotFound: "资源不存在",
	InvalidParameter: "参数错误",
}
