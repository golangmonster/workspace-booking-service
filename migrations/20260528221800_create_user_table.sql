-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR(20) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    phone VARCHAR(16) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX unique_login
    ON users (login)
    WHERE is_deleted = false;

CREATE UNIQUE INDEX unique_phone
    ON users (phone)
    WHERE is_deleted = false;

COMMENT ON TABLE users IS 'Пользователи';

COMMENT ON COLUMN users.login IS 'Логин пользователя в системе';
COMMENT ON COLUMN users.full_name IS 'Имя пользователя';
COMMENT ON COLUMN users.phone IS 'Телефон в формате +79001234567 ';
COMMENT ON COLUMN users.created_at IS 'Дата создания пользователя';
COMMENT ON COLUMN users.updated_at IS 'Дата последнего обновления пользователя';
COMMENT ON COLUMN users.is_deleted IS 'Флаг удаления';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
