package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
}

func main() {
	// Так сделано чтобы можно было использовать код возврата через os.Exit
	// и при этом отрабатывали deffer-ы
	ret := retMain()
	os.Exit(ret)
}

func retMain() int {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return 1
	}

	// Загружаем конфигурацию
	config := NewConfig()
	if err := config.Parse(configFile); err != nil {
		fmt.Println(err)
		return 1
	}

	// Настраиваем логгер
	logg := logger.New(config.Logger.Level)

	// Выбираем и настраиваем хранилище
	var storage app.Storage
	switch config.Storage.Storage {
	case "memory":
		storage = memorystorage.New()
	case "sql":
		sqlStorage := sqlstorage.New(config.SQL.DSN)
		err := sqlStorage.Connect(context.Background())
		if err != nil {
			fmt.Println(err)
			return 1
		}
		defer sqlStorage.Close(context.Background())
		storage = sqlStorage
	default:
		fmt.Println("unknown storage")
		return 1
	}

	// Запускаем приложение
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	// Останавливалка сервера
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done() // получение сигнала

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		fmt.Println("terminating...")
		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	// Запускалка сервера
	logg.Info("calendar is running...")

	if err := server.Start(ctx, config.Web.Host, config.Web.Port); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		return 1
	}

	return 0
}
