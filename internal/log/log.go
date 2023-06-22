package log

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Module = fx.Provide(
	NewLogger,
)

type Logger = *zap.SugaredLogger

const LogCtxKey string = "logCtxKey"

func NewLogger() Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	return logger.Sugar()
}

func FromContext(ctx context.Context) Logger {
	l, _ := ctx.Value(LogCtxKey).(Logger)
	return l
}

func With(ctx context.Context, fields ...Field) context.Context {
	logger := getLogger(ctx).With(fields)
	return context.WithValue(ctx, LogCtxKey, logger)
}
