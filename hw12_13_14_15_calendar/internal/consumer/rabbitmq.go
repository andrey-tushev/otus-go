package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/logger"

	conf "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/config/sender"
)

func GetMessagesFromRMQ(ctx context.Context, config conf.RabbitMQConf, logg *logger.Logger) (<-chan Message, error) {
	conn, err := amqp.Dial(config.URI)
	if err != nil {
		logg.Error("rabbitmq dial error: " + err.Error())
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	go func() {
		<-ctx.Done()
		if err := channel.Close(); err != nil {
			log.Println(err)
		}
	}()

	deliveries, err := channel.Consume(config.Queue, config.Consumer,
		false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("start consuming: %w", err)
	}

	messages := make(chan Message)

	go func() {
		defer func() {
			close(messages)
			logg.Info("close messages channel")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case delivery := <-deliveries:
				if err := delivery.Ack(false); err != nil {
					logg.Error("acknowledge error: " + err.Error())
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
