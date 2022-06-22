package helper

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/largezhou/gin_starter/app/app_const"
	"github.com/largezhou/gin_starter/app/app_error"
	"github.com/largezhou/gin_starter/app/config"
	"gorm.io/gorm"
)

// CheckAppKey 检查 app key
func CheckAppKey() {
	if len(config.Config.App.Key) < 32 {
		panic("APP_KEY 长度至少为 32 位")
	}
}

// ModelNotFound 处理模型未找到
func ModelNotFound(err error, msg string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return app_error.New(msg).SetCode(app_error.ResourceNotFound)
	} else {
		return err
	}
}

// NewTraceIdContext 返回一个新的带链路追踪 ID 的 context
func NewTraceIdContext() context.Context {
	return context.WithValue(context.Background(), app_const.TraceIdKey, uuid.NewString())
}
