-- +goose Up
-- +goose StatementBegin
create table users
(
    id            serial primary key,
    username      text      not null,
    password_hash text      not null,
    password_salt text      not null,
    created_at    timestamp not null default now(),
    updated_at    timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
