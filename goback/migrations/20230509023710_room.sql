-- +goose Up
-- +goose StatementBegin
create table if not exists rooms
(
    id         serial4 primary key,
    name text not null,
    author     text        not null,
    is_open    bool        not null default false,
    updated_at timestamptz not null default now(),
    created_at timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table rooms;
-- +goose StatementEnd
