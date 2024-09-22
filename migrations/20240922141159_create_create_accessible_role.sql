-- +goose Up
-- +goose StatementBegin
create table accessible_role
(
    accessible_role_id bigserial primary key,
    role               text      not null,
    endpoint_address   text      not null,
    created_at         timestamp not null default now(),
    CONSTRAINT accessible_role_unique UNIQUE (role, endpoint_address)
);

INSERT INTO accessible_role (role, endpoint_address) VALUES ('admin', '/chat_v1.ChatV1/Create');
INSERT INTO accessible_role (role, endpoint_address) VALUES ('admin', '/chat_v1.ChatV1/Delete');
INSERT INTO accessible_role (role, endpoint_address) VALUES ('admin', '/chat_v1.ChatV1/Get');
INSERT INTO accessible_role (role, endpoint_address) VALUES ('user', '/chat_v1.ChatV1/Get');
INSERT INTO accessible_role (role, endpoint_address) VALUES ('admin', '/chat_v1.ChatV1/SendMessage');
INSERT INTO accessible_role (role, endpoint_address) VALUES ('user', '/chat_v1.ChatV1/SendMessage');
INSERT INTO accessible_role (role, endpoint_address) VALUES ('admin', '/chat_v1.ChatV1/GetMessages');
INSERT INTO accessible_role (role, endpoint_address) VALUES ('user', '/chat_v1.ChatV1/GetMessages');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table accessible_role;
-- +goose StatementEnd
