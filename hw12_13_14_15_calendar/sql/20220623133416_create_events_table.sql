-- +goose Up
-- +goose StatementBegin
create table events
(
    id        varchar   not null    constraint events_pk    primary key,
    title     varchar   not null,
    date_time timestamp not null,
    duration  integer   not null,
    text      text      not null,
    user_id   integer   not null,
    remind    integer   not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events
-- +goose StatementEnd
