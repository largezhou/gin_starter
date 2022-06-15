package logger

import (
	"context"
	"github.com/largezhou/gin_starter/app/app_const"
	"go.uber.org/zap"
	"strings"
)

type GLogger struct {
	logger  *zap.Logger
	channel string
}

func (l *GLogger) getTraceId(ctx context.Context) string {
	if traceId, ok := ctx.Value(app_const.TraceIdKey).(string); ok {
		return traceId
	} else {
		return ""
	}
}

func (l *GLogger) getFields(ctx context.Context, fields []zap.Field) []zap.Field {
	traceId := l.getTraceId(ctx)
	if traceId != "" {
		fields = append(fields, zap.String(app_const.TraceIdKey, traceId))
	}
	if channel := strings.TrimSpace(l.channel); channel != "" {
		fields = append(fields, zap.String("channel", channel))
	}
	return fields
}

func (l *GLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Debug(msg, l.getFields(ctx, fields)...)
}

func (l *GLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Info(msg, l.getFields(ctx, fields)...)
}

func (l *GLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Warn(msg, l.getFields(ctx, fields)...)
}
func (l *GLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Error(msg, l.getFields(ctx, fields)...)
}
func (l *GLogger) DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, l.getFields(ctx, fields)...)
}
func (l *GLogger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Panic(msg, l.getFields(ctx, fields)...)
}
func (l *GLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, l.getFields(ctx, fields)...)
}
