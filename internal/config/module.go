package config

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewConfig,
)

func NewConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("default")
	v.AddConfigPath("configs")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return v
}
