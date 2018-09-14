package logger

import (
	"log"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

// zapLogger is the log instance to be used in application code
var zapLogger logr.Logger

func init() {
	loggerConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stdout"},
		InitialFields:    map[string]interface{}{"from": "VarnishService-operator"},
	}

	z, err := loggerConfig.Build(zap.AddCallerSkip(2))
	if err != nil {
		log.Panicf("could not initialize zap logger: %v", err)
	}
	zapLogger = zapr.NewLogger(z)
}

// Info is exactly the same as zapLogger.Info
func Info(msg string, keysAndValues ...interface{}) {
	zapLogger.Info(msg, keysAndValues...)
}

// Error is exactly the same as zapLogger.Error
func Error(err error, msg string, keysAndValues ...interface{}) {
	zapLogger.Error(err, msg, keysAndValues...)
}

// RError is the same as Error, except it also returns the error value
func RError(err error, msg string, keysAndValues ...interface{}) error {
	zapLogger.Error(err, msg, keysAndValues...)
	return err
}
