package logger

import (
	"hilo-api/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap(cfg config.Core) (*zap.Logger, error) {
	var err error

	var logCfg zap.Config
	if cfg.IsReleaseMode {
		logCfg = zap.NewProductionConfig()
	} else {
		logCfg = zap.NewDevelopmentConfig()
	}

	logCfg.EncoderConfig.EncodeName = zapcore.FullNameEncoder
	if len(cfg.LogLevel) > 0 {
		var lv zapcore.Level
		err := lv.Set(cfg.LogLevel)
		if err == nil {
			logCfg.Level.SetLevel(lv)
		}
	}
	logger, err := logCfg.Build(zap.Fields(zap.String("system", cfg.SystemName)))
	if err != nil {
		return nil, err
	}
	logger = logger.WithOptions(zap.AddStacktrace(zap.DPanicLevel))

	zap.ReplaceGlobals(logger)
	return logger, nil
}
