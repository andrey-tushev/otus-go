package sqlstorage

import (
	"context"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/andrey-tushev/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	dsn string
	db  *sqlx.DB
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("pgx", "postgres://calendar:calendar@localhost/calendar")

	if err != nil {
		return err
	}
	s.db = db

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.db.Close()

	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	//TODO implement me
	panic("implement me")
}
