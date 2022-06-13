package helper

import (
	"errors"
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
