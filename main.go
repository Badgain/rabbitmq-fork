package main

import (
	"rabbit-mq-fork/fork"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fork.ForkModule,
	)
	app.Run()
}
