-- +goose Up
-- +goose StatementBegin
create table users
(
    id            serial primary key,
    first_name    text        not null,
    last_name     text        not null,
    email         text        not null,
    oauth_id      text        not null,
    password_hash text        not null,
    password_salt text        not null,
    created_at    timestamptz not null default now(),
    updated_at    timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
