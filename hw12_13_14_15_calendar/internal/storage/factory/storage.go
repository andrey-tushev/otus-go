package factory

import (
	"context"
	"errors"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/app"
	conf "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

func GetStorage(config conf.Config) (app.Storage, error) {
	switch config.Storage.Storage {
	case "memory":
		return memorystorage.New(), nil

	case "sql":
		sqlStorage := sqlstorage.New(config.SQL.DSN)
		err := sqlStorage.Connect(context.Background())
		if err != nil {
			return nil, err
		}

		return sqlStorage, nil
	}

	return nil, errors.New("unknown storage")
}
