package helper

var shutdownFuncList []func()

// RegisterShutdownFunc 注册一个服务关闭时的回调函数
func RegisterShutdownFunc(f func()) {
	shutdownFuncList = append(shutdownFuncList, f)
}

// CallShutdownFunc 服务关闭时，执行所有回调函数
func CallShutdownFunc() {
	for _, f := range shutdownFuncList {
		f()
	}
}
