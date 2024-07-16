package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
   This file is based off the following:
   - https://betterstack.com/community/guides/logging/go/zap/
   - https://stackoverflow.com/a/70974263
   - logurus (uses a similar pattern)
*/

var logger *zap.Logger

// Logger is an abstraction so that we can swap out the logger implementation if needed
type Logger struct {
	logger *zap.Logger
}

func init() {
	NewLogger()
}

func NewLogger() *Logger {
	if logger == nil {
		var err error
		config := zap.NewProductionConfig()
		enccoderConfig := zap.NewProductionEncoderConfig()
		zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
		enccoderConfig.StacktraceKey = "" // to hide stacktrace info
		config.EncoderConfig = enccoderConfig

		logger, err = config.Build(
			zap.AddCallerSkip(1), // to skip the caller of the logger (functions in this file)
		)
		if err != nil {
			panic(err)
		}
	}

	return &Logger{
		logger: logger,
	}
}

// DEBUG (-1): for recording messages useful for debugging.
func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// INFO (0): for messages describing normal application operations.
func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// WARN (1): for recording messages indicating something unusual happened that may need attention before it escalates to a more severe issue.
func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(msg, fields...)
}

// ERROR (2): for recording unexpected error conditions in the program.
func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// DPANIC (3): for recording severe error conditions in development. It behaves like PANIC in development and ERROR in production.
func (l *Logger) DPanic(msg string, fields ...zapcore.Field) {
	l.logger.DPanic(msg, fields...)
}

// PANIC (4): calls panic() after logging an error condition.
func (l *Logger) Panic(msg string, fields ...zapcore.Field) {
	l.logger.Panic(msg, fields...)
}

// FATAL (5): calls os.Exit(1) after logging an error condition.
func (l *Logger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}
