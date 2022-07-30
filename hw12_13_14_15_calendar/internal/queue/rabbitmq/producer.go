package rabbitmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"

	conf "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/queue"
)

type Producer struct {
	exchange string
	channel  *amqp.Channel
}

func New(ctx context.Context, config conf.Config, logg *logger.Logger) (*Producer, error) {
	conn, err := amqp.Dial(config.RabbitMQ.URI)
	if err != nil {
		logg.Error("rabbitmq dial error: " + err.Error())
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	if err := channel.ExchangeDeclare(
		config.RabbitMQ.Exchange,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	return &Producer{
		channel:  channel,
		exchange: config.RabbitMQ.Exchange,
	}, nil
}

func (p *Producer) Publish(message queue.Message) error {
	err := p.channel.Publish(
		p.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{"a-header": "a-value"},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            message.Data,
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
