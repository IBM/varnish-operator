package logger

import (
	"context"
	"log"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

const (
	FieldComponent       = "component"
	FieldComponentName   = "component_name"
	FieldVarnishService  = "varnish_service"
	FieldOperatorVersion = "operator_version"
	FieldKwatcherVersion = "kwatcher_version"
	FieldFilePath        = "file_path"
	FieldPodName         = "pod_name"
	FieldNamespace       = "namespace"
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

type contextKey int

const loggerCtxKey contextKey = iota

func ToContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

func FromContext(ctx context.Context) *Logger {
	ctxValue := ctx.Value(loggerCtxKey)
	if ctxValue == nil {
		return &Logger{zap.NewNop().Sugar()}
	}

	logr, ok := ctxValue.(*Logger)
	if !ok {
		return &Logger{zap.NewNop().Sugar()}
	}

	return logr
}

func NewNopLogger() *Logger {
	return &Logger{SugaredLogger: zap.NewNop().Sugar()}
}
