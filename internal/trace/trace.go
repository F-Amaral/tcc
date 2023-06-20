package trace

import (
	"github.com/helios/go-sdk/sdk"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewTracer,
)

func NewTracer() *trace.TracerProvider {
	tracer, err := sdk.Initialize("tcc", "4eb2ae6ae3cb2bacbb8e", sdk.WithEnvironment("local"))
	if err != nil {
		panic(err)
	}
	return tracer
}
