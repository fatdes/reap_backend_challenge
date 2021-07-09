package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	*zap.SugaredLogger
}

func GetLogger(name string) *Log {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), os.Stdout, zap.InfoLevel)

	return &Log{zap.New(config).Sugar().With("logger", string(name))}
}

func (l *Log) Add(key string, value interface{}) *Log {
	l.SugaredLogger = l.SugaredLogger.With(key, value)
	return l
}
