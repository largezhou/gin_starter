package logger

import (
	"context"
	"strings"

	"github.com/largezhou/gin_starter/app/appconst"
	"go.uber.org/zap"
)

type AppLogger struct {
	logger  *zap.Logger
	channel string
}

func (l *AppLogger) getTraceId(ctx context.Context) string {
	if traceId, ok := ctx.Value(appconst.TraceIdKey).(string); ok {
		return traceId
	} else {
		return ""
	}
}

func (l *AppLogger) getFields(ctx context.Context, fields []zap.Field) []zap.Field {
	traceId := l.getTraceId(ctx)
	if traceId != "" {
		fields = append(fields, zap.String(appconst.TraceIdKey, traceId))
	}
	if channel := strings.TrimSpace(l.channel); channel != "" {
		fields = append(fields, zap.String("channel", channel))
	}
	return fields
}

func (l *AppLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Debug(msg, l.getFields(ctx, fields)...)
}

func (l *AppLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Info(msg, l.getFields(ctx, fields)...)
}

func (l *AppLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Warn(msg, l.getFields(ctx, fields)...)
}
func (l *AppLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Error(msg, l.getFields(ctx, fields)...)
}
func (l *AppLogger) DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, l.getFields(ctx, fields)...)
}
func (l *AppLogger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Panic(msg, l.getFields(ctx, fields)...)
}
func (l *AppLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, l.getFields(ctx, fields)...)
}
