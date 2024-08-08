package fork

import (
	"rabbit-mq-fork/config"
	"rabbit-mq-fork/fork/service"

	"go.uber.org/fx"
)

var ForkModule = fx.Module(
	"fork-service-module",
	fx.Provide(
		service.NewForkService,
		config.NewConfig,
	),
	fx.Invoke(config.NewConfig, service.NewForkService),
)
