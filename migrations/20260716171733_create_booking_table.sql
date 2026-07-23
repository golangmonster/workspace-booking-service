-- +goose Up
-- +goose StatementBegin
CREATE TABLE booking (
    id BIGSERIAL PRIMARY KEY,
    workspace_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX booking_unique_user_workspace_idx ON booking(user_id, workspace_id) WHERE status <> 'CANCELLED';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE booking;
-- +goose StatementEnd
