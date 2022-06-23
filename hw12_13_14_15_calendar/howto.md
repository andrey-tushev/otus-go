
goose postgres "postgres://calendar:calendar@localhost/calendar?sslmode=disable" status

goose create add_some_column sql