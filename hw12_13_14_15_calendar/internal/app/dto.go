package app

import "time"

type Event struct {
	ID       string
	Title    string
	DateTime time.Time
	Duration time.Duration
	Text     string
	UserId   string
	Remind   time.Duration
}
