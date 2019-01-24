package logger

import (
	"log"

	"github.com/juju/errors"

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
		loggerConfig.DisableCaller = true
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.Level = zap.NewAtomicLevelAt(level)

	zaplog, err := loggerConfig.Build()
	if err != nil {
		log.Panicf("Could not initialize zap logger: %v", err)
	}
	logger := zaplog.Sugar()

	return &Logger{SugaredLogger: logger}
}

func (l *Logger) RErrorw(err error, msg string, keysAndValues ...interface{}) error {
	var wrapped *errors.Err
	switch e := err.(type) {
	case *errors.Err:
		wrapped = e
	default:
		errWithCause := errors.NewErrWithCause(e, msg)
		wrapped = &errWithCause
	}
	wrapped.SetLocation(2)
	l.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar().Errorw(msg, "stacktrace", errors.ErrorStack(wrapped))
	return err
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
