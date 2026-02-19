package broker

import (
	"context"
	"directory-viewing-service/internal/domain/services"
	"directory-viewing-service/internal/infrastructure/dto"
	"directory-viewing-service/pkg"
	"encoding/json"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Receiver struct {
	channel *amqp.Channel
}

func NewReceiver(ch *amqp.Channel) *Receiver {
	return &Receiver{channel: ch}
}

func (r *Receiver) Receive(
	ctx context.Context,
	out chan<- *dto.FileDataMessage,
	s services.FileTaskService,
) error {
	msgs, err := r.channel.Consume(
		"files_messages",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return pkg.PackageError(packageName, "consume messages", err)
	}

	go func() {
		for m := range msgs {
			select {
			case <-ctx.Done():
				close(out)
				return
			default:
				r.processMessage(ctx, m, out, s)
			}
		}
	}()
	return nil
}

func (r *Receiver) processMessage(
	ctx context.Context,
	msg amqp.Delivery,
	out chan<- *dto.FileDataMessage,
	s services.FileTaskService,
) {
	var fm dto.FileDataMessage
	if err := json.Unmarshal(msg.Body, &fm); err != nil {
		msg.Nack(false, false)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := s.ChangeStatus(ctx, fm.ID, services.StatusProcessing); err != nil {
		slog.Warn("update file task status", "error", err)
		msg.Nack(false, true)
		return
	}

	select {
	case <-ctx.Done():
	case out <- &fm:
		msg.Ack(false)
	}
}
