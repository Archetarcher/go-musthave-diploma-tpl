-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_accrual
(
    id serial PRIMARY KEY,
    user_id INT NOT NULL,
    order_id BIGINT NOT NULL UNIQUE,
    status VARCHAR(10) NOT NULL,
    amount BIGINT NOT NULL,

    uploaded_at timestamp without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
    processing_started_at timestamp without time zone  NULL,
    processed_at timestamp without time zone  NULL,
    invalidated_at timestamp without time zone  NULL,


    CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES users(id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_accrual;
-- +goose StatementEnd
