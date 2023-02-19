package main

import (
	"git.jetbrains.space/philldev/tcc/internal/log"
	"git.jetbrains.space/philldev/tcc/internal/wireup"
	"git.jetbrains.space/philldev/tcc/pkg/ping"
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
