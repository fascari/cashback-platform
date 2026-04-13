package logger

import "go.uber.org/zap"

var instance *zap.Logger

func Init() {
	instance = zap.Must(zap.NewDevelopment())
}

func Zap() *zap.Logger {
	return instance
}

func Info(msg string, args ...any) {
	if instance != nil {
		instance.Sugar().Infow(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if instance != nil {
		instance.Sugar().Errorw(msg, args...)
	}
}

func Debug(msg string, args ...any) {
	if instance != nil {
		instance.Sugar().Debugw(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if instance != nil {
		instance.Sugar().Warnw(msg, args...)
	}
}
