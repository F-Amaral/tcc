package main

import (
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/F-Amaral/tcc/internal/wireup"
	"github.com/F-Amaral/tcc/pkg/ping"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		log.Module,
		wireup.Module,
		ping.Module,
	)

	app.Run()
}
