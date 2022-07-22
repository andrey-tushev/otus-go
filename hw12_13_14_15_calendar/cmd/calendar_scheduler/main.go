package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/app"
	conf "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/storage/factory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/scheduler_config.toml", "Path to configuration file")
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

	_ = ctx

	// Выбираем и настраиваем хранилище
	storage, err := factory.GetStorage(config)
	if err != nil {
		logg.Error(err.Error())
		return 1
	}

	// Запускаем приложение
	calendar := app.New(logg, storage)

	_ = calendar

	logg.Info("started")

	return 0
}
