package logger

import (
	"go.uber.org/zap"
)

// New TODO: убрать жесткую привязку к zap
func New(level string) (*zap.Logger, error) {
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

	log := zl
	return log, nil
}
