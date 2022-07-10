package sqlstorage

//nolint:gci
import (
	"context"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/storage"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib" //nolint:golint
	"github.com/jmoiron/sqlx"
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
	db, err := sqlx.Open("pgx", s.dsn)
	if err != nil {
		return err
	}
	s.db = db

	return nil
}

func (s *Storage) Exec(ctx context.Context, query string) error {
	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *Storage) Close(ctx context.Context) error {
	s.db.Close()

	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	id := uuid.New().String()
	event.ID = id
	query := `
		INSERT INTO events(id, title, date_time, duration, text, user_id, remind)
		VALUES (:id, :title, :date_time, :duration, :text, :user_id, :remind)
	`
	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	query := `
		UPDATE events
		SET 
		    title 		= :title, 
		    date_time 	= :date_time, 
		    duration 	= :duration, 
		    text 		= :text, 
		    user_id 	= :user_id, 
		    remind 		= :remind
		WHERE id = :id		
	`
	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	query := `
		DELETE FROM events		
		WHERE id = $1		
	`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	res := []storage.Event{}

	query := `SELECT * FROM events`
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := storage.Event{}
	for rows.Next() {
		if err := rows.StructScan(&list); err != nil {
			return nil, err
		}
		res = append(res, list)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
