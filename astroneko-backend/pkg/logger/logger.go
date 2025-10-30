package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

type gcpZapLogger struct {
	zapLogger *zap.Logger
}

func NewDualLogger(z *zap.Logger) Logger {
	return &gcpZapLogger{
		zapLogger: z,
	}
}

func (l *gcpZapLogger) Info(msg string, fields ...Field) {
	zapFields := convertZapFields(fields)
	l.zapLogger.Info(msg, zapFields...)
}

func (l *gcpZapLogger) Warn(msg string, fields ...Field) {
	zapFields := convertZapFields(fields)
	l.zapLogger.Warn(msg, zapFields...)
}

func (l *gcpZapLogger) Error(msg string, fields ...Field) {
	zapFields := convertZapFields(fields)
	l.zapLogger.Error(msg, zapFields...)
}

func convertZapFields(fields []Field) []zap.Field {
	zfs := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zfs = append(zfs, zap.Any(f.Key, f.Value))
	}
	return zfs
}
