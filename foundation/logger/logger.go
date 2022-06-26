package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(service string) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.DisableStacktrace = true
	cfg.InitialFields = map[string]interface{}{
		"service": service,
	}

	sl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return sl.Sugar(), nil
}
