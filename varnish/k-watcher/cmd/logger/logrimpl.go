package logger

import (
	"log"

	"github.com/juju/errors"
	"go.uber.org/zap"
)

var infoLogger *zap.SugaredLogger
var errorLogger *zap.SugaredLogger

func init() {
	infoConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stdout"},
		InitialFields:    map[string]interface{}{"from": "k-watcher"},
	}
	errorConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoding:          "json",
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields:     map[string]interface{}{"from": "k-watcher"},
	}

	iz, err := infoConfig.Build(zap.AddCallerSkip(2))
	if err != nil {
		log.Panicf("could not initialize zap logger: %v", err)
	}
	infoLogger = iz.Sugar()
	ez, err := errorConfig.Build()
	if err != nil {
		log.Panicf("could not initialize zap logger: %v", err)
	}
	errorLogger = ez.Sugar()
}

// Info is exactly the same as zapLogger.Infow
func Info(msg string, keysAndValues ...interface{}) {
	infoLogger.Infow(msg, keysAndValues...)
}

// references call site for logError/logAndPanic
func generateErrorStack(err error, msg string) string {
	wrapped := errors.NewErrWithCause(err, msg)
	wrapped.SetLocation(2)
	return errors.ErrorStack(&wrapped)
}

// Error logs the err and message
func Error(err error, msg string, keysAndValues ...interface{}) {
	errorLogger.Errorw(generateErrorStack(err, msg), keysAndValues...)
}

// Panic logs the err and message, then panics (exits the program)
func Panic(err error, msg string, keysAndValues ...interface{}) {
	errorLogger.Panicw(generateErrorStack(err, msg), keysAndValues...)
}

// Sync just calls sync on the zap loggers
func Sync() {
	infoLogger.Sync()
	errorLogger.Sync()
}
