package main

import (
	"github.com/F-Amaral/tcc/internal/config"
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/F-Amaral/tcc/internal/trace"
	"github.com/F-Amaral/tcc/internal/wireup"
	"github.com/F-Amaral/tcc/pkg/ping"
	"github.com/F-Amaral/tcc/pkg/tree"
	"github.com/F-Amaral/tcc/pkg/tree/http/ppt"
	"github.com/F-Amaral/tcc/pkg/tree/repository"
	"github.com/F-Amaral/tcc/pkg/tree/service"
	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		trace.Module,
		config.Module,
		log.Module,
		wireup.Module,
		repository.Module,
		service.Module,
		ppt.Module,
		tree.Module,
		ping.Module,
	)

	app.Run()
}
