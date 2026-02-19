package broker

import (
	"directory-viewing-service/internal/config"
	"directory-viewing-service/pkg"

	amqp "github.com/rabbitmq/amqp091-go"
)

const packageName = "infrastructure/broker"

type RabbitMQ struct {
	conn *amqp.Connection
	Sch  *amqp.Channel
	Rch  *amqp.Channel
}

func NewRabbitMQ(c config.DSNConfig) (*RabbitMQ, error) {
	var rabbit RabbitMQ
	conn, err := amqp.Dial(c.DSN())
	if err != nil {
		return nil, pkg.PackageError(packageName, "open connection", err)
	}
	rabbit.conn = conn
	sch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, pkg.PackageError(packageName, "open connection", err)

	}
	rch, err := conn.Channel()
	if err != nil {
		sch.Close()
		conn.Close()
		return nil, pkg.PackageError(packageName, "open connection", err)
	}
	rabbit.Sch, rabbit.Rch = sch, rch
	if err := rabbit.setupQueues(); err != nil {
		sch.Close()
		rch.Close()
		conn.Close()
	}
	return &rabbit, nil
}

func (r *RabbitMQ) setupQueues() error {
	_, err := r.Sch.QueueDeclare(
		"files_messages",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return pkg.PackageError(packageName, "declare queue", err)
	}
	return nil
}
