package logger

import (
	"context"
	"os"

	"github.com/largezhou/gin_starter/app/appconst"
	"github.com/largezhou/gin_starter/app/config"
	"github.com/largezhou/gin_starter/app/shutdown"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var cfg = config.Config.Log
var intLevelMap = map[string]zapcore.Level{
	"debug":  zap.DebugLevel,
	"info":   zap.InfoLevel,
	"warn":   zap.WarnLevel,
	"error":  zap.ErrorLevel,
	"dPanic": zap.DPanicLevel,
	"panic":  zap.PanicLevel,
	"fatal":  zap.FatalLevel,
}

// Logger 包级别默认日志
var Logger = WithChannel(appconst.LogDefault)

// callerSkip 需要跳过的 堆栈 数，由于 logger 方法可能会被封装，所以需要跳过封装的层数
const callerSkip = 2

var zapLogger = new(zap.Logger)
var Loggers = make(map[string]*AppLogger)

func init() {
	level, ok := intLevelMap[cfg.Level]
	if !ok {
		level = zap.InfoLevel
	}

	cores := []zapcore.Core{
		zapcore.NewCore(zapcore.NewJSONEncoder(createEncodeConfig()), createFileWriter(), level),
	}
	if cfg.Stdout {
		cores = append(
			cores,
			zapcore.NewCore(zapcore.NewConsoleEncoder(createEncodeConfig()), createConsoleWriter(), level),
		)
	}
	*zapLogger = *zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(callerSkip),
	)

	shutdown.RegisterShutdownFunc(func(ctx context.Context) {
		_ = zapLogger.Sync()
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

func WithChannel(channel string) *AppLogger {
	if l := Loggers[channel]; l == nil {
		l = &AppLogger{
			logger:  zapLogger,
			channel: channel,
		}
		Loggers[channel] = l
	}

	return Loggers[channel]
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Debug(ctx, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Info(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Warn(ctx, msg, fields...)
}
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Error(ctx, msg, fields...)
}
func DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.DPanic(ctx, msg, fields...)
}
func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Panic(ctx, msg, fields...)
}
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	Logger.Fatal(ctx, msg, fields...)
}
