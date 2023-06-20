package trace

import (
	"github.com/F-Amaral/tcc/constants"
	"github.com/helios/go-sdk/sdk"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewTracer,
)

func NewTracer(viper *viper.Viper) *trace.TracerProvider {
	apiName := viper.GetString(constants.ApiNameKey)
	environment := viper.GetString(constants.EnvironmentKey)
	apiToken := viper.GetString(constants.HeliosApiKey)

	tracer, err := sdk.Initialize(apiName, apiToken, sdk.WithEnvironment(environment))
	if err != nil {
		panic(err)
	}
	return tracer
}
