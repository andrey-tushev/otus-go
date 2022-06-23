package storage

import "time"

type Event struct {
	ID       string    `db:"id"`
	Title    string    `db:"title"`
	DateTime time.Time `db:"date_time"`
	Duration int       `db:"duration"` // хотел использовать time.Duration, но sqlx не умеет конвертировать такое :-(
	Text     string    `db:"text"`
	UserId   int       `db:"user_id"`
	Remind   int       `db:"remind"`
}
