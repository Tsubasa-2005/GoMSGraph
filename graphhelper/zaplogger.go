package graphhelper

import (
	"go.uber.org/zap"
)

type Logger struct {
	sugar *zap.SugaredLogger
}

func NewDefaultLogger(logger *zap.Logger) *Logger {
	return &Logger{sugar: logger.Sugar()}
}

func (zl *Logger) Debugf(format string, v ...interface{}) {
	zl.sugar.Debugf(format, v...)
}

func (zl *Logger) Infof(format string, v ...interface{}) {
	zl.sugar.Infof(format, v...)
}

func (zl *Logger) Warnf(format string, v ...interface{}) {
	zl.sugar.Warnf(format, v...)
}

func (zl *Logger) Errorf(format string, v ...interface{}) {
	zl.sugar.Errorf(format, v...)
}
