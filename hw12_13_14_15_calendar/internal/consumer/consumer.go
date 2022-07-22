package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type RMQConnection interface {
	Channel() (*amqp.Channel, error)
}

type Consumer struct {
	name string
	conn RMQConnection
	logg *logger.Logger
}

func New(name string, conn RMQConnection, logg *logger.Logger) *Consumer {
	return &Consumer{
		name: name,
		conn: conn,
		logg: logg,
	}
}

type Message struct {
	Data []byte
}

func (c *Consumer) Consume(ctx context.Context, queue string) (<-chan Message, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	go func() {
		<-ctx.Done()
		if err := ch.Close(); err != nil {
			log.Println(err)
		}
	}()

	deliveries, err := ch.Consume(queue, c.name, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("start consuming: %w", err)
	}

	messages := make(chan Message)

	go func() {
		defer func() {
			close(messages)
			c.logg.Info("close messages channel")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case delivery := <-deliveries:
				if err := delivery.Ack(false); err != nil {
					c.logg.Error("acknowledge error: " + err.Error())
				}

				message := Message{
					Data: delivery.Body,
				}

				select {
				case <-ctx.Done():
					return
				case messages <- message:
				}
			}
		}
	}()

	return messages, nil
}
