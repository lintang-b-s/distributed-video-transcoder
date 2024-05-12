package rabbitmq

import (
	"context"
	"lintang/video-transcoder-api/config"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(cfg *config.Config) *RabbitMQ {
	conn, err := amqp.Dial(cfg.RabbitMQ.RabbitURL)
	if err != nil {
		zap.L().Fatal("amqp.Dial", zap.Error(err))
	}

	channel, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("conn.Channel()", zap.Error(err))
	}

	err = channel.ExchangeDeclare(
		"metadata",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal(" channel.ExchangeDeclare", zap.Error(err))
	}
	if err := channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		zap.L().Fatal("channel.Qos", zap.Error(err))

	}

	return &RabbitMQ{
		conn, channel,
	}
}

func (r *RabbitMQ) Close( context.Context)   {
	 r.Connection.Close()
	
}
