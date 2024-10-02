package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Initialize(level string) (*zap.SugaredLogger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig.TimeKey = "time"

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return zl.Sugar(), nil
}
