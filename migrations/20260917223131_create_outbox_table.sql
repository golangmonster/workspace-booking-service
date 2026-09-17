-- +goose Up
-- +goose StatementBegin
CREATE TABLE outbox (
    id BIGSERIAL PRIMARY KEY,
    topic  TEXT      NOT NULL,
    aggregate_id TEXT NOT NULL,
    message_value  TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE outbox;
-- +goose StatementEnd
