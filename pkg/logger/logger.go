package logger

import (
	"icm-varnish-k8s-operator/pkg/config"
	"log"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

// InternalZapLogger wrapper for sugared logger
type InternalZapLogger struct {
	*zap.SugaredLogger
}

// RErrorw extends error logger
func (izl *InternalZapLogger) RErrorw(err error, msg string, keysAndValues ...interface{}) error {
	izl.Errorw(msg, append(keysAndValues, zap.Error(err))...)
	return err
}

// Infoc provides conditional logging based on provided loglevel
func (izl *InternalZapLogger) Infoc(msg string, keysAndValues ...interface{}) {
	if config.GlobalConf.LogLevel == zapcore.DebugLevel {
		izl.Debugw(msg, keysAndValues...)
	} else {
		izl.Infow(msg)
	}
}

// With extends logger.With
func (izl *InternalZapLogger) With(fields ...interface{}) *InternalZapLogger {
	return &InternalZapLogger{izl.SugaredLogger.With(fields...)}
}

var logger *InternalZapLogger

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

	logger = &InternalZapLogger{zaplog.Sugar()}
}

// Infow is exactly the same as logger.Infow
func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

// Debugw is exactly the same as logger.Debugw
func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

// Errorw is exactly the same as logger.Errorw
func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}

// RErrorw is logs an error and then returns it
func RErrorw(err error, msg string, keysAndValues ...interface{}) error {
	return logger.RErrorw(err, msg, keysAndValues...)
}

// With wraps logger.With
func With(fields ...interface{}) *InternalZapLogger {
	return logger.With(fields...)
}
