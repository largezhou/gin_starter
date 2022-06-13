package shutdown

import (
	"context"
)

type HandlerFunc func(ctx context.Context)

var shutdownFuncList []HandlerFunc

// RegisterShutdownFunc 注册一个服务关闭时的回调函数
func RegisterShutdownFunc(f HandlerFunc) {
	shutdownFuncList = append(shutdownFuncList, f)
}

// CallShutdownFunc 服务关闭时，执行所有回调函数
func CallShutdownFunc(ctx context.Context) {
	for _, f := range shutdownFuncList {
		f(ctx)
	}
}
