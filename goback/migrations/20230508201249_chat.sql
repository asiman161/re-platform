-- +goose Up
-- +goose StatementBegin
create table if not exists chat
(
    id         serial4 primary key,
    chat_id    text not null,
    content    text not null,
    author text not null,
    created_at timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chat;
-- +goose StatementEnd
