package worker

import (
	"context"

	"github.com/Badgain/rabbit/listener"
	"github.com/Badgain/rabbit/producer"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Worker interface {
	Start() error
	Stop() error
}

type impl struct {
	queue string
	ctx   context.Context
	lg    *zap.Logger
	ln    listener.Listener
	p     producer.Producer
}

func NewWorker(ctx context.Context, queue string, ln listener.Listener, p producer.Producer, lg *zap.Logger) Worker {
	return &impl{
		queue: queue,
		lg:    lg,
		ctx:   ctx,
		ln:    ln,
		p:     p,
	}
}

func (w *impl) Start() error {
	return w.ln.Consume(w.ctx, w.pipe)
}

func (w *impl) Stop() error {
	w.ln.Stop()
	return w.p.Stop()
}

func (w *impl) pipe(msg amqp.Delivery) error {
	w.lg.Info("Processing message",
		zap.String("from", w.queue),
		zap.String("to", msg.Exchange),
		zap.String("key", msg.RoutingKey),
		zap.String("userId", msg.UserId),
		zap.Time("timestamp", msg.Timestamp),
	)
	return w.p.Produce(w.ctx, msg.Body)
}
