package telemetry

import (
	"context"
	"github.com/F-Amaral/tcc/constants"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"os"
	"time"
)

var Module = fx.Provide(NewTracer)

type Telemetry = *newrelic.Application

const TelemetryCtxKey string = "telemetryCtxKey"

type wrapper struct {
	tracer *newrelic.Application
}

func With(ctx context.Context) wrapper {
	return wrapper{
		tracer: getTracer(ctx),
	}
}

func (w wrapper) StartTransaction(name string, options ...newrelic.TraceOption) *newrelic.Transaction {
	tx := w.tracer.StartTransaction(name, options...)

	return tx
}

func FromContext(ctx context.Context) Telemetry {
	l, _ := ctx.Value(TelemetryCtxKey).(Telemetry)
	return l
}

func NewTracer(viper *viper.Viper) Telemetry {
	nrAppName := viper.GetString(constants.NRAppNameKey)
	nrLicense := viper.GetString(constants.NRAppLicenseKey)

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(nrAppName),
		newrelic.ConfigLicense(nrLicense),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDebugLogger(os.Stdout),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		panic(err)
	}

	app.WaitForConnection(5 * time.Second)
	return app
}

func getTracer(ctx context.Context) Telemetry {
	logger := FromContext(ctx)
	if logger != nil {
		return logger
	}

	app, _ := newrelic.NewApplication()
	return app
}
