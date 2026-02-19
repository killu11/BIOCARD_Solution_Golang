package broker

import (
	"context"
	"directory-viewing-service/internal/infrastructure/dto"
	"directory-viewing-service/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrSendTimeout = errors.New("send timeout")

type Sender struct {
	channel *amqp.Channel
}

func NewSender(ch *amqp.Channel) *Sender {
	return &Sender{channel: ch}
}

func (s *Sender) Send(fm *dto.FileDataMessage) error {
	body, err := json.Marshal(fm)
	if err != nil {
		return pkg.PackageError(packageName, "marshal message", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = s.channel.PublishWithContext(
		ctx,
		"",
		"files_messages",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: %s", ErrSendTimeout, fm.Filename)
		}
		return pkg.PackageError(packageName, "publish message", err)
	}
	return nil
}
