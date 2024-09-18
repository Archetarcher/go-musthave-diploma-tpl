-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_withdrawal
(
    id serial PRIMARY KEY,
    user_id INT NOT NULL,
    order_id BIGINT NOT NULL UNIQUE,
    amount DOUBLE PRECISION NOT NULL,

    created_at timestamp without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),

    CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES users(id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_withdrawal;
-- +goose StatementEnd
