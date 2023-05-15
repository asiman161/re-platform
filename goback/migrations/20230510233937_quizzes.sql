-- +goose Up
-- +goose StatementBegin
create table if not exists quizzes
(
    id         serial4 primary key,
    room_id    text        not null,
    author     text        not null,
    content    text        not null,
    variants   jsonb       not null default '[]',
    answers    jsonb       not null default '[]',
    is_open    bool        not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table quizzes;
-- +goose StatementEnd
