package logger

import (
	"log"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(format string, level zapcore.Level) *Logger {
	var loggerConfig zap.Config

	if format == "json" {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}
	loggerConfig.DisableStacktrace = true //stack traces are shown by github.com/pkg/errors package
	loggerConfig.Level = zap.NewAtomicLevelAt(level)

	zaplog, err := loggerConfig.Build()
	if err != nil {
		log.Panicf("Could not initialize zap logger: %v", err)
	}
	logger := zaplog.Sugar()

	return &Logger{SugaredLogger: logger}
}

// Infoc provides conditional logging based on provided loglevel
func (l *Logger) Infoc(msg string, keysAndValues ...interface{}) {
	desugared := l.Desugar().WithOptions(zap.AddCallerSkip(1))
	if desugared.Core().Enabled(zapcore.DebugLevel) {
		desugared.Sugar().Debugw(msg, keysAndValues...)
	} else {
		desugared.Sugar().Infow(msg)
	}
}

func (l *Logger) With(fields ...interface{}) *Logger {
	return &Logger{l.SugaredLogger.With(fields...)}
}
