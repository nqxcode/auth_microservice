-- +goose Up
-- +goose StatementBegin
alter table "user" add constraint email_unique UNIQUE (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table "user" drop constraint email_unique;
-- +goose StatementEnd
