package log

import (
	"context"
	"github.com/F-Amaral/tcc/internal/telemetry"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Module = fx.Provide(
	NewLogger,
)

type Logger = *zap.SugaredLogger

const LogCtxKey string = "logCtxKey"

func NewLogger(telemetry telemetry.Telemetry) Logger {
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zap.InfoLevel)
	backgroundCore, err := nrzap.WrapBackgroundCore(core, telemetry)
	if err != nil && err != nrzap.ErrNilApp {
		panic(err)
	}
	return zap.New(backgroundCore).Sugar()
}

func NewCore() zapcore.Core {
	return zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zap.InfoLevel)
}

func FromContext(ctx context.Context) Logger {
	l, _ := ctx.Value(LogCtxKey).(Logger)
	return l
}

func With(ctx context.Context, fields ...Field) context.Context {
	logger := getLogger(ctx).With(fields)
	return context.WithValue(ctx, LogCtxKey, logger)
}

func WrapTransaction(tx *newrelic.Transaction) Logger {
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zap.InfoLevel)
	backgroundCore, err := nrzap.WrapTransactionCore(core, tx)
	if err != nil && err != nrzap.ErrNilApp {
		return zap.New(core).Sugar()
	}
	return zap.New(backgroundCore).Sugar()
}
