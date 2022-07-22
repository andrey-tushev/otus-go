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
	conf "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
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
	config := conf.New()
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

	webServer := internalhttp.NewServer(logg, calendar)
	grpcServer := internalgrpc.NewServer(logg, calendar)

	// Останавливалка серверов по сигналу
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done() // получение сигнала
		logg.Info("got terminating signal")

		// На остановку выделяем не более 3 секунд
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		logg.Info("terminating...")

		// Останавливаем web-сервер
		go func() {
			logg.Info("terminating web-server")
			if err := webServer.Stop(ctx); err != nil {
				logg.Error("failed to stop web-server: " + err.Error())
			}
		}()

		// Останавливаем grpc-сервер
		go func() {
			logg.Info("terminating grpc-server")
			if err := grpcServer.Stop(ctx); err != nil {
				logg.Error("failed to stop grpc-server: " + err.Error())
			}
		}()
	}()

	// Запускаем оба web и grpc сервер.
	// Программа завершиться когда завершатся оба сервера
	logg.Info("running web and grpc servers...")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := webServer.Start(ctx, config.Web.Host, config.Web.Port); err != nil {
			logg.Error("failed to start http-server: " + err.Error())
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Start(ctx, config.GRPC.Port); err != nil {
			logg.Error("failed to start grpc-server: " + err.Error())
			cancel()
		}
	}()

	logg.Info("waiting for all servers finished")
	wg.Wait()
	logg.Info("all servers are finished")

	return 0
}
