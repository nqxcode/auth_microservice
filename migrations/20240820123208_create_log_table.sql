-- +goose Up
-- +goose StatementBegin
create table log
(
    log_id     bigserial primary key,
    message    text      not null,
    payload    jsonb,
    ip         text      not null,
    created_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table log;
-- +goose StatementEnd
