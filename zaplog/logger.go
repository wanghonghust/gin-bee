package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger = GetLogger()

func GetLogger() *zap.SugaredLogger {
	lg, _ := NewLoggerConfig().Build()
	return lg.Sugar()
}

func NewLoggerConfig() zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    NewLoggerEncoderConfig(),
		OutputPaths:      []string{"stdout", "./zaplog/gin.log"},
		ErrorOutputPaths: []string{"stderr", "./zaplog/zap-error.log"},
	}
}

func NewLoggerEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "Date",
		LevelKey:       "Level",
		NameKey:        "Name",
		CallerKey:      "Caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "Msg",
		StacktraceKey:  "Stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
