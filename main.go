package main

import (
	"github.com/F-Amaral/tcc/internal/config"
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/F-Amaral/tcc/internal/telemetry"
	"github.com/F-Amaral/tcc/internal/wireup"
	"github.com/F-Amaral/tcc/pkg/ping"
	"github.com/F-Amaral/tcc/pkg/tree"
	"github.com/F-Amaral/tcc/pkg/tree/http/nested"
	"github.com/F-Amaral/tcc/pkg/tree/http/ppt"
	"github.com/F-Amaral/tcc/pkg/tree/repository"
	"github.com/F-Amaral/tcc/pkg/tree/service"
	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		config.Module,
		telemetry.Module,
		log.Module,
		wireup.Module,
		repository.Module,
		service.Module,
		ppt.Module,
		nested.Module,
		tree.Module,
		ping.Module,
	)

	app.Run()
}
