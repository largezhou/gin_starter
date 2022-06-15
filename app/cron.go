package app

import (
	"context"
	"github.com/largezhou/gin_starter/app/logger"
	"go.uber.org/zap"
)

type CronLogger struct {
}

func (c CronLogger) Info(msg string, keysAndValues ...interface{}) {
	logger.Info(context.Background(), msg, zap.Any("keysAndValues", keysAndValues))
}

func (c CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	logger.Error(context.Background(), msg, zap.Error(err), zap.Any("keysAndValues", keysAndValues))
}
