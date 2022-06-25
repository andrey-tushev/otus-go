package app

import "time"

type Event struct {
	ID       string
	Title    string
	DateTime time.Time
	Duration int
	Text     string
	UserId   int
	Remind   int
}
