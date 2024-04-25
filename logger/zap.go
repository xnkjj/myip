package logger

import "go.uber.org/zap"

var l *zap.Logger

func init() {
	l, _ = zap.NewProduction()
}

func Info(msg string, fields ...zap.Field) {
	l.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	l.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	l.Panic(msg, fields...)
}
