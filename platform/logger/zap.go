package logger

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() *ZapLogger {
	// For production, use zap.NewProduction()
	// For development, use zap.NewDevelopment()
	l, _ := zap.NewProduction()
	return &ZapLogger{logger: l}
}

func (z *ZapLogger) Info(msg string, fields ...interface{}) {
	z.logger.Sugar().Infow(msg, fields...)
}

func (z *ZapLogger) Error(msg string, fields ...interface{}) {
	z.logger.Sugar().Errorw(msg, fields...)
}

func (z *ZapLogger) Debug(msg string, fields ...interface{}) {
	z.logger.Sugar().Debugw(msg, fields...)
}

func (z *ZapLogger) Warn(msg string, fields ...interface{}) {
	z.logger.Sugar().Warnw(msg, fields...)
}

func (z *ZapLogger) Fatal(msg string, fields ...interface{}) {
	z.logger.Sugar().Fatalw(msg, fields...)
}
