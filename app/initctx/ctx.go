package initctx

import (
	"context"

	"github.com/google/uuid"
	"github.com/largezhou/gin_starter/app/appconst"
)

// Ctx 应用启动时的 ctx
// 不依赖其他包，尽量避免循环依赖，所以不使用 helper.NewTraceIdContext
var Ctx = context.WithValue(context.Background(), appconst.TraceIdKey, uuid.NewString())
