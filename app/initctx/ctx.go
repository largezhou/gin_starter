package initctx

import (
	"github.com/largezhou/gin_starter/app/helper"
)

// Ctx 应用启动时的 ctx
var Ctx = helper.NewTraceIdContext()
