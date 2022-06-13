package init_ctx

import (
	"context"
	"github.com/google/uuid"
	"github.com/largezhou/gin_starter/app/app_const"
)

// Ctx 应用启动时的 ctx
var Ctx = context.WithValue(context.Background(), app_const.RequestIdKey, uuid.NewString())
