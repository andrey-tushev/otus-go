package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/consumer"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender_config.toml", "Path to configuration file")
}

func main() {
	ret := retMain()
	os.Exit(ret)
}

func retMain() int {
	flag.Parse()

	config := NewConfig()
	if err := config.Parse(configFile); err != nil {
		fmt.Println(err)
		return 1
	}

	logg := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	messages, err := getMessages(ctx, config, logg)
	if err != nil {
		logg.Error("message chan error: " + err.Error())
		return 1
	}

	logg.Info("start consuming")
	for message := range messages {
		logg.Info("mock remind: " + string(message.Data))
	}
	logg.Info("finish consuming")

	return 0
}

// getMessages возвращает канал сообщениями
func getMessages(ctx context.Context, config *Config, logg *logger.Logger) (<-chan consumer.Message, error) {
	conn, err := amqp.Dial(config.RabbitMQ.URI)
	if err != nil {
		logg.Error("rabbitmq dial error: " + err.Error())
		return nil, err
	}

	consumer := consumer.New(config.RabbitMQ.Consumer, conn, logg)
	messages, err := consumer.Consume(ctx, config.RabbitMQ.Queue)
	if err != nil {
		logg.Error("rabbitmq dial error: " + err.Error())
		return nil, err
	}

	return messages, nil
}
