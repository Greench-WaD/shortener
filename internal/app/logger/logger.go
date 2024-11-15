package logger

import (
	"go.uber.org/zap"
)

func New(level string) (*zap.Logger, error) {
	log := zap.NewNop()
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	log = zl
	return log, nil
}
