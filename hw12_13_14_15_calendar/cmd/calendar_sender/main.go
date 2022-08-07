package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	conf "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/queue/rabbitmq"
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

	config := conf.New()
	if err := config.Parse(configFile); err != nil {
		fmt.Println(err)
		return 1
	}

	logg := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	messages, err := rabbitmq.GetMessages(ctx, config.RabbitMQ, logg)
	if err != nil {
		logg.Error("message chan error: " + err.Error())
		return 1
	}

	logg.Info("start consuming")
	for message := range messages {
		logg.Info("mock remind: " + string(message.String()))
	}
	logg.Info("finish consuming")

	return 0
}
