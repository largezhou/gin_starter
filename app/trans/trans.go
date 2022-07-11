package trans

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/largezhou/gin_starter/app/config"
	"github.com/largezhou/gin_starter/app/initctx"
	"github.com/largezhou/gin_starter/app/logger"
	"go.uber.org/zap"
)

var Translator ut.Translator

var cfg = config.Config.App

var localeTranslatorMap = map[string]locales.Translator{
	"zh": zh.New(),
	"en": en.New(),
	// other
}
var localeTranslationFuncMap = map[string]func(*validator.Validate, ut.Translator) error{
	"zh": zhTrans.RegisterDefaultTranslations,
	"en": enTrans.RegisterDefaultTranslations,
	// other
}

func init() {
	if err := initTranslator(); err != nil {
		logger.Error(initctx.Ctx, "初始化验证器语言失败", zap.Error(err))
	}
}

func initTranslator() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("验证器类型错误")
	}

	lt, ok := localeTranslatorMap[cfg.Locale]
	if !ok {
		return fmt.Errorf("语言 [ %s ] 不存在", cfg.Locale)
	}
	fallbackLt, ok := localeTranslatorMap[cfg.FallbackLocale]
	if !ok {
		return fmt.Errorf("备用语言 [ %s ] 不存在", cfg.FallbackLocale)
	}

	uni := ut.New(fallbackLt, lt)
	Translator, _ = uni.GetTranslator(cfg.Locale)

	transFunc, ok := localeTranslationFuncMap[cfg.Locale]
	if !ok {
		return fmt.Errorf("语言 [ %s ] 不存在", cfg.Locale)
	}
	fallbackTransFunc, ok := localeTranslationFuncMap[cfg.FallbackLocale]
	if !ok {
		return fmt.Errorf("备用语言 [ %s ] 不存在", cfg.FallbackLocale)
	}

	if err := transFunc(v, Translator); err != nil {
		logger.Error(initctx.Ctx, "注册语言失败，使用备用语言", zap.Error(err))
		return fallbackTransFunc(v, Translator)
	}

	return nil
}
