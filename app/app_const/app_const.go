package app_const

const (
	// TraceIdHeaderKey 请求头中链路追踪 ID 的键名
	TraceIdHeaderKey = "Trace-Id"
	// CliKey cli 命令行的入口
	CliKey = "cli"
	// TraceIdKey 链路追踪 ID 的键名
	TraceIdKey = "traceId"
)

// 日志渠道
const (
	LogDefault = "default"
	LogAccess  = "access"
	LogCron    = "cron"
	LogSql     = "sql"
)
