-- +goose Up
create table "user"
(
    user_id    serial primary key,
    name       text      not null,
    email      text      not null,
    role       text      not null,
    password   text      not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);


-- +goose Down
drop table "user";

