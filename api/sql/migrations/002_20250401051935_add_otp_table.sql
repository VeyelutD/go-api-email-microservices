-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_login_codes
(
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT NOT NULL,
    code       TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (email) REFERENCES users (email) ON DELETE CASCADE,
    UNIQUE (email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    DROP TABLE user_login_codes;
-- +goose StatementEnd
