package logger

import (
	"context"
	"github.com/largezhou/gin_starter/app/app_const"
	"github.com/largezhou/gin_starter/app/config"
	"github.com/largezhou/gin_starter/app/shutdown"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var c = config.Config.Log
var Logger *zap.Logger

// callerSkip 需要跳过的 堆栈 数，由于 logger 方法被包装了 1 层，所以需要跳过 1 层
const callerSkip = 1

var intLevelMap = map[string]zapcore.Level{
	"debug":     zap.DebugLevel,
	"info":      zap.InfoLevel,
	"warn":      zap.WarnLevel,
	"app_error": zap.ErrorLevel,
	"dPanic":    zap.DPanicLevel,
	"panic":     zap.PanicLevel,
	"fatal":     zap.FatalLevel,
}

func init() {
	level, ok := intLevelMap[c.Level]
	if !ok {
		level = zap.InfoLevel
	}

	cores := []zapcore.Core{
		zapcore.NewCore(zapcore.NewJSONEncoder(createEncodeConfig()), createFileWriter(), level),
	}
	if c.Stdout {
		cores = append(
			cores,
			zapcore.NewCore(zapcore.NewConsoleEncoder(createEncodeConfig()), createConsoleWriter(), level),
		)
	}
	Logger = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(callerSkip),
	)

	shutdown.RegisterShutdownFunc(func(ctx context.Context) {
		_ = Logger.Sync()
	})
}

func createConsoleWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func createFileWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename: "./logs/log.log",
		MaxSize:  100,
		MaxAge:   14,
		Compress: false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func createEncodeConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000Z0700"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getRequestId(ctx context.Context) string {
	if requestId, ok := ctx.Value(app_const.RequestIdKey).(string); ok {
		return requestId
	} else {
		return ""
	}
}

func getFields(ctx context.Context, fields []zap.Field) []zap.Field {
	requestId := getRequestId(ctx)
	if requestId != "" {
		return append(fields, zap.String(app_const.RequestIdKey, requestId))
	}
	return fields
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Debug(msg, getFields(ctx, fields)...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Info(msg, getFields(ctx, fields)...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Warn(msg, getFields(ctx, fields)...)
}
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Error(msg, getFields(ctx, fields)...)
}
func DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.DPanic(msg, getFields(ctx, fields)...)
}
func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Panic(msg, getFields(ctx, fields)...)
}
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Fatal(msg, getFields(ctx, fields)...)
}
