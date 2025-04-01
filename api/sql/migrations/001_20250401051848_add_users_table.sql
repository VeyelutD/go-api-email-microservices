-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id           BIGSERIAL PRIMARY KEY,
    email        TEXT                    NOT NULL,
    password     TEXT,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_confirmed BOOLEAN   DEFAULT FALSE NOT NULL,
    UNIQUE (email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    DROP TABLE users CASCADE;
-- +goose StatementEnd
