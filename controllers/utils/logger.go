package utils

import (
	"log/slog"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ComponentKey = "cmp"
	ResourceKey  = "res"
)

func Logger(name string) logr.Logger {
	cfg := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			NameKey:     "logger",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey: "time",
			EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(time.Format("2006-01-02T15:04:05.999"))
			},
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic("Can not initialize logger")
	}
	return zapr.NewLogger(logger).WithName(name)
}

func Warn(log logr.Logger, msg string, keysAndValues ...interface{}) {
	slog.New(logr.ToSlogHandler(log)).Warn(msg, keysAndValues...)
}
