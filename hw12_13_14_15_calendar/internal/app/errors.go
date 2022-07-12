package app

import "errors"

var (
	ErrDateBusy = errors.New("date busy")
	ErrNotFound = errors.New("event not found")
)
