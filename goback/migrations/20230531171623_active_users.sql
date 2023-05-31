-- +goose Up
-- +goose StatementBegin
create table if not exists active_room_users
(
    id         serial8     not null primary key,
    room_id    text        not null,
    email      text        not null,
    connected  bool        not null,
    active     bool        not null,
    created_at timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists active_room_users;
-- +goose StatementEnd
