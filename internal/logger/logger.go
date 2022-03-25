package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

func init() {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	accessEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "timestamp",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	})
	errorEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "timestamp",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	})

	errorFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		msg := fmt.Sprintf("error when createing error.log file: %v", err)
		panic(msg)
	}
	errorSync := zapcore.AddSync(errorFile)

	accessFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		msg := fmt.Sprintf("error when createing access.log file: %v", err)
		panic(msg)
	}
	accessSync := zapcore.AddSync(accessFile)

	core := zapcore.NewTee(
		zapcore.NewCore(accessEncoder, accessSync, lowPriority),
		zapcore.NewCore(errorEncoder, errorSync, highPriority),
	)

	ret := zap.New(core)

	logger = ret
}

func GetLogger() *zap.Logger {
	return logger
}
