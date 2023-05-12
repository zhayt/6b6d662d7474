package logger

import (
	"github.com/zhayt/kmf-tt/config"
	"go.uber.org/zap"
)

func Init(cfg *config.Config) (*zap.Logger, error) {
	switch cfg.App.AppMode {
	case "dev":
		return zap.NewDevelopment()
	default:
		return zap.NewProduction()
	}
}
