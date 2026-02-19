package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type RabbitMQConfig struct {
	User     string `env:"RABBITMQ_USER" envDefault:"guest"`
	Password string `env:"RABBITMQ_PASSWORD" envDefault:"guest"`
	Host     string `env:"RABBITMQ_HOST" envDefault:"localhost"`
	Port     int    `env:"RABBITMQ_PORT" envDefault:"5672"`
}

func NewRabbitMQConfig() (*RabbitMQConfig, error) {
	var r RabbitMQConfig
	if err := env.Parse(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *RabbitMQConfig) DSN() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}
