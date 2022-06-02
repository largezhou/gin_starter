package logger

import (
	"github.com/largezhou/gin_starter/lib/config"
	"github.com/largezhou/gin_starter/lib/helper"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var c = config.Config.Log
var Logger *zap.Logger

var intLevelMap = map[string]zapcore.Level{
	"debug":  zap.DebugLevel,
	"info":   zap.InfoLevel,
	"warn":   zap.WarnLevel,
	"error":  zap.ErrorLevel,
	"dPanic": zap.DPanicLevel,
	"panic":  zap.PanicLevel,
	"fatal":  zap.FatalLevel,
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
	)

	helper.RegisterShutdownFunc(func() {
		_ = Logger.Sync()
	})
}

func createConsoleWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func createFileWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename: "./logs/log.log",
		MaxSize:  1,
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

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}
func DPanic(msg string, fields ...zap.Field) {
	Logger.DPanic(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
