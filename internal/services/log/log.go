package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger := zap.Must(zap.NewProduction())
	zap.ReplaceGlobals(logger)
}

func Debug(msg string, fields ...zapcore.Field) {
	zap.L().Debug(msg, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	zap.L().Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	zap.L().Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	zap.L().Error(msg, fields...)
}

func DPanic(msg string, fields ...zapcore.Field) {
	zap.L().DPanic(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	zap.L().Fatal(msg, fields...)
}
