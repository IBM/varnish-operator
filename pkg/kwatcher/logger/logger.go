package logger

import (
	"icm-varnish-k8s-operator/pkg/kwatcher/config"
	"log"

	"github.com/juju/errors"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger
var errLogger *zap.SugaredLogger

func init() {
	var loggerConfig zap.Config

	if config.GlobalConf.LogFormat == "json" {
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.DisableCaller = true
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.Level = zap.NewAtomicLevelAt(config.GlobalConf.LogLevel)

	zaplog, err := loggerConfig.Build(zap.AddCallerSkip(2))
	if err != nil {
		log.Panicf("Could not initialize zap logger: %v", err)
	}
	logger = zaplog.Sugar()

	errZapLog, err := loggerConfig.Build(zap.AddCallerSkip(3))
	if err != nil {
		log.Panicf("Could not initialize error zap logger: %v", err)
	}
	errLogger = errZapLog.Sugar()
}

func WrappedError(err error) {
	var wrapped *errors.Err
	switch e := err.(type) {
	case *errors.Err:
		wrapped = e
	default:
		errWithCause := errors.NewErrWithCause(e, "")
		wrapped = &errWithCause
	}
	wrapped.SetLocation(2)
	errLogger.Error(errors.ErrorStack(wrapped))
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}
