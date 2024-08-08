package service

import (
	"context"
	"fmt"
	"rabbit-mq-fork/config"
	"rabbit-mq-fork/fork/worker"

	rmqconf "github.com/Badgain/rabbit/config"
	"github.com/Badgain/rabbit/listener"
	"github.com/Badgain/rabbit/producer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ForkService interface {
}

type service struct {
	lg      *zap.Logger
	workers map[string]worker.Worker
}

func NewForkService(cfg config.ForkConfig, lifecycle fx.Lifecycle) ForkService {
	lg, _ := zap.NewProduction()
	s := &service{
		lg:      lg.With(zap.String("layer", "service")),
		workers: make(map[string]worker.Worker),
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			for _, mapping := range cfg.GetMapping() {
				ln, err := listener.NewListener(rmqconf.ListenerConfig{
					Name:         fmt.Sprintf("pipe-%s-to-%s", mapping.Queue.Queue, mapping.Exchange.Name),
					QueueConfig:  mapping.Queue,
					ServerConfig: cfg.GetServerInfo(),
				}, lg)
				if err != nil {
					lg.Sugar().Error("Failed to create listener", err)
					return err
				}

				p, err := producer.NewProducer(rmqconf.PublisherConfig{
					Exchange:     mapping.Exchange,
					MessageType:  mapping.MessageType,
					ServerConfig: cfg.GetServerInfo(),
				}, lg)
				if err != nil {
					lg.Sugar().Error("Failed to create publisher", err)
					return err
				}

				w := worker.NewWorker(ctx, mapping.Queue.Queue, ln, p, lg)
				if err = w.Start(); err != nil {
					lg.Sugar().Error("Failed to start worker", err)
					return err
				}
				s.workers[mapping.Queue.Queue] = w
			}

			s.lg.Info("Starting service")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.lg.Info("Stopping service")
			for k, v := range s.workers {
				s.lg.Info("Stopping worker", zap.String("queue", k))
				if err := v.Stop(); err != nil {
					s.lg.Sugar().Error("Failed to stop worker", err)
					return err
				}
				delete(s.workers, k)
			}
			return nil
		},
	})

	return s
}
