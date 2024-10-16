-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id serial PRIMARY KEY,
    login VARCHAR (100) NOT NULL unique,
    hash VARCHAR (100) NOT NULL,
    balance DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at timestamp without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
    updated_at timestamp without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC')
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
