package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Log *zap.Logger

func InitLogger() {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	// Choose JSON encoder for structured logging
	encoder := zapcore.NewJSONEncoder(config)

	// Output to stdout
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)

	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

func Info(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(msg, fields...)
	}
}
