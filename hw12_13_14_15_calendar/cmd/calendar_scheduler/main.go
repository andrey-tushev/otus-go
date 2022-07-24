package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/app"
	conf "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/config"
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
	defer storage.Close(context.Background())

	// Запускаем приложение
	calendar := app.New(logg, storage)

	list, _ := calendar.ListEvents(context.Background())
	fmt.Println(list)

	// Подготовим тикеры
	cleanInterval, _ := time.ParseDuration(config.Scheduler.CleanInterval)
	if cleanInterval <= 0 {
		logg.Error("bad clean interval")
		return 1
	}
	cleanTicker := time.NewTicker(cleanInterval)

	remindInterval, _ := time.ParseDuration(config.Scheduler.RemindInterval)
	if remindInterval <= 0 {
		logg.Error("bad remind interval")
		return 1
	}
	remindTicker := time.NewTicker(remindInterval)

	wg := sync.WaitGroup{}

	// Очистка старых событий
	wg.Add(1)
	go func() {
		calendar.Clean(context.Background())

		for {
			select {
			case <-cleanTicker.C:
				calendar.Clean(context.Background())

			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}()

	// Напоминалки
	wg.Add(1)
	go func() {
		calendar.Remind(context.Background())

		for {
			select {
			case <-remindTicker.C:
				calendar.Remind(context.Background())

			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
	logg.Info("finished")

	return 0
}
